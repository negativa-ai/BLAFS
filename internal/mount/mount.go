package mount

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MountType = "fuse.debloated_fs"

type Mount struct {
	exePath    string
	mountPoint string
	kvArgs     map[string]string
	flagArgs   []string
}

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
	fmt.Println(cmd)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}

func (m Mount) Unmount() {
	// umount /tmp/mnt5
	cmd := exec.Command("umount", m.mountPoint)
	fmt.Println(cmd)
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
