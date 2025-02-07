package mount

import (
	"fmt"
	"os/exec"
)

type Mount struct {
	exePath    string
	mountPoint string
	kvArgs     map[string]string
	flagArgs   []string
}

func (m Mount) Mount() {
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
