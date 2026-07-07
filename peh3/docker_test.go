package peh3_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pndlm/peh/peh3"
	"github.com/pndlm/peh/peh3/docker"
)

const testStackComposeFile = `version: '3.8'
services:
  sleeper:
    image: busybox:latest
    command: sleep 3600
`

// requireSwarm skips the test unless a local docker daemon is reachable
// and swarm mode is active, since StackUp/StackDown shell out to the
// real "docker stack" CLI against whatever daemon is on PATH.
func requireSwarm(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping docker integration test in -short mode")
	}
	out, err := exec.Command("docker", "info", "--format", "{{.Swarm.LocalNodeState}}").Output()
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	if strings.TrimSpace(string(out)) != "active" {
		t.Skip("docker swarm is not active on this host")
	}
}

func waitUntil(t *testing.T, timeout time.Duration, desc string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		if cond() {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for: %s", desc)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func stackServiceNames(t *testing.T, stackName string) []string {
	t.Helper()
	out, err := exec.Command("docker", "stack", "services", stackName, "--format", "{{.Name}}").Output()
	if err != nil {
		return nil
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}

func containerState(t *testing.T, id string) (state string, exists bool) {
	t.Helper()
	out, err := exec.Command("docker", "inspect", "--format", "{{.State.Status}}", id).Output()
	if err != nil {
		return "", false
	}
	return strings.TrimSpace(string(out)), true
}

func stackExists(t *testing.T, stackName string) bool {
	t.Helper()
	out, err := exec.Command("docker", "stack", "ls", "--format", "{{.Name}}").Output()
	if err != nil {
		t.Fatalf("docker stack ls: %v", err)
	}
	for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if name == stackName {
			return true
		}
	}
	return false
}

// networkExists checks the daemon directly rather than relying on
// "docker stack ls", since overlay network teardown is asynchronous and
// can still be in progress after the stack itself stops being listed -
// racing a same-named network's removal against a subsequent test run's
// creation of it is what makes StackUp flaky back-to-back.
func networkExists(t *testing.T, networkName string) bool {
	t.Helper()
	out, err := exec.Command("docker", "network", "ls", "--filter", "name=^"+networkName+"$", "--format", "{{.Name}}").Output()
	if err != nil {
		t.Fatalf("docker network ls: %v", err)
	}
	return strings.TrimSpace(string(out)) != ""
}

func waitForStackTeardown(t *testing.T, proj *peh3.Project) {
	t.Helper()
	waitUntil(t, 30*time.Second, "stack and its default network to be removed", func() bool {
		return !stackExists(t, proj.Name) && !networkExists(t, proj.Name+"_default")
	})
}

func TestStackUpAndDown(t *testing.T) {
	requireSwarm(t)

	composeFile := filepath.Join(t.TempDir(), "docker-compose.yml")
	if err := os.WriteFile(composeFile, []byte(testStackComposeFile), 0o644); err != nil {
		t.Fatalf("write compose file: %v", err)
	}

	proj := &peh3.Project{Name: "peh3-test-stack"}

	// guard against leftover state from a previous interrupted run racing
	// this one's StackUp (overlay network teardown is asynchronous)
	waitForStackTeardown(t, proj)

	// make sure we don't leave the stack behind even if an assertion fails
	defer func() {
		proj.StackDown()
		waitForStackTeardown(t, proj)
	}()

	proj.StackUp(composeFile)

	waitUntil(t, 60*time.Second, "sleeper service to appear", func() bool {
		names := stackServiceNames(t, proj.Name)
		for _, n := range names {
			if strings.HasSuffix(n, "_sleeper") {
				return true
			}
		}
		return false
	})

	// RunningServiceContainers/RunningServiceContainer only see running
	// containers, which can lag behind the service being listed, so poll
	// using the function under test itself before asserting on it.
	var container docker.Container
	waitUntil(t, 60*time.Second, "sleeper container to be running", func() bool {
		containers := proj.RunningServiceContainers("sleeper")
		if len(containers) != 1 {
			return false
		}
		container = containers[0]
		return true
	})
	if name := container.Labels["com.docker.swarm.service.name"]; !strings.HasSuffix(name, "_sleeper") {
		t.Fatalf("expected container label to end with _sleeper, got %q", name)
	}

	if got := proj.RunningServiceContainer("sleeper"); got.ID != container.ID {
		t.Fatalf("RunningServiceContainer returned ID %q, expected %q", got.ID, container.ID)
	}

	if unmatched := proj.RunningServiceContainers("does-not-exist"); len(unmatched) != 0 {
		t.Fatalf("expected no containers for unmatched service name, got %d", len(unmatched))
	}

	proj.StopServiceContainers("sleeper")

	// swarm should reschedule a replacement task to maintain the service's
	// desired replica count, so a new container ID should appear.
	waitUntil(t, 60*time.Second, "swarm to respawn the sleeper container", func() bool {
		containers := proj.RunningServiceContainers("sleeper")
		if len(containers) != 1 {
			return false
		}
		return containers[0].ID != container.ID
	})

	// the container StopServiceContainers stopped should now be sitting
	// around exited (swarm keeps stopped task containers around for
	// history) - confirm that before relying on it as fixture state.
	waitUntil(t, 30*time.Second, "old container to reach exited state", func() bool {
		state, exists := containerState(t, container.ID)
		return exists && state == "exited"
	})

	proj.DeleteExitedContainers()

	if _, exists := containerState(t, container.ID); exists {
		t.Fatalf("expected exited container %s to be removed by DeleteExitedContainers", container.ID)
	}

	proj.StackDown()

	waitForStackTeardown(t, proj)
}
