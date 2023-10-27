package cmd

import (
	"fmt"
	"os"

	"github.com/pndlm/peh/peh3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "peh",
	Short: "peh",
}

var proj *peh3.Project

func SetProject(p *peh3.Project) {
	proj = p
	rootCmd.Short = fmt.Sprintf("%s peh", proj.Name)
}

func AddNamedCommand(use string, cmd *cobra.Command) {
	cmd2 := *cmd
	cmd2.Use = use
	rootCmd.AddCommand(&cmd2)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
