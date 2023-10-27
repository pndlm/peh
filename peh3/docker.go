package peh3

import (
	"context"
	"fmt"
	"os"
	"strings"

	dkrtypes "github.com/docker/docker/api/types"
	dkrcontainertypes "github.com/docker/docker/api/types/container"
	mobyclient "github.com/moby/moby/client"
)

func (proj *Project) DockerClient() *mobyclient.Client {
	if proj.mobyclient == nil {
		mcli, err := mobyclient.NewClientWithOpts(mobyclient.FromEnv)
		if err != nil {
			panic(err)
		}
		proj.mobyclient = mcli
	}
	return proj.mobyclient
}

func (proj *Project) DeleteExitedContainers() {
	containers, err := proj.DockerClient().ContainerList(
		context.TODO(),
		dkrtypes.ContainerListOptions{
			All: true,
		},
	)
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if container.State == "exited" {
			fmt.Fprintf(os.Stderr, "Removing %s\n", container.Labels["com.docker.swarm.service.name"])
			err := proj.DockerClient().ContainerRemove(context.TODO(), container.ID, dkrtypes.ContainerRemoveOptions{})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (proj *Project) GetServiceContainerShell(serviceName string) {
	containers := proj.RunningServiceContainers(serviceName)
	if len(containers) < 1 {
		fmt.Fprintf(os.Stderr, "Service %s has no containers", serviceName)
		return
	}
	if len(containers) > 1 {
		fmt.Fprintf(os.Stderr, "Service %s has more than 1 container", serviceName)
		return
	}
	cmd := StdStreamCommand("docker", "exec", "-it", containers[0].ID, "/bin/bash")
	cmd.Run()
}

func (proj *Project) RunningServiceContainers(serviceName string) []dkrtypes.Container {
	matches := []dkrtypes.Container{}
	containers, err := proj.DockerClient().ContainerList(context.TODO(), dkrtypes.ContainerListOptions{})
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
		err := proj.DockerClient().ContainerStop(context.TODO(), container.ID, dkrcontainertypes.StopOptions{})
		if err != nil {
			panic(err)
		}
	}
}

func (proj *Project) TailServiceContainer(serviceName string) {
	containers := proj.RunningServiceContainers(serviceName)
	if len(containers) < 1 {
		fmt.Fprintf(os.Stderr, "Service %s has no containers", serviceName)
		return
	}
	if len(containers) > 1 {
		fmt.Fprintf(os.Stderr, "Service %s has more than 1 container", serviceName)
		return
	}
	cmd := StdStreamCommand("docker", "logs", "-f", containers[0].ID)
	cmd.Run()
}
