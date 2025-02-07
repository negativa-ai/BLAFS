package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLayerGetLayerSize(t *testing.T) {
	layerInfo := NewLayerInfo("/var/lib/docker/overlay2/474f8a8095314d6f8a42d484b7342c3f3dd1736f503da0cd929ea019abf45090/")

	size := layerInfo.GetLayerSize()

	assert.Greater(t, size, int64(100))
}

func TestNewLayerInfo(t *testing.T) {
	layerInfo := NewLayerInfo("/var/lib/docker/overlay2/474f8a8095314d6f8a42d484b7342c3f3dd1736f503da0cd929ea019abf45090/")

	assert.NotEmpty(t, layerInfo.diffPath)
	assert.NotEmpty(t, layerInfo.linkPath)
}
