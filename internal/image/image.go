package image

import (
	"crypto/sha256"
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
	realPath     string
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

// SetLowers replace existing lowers to new lowers
func (l *Layer) SetLowers(newLowers string) {
	l.lowerContent = newLowers
}

func (l *Layer) SetLayerSize(size string) {
	l.size = size
}

func (l *Layer) SetCacheId(cacheId string) {
	l.cacheid = cacheId
}

func (l *Layer) GetLayerSize() int64 {
	var size int64
	var path string

	// the files exist either absolute_real_path(shadow layer) or absolute_diff_path(original layer)
	if l.realPath != "" {
		path = l.realPath
	} else {
		path = l.diffPath
	}
	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if info, err := d.Info(); err != nil {
				panic(err)
			} else {
				size += info.Size()
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return size
}

func (l *Layer) DumpLayerSize(size string) {
	if err := os.WriteFile(l.sizePath, []byte(size), 0600); err != nil {
		panic(err)
	}
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

// Construct a shadow layer in memory, but not create in the filesystem
func (l *Layer) Shadow() Layer {
	layerName := "shadow_" + l.layerName
	parentPath := l.layerPath[:len(l.layerPath)-len(l.layerName)]
	shadow := Layer{
		layerPath:    filepath.Join(parentPath, layerName),
		diffPath:     filepath.Join(parentPath, layerName, "diff"),
		linkPath:     filepath.Join(parentPath, layerName, "link"),
		realPath:     filepath.Join(parentPath, layerName, "real"),
		linkContent:  "shadow_" + l.linkContent,
		lowerContent: "",
		layerName:    layerName,
		metaPath:     l.metaPath,
		cacheidPath:  l.cacheidPath,
		sizePath:     l.sizePath,
		cacheid:      layerName,
		size:         l.size,
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

// Construct an original layer from a shadow layer, not create it in the filesystem
func (l *Layer) Original() Layer {
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

func (l *Layer) Restore() {
	if !strings.Contains(l.layerName, "shadow") {
		panic("Only shadow layers can be restored")
	}
	original := l.Original()
	bakCacheIdData, err := os.ReadFile(original.cacheidPath + ".bak")
	if err != nil {
		panic(err)
	}
	if err = os.WriteFile(original.cacheidPath, bakCacheIdData, 0600); err != nil {
		panic(err)
	}
}

// // Construct a cache layer from an original/shadow layer
// func (l *Layer) Cache() Layer {
// 	layerName := "cache_" + l.layer_name
// 	parentPath := l.absolute_layer_path[:len(l.absolute_layer_path)-len(l.layer_name)]
// 	cache := Layer{
// 		absolute_layer_path: filepath.Join(parentPath, layerName),
// 		absolute_diff_path:  filepath.Join(parentPath, layerName, "diff"),
// 		absolute_link_path:  filepath.Join(parentPath, layerName, "link"),
// 		absolute_real_path:  filepath.Join(parentPath, layerName, "real"),
// 		link_content:        "link_" + layerName,
// 		lower_content:       "",
// 		layer_name:          layerName,
// 	}
// 	cache.absolute_l_link_path = filepath.Join(l.absolute_l_link_path[:len(l.absolute_l_link_path)-len(l.link_content)], cache.link_content)
// 	return cache

// }

// // Create dir and files according to the layer
// // Only shadow layers can call this function
// func (l *Layer) Dump() {
// 	var mode fs.FileMode = 0755
// 	// create layer
// 	if err := os.Mkdir(l.absolute_layer_path, mode); err != nil {
// 		if errors.Is(err, os.ErrExist) {
// 			fmt.Println("Layer diff dir already exists")
// 		} else {
// 			panic(err)
// 		}
// 	}

// 	// create diff dir
// 	if err := os.Mkdir(l.absolute_diff_path, mode); err != nil {
// 		if errors.Is(err, os.ErrExist) {
// 			fmt.Println("Layer diff dir already exists")
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	//create link file
// 	link_file, err := os.Create(l.absolute_link_path)
// 	if err != nil {
// 		if errors.Is(err, os.ErrExist) {
// 			fmt.Println("Layer link file already exists")
// 		} else {
// 			panic(err)
// 		}
// 	}
// 	link_file.WriteString(l.link_content)

// 	// create lower file, optionally
// 	if l.absolute_lower_path != "" {
// 		if lower_file, err := os.Create(l.absolute_lower_path); err != nil {
// 			if errors.Is(err, os.ErrExist) {
// 				fmt.Println("Layer lower file already exists")
// 			} else {
// 				panic(err)
// 			}
// 		} else {
// 			// set lowers
// 			lower_file.WriteString(l.lower_content)
// 		}
// 	}

// 	// create real dir
// 	if l.absolute_real_path != "" {
// 		if err := os.Mkdir(l.absolute_real_path, mode); err != nil {
// 			if errors.Is(err, os.ErrExist) {
// 				fmt.Println("Layer real dir already exists")
// 			} else {
// 				panic(err)
// 			}
// 		}
// 	}

// 	// create l/link_{layer_name}
// 	if err := os.Symlink(filepath.Join("../", l.layer_name, "diff"), l.absolute_l_link_path); err != nil {
// 		if errors.Is(err, os.ErrExist) {
// 			fmt.Println("l/link file already exists")
// 		} else {
// 			panic(err)
// 		}
// 	}

// 	if l.absolute_cacheid_path != "" {
// 		// set cache-id
// 		// back up the original cache-id firstly
// 		if !file_exist(l.absolute_cacheid_path + ".bak") {
// 			copy(l.absolute_cacheid_path, l.absolute_cacheid_path+".bak")
// 		}
// 		if err := os.WriteFile(l.absolute_cacheid_path, []byte(l.cacheid), 0600); err != nil {
// 			panic(err)
// 		}

// 	}
// }

// // archive files under diff directory
// func (l *Layer) TarDiff(destFile string) {
// 	TarFiles(l.absolute_diff_path, destFile)
// }

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
	absRealPath := filepath.Join(layerPath, "real")
	if !util.PathExist(absRealPath) {
		layer.realPath = ""
	} else {
		layer.realPath = absRealPath
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
