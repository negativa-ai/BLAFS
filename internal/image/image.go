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
package image

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/negativa-ai/BLAFS/internal/util"
	log "github.com/sirupsen/logrus"
)

/*
A LayerInfo represents the file structure of an original overlay layer of a docker image.
An original overlay layer is organized as follows:
{overlay_root_dir}/
	{layer_name}/
		diff/
		link: l/link_{layer_name}
		lower
		real/
	l/
		link_{layer_name} -> ../{layer_name}/diff

All members of LayerInfo are absolute paths.
None of the methods of LayerInfo change the filesystem.
*/

type LayerInfo struct {
	// files under /var/lib/docker/overlay2
	// all path should be absolute path
	layerPath    string
	diffPath     string
	linkPath     string
	lowerPath    string
	linkContent  string
	lowerContent string
	lLinkPath    string
	layerName    string
	// files under /var/lib/docker/image/overlay2/layerdb/sha256
	metaPath    string
	cacheidPath string
	sizePath    string
	cacheid     string
	size        string
}

func (l *LayerInfo) GetLayerSize() int64 {
	return util.GetDirSize(l.diffPath)
}

func (l *LayerInfo) GetLayerPath() string {
	return l.layerPath
}

func (l *LayerInfo) GetDiffPath() string {
	return l.diffPath
}

func (l *LayerInfo) GetlinkContent() string {
	return l.linkContent
}

func (l *LayerInfo) SetMetaPath(metaPath string) {
	l.metaPath = metaPath
}

func (l *LayerInfo) SetCacheIdPath(cacheidPath string) {
	l.cacheidPath = cacheidPath
}

func (l *LayerInfo) SetCacheId(cacheId string) {
	l.cacheid = cacheId
}

// Set size path, this will also set the size field
func (l *LayerInfo) SetSizePath(sizePath string) {
	l.sizePath = sizePath
	if size, err := os.ReadFile(sizePath); err != nil {
		panic(err)
	} else {
		l.size = string(size)
	}
}

// NewLayerInfo creates a new LayerInfo from a layer path, with some very basic fields
func NewLayerInfo(layerPath string) *LayerInfo {
	var layer LayerInfo
	layer.layerPath = layerPath
	absDiffPath := filepath.Join(layerPath, "diff")
	if !util.PathExist(absDiffPath) {
		log.Error("Diff dir not exist: ", absDiffPath)
	}
	layer.diffPath = absDiffPath

	absLinkPath := filepath.Join(layerPath, "link")
	if !util.PathExist(absLinkPath) {
		log.Error("Link file not exist: ", absLinkPath)
	}
	layer.linkPath = absLinkPath

	linkContent, err := os.ReadFile(absLinkPath)
	if err != nil {
		log.Error("Read link file failed: ", absLinkPath)
	}
	layer.linkContent = string(linkContent)

	absLowerPath := filepath.Join(layerPath, "lower")
	if !util.PathExist(absLowerPath) {
		layer.lowerPath = ""
		layer.lowerContent = ""
	} else {
		layer.lowerPath = absLowerPath
		lowerData, err := os.ReadFile(absLowerPath)
		if err != nil {
			panic(err)
		}
		layer.lowerContent = string(lowerData)
	}

	layer.layerName = filepath.Base(layerPath)
	layer.lLinkPath = filepath.Join(filepath.Dir(layerPath), "l", layer.linkContent)

	return &layer
}

// An OriginalLayer represents an original overlay layer of a docker image.
type OriginalLayer struct {
	LayerInfo
}

// Shadow creates a ShadowLayer from an original layer in memory, not create it in the filesystem.
// It is the counterpart of ShadowLayer.Original.
func (l *OriginalLayer) Shadow() ShadowLayer {
	layerName := "shadow_" + l.layerName
	parentPath := l.layerPath[:len(l.layerPath)-len(l.layerName)]
	shadow := ShadowLayer{
		LayerInfo: LayerInfo{
			layerPath:    filepath.Join(parentPath, layerName),
			diffPath:     filepath.Join(parentPath, layerName, "diff"),
			linkPath:     filepath.Join(parentPath, layerName, "link"),
			linkContent:  "shadow_" + l.linkContent,
			lowerContent: "",
			layerName:    layerName,
			metaPath:     l.metaPath,
			cacheidPath:  l.cacheidPath,
			sizePath:     l.sizePath,
			cacheid:      layerName,
			size:         l.size,
		},
		realPath: filepath.Join(parentPath, layerName, "real"),
	}

	// lower file is optional if it's the bottom layer
	if l.lowerPath == "" {
		shadow.lowerPath = ""

	} else {
		shadow.lowerPath = filepath.Join(parentPath, layerName, "lower")
	}

	shadow.lLinkPath = filepath.Join(l.lLinkPath[:len(l.lLinkPath)-len(l.linkContent)], shadow.linkContent)

	return shadow
}

