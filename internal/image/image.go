package image

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/jzh18/baffs/internal/util"
)

/*
Create an overlay layer.
This layer is in memory, not in the filesystem.
Until Dump() is called, the layer will be created in the filesystem.
{overlay_root_dir}/

	{layer_name}/
		diff/
		link: l/link_{layer_name}
		lower
		real/
	l/
		link_{layer_name} -> ../{layer_name}/diff
*/
type Layer struct {
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

func (l *Layer) GetLayerSize() int64 {
	return util.GetDirSize(l.diffPath)
}

// Truncate diff layer
func (l *Layer) TruncateDiff() {
	// this will remove the diff layer too
	if err := os.RemoveAll(l.diffPath); err != nil {
		panic(err)
	}
	// create an empty diff layer
	if err := os.Mkdir(l.diffPath, 0755); err != nil {
		panic(err)
	}
}

// archive files under diff directory
// func (l *Layer) TarDiff(destFile string) {
// 	TarFiles(l.absolute_diff_path, destFile)
// }

// Construct a shadow layer in memory, but not create in the filesystem
func (l *Layer) Shadow() ShadowLayer {
	layerName := "shadow_" + l.layerName
	parentPath := l.layerPath[:len(l.layerPath)-len(l.layerName)]
	shadow := ShadowLayer{
		Layer: Layer{
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

type ShadowLayer struct {
	Layer
	realPath string
}

// SetLowers replace existing lowers to new lowers
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

func (l *ShadowLayer) DumpLayerSize(size string) {
	if err := os.WriteFile(l.sizePath, []byte(size), 0600); err != nil {
		panic(err)
	}
}

// Construct an original layer from a shadow layer, not create it in the filesystem
func (l *ShadowLayer) Original() Layer {
	layerName := l.layerName[len("shadow_"):]
	parentPath := l.layerPath[:len(l.layerPath)-len(l.layerName)]
	original := NewLayer(filepath.Join(parentPath, layerName))
	original.layerName = layerName
	original.metaPath = l.metaPath
	original.cacheidPath = l.cacheidPath
	original.sizePath = l.sizePath

	original.lLinkPath = filepath.Join(l.lLinkPath[:len(l.lLinkPath)-len(l.linkContent)], original.linkContent)
	return *original
}

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

// Create dir and files according to the layer
// Only shadow layers can call this function
func (l *ShadowLayer) Dump() {
	var mode fs.FileMode = 0755
	// create layer
	if err := os.Mkdir(l.layerPath, mode); err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("Layer diff dir already exists")
		} else {
			panic(err)
		}
	}

	// create diff dir
	if err := os.Mkdir(l.diffPath, mode); err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("Layer diff dir already exists")
		} else {
			panic(err)
		}
	}
	//create link file
	linkFile, err := os.Create(l.linkPath)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("Layer link file already exists")
		} else {
			panic(err)
		}
	}
	linkFile.WriteString(l.linkContent)

	// create lower file, optionally
	if l.lowerPath != "" {
		if lower_file, err := os.Create(l.lowerPath); err != nil {
			if errors.Is(err, os.ErrExist) {
				fmt.Println("Layer lower file already exists")
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
				fmt.Println("Layer real dir already exists")
			} else {
				panic(err)
			}
		}
	}

	// create l/link_{layer_name}
	if err := os.Symlink(filepath.Join("../", l.layerName, "diff"), l.lLinkPath); err != nil {
		if errors.Is(err, os.ErrExist) {
			fmt.Println("l/link file already exists")
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

func NewLayer(layerPath string) *Layer {
	var layer Layer
	layer.layerPath = layerPath
	absDiffPath := filepath.Join(layerPath, "diff")
	if !util.PathExist(absDiffPath) {
		panic("Diff dir not exist: " + absDiffPath)
	}
	layer.diffPath = absDiffPath

	absLinkPath := filepath.Join(layerPath, "link")
	if !util.PathExist(absLinkPath) {
		panic("Link file not exist: " + absLinkPath)
	}
	layer.linkPath = absLinkPath

	linkContent, err := os.ReadFile(absLinkPath)
	if err != nil {
		panic(err)
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
	return &layer
}

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
func ExtractLayers(imgInfo *types.ImageInspect, overlayPath string, dockerRootDir string) []Layer {
	layerNames := extractLayerNames(imgInfo.GraphDriver)
	// there should only return image inspect dirs
	// after modified, frist test cold run, then warm run
	allOriginalLayers := []Layer{}
	for _, p := range layerNames {
		l := NewLayer(filepath.Join(overlayPath, p))
		l.layerName = p
		l.linkPath = filepath.Join(overlayPath, "l", l.linkContent)
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
			if allOriginalLayers[i].layerPath == expectedAbsDir {
				break
			}
		}
		allOriginalLayers[i].metaPath = layerDir
		allOriginalLayers[i].cacheidPath = cacheIdDir
		allOriginalLayers[i].cacheid = cacheIdStr
		allOriginalLayers[i].sizePath = filepath.Join(layerDir, "size")
		if size, err := os.ReadFile(allOriginalLayers[i].sizePath); err != nil {
			panic(err)
		} else {
			allOriginalLayers[i].size = string(size)
		}

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
