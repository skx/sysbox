package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Structure for our options and state.
type todoCommand struct {

	// The current date/time
	now time.Time

	// regular expression to find (nn/NN...)
	reg *regexp.Regexp

	// silent?
	silent bool

	// verbose?
	verbose bool
}

// Arguments adds per-command args to the object.
func (t *todoCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&t.silent, "silent", false, "Should we be silent in the case of permission-errors?")
	f.BoolVar(&t.verbose, "verbose", false, "Should we report on what we're doing?")
}

// Info returns the name of this subcommand.
func (t *todoCommand) Info() (string, string) {
	return "todo", `Flag TODO-notes past their expiry date.

Details:

This command recursively examines files beneath the current directory,
or the named directory, and outputs any comments which have an associated
date which is in the past.

Two comment-types are supported 'TODO' and 'FIXME' - these must occur
literally, and in upper-case only.   To find comments which should be
reported the line must also contain a date, enclosed in parenthesis.

The following examples show the kind of comments that will be reported
when the given date(s) are in the past:

    // TODO (10/03/2022) - Raise this after 10th March 2022.
    // TODO (03/2022)    - Raise this after March 2022.
    // TODO (02/06/2022) - Raise this after 2nd June 2022.
    // FIXME - This will break at the end of the year (2023).
    // FIXME - RootCA must be renewed & replaced before (10/2025).

Usage:

   $ sysbox todo
   $ sysbox todo ~/Projects/

`
}

// Process all the files beneath the given path
func (t *todoCommand) scanPath(path string) error {

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if !os.IsPermission(err) {
					return err
				}

				if !t.silent {
					fmt.Fprintf(os.Stderr, "permission denied: %s\n", path)
				}
				return nil
			}

			// We only want to read files
			isDir := info.IsDir()

			if !isDir {
				err := t.processFile(path)
				return err
			}

			return nil
		})

	return err
}

// processFile opens a file and reads line by line for a date.
func (t *todoCommand) processFile(path string) error {

	if t.verbose {
		fmt.Printf("examining %s\n", path)
	}

	// open the file
	file, err := os.Open(path)
	if err != nil {

		// error - is it permission-denied?  If so we can swallow that
		if os.IsPermission(err) {
			if !t.silent {
				fmt.Fprintf(os.Stderr, "permission denied opening: %s\n", path)
			}
			return nil
		}

		// ok another error
		return fmt.Errorf("failed to scan file %s:%s", path, err)
	}
	defer file.Close()

	// prepare to read the file
	scanner := bufio.NewScanner(file)

	// 64k is the default max length of the line-buffer - double it.
	const maxCapacity int = 128 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Process each line
	for scanner.Scan() {

		// The line we're operating upon.
		line := scanner.Text()

		// Does this line contain TODO, or FIXME?
		if strings.Contains(line, "TODO") || strings.Contains(line, "FIXME") {

			// remove leading/trailing space
			line = strings.TrimSpace(line)

			// Does it contain a date?
			match := t.reg.FindStringSubmatch(line)

			// OK we have a date.
			if len(match) >= 2 {

				// The date we've found
				date := match[1]

				var found time.Time

				// Split by "/" to find the number
				// of values we've got:
				//
				//   "DD/MM/YYYY"
				//   "MM/YYYY"
				//   "YYYY"
				parts := strings.Split(date, "/")

				switch len(parts) {
				case 3:
					found, err = time.Parse("02/01/2006", date)
					if err != nil {
						return fmt.Errorf("failed to parse %s:%s", date, err)
					}
				case 2:
					found, _ = time.Parse("01/2006", date)
					if err != nil {
						return fmt.Errorf("failed to parse %s:%s", date, err)
					}
				case 1:
					found, _ = time.Parse("2006", date)
					if err != nil {
						return fmt.Errorf("failed to parse %s:%s", date, err)
					}
				default:
					return fmt.Errorf("unknown date-format %s", date)
				}

				// If the date we've parsed is before today
				// then we alert on the line.
				if found.Before(t.now) {
					fmt.Printf("%s:%s\n", path, line)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Execute is invoked if the user specifies `todo` as the subcommand.
func (t *todoCommand) Execute(args []string) int {

	// Save today's date/time which we'll use for comparison.
	t.now = time.Now()

	// Create the capture regexp
	var err error
	t.reg, err = regexp.Compile(`\(([0-9/]+)\)`)
	if err != nil {
		fmt.Printf("internal error compiling regular expression:%s\n", err)
		return 1
	}

	// If we got any directories ..
	if len(args) > 0 {

		failed := false

		// process each path
		for _, path := range args {

			// error? then report it, but continue
			err = t.scanPath(path)
			if err != nil {
				fmt.Printf("error handling %s: %s\n", path, err)
				failed = true
			}
		}

		// exit-code will reveal errors
		if failed {
			return 1
		}
		return 0
	}

	// No named directory/directories - just handle the PWD
	err = t.scanPath(".")
	if err != nil {
		fmt.Printf("error handling search:%s\n", err)
		return 1
	}

	return 0
}
