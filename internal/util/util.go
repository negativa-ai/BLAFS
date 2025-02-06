package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func PathExist(absPath string) bool {
	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CopyFile(src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		panic(err)
	}
	// Write data to dst
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		panic(err)
	}
}

func GetDirSize(path string) int64 {
	var size int64
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
