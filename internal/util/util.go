package util

import (
	"errors"
	"os"
)

func PathExist(absPath string) bool {
	if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
