package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// Structure for our options and state.
type timeoutCommand struct {
	duration int
}

// Arguments adds per-command args to the object.
func (t *timeoutCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&t.duration, "timeout", 300, "The number of seconds to let the command run for")

}

// Info returns the name of this subcommand.
func (t *timeoutCommand) Info() (string, string) {
	return "timeout", `Run a command, with a timeout.

Details:

This command allows you to execute an arbitrary command, but terminate it
after the given number of seconds.

The command is launched with a PTY to allow interactive commands to work
as expected, for example

$ sysbox timeout -duration=10 top`
}

// Execute is invoked if the user specifies `timeout` as the subcommand.
func (t *timeoutCommand) Execute(args []string) int {

	if len(args) <= 0 {
		fmt.Printf("Usage: timeout command [arg1] [arg2] ..[argN]\n")
		return 1
	}

	// Create a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.duration)*time.Second)
	defer cancel()

	// Create the command, using our context.
	c := exec.CommandContext(ctx, args[0], args[1:]...)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		fmt.Printf("Failed to launch %s\n", err.Error())
		return 1
	}

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }()

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout/stderr.
	//
	// If any of the copy commands complete kill our context, which
	// will let us stop awaiting completion.
	go func() {
		io.Copy(ptmx, os.Stdin)
		cancel()
	}()
	go func() {
		io.Copy(os.Stdout, ptmx)
		cancel()
	}()
	go func() {
		io.Copy(os.Stderr, ptmx)
		cancel()
	}()

	//
	// Wait for our command to complete.
	//
	<-ctx.Done()

	return 0
}
