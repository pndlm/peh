package main

import "github.com/pndlm/peh/peh3"

func main() {
	pehdkr := peh3.MustGetPehDockerClient()
	pehdkr.DeleteExitedContainers()
}
