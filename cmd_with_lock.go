package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nightlyone/lockfile"
)

// Structure for our options and state.
type withLockCommand struct {

	// prefix is the directory-root beneath which we write our lockfile.
	prefix string
}

// Arguments adds per-command args to the object.
func (wl *withLockCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&wl.prefix, "prefix", "/var/tmp", "The location beneath which to write our lockfile")
}

// Info returns the name of this subcommand.
func (wl *withLockCommand) Info() (string, string) {
	return "with-lock", `Execute a process, with a lock.

Details:

This command allows you to execute a command, with arguments,
using a lockfile.  This will prevent multiple concurrent executions
of the same command.

The expected use-case is to prevent overlapping executions of cronjobs,
etc.

Implementation:

A filename is constructed based upon the command to be executed, and
this is used to prevent the concurrent execution.  The command, and
arguments, to be executed are passed through a SHA1 hash for consistency.`
}

// Execute is invoked if the user specifies `with-lock` as the subcommand.
func (wl *withLockCommand) Execute(args []string) int {

	//
	// Ensure we have an argument
	//
	if len(args) < 1 {
		fmt.Printf("You must specify the command to execute\n")
		return 1
	}

	//
	// Generate a lockfile
	//
	h := sha1.New()
	for i, arg := range args {
		h.Write([]byte(fmt.Sprintf("%d:%s", i, arg)))
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))

	//
	// Create the lockfile
	//
	lock, err := lockfile.New(filepath.Join(wl.prefix, string(hash)))
	if err != nil {
		fmt.Printf("Cannot init lock. reason: %v", err)
		return 1
	}

	// Error handling is essential, as we only try to get the lock.
	if err = lock.TryLock(); err != nil {
		fmt.Printf("Cannot lock %q, reason: %v", lock, err)
		return 1
	}

	defer func() {
		if errr := lock.Unlock(); errr != nil {
			fmt.Printf("Cannot unlock %q, reason: %v", lock, errr)
			os.Exit(1)
		}
	}()

	//
	// Run the command.
	//
	cmd := exec.Command(args[0], args[1:]...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()

	if len(stdout.String()) > 0 {
		fmt.Print(stdout.String())
	}
	if len(stderr.String()) > 0 {
		fmt.Print(stderr.String())
	}
	if err != nil {
		fmt.Printf("Error running command:%s\n", err.Error())
		return 1
	}

	return 0
}
