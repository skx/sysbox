//go:build go1.18
// +build go1.18

package main

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/skx/subcommands"
)

var (
	version = "unreleased"
)

// Structure for our options and state.
type versionCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (t *versionCommand) Info() (string, string) {
	return "version", `Show the version of this binary.

Details:

This reports upon the version of the application.
`
}

// Execute is invoked if the user specifies `version` as the subcommand.
func (t *versionCommand) Execute(args []string) int {

	fmt.Printf("%s\n", version)

	info, ok := debug.ReadBuildInfo()

	if ok {
		for _, settings := range info.Settings {
			if strings.Contains(settings.Key, "vcs") {
				fmt.Printf("%s: %s\n", settings.Key, settings.Value)
			}
		}
	}

	return 0
}
