package peh3

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pndlm/peh/peh3/docker"
)

const (
	swarmServiceNameLabel    = "com.docker.swarm.service.name"
	composeProjectLabel      = "com.docker.compose.project"
	composeServiceLabel      = "com.docker.compose.service"
	composeVolumeLabel       = "com.docker.compose.volume"
	swarmStackNamespaceLabel = "com.docker.stack.namespace"
)

func (proj *Project) DockerClient() *docker.Client {
	return docker.NewClient()
}

// containerBelongsToProject reports whether a container was launched as
// part of this project's stack/compose deployment.
func (proj *Project) containerBelongsToProject(container docker.Container) bool {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		return container.Labels[composeProjectLabel] == proj.Name
	case ProjectApparatusSwarm:
		return strings.HasPrefix(container.Labels[swarmServiceNameLabel], proj.Name+"_")
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

// containerMatchesService reports whether a container is a running instance
// of the given service within this project's deployment.
func (proj *Project) containerMatchesService(container docker.Container, serviceName string) bool {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		return container.Labels[composeProjectLabel] == proj.Name &&
			container.Labels[composeServiceLabel] == serviceName
	case ProjectApparatusSwarm:
		return strings.HasSuffix(container.Labels[swarmServiceNameLabel], "_"+serviceName)
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

// containerServiceLabel returns a human-readable service identifier for
// log/status output, in whichever form the active apparatus labels it with.
func (proj *Project) containerServiceLabel(container docker.Container) string {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		return container.Labels[composeServiceLabel]
	case ProjectApparatusSwarm:
		return container.Labels[swarmServiceNameLabel]
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

// volumeMatchesName reports whether a volume is this project's named volume
// with the given short name (as declared in the compose file), however the
// active apparatus happens to name/label it. Compose labels the short name
// directly; a Swarm stack only labels its namespace and instead bakes the
// short name into the volume's literal "<stack>_<volume>" name.
func (proj *Project) volumeMatchesName(volume docker.Volume, name string) bool {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		return volume.Labels[composeProjectLabel] == proj.Name &&
			volume.Labels[composeVolumeLabel] == name
	case ProjectApparatusSwarm:
		return volume.Labels[swarmStackNamespaceLabel] == proj.Name &&
			volume.Name == proj.Name+"_"+name
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

// DeleteVolume removes this project's named volume (e.g. "db") so that
// bringing the stack back up recreates it empty. The stack/compose project
// must already be down: Docker refuses to remove a volume still in use by a
// container, and that's deliberate here rather than forced, since a
// force-remove upending a volume that's actually still mounted would be a
// good way to lose data.
func (proj *Project) DeleteVolume(name string) {
	volumes, err := proj.DockerClient().VolumeList()
	if err != nil {
		panic(err)
	}
	found := false
	for _, volume := range volumes {
		if !proj.volumeMatchesName(volume, name) {
			continue
		}
		found = true
		fmt.Fprintf(os.Stderr, "Removing volume %s\n", volume.Name)
		if err := proj.DockerClient().VolumeRemove(volume.Name); err != nil {
			panic(err)
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "Volume %s not found\n", name)
	}
}

func (proj *Project) DeleteExitedContainers() {
	containers, err := proj.DockerClient().ContainerList(docker.ContainerListOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if container.State == "exited" && proj.containerBelongsToProject(container) {
			fmt.Fprintf(os.Stderr, "Removing %s\n", proj.containerServiceLabel(container))
			err := proj.DockerClient().ContainerRemove(container.ID)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (proj *Project) GetServiceContainerShell(serviceName string) {
	container := proj.RunningServiceContainer(serviceName)
	for _, sh := range []string{"/bin/bash", "/bin/sh"} {
		rChk := exec.Command("docker", "exec", container.ID, "test", "-x", sh).Run()
		if rChk != nil && rChk.Error() != "" {
			continue
		}
		cmd := StdStreamCommand("docker", "exec", "-it", container.ID, sh)
		cmd.Run()
		return
	}
	fmt.Fprintf(os.Stderr, "Service %s has no suitable shell", serviceName)
	os.Exit(1)
}

func (proj *Project) RunningServiceContainers(serviceName string) []docker.Container {
	var matches []docker.Container
	containers, err := proj.DockerClient().ContainerList(docker.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		if proj.containerMatchesService(container, serviceName) {
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

// StackUpCmd builds the apparatus-appropriate "bring it up" command without
// running it, so callers can apply env vars or other tweaks (e.g. via
// ApplyCmdEnv) before executing it themselves.
//
// https://github.com/docker/cli/tree/master/cli/command/stack
// https://docs.docker.com/reference/cli/docker/compose/up/
func (proj *Project) StackUpCmd(composeFile string) *exec.Cmd {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		return StdStreamCommand("docker", "compose", "-p", proj.Name, "-f", composeFile, "up", "-d")
	case ProjectApparatusSwarm:
		return StdStreamCommand("docker", "stack", "up", "-c", composeFile, proj.Name)
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

func (proj *Project) StackUp(composeFile string) {
	proj.StackUpCmd(composeFile).Run()
}

func (proj *Project) StackDown() {
	switch proj.Apparatus {
	case ProjectApparatusCompose:
		cmd := StdStreamCommand("docker", "compose", "-p", proj.Name, "down")
		cmd.Run()
	case ProjectApparatusSwarm:
		cmd := StdStreamCommand("docker", "stack", "down", proj.Name)
		cmd.Run()
	default:
		panic(fmt.Errorf("unsupported apparatus '%s'", proj.Apparatus))
	}
}

func (proj *Project) StopServiceContainers(serviceName string) {
	containers := proj.RunningServiceContainers(serviceName)
	for _, container := range containers {
		fmt.Fprintf(os.Stderr, "Stopping %s\n", proj.containerServiceLabel(container))
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
