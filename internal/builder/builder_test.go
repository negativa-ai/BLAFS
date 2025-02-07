package builder

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestExtractLayerName(t *testing.T) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	imgInfo, _, _ := cli.ImageInspectWithRaw(ctx, "hello-world")

	layerNames := extractLayerNames(imgInfo.GraphDriver)

	assert.Equal(t, len(layerNames), 1)
}

func TestExtractLayers(t *testing.T) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	imgInfo, _, _ := cli.ImageInspectWithRaw(ctx, "hello-world")

	layers := ExtractLayersInfo(&imgInfo, "/var/lib/docker/overlay2/", "/var/lib/docker/")

	assert.Equal(t, len(layers), 1)
}
