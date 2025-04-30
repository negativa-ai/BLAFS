// MIT License

// Copyright (c) [2025] [jzh18]

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package builder

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/negativa-ai/BLAFS/internal/image"
	"github.com/negativa-ai/BLAFS/internal/mount"
	"github.com/negativa-ai/BLAFS/internal/util"
	log "github.com/sirupsen/logrus"
)

// extractLayerNames extracts layer names from graph driver.
// It returns a list of layer names.
func extractLayerNames(graphDriver types.GraphDriverData) []string {
	var allLayers []string
	upper := graphDriver.Data["UpperDir"]
	allLayers = append(allLayers, upper)
	if lower, ok := graphDriver.Data["LowerDir"]; ok {
		allLowers := strings.Split(lower, ":")
		allLayers = append(allLayers, allLowers...)
	}

	var layerNames []string
	for _, layer := range allLayers {
		tmp := strings.Split(layer, "/")
		layerName := tmp[len(tmp)-2]
		layerNames = append(layerNames, layerName)
	}
	return layerNames
}

// generateChainId generates a chain id from previous chain id and diff id.
// See https://www.baeldung.com/linux/docker-image-storage-host
// It returns a chain id.
func generateChainId(preChainId string, diffId string) string {
	str := preChainId + " " + diffId
	chainId := fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
	return chainId

}

// ExtractLayersInfo extracts layer info from image inspect, from top to bottom.
// It returns a list of layer info.
func ExtractLayersInfo(imgInfo *types.ImageInspect, overlayPath string, dockerRootDir string) []image.LayerInfo {
	layerNames := extractLayerNames(imgInfo.GraphDriver)
	// there should only return image inspect dirs
	// after modified, frist test cold run, then warm run
	layerInfos := []image.LayerInfo{}
	for _, p := range layerNames {
		l := image.NewLayerInfo(filepath.Join(overlayPath, p))
		layerInfos = append(layerInfos, *l)
	}

	rootfsLayers := imgInfo.RootFS.Layers
	chainId := rootfsLayers[0]

	// generate layer meta info for each layer
	count := 1
	for {
		dirName := strings.Split(chainId, ":")[1]
		// read layer size
		layerDir := filepath.Join(dockerRootDir, "image/overlay2/layerdb/sha256", dirName)
		cacheIdDir := filepath.Join(layerDir, "cache-id")
		cacheId, err := os.ReadFile(cacheIdDir)
		if err != nil {
			panic(err)
		}
		cacheIdStr := string(cacheId)
		expectedAbsDir := filepath.Join(overlayPath, cacheIdStr)

		// find corresponsding original layer
		i := 0
		for ; i < len(layerInfos); i++ {
			log.Debug("layer path: ", layerInfos[i].GetLayerPath(), " expected path: ", expectedAbsDir)
			if layerInfos[i].GetLayerPath() == expectedAbsDir {
				break
			}
		}

		layerInfos[i].SetMetaPath(layerDir)
		layerInfos[i].SetCacheIdPath(cacheIdDir)
		layerInfos[i].SetCacheId(cacheIdStr)
		layerInfos[i].SetSizePath(filepath.Join(layerDir, "size"))

		if count >= len(rootfsLayers) {
			break
		}
		// generate new chain_id
		diffId := rootfsLayers[count]
		chainId = generateChainId(chainId, diffId)
		chainId = "sha256:" + chainId
		count++
	}
	return layerInfos
}

// checkIfShadowed checks if the image is already shadowed.
func checkIfShadowed(graphDriver types.GraphDriverData) bool {
	return strings.Contains(graphDriver.Data["UpperDir"], "shadow")
}

// generateTarFileName generates a tar file name from image name.
func generateTarFileName(imgName string) string {
	s := strings.ReplaceAll(imgName, ":", "_")
	s = strings.ReplaceAll(s, "/", "_") + ".tar"
	return s

}

