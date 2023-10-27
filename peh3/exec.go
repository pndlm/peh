package peh3

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func StdStreamCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func ReadEnv(envPath string) map[string]string {
	env, err := godotenv.Read(envPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s probably does not exist; please create it from the example file\n\n", envPath)
		// panic(err)
		os.Exit(1)
	}
	return env
}

func ApplyCmdEnv(cmd *exec.Cmd, envPath string) {
	env := ReadEnv(envPath)
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}
