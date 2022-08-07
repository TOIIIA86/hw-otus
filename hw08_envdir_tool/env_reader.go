package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type Environment map[string]EnvVal

// EnvVal helps to distinguish between empty files and files with the first empty line.
type EnvVal struct {
	Value    string
	UnsetVal bool
}

var (
	ErrInvalidFileName  = errors.New("filename contains a '=' symboll")
	ErrFileIsADirectory = errors.New("file is a directory")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envVars := make(Environment, len(filesInfo))
	for _, fInfo := range filesInfo {
		if err := validateFileInfo(fInfo); err != nil {
			return nil, err
		}

		file, err := os.Open(filepath.Join(dir, fInfo.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		envVars[fInfo.Name()] = filterEnvVal(scanner.Text())
	}

	return envVars, nil
}

func validateFileInfo(fInfo os.FileInfo) error {
	if fInfo.IsDir() {
		return ErrFileIsADirectory
	}

	if strings.Contains(fInfo.Name(), "=") {
		return ErrInvalidFileName
	}

	return nil
}

func filterEnvVal(s string) EnvVal {
	var envVal EnvVal

	s = strings.TrimRightFunc(s, unicode.IsSpace)
	s = strings.ReplaceAll(s, string([]byte{'\x00'}), "\n")
	if len(s) == 0 {
		envVal.UnsetVal = true
	}

	envVal.Value = s

	return envVal
}
