package utils

import (
	"os"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
