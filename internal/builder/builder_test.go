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
