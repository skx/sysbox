package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Structure for our options and state.
type treeCommand struct {

	// show only directories?
	directories bool

	// show all files?
	all bool
}

// Arguments adds per-command args to the object.
func (t *treeCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&t.directories, "d", false, "Show only directories.")
	f.BoolVar(&t.all, "a", false, "Show all files, including dotfiles.")

}

// Info returns the name of this subcommand.
func (t *treeCommand) Info() (string, string) {
	return "tree", `Show filesystem contents as a tree.

Details:

This is a minimal reimplementation of the standard 'tree' command, it
supports showing a directory tree.

Usage:

   $ sysbox tree /etc/

To show only directory entries:

   $ sysbox tree -d /opt

If there were any errors encountered then the return-code will be 1, otherwise 0.`
}

// Execute is invoked if the user specifies `tree` as the subcommand.
func (t *treeCommand) Execute(args []string) int {

	//
	// Starting directory defaults to the current working directory
	//
	start := "."

	//
	// But can be changed
	//
	if len(args) > 0 {
		start = args[0]
	}

	type Entry struct {
		name      string
		error     string
		directory bool
	}

	//
	// Keep track of directory entries here.
	//
	entries := []*Entry{}

	//
	// Find the contents
	//
	filepath.Walk(start,
		func(path string, info os.FileInfo, err error) error {

			// Null info?  That probably means that the
			// destination we're trying to walk doesn't exist.
			if info == nil {
				return nil
			}

			entry := &Entry{name: path}

			if err == nil {
				switch mode := info.Mode(); {
				case mode.IsDir():
					entry.directory = true
				}
			} else {
				entry.error = err.Error()
			}
			entries = append(entries, entry)
			return nil
		})

	//
	// Did we hit an error?
	//
	error := false

	//
	// Show the entries
	//
	for _, ent := range entries {

		// showing only directories?  Then skip this
		// entry unless it is a directory
		if t.directories && !ent.directory {
			continue
		}

		// skip dotfiles by default
		if (strings.Contains(ent.name, "/.") || strings.HasPrefix(ent.name, ".")) && !t.all {
			continue
		}

		if ent.error != "" {
			fmt.Printf("%s - %s\n", ent.name, ent.error)
			error = true
			continue
		}
		fmt.Printf("%s\n", ent.name)
	}

	if error {
		return 1
	}
	return 0
}
