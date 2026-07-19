package cmd

import (
	"github.com/pndlm/peh/peh3"
	"github.com/spf13/cobra"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Run mysql CLI",
	Run: func(cmd *cobra.Command, args []string) {
		env := peh3.ReadEnv(proj.RelPath(".env"))
		peh3.RequireEnv(env, []string{
			"MYSQL_ROOT_PASSWORD",
			"MYSQL_DATABASE",
		})
		container := proj.RunningServiceContainer("db")
		cmd1 := peh3.StdStreamCommand(
			"docker", "exec", "-it", container.ID,
			"mysql", "-p"+env["MYSQL_ROOT_PASSWORD"], "-uroot",
			env["MYSQL_DATABASE"],
		)
		cmd1.Run()
	},
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
}
