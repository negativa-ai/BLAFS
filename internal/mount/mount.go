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
package mount

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

const MountType = "fuse.debloated_fs"

// A Mount represents a BAFFS mount.
type Mount struct {
	exePath    string            // Path to the debloated_fs executable
	mountPoint string            // Path to the mount point
	kvArgs     map[string]string // Key-value arguments
	flagArgs   []string          // Flag arguments
}

// Mounts a BAFFS mount.
// A mount is skipped if the mount point is already mounted with the same mount type.
func (m Mount) Mount() {
	if IsMountedWithType(m.mountPoint, MountType) {
		return
	}
	// debloated_fs -s -d  --realdir=/tmp/real5 --lowerdir=/tmp/lower5  /tmp/mnt5
	args := []string{}
	args = append(args, m.flagArgs...)
	args = append(args, "--realdir="+m.kvArgs["--realdir"])
	args = append(args, "--lowerdir="+m.kvArgs["--lowerdir"])
	args = append(args, "--optimize="+m.kvArgs["--optimize"])
	args = append(args, m.mountPoint)

	cmd := exec.Command(m.exePath, args...)
	log.Debug(cmd)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

// Unmounts a BAFFS mount.
func (m Mount) Unmount() {
	// umount /tmp/mnt5
	cmd := exec.Command("umount", m.mountPoint)
	log.Debug(cmd)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

func NewMount(exePath string, mountPoint string, kvArgs map[string]string, flagArgs []string) Mount {
	return Mount{
		exePath:    exePath,
		mountPoint: mountPoint,
		kvArgs:     kvArgs,
		flagArgs:   flagArgs,
	}
}

// IsMountedWithType returns true if the given path is a mount point using the specified mount type.
func IsMountedWithType(path, mountType string) bool {
	// Open /proc/self/mounts to read mount information.
	file, err := os.Open("/proc/self/mounts")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Each line is expected to have at least 3 fields:
		// device mountpoint fstype [options dump pass]
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue // skip malformed lines
		}

		mountPoint := fields[1]
		fsType := fields[2]

		// Check if the mount point and filesystem type match.
		if mountPoint == path && fsType == mountType {
			return true
		}
	}

	// Check for scanner errors.
	if err := scanner.Err(); err != nil {
		return false
	}

	return false
}
