package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Structure for our options and state.
type execSTDINCommand struct {

	// testing the command
	dryRun bool

	// verbose flag
	verbose bool

	// field separator
	split string
}

// Arguments adds per-command args to the object.
func (es *execSTDINCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&es.dryRun, "dry-run", false, "Don't run the command.")
	f.BoolVar(&es.verbose, "verbose", false, "Be verbose")
	f.StringVar(&es.split, "split", "", "Split on a different character")

}

// Info returns the name of this subcommand.
func (es *execSTDINCommand) Info() (string, string) {
	return "exec-stdin", `Execute a command for each line of STDIN

Details:

This command reads lines from STDIN, and executes the named command with
that line as input.

The line read from STDIN will be available as '{}' and each space-separated
field will be available as {1}, {2}, etc.

Examples:

  $ echo -e "foo\tbar\nbar\tSteve" | sysbox exec-stdin echo {1}
  foo
  bar

Here you see that STDIN would contain:

  foo bar
  bar Steve

However only the first field was displayed, because {1} means the first field.

To show all input you'd run:

  $ echo -e "foo\tbar\nbar\tSteve" | sysbox exec-stdin echo {}
  foo bar
  bar Steve

Flags:

If you prefer you can split fields on specific characters, which is useful
for operating upon CSV files, or in case you wish to split '/etc/passwd' on
':' to work on usernames:

  $ sysbox exec-stdin -split=: groups {1}

The only other flag is '-verbose', to show the command that would be
executed and 'dry-run' to avoid running anything.`
}

// Execute is invoked if the user specifies `exec-stdin` as the subcommand.
func (es *execSTDINCommand) Execute(args []string) int {

	//
	// Join all arguments, in case we have been given "{1}", "{2}", etc.
	//
	cmd := ""

	for _, arg := range args {
		cmd += arg
		cmd += " "
	}

	//
	// If we have no arguments then we're in the repl.
	//
	// Otherwise we process the input.
	//
	if cmd == "" {
		fmt.Printf("Usage: sysbox exec-stdin command .. args {}..\n")
		return 1
	}

	//
	// Prepare to read line-by-line
	//
	scanner := bufio.NewReader(os.Stdin)

	//
	// Read a line
	//
	line, err := scanner.ReadString(byte('\n'))
	for err == nil && line != "" {

		//
		// Remove any leading/trailing whitespace
		//
		line = strings.TrimSpace(line)

		//
		// We're now going to build up the command
		// to execute.
		//
		//   {} -> The complete line read from STDIN
		//
		//   {1} -> The first field of the input.
		//
		//   {2} -> The second field of the input.
		//
		//   {N} -> The Nth field of the input.
		//
		fields := strings.Fields(line)

		//
		// Different split character?
		//
		if es.split != "" {
			fields = strings.Split(line, es.split)
		}

		//
		// Copy the string
		//
		sh := cmd

		//
		// Look for {NNNN}
		//
		reg := regexp.MustCompile("({[0-9]+})")
		matches := reg.FindAllStringSubmatch(sh, -1)

		//
		// For each match, perform the expansion
		//
		for _, v := range matches {

			//
			// Copy the match and remove the {}
			//
			// So we just have "1", "3", etc.
			//
			match := v[1]
			match = strings.ReplaceAll(match, "{", "")
			match = strings.ReplaceAll(match, "}", "")

			//
			// Convert the string to a number.
			//
			num, err := strconv.Atoi(match)
			if err != nil {
				fmt.Printf("failed to convert %s to number: %s", match, err.Error())
				return 1
			}

			//
			// If the field matches then we can replace it
			//
			if num >= 1 && num <= len(fields) {
				sh = strings.ReplaceAll(sh, v[1], fields[num-1])
			}

		}

		//
		// Replace "{}" with the complete input.
		//
		sh = strings.ReplaceAll(sh, "{}", line)

		//
		// Show command if being verbose
		//
		if es.verbose || es.dryRun {
			fmt.Printf("%s\n", sh)
		}

		//
		// Run, unless we're not supposed to
		//
		if !es.dryRun {

			pieces := strings.Fields(sh)
			cmd := exec.Command(pieces[0], pieces[1:]...)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error running '%s': %s\n", sh, err.Error())
				return 1
			}

			//
			// Show the output
			//
			fmt.Printf("%s", out)
		}

		//
		// Loop again
		//
		line, err = scanner.ReadString(byte('\n'))
	}

	return 0
}