// NewOriginalLayer creates a new OriginalLayer from a LayerInfo
func NewOriginalLayer(layerInfo LayerInfo) OriginalLayer {
	return OriginalLayer{LayerInfo: layerInfo}
}

// A ShadowLayer represents a shadowed overlay layer of a docker image.
// Each original layer will have a corresponding shadow layer.
// Files used in an original layer will be copied to the shadow layer.
type ShadowLayer struct {
	LayerInfo
	realPath string
}

func (l *ShadowLayer) GetRealPath() string {
	return l.realPath
}

// SetLowers replaces existing lowers to new lowers
func (l *ShadowLayer) SetLowers(newLowers string) {
	l.lowerContent = newLowers
}

func (l *ShadowLayer) SetLayerSize(size string) {
	l.size = size
}

func (l *ShadowLayer) SetCacheId(cacheId string) {
	l.cacheid = cacheId
}

func (l *ShadowLayer) GetLayerSize() int64 {
	return util.GetDirSize(l.realPath)
}

// DumpLayerSize writes the size of the layer to the size file.
// It modifies the filesystem.
func (l *ShadowLayer) DumpLayerSize(size string) {
	if err := os.WriteFile(l.sizePath, []byte(size), 0600); err != nil {
		panic(err)
	}
}

// TarDiff archives the diff directory of the shadow layer to a tar file.
func (l *ShadowLayer) TarDiff(destFile string) {
	util.TarFiles(l.diffPath, destFile)
}

// Original returns the original layer from a shadow layer in memory, not create it in the filesystem
// It is the counterpart of OriginalLayer.Shadow.
func (l *ShadowLayer) Original() OriginalLayer {
	layerName := l.layerName[len("shadow_"):]
	parentPath := l.layerPath[:len(l.layerPath)-len(l.layerName)]
	original := NewLayerInfo(filepath.Join(parentPath, layerName))
	original.layerName = layerName
	original.metaPath = l.metaPath
	original.cacheidPath = l.cacheidPath
	original.sizePath = l.sizePath

	original.lLinkPath = filepath.Join(l.lLinkPath[:len(l.lLinkPath)-len(l.linkContent)], original.linkContent)
	return OriginalLayer{LayerInfo: *original}
}

// Restore restores the filesystem of the docker image to the state before shadowing.
// It modifies the filesystem.
// It is the counterpart of Dump.
func (l *ShadowLayer) Restore() {
	original := l.Original()
	bakCacheIdData, err := os.ReadFile(original.cacheidPath + ".bak")
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(original.cacheidPath, bakCacheIdData, 0600); err != nil {
		panic(err)
	}
}

// Dump creates the shadow layer in the filesystem.
// It is the counterpart of Restore.
func (l *ShadowLayer) Dump() {
	var mode fs.FileMode = 0755
	// create layer
	if err := os.Mkdir(l.layerPath, mode); err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Debug("Layer diff dir already exists")
		} else {
			panic(err)
		}
	}

	// create diff dir
	if err := os.Mkdir(l.diffPath, mode); err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Debug("Layer diff dir already exists")
		} else {
			panic(err)
		}
	}
	//create link file
	linkFile, err := os.Create(l.linkPath)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Debug("Layer link file already exists")
		} else {
			panic(err)
		}
	}
	linkFile.WriteString(l.linkContent)

	// create lower file, optionally
	if l.lowerPath != "" {
		if lower_file, err := os.Create(l.lowerPath); err != nil {
			if errors.Is(err, os.ErrExist) {
				log.Debug("Layer lower file already exists")
			} else {
				panic(err)
			}
		} else {
			// set lowers
			lower_file.WriteString(l.lowerContent)
		}
	}

	// create real dir
	if l.realPath != "" {
		if err := os.Mkdir(l.realPath, mode); err != nil {
			if errors.Is(err, os.ErrExist) {
				log.Debug("Layer real dir already exists")
			} else {
				panic(err)
			}
		}
	}

	// create l/link_{layer_name}
	if err := os.Symlink(filepath.Join("../", l.layerName, "diff"), l.lLinkPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			log.Debug("l/link file already exists")
		} else {
			panic(err)
		}
	}

	if l.cacheidPath != "" {
		// set cache-id
		// back up the original cache-id firstly
		if !util.PathExist(l.cacheidPath + ".bak") {
			util.CopyFile(l.cacheidPath, l.cacheidPath+".bak")
		}
		if err := os.WriteFile(l.cacheidPath, []byte(l.cacheid), 0600); err != nil {
			panic(err)
		}
	}
}

