package mount

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMountedWithType(t *testing.T) {

	mounted := IsMountedWithType("/home/ubuntu/repos/BAFFS/build/mount", "fuse.debloated_fs")

	assert.True(t, mounted)
}
