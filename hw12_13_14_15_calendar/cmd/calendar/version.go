package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "develop"
	buildDate = "2022-08-04T06:23:11"
	gitHash   = "4bf46c9"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
