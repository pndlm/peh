package cmd

import (
	"os"
	"path/filepath"

	"github.com/pndlm/peh/peh3"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Bring the stack up",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(filepath.Join(proj.Dir, "docker", "working", "letsencrypt"), 0750)
		if err != nil {
			panic(err)
		}
		cmd2 := peh3.StdStreamCommand("docker", "stack", "up", "-c", filepath.Join("docker", "docker-compose.yaml"), proj.Name)
		peh3.ApplyCmdEnv(cmd2, filepath.Join(proj.Dir, "docker", ".env"))
		cmd2.Run()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
