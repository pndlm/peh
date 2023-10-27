package peh3

import (
	"fmt"
	"os"
	"path/filepath"

	mobyclient "github.com/moby/moby/client"
)

type Project struct {
	Dir        string
	Name       string
	mobyclient *mobyclient.Client
}

func ProjectAtCwd(name string) *Project {
	var projectDir string
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// traverse up from cwd
	for dir != string(filepath.Separator) {
		pehDir := filepath.Join(dir, "peh")
		if s, err := os.Stat(pehDir); err == nil && s.IsDir() {
			projectDir = dir
			break
		}
		dir = filepath.Dir(dir)
	}
	if projectDir == "" {
		fmt.Fprintf(os.Stderr, "Enclosing project directory not found\n")
		os.Exit(1)
	}
	return &Project{
		Dir:  projectDir,
		Name: name,
	}
}
