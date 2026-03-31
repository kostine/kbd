package main

import (
	"github.com/kostine/kbd/cmd"
	"github.com/kostine/kbd/internal/logging"
)

// Set by goreleaser ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersion(version, commit, date)
	logging.Init()
	defer logging.Close()
	defer logging.RecoverPanic()
	cmd.Execute()
}
