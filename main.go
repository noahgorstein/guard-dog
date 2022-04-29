package main

import (
	"os"

	"github.com/noahgorstein/guard-dog/cmd"
)

func main() {
	cmd := cmd.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
