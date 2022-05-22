package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// Structure for our options and state.
type findCommand struct {

	// Show the names of matching directories?
	directories bool

	// Show the names of matching files?
	files bool

	// Starting path
	path string

	// Ignore errors?
	silent bool
}

// Info returns the name of this subcommand.
func (fc *findCommand) Info() (string, string) {
	return "find", `Trivial file-finder.

Details:

This command is a minimal 'find' replacement, allowing you to find
files by regular expression.  Basic usage defaults to finding files,
by regular expression:

    $ sysbox find 'foo*.go' '_test'

To find directories instead of files:

    $ sysbox find -files=false -directories=true 'blah$'

Or both:

    $ sysbox find -path=/etc -files=true -directories=true 'blah$'
`
}

// Arguments adds per-command args to the object.
func (fc *findCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&fc.files, "files", true, "Show the names of matching files?")
	f.BoolVar(&fc.directories, "directories", false, "Show the names of matching directories?")
	f.BoolVar(&fc.silent, "silent", true, "Ignore permission-denied errors when recursing into unreadable directories?")
	f.StringVar(&fc.path, "path", ".", "Starting path for search.")
}

// find runs the find operation
func (fc *findCommand) find(patterns []string) error {

	// build up a list of regular expressions
	regs := []*regexp.Regexp{}

	for _, pattern := range patterns {

		reg, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("failed to compile %s:%s", pattern, err)
		}

		regs = append(regs, reg)
	}

	//
	// Walk the filesystem
	//
	err := filepath.Walk(fc.path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if !os.IsPermission(err) {
					return err
				}

				if !fc.silent {
					fmt.Fprintln(os.Stderr, "permission denied handling "+path)
				}
			}

			// We have a path.
			//
			// If it doesn't match any of our regexps then we return.
			//
			// i.e. We must match ALL supplied patterns, not just
			// one of them.
			//
			for _, r := range regs {
				if !r.MatchString(path) {
					return nil
				}
			}

			// is it a file?
			isDir := info.IsDir()
			isFile := !isDir

			if (isDir && fc.directories) ||
				(isFile && fc.files) {
				fmt.Printf("%s\n", path)
			}
			return nil
		})

	if err != nil {
		return fmt.Errorf("error walking filesystem %s", err)
	}

	return nil
}

// Execute is invoked if the user specifies `find` as the subcommand.
func (fc *findCommand) Execute(args []string) int {

	//
	// Build up the list of patterns
	//
	patterns := []string{}
	patterns = append(patterns, args...)

	//
	// Ensure we have a least one.
	//
	if len(patterns) < 1 {
		fmt.Printf("Usage: sysbox find pattern1 [pattern2..]\n")
		return 1
	}

	//
	// Run the find.
	//
	err := fc.find(patterns)
	if err != nil {
		fmt.Printf("%s\n", err)
		return 1
	}

	//
	// All done
	//
	return 0
}
