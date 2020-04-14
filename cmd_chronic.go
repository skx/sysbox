package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type chronicCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (c *chronicCommand) Info() (string, string) {
	return "chronic", `Run a command quietly, if it succeeds.

Details:

The chronic command allows you to execute a program, and hide the output
if the command succeeds.

The ideal use-case is for wrapping cronjobs, where you don't care about the
output unless the execution fails.

Example:

Compare the output of these two commands:

$ sysbox chronic ls
$

$ sysbox chronic ls /missing/dir
ls: cannot access '/missing/file': No such file or directory
`
}

// RunCommand is a helper to run a command, returning output and the exit-code.
func (c *chronicCommand) RunCommand(command []string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not
			// available in $PATH, in this situation, exit
			// code could not be get, and stderr will be
			// empty string very likely, so we use the default
			// fail code, and format err to string and set to stderr
			exitCode = 1
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return stdout, stderr, exitCode
}

// Execute is invoked if the user specifies `chronic` as the subcommand.
func (c *chronicCommand) Execute(args []string) int {

	if len(args) <= 0 {
		fmt.Printf("Usage: chronic command to execute ..\n")
		return 1
	}

	stdout, stderr, exit := c.RunCommand(args)
	if exit == 0 {
		return 0
	}

	fmt.Printf("%q exited with status code %d\n", args, exit)
	if len(stdout) > 0 {
		fmt.Println(stdout)
	}
	if len(stderr) > 0 {
		fmt.Println(stderr)
	}
	return exit
}
