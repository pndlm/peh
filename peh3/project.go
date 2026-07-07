package peh3

import (
	"fmt"
	"os"
	"path/filepath"
)

type ProjectApparatus string

const (
	ProjectApparatusCompose ProjectApparatus = "compose"
	ProjectApparatusSwarm   ProjectApparatus = "swarm"
)

type Project struct {
	Dir       string
	Name      string
	Apparatus ProjectApparatus
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
		// the old peh default, clunkier for dev
		// but technically better for a quick vm-based deploy
		Apparatus: ProjectApparatusSwarm,
	}
}

// Get an OS-agnostic path relative to the project directory
// using filepath.FromSlash
func (proj *Project) RelPath(path string) string {
	return filepath.Join(proj.Dir, filepath.FromSlash(path))
}
