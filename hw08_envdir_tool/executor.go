package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	exitCodeOk             = 0
	exitCodeErr            = 1
	exitCodeCannotUnsetEnv = 100
	exitCodeCannotSetEnv   = 101
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	name, args := cmd[0], cmd[1:]

	exCmd := exec.Command(name, args...)

	code := fillEnv(env)
	if code != 0 {
		return code
	}

	exCmd.Stdin = os.Stdin
	exCmd.Stdout = os.Stdout
	exCmd.Stderr = os.Stderr
	if err := exCmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return exitCodeErr
	}
	return exitCodeOk
}

func fillEnv(env Environment) int {
	for key, elem := range env {
		if elem.UnsetVal {
			err := os.Unsetenv(key)
			if err != nil {
				return exitCodeCannotUnsetEnv
			}
			continue
		}

		err := os.Setenv(key, elem.Value)
		if err != nil {
			return exitCodeCannotSetEnv
		}
	}
	return exitCodeOk
}
