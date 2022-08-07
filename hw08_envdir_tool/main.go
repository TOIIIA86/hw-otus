package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Not enough args")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Env reading err: %v", err)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
