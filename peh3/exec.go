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
		fmt.Fprintf(os.Stderr, "Env file '%s' is not readable; ensure it has been created from example (see README)\n\n", envPath)
		// panic(err)
		os.Exit(1)
	}
	return env
}

func RequireEnv(env map[string]string, keys []string) {
	for _, key := range keys {
		if env[key] == "" {
			fmt.Fprintf(os.Stderr, "Env value '%s' does not exist; see example file and README\n\n", key)
			os.Exit(1)
		}
	}
}

func ApplyCmdEnv(cmd *exec.Cmd, envPath string) {
	ApplyCmdEnvVal(cmd, ReadEnv(envPath))
}

func ApplyCmdEnvVal(cmd *exec.Cmd, env map[string]string) {
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}