// generateLowers generates lowers for each layer.
// `layers[0]` is the top layer.
// It returns a list of lowers.
func generateLowers(layers []image.ShadowLayer) []string {
	allLowers := []string{""}
	curLower := "l/" + layers[len(layers)-1].GetlinkContent()
	for i := len(layers) - 2; i >= 0; i-- {
		allLowers = append([]string{curLower}, allLowers...)
		curLower = "l/" + layers[i].GetlinkContent() + ":" + curLower
	}
	return allLowers
}

// saveImage saves the original image to a tar file.
func saveImage(workDir string, cli *client.Client, ctx *context.Context, imgName string) {
	// bak up the original image
	reader, err := cli.ImageSave(*ctx, []string{imgName})
	if err != nil {
		panic(err)
	}
	imgTarPath := filepath.Join(workDir, generateTarFileName(imgName))
	log.Debug("original img backup path: ", imgTarPath)
	out, err := os.OpenFile(imgTarPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(out, reader)
	if err != nil {
		panic(err)
	}
	reader.Close()
}

// ShadowImage shadows the image. For each original layer, it creates a shadow layer in memory.
// It does not create anything on the filesystem.
// It returns if shadowed, original layers, shadow layers.
func ShadowImage(imgName string, workDir string, overlayPath string,
	dockerRootDir string, cli *client.Client, ctx *context.Context, optimize string) (bool, []image.OriginalLayer, []image.ShadowLayer) {
	imgInfo, _, err := cli.ImageInspectWithRaw(*ctx, imgName)
	if err != nil {
		panic(err)
	}
	var originalLayers []image.OriginalLayer
	var shadowLayers []image.ShadowLayer
	layerInfos := ExtractLayersInfo(&imgInfo, overlayPath, dockerRootDir)
	shadowed := checkIfShadowed(imgInfo.GraphDriver)
	if !shadowed {
		log.Debug("shadowing container")
		saveImage(workDir, cli, ctx, imgName)

		for _, l := range layerInfos {
			originalLayer := image.OriginalLayer{LayerInfo: l}
			originalLayers = append(originalLayers, originalLayer)
			shadowLayers = append(shadowLayers, originalLayer.Shadow())
		}
		allShadowLowers := generateLowers(shadowLayers)
		// we should not set lowers for the bottom layer
		for i := 0; i < len(shadowLayers)-1; i++ {
			shadowLayers[i].SetLowers(allShadowLowers[i])
		}
	} else {
		for _, l := range layerInfos {
			shadowLayer := image.NewShadowLayer(l)
			shadowLayers = append(shadowLayers, shadowLayer)
			originalLayers = append(originalLayers, shadowLayer.Original())
		}
	}
	log.Debug("total layers: ", len(originalLayers))
	return shadowed, originalLayers, shadowLayers
}

// createMount creates a mount in memory abstraction, not create anything on the filesystem.
func createMount(fsExePath string, originalLaye image.OriginalLayer, shadowLayer image.ShadowLayer) mount.Mount {
	mountPoint := shadowLayer.GetDiffPath()
	kvArgs := map[string]string{
		"--realdir":  shadowLayer.GetRealPath(),
		"--lowerdir": originalLaye.GetDiffPath(),
		"--optimize": "",
	}
	flagArgs := []string{"-s"} // -s for silent output
	return mount.NewMount(fsExePath, mountPoint, kvArgs, flagArgs)
}

// CreateMounts creates mounts for each layer. It does not create anything on the filesystem.
func CreateMounts(fsExePath string, originalLayers []image.OriginalLayer, shadowLayers []image.ShadowLayer) []mount.Mount {
	var mounts []mount.Mount
	for i := 0; i < len(originalLayers); i++ {
		mounts = append(mounts, createMount(fsExePath, originalLayers[i], shadowLayers[i]))
	}
	return mounts
}

func umount(mountPoint string, mountType string) {
	cmd := exec.Command("umount", "-f", "-A", "-t", mountType, mountPoint)
	log.Debug("umount layer: ", cmd)
	cmd.Output()
}

func umountAllLayers(graphDriver types.GraphDriverData, mountType string) {
	allLowers := strings.Split(graphDriver.Data["LowerDir"], ":")
	// umount all lower shadow layers
	for _, lower := range allLowers {
		umount(lower, mountType)
	}
	shadowDiffPath := graphDriver.Data["UpperDir"]
	umount(shadowDiffPath, mountType)

}

// ExportImg exports the debloated image to a tar file.
// It returns if exported, target tar path, shadow layers.
func ExportImg(imgName string, workDir string, overlayPath string, dockerRootDir string, cli *client.Client, ctx *context.Context, topN int) (bool, string, []image.ShadowLayer) {
	imgInfo, _, err := cli.ImageInspectWithRaw(*ctx, imgName)
	if err != nil {
		panic(err)
	}

	if !checkIfShadowed(imgInfo.GraphDriver) {
		log.Info("Container not shadowed, cannot perform debloating")
		return false, "", make([]image.ShadowLayer, 0)
	}

	// export original image to reuse the structure
	tarName := generateTarFileName(imgName)
	originalImgPath := filepath.Join(workDir, tarName)
	untarPath := filepath.Join("/tmp", tarName)
	if err := os.Mkdir(untarPath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}
	cmd := exec.Command("tar", "-xf", originalImgPath, "-C", untarPath)
	log.Debug("untar image tar file: ", cmd)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	imgsTarFs := image.ParseImgTarFs(untarPath)
	for _, l := range imgsTarFs.GetLayers() {
		log.Debug("layer tar path: ", l.GetLayerTarPath())
		l.RmLayerTar()
	}

	// copy file from real path to diff path
	layerInfos := ExtractLayersInfo(&imgInfo, overlayPath, dockerRootDir)
	shadowLayers := []image.ShadowLayer{}
	for _, l := range layerInfos {
		shadowLayer := image.NewShadowLayer(l)
		shadowLayers = append(shadowLayers, shadowLayer)
	}
	umountAllLayers(imgInfo.GraphDriver, "fuse.debloated_fs")
	time.Sleep(1 * time.Second)
	log.Debug("Total layers: ", len(shadowLayers))
	for _, l := range shadowLayers {
		if !util.PathExist(l.GetRealPath()) {
			log.Debug("real path not exist, this layer might already be exported: ", l.GetRealPath())
			continue
		} else {
			if err := os.RemoveAll(l.GetDiffPath()); err != nil {
				panic(err)
			} else {
				if err := util.Move(l.GetRealPath(), l.GetDiffPath()); err != nil {
					log.Error("failed to move from: ", l.GetRealPath(), " to: ", l.GetDiffPath())
					panic(err)
				}
			}
		}

	}

	// tar diff file to untarpath & update layer diff ids
	if len(shadowLayers) != len(imgsTarFs.GetLayers()) {
		panic("number of shadow layers should be equal to img tar fs layers")
	}
	layerLen := len(shadowLayers)
	count := 0
	for i := 0; i < len(shadowLayers); i++ {
		shadow := shadowLayers[i]
		tarFsLayer := imgsTarFs.GetLayers()[layerLen-1-i]
		shadow.TarDiff(tarFsLayer.GetLayerTarPath())
		imgsTarFs.GetImageJson().Rootfs.DiffIds[layerLen-1-i] = "sha256:" + tarFsLayer.LayerTarSha256Sum()
		count++
		if topN != -1 && count >= topN {
			log.Debug("Only export top ", topN, " layers.")
			break
		}
	}
	imgsTarFs.DumpImgJson()

	// set tag
	imgsTarFs.GetManifest()[0].RepoTags[0] = imgsTarFs.GetManifest()[0].RepoTags[0] + "-baffs"
	imgsTarFs.DumpManifest()

	// tar the image fs
	targetTarPath := filepath.Join("/tmp/", tarName+".debloated")
	log.Debug("target tar path: ", targetTarPath)
	imgsTarFs.TarWholeFs(targetTarPath)

	return true, targetTarPath, shadowLayers
}

// LoadImage loads the generated image tar file.
func LoadImage(imgTarPath string, cli *client.Client) {
	// load the generated image tar file
	imageFile, err := os.Open(imgTarPath)
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()
	imageResponse, err := cli.ImageLoad(context.Background(), imageFile, false)
	if err != nil {
		panic(err)
	}
	defer imageResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageResponse.Body)
	if err != nil {
		panic(err)
	}
}
