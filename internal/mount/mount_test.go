package mount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMountedWithType(t *testing.T) {

	mounted := IsMountedWithType("/tmp", "fuse.debloated_fs")

	assert.False(t, mounted)
}
