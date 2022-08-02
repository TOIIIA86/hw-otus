package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidName = errors.New("name contains invalid characters")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		// TODO: как получить dir из os.DirEnv?
		err = registerVar(env, dir, entry)
		if err != nil {
			return nil, err
		}
	}

	return env, nil
}

func registerVar(env Environment, dir string, entry os.DirEntry) error {
	name := entry.Name()
	if strings.ContainsAny(name, "=") {
		return ErrInvalidName
	}

	// case: remove
	info, err := entry.Info()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		env[name] = EnvValue{NeedRemove: true}
		return nil
	}

	path := filepath.Join(dir, name)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	value = strings.TrimRight(value, " \t\n")
	value = string(bytes.ReplaceAll([]byte(value), []byte{0x0}, []byte("\n")))
	env[name] = EnvValue{Value: value}

	return nil
}
