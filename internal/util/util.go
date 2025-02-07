package util

import (
	"archive/tar"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
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

// mv, works for directory and file
func Move(src string, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
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

func Sha256Sum(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func TarFiles(sourceDir string, destFile string) {
	// Open the destination file for writing
	dest, err := os.Create(destFile)
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	// Create a new tar writer using the destination file
	tw := tar.NewWriter(dest)
	defer tw.Close()

	// Walk through the source directory and add each file to the tar archive
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a new header for the file
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Set the name of the header to the relative path of the file within the source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		// Handle symbolic links
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			header.Linkname = link
			header.Typeflag = tar.TypeSymlink
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			return nil
		}

		// Write the header to the tar archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If the file is a regular file or a directory, write its contents to the tar archive
		if info.Mode().IsRegular() {
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()

			if _, err := io.Copy(tw, src); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
