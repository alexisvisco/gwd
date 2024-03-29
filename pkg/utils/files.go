package utils

import (
	"os"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func PathIsDir(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}
