package main

import (
	"os"

	"github.com/skx/subcommands"
)

//
// Register the subcommands, and run the one the user chose.
//
func main() {

	//
	// Register each of our subcommands.
	//
	subcommands.Register(&collapseCommand{})
	subcommands.Register(&envTemplateCommand{})
	subcommands.Register(&httpdCommand{})
	subcommands.Register(&installCommand{})
	subcommands.Register(&ipsCommand{})
	subcommands.Register(&passwordCommand{})
	subcommands.Register(&peerdCommand{})
	subcommands.Register(&runDirectoryCommand{})
	subcommands.Register(&splayCommand{})
	subcommands.Register(&SSLExpiryCommand{})
	subcommands.Register(&validateJSONCommand{})
	subcommands.Register(&validateYAMLCommand{})
	subcommands.Register(&withLockCommand{})

	//
	// Execute the one the user chose.
	//
	os.Exit(subcommands.Execute())
}
