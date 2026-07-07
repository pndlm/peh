package peh3

import (
	"fmt"
	"os"
	"strings"

	"github.com/pndlm/peh/peh3/docker"
)

func (proj *Project) DockerClient() *docker.Client {
	return docker.NewClient()
}

func (proj *Project) DeleteExitedContainers() {
	containers, err := proj.DockerClient().ContainerList(docker.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		name := container.Labels["com.docker.swarm.service.name"]
		if container.State == "exited" && strings.HasPrefix(name, proj.Name+"_") {
			fmt.Fprintf(os.Stderr, "Removing %s\n", name)
			err := proj.DockerClient().ContainerRemove(container.ID)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (proj *Project) GetServiceContainerShell(serviceName string) {
	container := proj.RunningServiceContainer(serviceName)
	cmd := StdStreamCommand("docker", "exec", "-it", container.ID, "/bin/bash")
	cmd.Run()
}

func (proj *Project) RunningServiceContainers(serviceName string) []docker.Container {
	var matches []docker.Container
	containers, err := proj.DockerClient().ContainerList(docker.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		name := container.Labels["com.docker.swarm.service.name"]
		if strings.HasSuffix(name, "_"+serviceName) {
			matches = append(matches, container)
		}
	}
	return matches
}

func (proj *Project) RunningServiceContainer(serviceName string) docker.Container {
	containers := proj.RunningServiceContainers(serviceName)
	if len(containers) < 1 {
		fmt.Fprintf(os.Stderr, "Service %s has no containers\n", serviceName)
		os.Exit(1)
	}
	if len(containers) > 1 {
		fmt.Fprintf(os.Stderr, "Service %s has more than 1 container\n", serviceName)
		os.Exit(1)
	}
	return containers[0]
}

// https://github.com/docker/cli/tree/master/cli/command/stack
func (proj *Project) StackUp(composeFile string) {
	cmd := StdStreamCommand("docker", "stack", "up", "-c", composeFile, proj.Name)
	cmd.Run()
}

func (proj *Project) StackDown() {
	cmd := StdStreamCommand("docker", "stack", "down", proj.Name)
	cmd.Run()
}

func (proj *Project) StopServiceContainers(serviceName string) {
	containers := proj.RunningServiceContainers(serviceName)
	for _, container := range containers {
		fmt.Fprintf(os.Stderr, "Stopping %s\n", container.Labels["com.docker.swarm.service.name"])
		err := proj.DockerClient().ContainerStop(container.ID)
		if err != nil {
			panic(err)
		}
	}
}

func (proj *Project) TailServiceContainer(serviceName string) {
	container := proj.RunningServiceContainer(serviceName)
	cmd := StdStreamCommand("docker", "logs", "-f", container.ID)
	cmd.Run()
}
