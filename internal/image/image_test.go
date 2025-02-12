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
package image

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func setUp() string {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli.ImagePull(ctx, "hello-world:latest", image.PullOptions{})
	time.Sleep(1 * time.Second)

	imgInfo, _, _ := cli.ImageInspectWithRaw(ctx, "hello-world:latest")
	return filepath.Dir(imgInfo.GraphDriver.Data["UpperDir"])
}

func tearDown() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cli.ImageRemove(ctx, "hello-world:latest", image.RemoveOptions{})
}

func TestLayerGetLayerSize(t *testing.T) {
	defer tearDown()
	layerPath := setUp()
	layerInfo := NewLayerInfo(layerPath)

	size := layerInfo.GetLayerSize()

	assert.Greater(t, size, int64(100))
}

func TestNewLayerInfo(t *testing.T) {
	defer tearDown()
	layerPath := setUp()

	layerInfo := NewLayerInfo(layerPath)

	assert.NotEmpty(t, layerInfo.diffPath)
	assert.NotEmpty(t, layerInfo.linkPath)
}
