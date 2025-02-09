package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/docker/docker/client"
	"github.com/jzh18/baffs/internal/builder"
	"github.com/jzh18/baffs/internal/util"
)

type ShadowCmd struct {
	Images string `arg:"-i,--images" help:"Images to shadow, separated by comma"`
}
type DebloatCmd struct {
	Images string `arg:"-i,--images" help:"Images to debloat separated by comma"`
	Top    int    `arg:"-t,--top" help:"Top N layers to debloat" default:"-1"`
}

var args struct {
	Shadow  *ShadowCmd  `arg:"subcommand:shadow" help:"Shadow images"`
	Debloat *DebloatCmd `arg:"subcommand:debloat" help:"Debloat images"`
}

func restartDocker() {
	cmd := exec.Command("systemctl", "restart", "docker")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(stdout))
}

func shadow(imgName string, workDir string, overlayPath string,
	dockerRootDir string, cli *client.Client, ctx *context.Context) {
	shadowed, originalLayers, shadowLayers := builder.ShadowImage(imgName, workDir, overlayPath, dockerRootDir, cli, ctx, "")
	if !shadowed {
		for _, l := range shadowLayers {
			l.Dump()
		}
		time.Sleep(3 * time.Second)
		restartDocker()
		mounts := builder.CreateMounts("/home/ubuntu/repos/BAFFS/build/debloated_fs", originalLayers, shadowLayers)
		for _, m := range mounts {
			m.Mount()
		}

	} else {
		fmt.Println("already shadowed")
	}

}

func debloat(imgName string, workDir string, overlayPath string, dockerRootDir string, cli *client.Client, ctx *context.Context, topN int) {
	shadowed, imgTarPath := builder.ExportImg(imgName, workDir, overlayPath, dockerRootDir, cli, ctx, topN)
	if shadowed {
		restartDocker()
		time.Sleep(3 * time.Second)
		builder.LoadImage(imgTarPath, cli)

	}
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	dockerInfo, _ := cli.Info(ctx)
	dockerRootDir := dockerInfo.DockerRootDir
	overlayPath := filepath.Join(dockerRootDir, "overlay2")
	workDir := "/usr/local/bafs"
	if !util.PathExist(workDir) {
		if err := os.Mkdir(workDir, 0755); err != nil {
			panic(err)
		}
	}

	switch {
	case args.Shadow != nil:
		images := strings.Split(args.Shadow.Images, ",")
		shadow(images[0], workDir, overlayPath, dockerRootDir, cli, &ctx)
	case args.Debloat != nil:
		images := strings.Split(args.Debloat.Images, ",")
		debloat(images[0], workDir, overlayPath, dockerRootDir, cli, &ctx, args.Debloat.Top)
	}
}
