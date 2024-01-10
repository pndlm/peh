package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var projectDir string
	var cmdName string
	var cmdArgs []string

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// traverse up from cwd
	for dir != string(filepath.Separator) {
		pehDir := filepath.Join(dir, "peh")
		if s, err := os.Stat(pehDir); err == nil && s.IsDir() {
			bin := filepath.Join(pehDir, "bin", "peh")
			if _, err := os.Stat(bin); err == nil {
				// run found peh bin
				projectDir = dir
				cmdName = "peh"
				break
			}
			srcMain := filepath.Join(pehDir, "src", "peh.go")
			if _, err := os.Stat(srcMain); err == nil {
				// run found peh source
				projectDir = dir
				cmdName = "go"
				cmdArgs = []string{
					"run",
					srcMain,
				}
				break
			}
		}
		shPath := filepath.Join(dir, "peh.sh")
		if _, err := os.Stat(shPath); err == nil {
			// run found peh.sh
			projectDir = dir
			cmdName = "./peh.sh"
			break
		}
		dir = filepath.Dir(dir)
	}
	if projectDir == "" {
		fmt.Fprintf(os.Stderr, "It doesn't look like you're in a peh project\n")
		os.Exit(1)
	}

	cmdArgs = append(cmdArgs, os.Args[1:]...)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = projectDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
	}
}
