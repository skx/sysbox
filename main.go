package main

import (
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
	subcommands.Register(&passwordCommand{})
	subcommands.Register(&splayCommand{})
	subcommands.Register(&SSLExpiryCommand{})
	subcommands.Register(&withLockCommand{})

	//
	// Execute the one the user chose.
	//
	subcommands.Execute()
}
