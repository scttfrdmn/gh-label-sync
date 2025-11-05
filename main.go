package main

import (
	"os"

	"github.com/scttfrdmn/gh-label-sync/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
