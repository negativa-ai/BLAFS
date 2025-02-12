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
