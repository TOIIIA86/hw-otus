package main

import (
	"log"
	"os"
)

func main() {
	l := len(os.Args)
	if l < 2 {
		log.Fatalln(
			"the path to env dir and command are not specified: go-envdir /path/to/env/dir command",
		)
	}
	if l < 3 {
		log.Fatalln("the command is not specified: go-envdir /path/to/env/dir command")
	}

	envdir, command := os.Args[1], os.Args[2:]

	env, err := ReadDir(envdir)
	if err != nil {
		log.Fatalln(err)
	}

	retCode := RunCmd(command, env)
	os.Exit(retCode)
}