// NewShadowLayer creates a new ShadowLayer from a LayerInfo
func NewShadowLayer(layerInfo LayerInfo) ShadowLayer {
	return ShadowLayer{
		LayerInfo: layerInfo,
		realPath:  filepath.Join(layerInfo.layerPath, "real"),
	}

}

// An ImgTarLayer represents a layer of a docker image in a tar file.
type ImgTarLayer struct {
	versionPath  string
	layerTarPath string
	jsonPath     string
}

func (l *ImgTarLayer) RmLayerTar() {
	if err := os.Remove(l.layerTarPath); err != nil {
		panic(err)
	}
}

// LayerTarSha256Sum returns the sha256 sum of the layer tar file.
func (l *ImgTarLayer) LayerTarSha256Sum() string {
	s, err := util.Sha256Sum(l.layerTarPath)
	if err != nil {
		panic(err)
	}
	return s
}

// GetLayerTarPath returns the path of the layer tar file.
func (l *ImgTarLayer) GetLayerTarPath() string {
	return l.layerTarPath
}

// we only care about Rootfs, so we set other fileds as interface{}

// An ImgJson represents the json file of a docker image in a tar file.
// See: https://github.com/moby/moby/blob/master/image/spec/v1.2.md
type ImgJson struct {
	Created      interface{} `json:"created"`
	Author       interface{} `json:"author"`
	Architecture interface{} `json:"architecture"`
	Os           interface{} `json:"os"`
	Config       interface{} `json:"config"`
	Rootfs       struct {
		DiffIds []string `json:"diff_ids"` // from bottom to top
		Type    string   `json:"type"`
	} `json:"rootfs"`
	History interface{} `json:"history"`
}

type ImgTarFs struct {
	basePath        string
	layers          []ImgTarLayer // from bottom to top
	manifestPath    string
	repoPath        string
	imgJsonPath     string
	manifestContent []Manifest
	imgJsonContent  ImgJson
}

type Manifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"` // 0->n : bottom->top
}

func (f *ImgTarFs) GetLayers() []ImgTarLayer {
	return f.layers
}

func (f *ImgTarFs) DumpImgJson() {
	data, err := json.Marshal(f.imgJsonContent)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(f.imgJsonPath, data, 0755)
	if err != nil {
		panic(err)
	}
}

func (f *ImgTarFs) GetImageJson() ImgJson {
	return f.imgJsonContent
}

func (f *ImgTarFs) GetManifest() []Manifest {
	return f.manifestContent
}

func (f *ImgTarFs) DumpManifest() {
	data, err := json.Marshal(f.manifestContent)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(f.manifestPath, data, 0755)
	if err != nil {
		panic(err)
	}
}

// TarWholeFs archives the whole filesystem of the image to a tar file.
func (f *ImgTarFs) TarWholeFs(dst string) {
	util.TarFiles(f.basePath, dst)
}

// ParseImgTarFs parses the filesystem of a docker image in a tar file.
func ParseImgTarFs(path string) ImgTarFs {
	imgTarFs := ImgTarFs{
		basePath: path,
	}
	imgTarFs.manifestPath = filepath.Join(path, "manifest.json")
	imgTarFs.repoPath = filepath.Join(path, "repositories")

	manifestFile, err := os.Open(imgTarFs.manifestPath)
	if err != nil {
		panic(err)
	}
	defer manifestFile.Close()

	var mf []Manifest
	err = json.NewDecoder(manifestFile).Decode(&mf)
	if err != nil {
		panic(err)
	}
	imgTarFs.manifestContent = mf

	manifestEle := imgTarFs.manifestContent[0]

	imgTarFs.imgJsonPath = filepath.Join(path, manifestEle.Config)
	imgJsonFile, err := os.Open(imgTarFs.imgJsonPath)
	if err != nil {
		panic(err)
	}
	defer imgJsonFile.Close()

	var imgJson ImgJson
	err = json.NewDecoder(imgJsonFile).Decode(&imgJson)
	if err != nil {
		panic(err)
	}
	imgTarFs.imgJsonContent = imgJson

	for _, layerTar := range manifestEle.Layers {
		imgTarLayer := ImgTarLayer{}
		layer := layerTar[0 : len(layerTar)-len("layer.tar")]
		imgTarLayer.jsonPath = filepath.Join(path, layer, "json")
		imgTarLayer.layerTarPath = filepath.Join(path, layerTar)
		imgTarLayer.versionPath = filepath.Join(path, layer, "VERSION")
		imgTarFs.layers = append(imgTarFs.layers, imgTarLayer)
	}
	return imgTarFs
}
