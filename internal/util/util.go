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

// PathExist checks if a path exists
func PathExist(absPath string) bool {
	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// CopyFile copies a file from src to dst
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

// Move moves a file from src to dst, works for directory and file
func Move(src string, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}

// GetDirSize returns the size of a directory, including all files and subdirectories
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

// Sha256Sum returns the sha256 sum of a file
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

// TarFiles creates a tar archive from a directory
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
