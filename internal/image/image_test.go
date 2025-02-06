package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayerGetLayerSize(t *testing.T) {
	layer := NewLayer("/var/lib/docker/overlay2/c786349027930120f67579f2bd2ecff92178e702bf8676219022320bfc223f1f/")

	size := layer.GetLayerSize()

	assert.Greater(t, size, int64(100))
}

func TestNewLayer(t *testing.T) {
	layer := NewLayer("/var/lib/docker/overlay2/c786349027930120f67579f2bd2ecff92178e702bf8676219022320bfc223f1f/")

	assert.NotEmpty(t, layer.diffPath)
	assert.NotEmpty(t, layer.linkPath)
}
