package peh3

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

func (proj *Project) CmdDown() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Take the stack down",
		Run: func(cmd *cobra.Command, args []string) {
			proj.StackDown()
		},
	}
}

func (proj *Project) CmdExitedRm() *cobra.Command {
	return &cobra.Command{
		Use:   "exitedrm",
		Short: "Clean up (delete) exited containers",
		Run: func(cmd *cobra.Command, args []string) {
			proj.DeleteExitedContainers()
		},
	}
}

func (proj *Project) CmdRestart() *cobra.Command {
	return &cobra.Command{
		Use:   "restart {serviceName}",
		Short: "Stop a service container (causing respawn)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			proj.StopServiceContainers(args[0])
		},
	}
}

func (proj *Project) CmdSh() *cobra.Command {
	return &cobra.Command{
		Use:   "sh {serviceName}",
		Short: "Get a shell to a service container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			proj.GetServiceContainerShell(args[0])
		},
	}
}

func (proj *Project) CmdTail() *cobra.Command {
	return &cobra.Command{
		Use:   "tail {serviceName}",
		Short: "Tail a service container's logs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			proj.TailServiceContainer(args[0])
		},
	}
}

func (proj *Project) CmdUp() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Bring the stack up",
		Run: func(cmd *cobra.Command, args []string) {
			proj.StackUp(filepath.Join("docker", "docker-compose.yaml"))
		},
	}
}
