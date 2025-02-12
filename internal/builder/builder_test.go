package builder

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestExtractLayerName(t *testing.T) {

	mocked := types.GraphDriverData{
		Data: map[string]string{
			"MergedDir": "/var/lib/docker/overlay2/02bbe40378a44f8a88293229da146a33578eedbb9d7947808503722480e00505/merged",
			"UpperDir":  "/var/lib/docker/overlay2/02bbe40378a44f8a88293229da146a33578eedbb9d7947808503722480e00505/diff",
			"WorkDir":   "/var/lib/docker/overlay2/02bbe40378a44f8a88293229da146a33578eedbb9d7947808503722480e00505/work",
		},
		Name: "overlay2",
	}

	layerNames := extractLayerNames(mocked)

	assert.Equal(t, len(layerNames), 1)
}
