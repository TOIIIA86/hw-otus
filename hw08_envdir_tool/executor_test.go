package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success case simple", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, exitCodeOk, exitCode)
	})

	t.Run("Success case with filled dir env", func(t *testing.T) {
		cmd := []string{"pwd", "-L"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_SET": EnvVal{
				Value: "TEST_SET",
			},
		})

		require.Equal(t, exitCodeOk, exitCode)
		require.Contains(t, os.Environ(), "TEST_SET=TEST_SET")
	})

	t.Run("Exec err case", func(t *testing.T) {
		cmd := []string{"pwd", "-RRR"}
		exitCode := RunCmd(cmd, Environment{
			"TEST_QWE": EnvVal{
				Value: "TEST_QWE",
			},
		})

		require.Equal(t, 1, exitCode)
	})
}

func TestFillEnv(t *testing.T) {
	t.Run("Success simple case", func(t *testing.T) {
		dirEnv := Environment{
			"UNSET": EnvVal{
				Value:    "",
				UnsetVal: true,
			},
			"SET_VAL": EnvVal{
				Value:    "VAL",
				UnsetVal: false,
			},
		}

		resCode := fillEnv(dirEnv)

		require.Equal(t, 0, resCode)
		require.Contains(t, os.Environ(), "SET_VAL=VAL")
	})
}
