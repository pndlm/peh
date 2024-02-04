package peh3

import (
	"io/fs"
	"os"
)

func MustMkdirAll(path string, perm fs.FileMode) {
	err := os.MkdirAll(path, perm)
	if err != nil {
		panic(err)
	}
}
