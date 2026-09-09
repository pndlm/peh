package cmd

import (
	"github.com/pndlm/peh/peh3"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring the stack up",
	Run: func(cmd *cobra.Command, args []string) {
		cmd2 := proj.StackUpCmd(proj.RelPath("docker-compose.yaml"))
		peh3.ApplyCmdEnv(cmd2, proj.RelPath(".env"), false)
		cmd2.Run()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
