package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type installCommand struct {

	// Path to the binary
	binary string

	// Installation directory
	directory string

	// force creation?
	force bool
}

// Arguments adds per-command args to the object.
func (i *installCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&i.binary, "binary", "sysbox", "The path to the sysbox-executable")
	f.StringVar(&i.directory, "directory", "/usr/local/bin", "The directory within which to create the symlinks.")
	f.BoolVar(&i.force, "force", false, "Force creation?")

}

// Info returns the name of this subcommand.
func (i *installCommand) Info() (string, string) {
	return "install", `Create symlinks for each known binary.

Details:

The sysbox-executable has support for running a variety of sub-commands,
which are specified as the first argument to the main binary (and are
then followed by command-specific options).

To save type you can also run the subcommand "foo" by creating a symlink
from the name "foo" to the sysbox executable.

Example:

     sysbox install -binary=$(pwd)/sysbox -directory=/usr/local/bin

This will output the commands to create the symlinks, which you can execute
like so:

     sysbox install -binary=$(pwd)/sysbox | sudo sh`

}

// Does the target file/directory exist?
func (i *installCommand) Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Execute is invoked if the user specifies `install` as the subcommand.
func (i *installCommand) Execute(args []string) int {

	// Ensure we have arguments
	if !i.Exists(i.binary) {
		fmt.Printf("Binary %s does not exist.", i.binary)
		return 1
	}
	if !i.Exists(i.directory) {
		fmt.Printf("The target directory %s does not exist", i.directory)
		return 1
	}

	// Force creation of the symlinks?
	force := ""
	if i.force {
		force = "f"
	}

	// For each command
	for _, cmd := range subcommands.Commands() {

		// Skip `help`, `install`
		if cmd == "help" || cmd == "install" {
			continue
		}

		// Work out the target
		target := filepath.Join(i.directory, cmd)

		fmt.Printf("ln -s%s %s %s\n", force, i.binary, target)
	}

	return 0
}
