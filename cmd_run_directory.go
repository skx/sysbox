package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Structure for our options and state.
type runDirectoryCommand struct {

	// Exit on error?
	exit bool

	// Be verbose?
	verbose bool
}

// Arguments adds per-command args to the object.
func (rd *runDirectoryCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&rd.exit, "exit", false, "Exit if any command terminates with a non-zero exit-code")
	f.BoolVar(&rd.verbose, "verbose", false, "Be verbose.")
}

// Info returns the name of this subcommand.
func (rd *runDirectoryCommand) Info() (string, string) {
	return "run-directory", `Run all the executables in a directory.

Details:

This command allows you to run each of the (executable) files in a given
directory.

Optionally you can terminate processing if any of the executables exit
with a non-zero exit-code.`
}

// IsExecutable returns true if the given path points to an executable file.
func (rd *runDirectoryCommand) IsExecutable(path string) bool {
	d, err := os.Stat(path)
	if err == nil {
		m := d.Mode()
		return !m.IsDir() && m&0111 != 0
	}
	return false
}

// RunCommand is a helper to run a command, returning output and the exit-code.
func (rd *runDirectoryCommand) RunCommand(command string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(command)
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

// RunParts runs all the executables in the given directory.
func (rd *runDirectoryCommand) RunParts(directory string) {

	//
	// Find the files beneath the named directory.
	//
	files, err := os.ReadDir(directory)
	if err != nil {
		fmt.Printf("error reading directory contents %s - %s\n", directory, err)
		os.Exit(1)
	}

	//
	// For each file we found.
	//
	for _, f := range files {

		//
		// Get the absolute path to the file.
		//
		path := filepath.Join(directory, f.Name())

		//
		// We'll skip any dotfiles.
		//
		if f.Name()[0] == '.' {
			if rd.verbose {
				fmt.Printf("Skipping dotfile: %s\n", path)
			}
			continue
		}

		//
		// We'll skip any non-executable files.
		//
		if !rd.IsExecutable(path) {
			if rd.verbose {
				fmt.Printf("Skipping non-executable %s\n", path)
			}
			continue
		}

		//
		// Show what we're doing.
		//
		if rd.verbose {
			fmt.Printf("%s - launching\n", path)
		}

		//
		// Run the command, capturing output and exit-code
		//
		stdout, stderr, exitCode := rd.RunCommand(path)

		//
		// Show STDOUT
		//
		if len(stdout) > 0 {
			fmt.Print(stdout)
		}

		//
		// Show STDERR
		//
		if len(stderr) > 0 {

			fmt.Print(stderr)
		}

		//
		// Show the duration, if we should
		//
		if rd.verbose {
			fmt.Printf("%s - completed\n", path)
		}

		//
		// If the exit-code was non-zero then we have to
		// terminate.
		//
		if exitCode != 0 {
			if rd.verbose {
				fmt.Printf("%s returned non-zero exit-code\n", path)
			}
			if rd.exit {
				os.Exit(1)
			}
		}

	}
}

// Execute is invoked if the user specifies `run-directory` as the subcommand.
func (rd *runDirectoryCommand) Execute(args []string) int {
	//
	// Ensure we have at least one argument.
	//
	if len(args) < 1 {
		fmt.Printf("Usage: run-directory <directory1> [directory2] .. [directoryN]\n")
		os.Exit(1)
	}

	//
	// Process each named directory
	//
	for _, entry := range args {
		rd.RunParts(entry)
	}

	return 0
}
