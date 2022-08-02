package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdName := cmd[0]
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	for k, v := range env {
		if v.NeedRemove {
			err := os.Unsetenv(k)
			if err != nil {
				return -1
			}
		}
		err := os.Setenv(k, v.Value)
		if err != nil {
			return -2
		}
	}

	ex := exec.Command(cmdName, args...)
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout
	ex.Stderr = os.Stderr
	err := ex.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		return -3
	}
	return
}
