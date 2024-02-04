package cmd

import (
	"github.com/pndlm/peh/peh3"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring the stack up",
	Run: func(cmd *cobra.Command, args []string) {
		peh3.MustMkdirAll(proj.RelPath("docker/working/letsencrypt"), 0750)
		cmd2 := peh3.StdStreamCommand("docker", "stack", "up", "-c", proj.RelPath("docker/docker-compose.yaml"), proj.Name)
		peh3.ApplyCmdEnv(cmd2, proj.RelPath("docker/.env"))
		cmd2.Run()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
