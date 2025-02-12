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
	"github.com/jzh18/baffs/internal/image"
	"github.com/jzh18/baffs/internal/mount"
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

func shadow(imgName []string, workDir string, overlayPath string,
	dockerRootDir string, cli *client.Client, ctx *context.Context) {
	var allShadowLayers [][]image.ShadowLayer
	var allImgMounts [][]mount.Mount
	for _, imgName := range imgName {

		shadowed, originalLayers, shadowLayers := builder.ShadowImage(imgName, workDir, overlayPath, dockerRootDir, cli, ctx, "")
		if !shadowed {
			allShadowLayers = append(allShadowLayers, shadowLayers)
			mounts := builder.CreateMounts("/home/ubuntu/repos/BAFFS/build/debloated_fs", originalLayers, shadowLayers)
			allImgMounts = append(allImgMounts, mounts)
		} else {
			fmt.Println("already shadowed")
		}
	}

	for _, shadowLayers := range allShadowLayers {
		for _, l := range shadowLayers {
			l.Dump()
		}
	}

	time.Sleep(3 * time.Second)
	restartDocker()
	for _, mounts := range allImgMounts {
		for _, m := range mounts {
			m.Mount()
		}
	}
}

func debloat(imgNames []string, workDir string, overlayPath string, dockerRootDir string, cli *client.Client, ctx *context.Context, topN int) {
	var imgPaths []string
	var allShadowLayers [][]image.ShadowLayer
	for _, imgName := range imgNames {
		shadowed, imgTarPath, shadowLayers := builder.ExportImg(imgName, workDir, overlayPath, dockerRootDir, cli, ctx, topN)
		if shadowed {
			imgPaths = append(imgPaths, imgTarPath)
			allShadowLayers = append(allShadowLayers, shadowLayers)
		}
	}

	for _, shadowLayers := range allShadowLayers {
		for _, l := range shadowLayers {
			l.Restore()
		}
	}

	restartDocker()
	time.Sleep(3 * time.Second)

	for _, imgTarPath := range imgPaths {
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
		shadow(images, workDir, overlayPath, dockerRootDir, cli, &ctx)
	case args.Debloat != nil:
		images := strings.Split(args.Debloat.Images, ",")
		debloat(images, workDir, overlayPath, dockerRootDir, cli, &ctx, args.Debloat.Top)
	}
}
