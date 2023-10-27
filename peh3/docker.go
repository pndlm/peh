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

type PehDockerClient struct {
	Mobyclient *mobyclient.Client
}

func MustGetPehDockerClient() *PehDockerClient {
	mcli, err := mobyclient.NewClientWithOpts(mobyclient.FromEnv)
	if err != nil {
		panic(err)
	}
	return &PehDockerClient{
		Mobyclient: mcli,
	}
}

func (pehdkr *PehDockerClient) RunningServiceContainers(serviceName string) []dkrtypes.Container {
	matches := []dkrtypes.Container{}
	containers, err := pehdkr.Mobyclient.ContainerList(context.TODO(), dkrtypes.ContainerListOptions{})
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

func (pehdkr *PehDockerClient) StopServiceContainers(serviceName string) {
	containers := pehdkr.RunningServiceContainers(serviceName)
	for _, container := range containers {
		fmt.Fprintf(os.Stderr, "Stopping %s\n", container.Labels["com.docker.swarm.service.name"])
		err := pehdkr.Mobyclient.ContainerStop(context.TODO(), container.ID, dkrcontainertypes.StopOptions{})
		if err != nil {
			panic(err)
		}
	}
}

func (pehdkr *PehDockerClient) DeleteExitedContainers() {
	containers, err := pehdkr.Mobyclient.ContainerList(
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
			err := pehdkr.Mobyclient.ContainerRemove(context.TODO(), container.ID, dkrtypes.ContainerRemoveOptions{})
			if err != nil {
				panic(err)
			}
		}
	}
}
