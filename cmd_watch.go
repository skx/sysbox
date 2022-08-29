package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type watchCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (w *watchCommand) Info() (string, string) {
	return "watch", `Watch the output of a command.

Details:

This command allows you execute a command every five seconds,
and see the output.

It is included because Mac OS does not include a watch-command
by default.

Notes:

Between executing the specified command the utility will
clear thes creen by executing 'cls' or 'clear', which is
a terrible approach.

In the future this command might be reimplemented using
a TUI instead, to avoid this, but for the moment it is a quick
hack.
`
}

// clearScreen clears the screen, in a horrid fashion.
func (w *watchCommand) clearScreen() {
	switch runtime.GOOS {
	case "windows":
		w.runCmd("cmd", "/c", "cls")
	default:
		w.runCmd("clear")
	}
}

// runCmd runs a command.
func (w *watchCommand) runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Execute is invoked if the user specifies `watch` as the subcommand.
func (w *watchCommand) Execute(args []string) int {

	if len(args) < 1 {
		fmt.Printf("Usage: watch cmd arg1 arg2 .. argN")
		return 1
	}

	// Run forever..
	for {

		// clear the screen, horridly
		w.clearScreen()

		// Run the command
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("cmd.Run() failed with %s\n", err)
			return 1
		}

		// Rerun after a delay of five seconds
		time.Sleep(5 * time.Second)
	}

	//
	// All done
	//
	return 0
}
