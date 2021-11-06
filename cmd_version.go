package main

//go:generate echo Hello, Go Generate!
import (
	"fmt"

	"github.com/skx/subcommands"
)

var (
	versionString = "unreleased"
)

// Structure for our options and state.
type versionCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (t *versionCommand) Info() (string, string) {
	return "version", `Show the version of the binary.

Details:

This reports upon the version of the sysbox application.

Usage:

   $ sysbox version

`
}

// Execute is invoked if the user specifies `version` as the subcommand.
func (t *versionCommand) Execute(args []string) int {

	fmt.Printf("%s\n", versionString)

	return 0
}
