package builder

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/jzh18/baffs/internal/image"
)

// Extract layer name from graph driver, from top to bottom
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

// https://www.baeldung.com/linux/docker-image-storage-host
func generateChainId(preChainId string, diffId string) string {
	str := preChainId + " " + diffId
	chainId := fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
	return chainId

}

// extract layers from image inspect, from top to bottom
func ExtractLayers(imgInfo *types.ImageInspect, overlayPath string, dockerRootDir string) []image.Layer {
	layerNames := extractLayerNames(imgInfo.GraphDriver)
	// there should only return image inspect dirs
	// after modified, frist test cold run, then warm run
	allOriginalLayers := []image.Layer{}
	for _, p := range layerNames {
		l := image.NewLayer(filepath.Join(overlayPath, p))
		allOriginalLayers = append(allOriginalLayers, *l)
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
		for ; i < len(allOriginalLayers); i++ {
			if allOriginalLayers[i].GetLayerPath() == expectedAbsDir {
				break
			}
		}

		allOriginalLayers[i].SetMetaPath(layerDir)
		allOriginalLayers[i].SetCacheIdPath(cacheIdDir)
		allOriginalLayers[i].SetCacheId(cacheIdStr)
		allOriginalLayers[i].SetSizePath(filepath.Join(layerDir, "size"))

		if count >= len(rootfsLayers) {
			break
		}
		// generate new chain_id
		diffId := rootfsLayers[count]
		chainId = generateChainId(chainId, diffId)
		chainId = "sha256:" + chainId
		count++
	}

	return allOriginalLayers
}
