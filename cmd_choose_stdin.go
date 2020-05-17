package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/skx/sysbox/chooseui"
)

// Structure for our options and state.
type chooseSTDINCommand struct {
	// Command to execute
	exec string

	// Filenames we'll let the user choose between
	stdin []string
}

// Arguments adds per-command args to the object.
func (cs *chooseSTDINCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&cs.exec, "execute", "", "Command to execute once a selection has been made")
}

func (cs *chooseSTDINCommand) Info() (string, string) {
	return "choose-stdin", `Choose an item from STDIN, interactively

Details:

This command presents a simple UI, showing all the lines read from STDIN.

You can navigate with the keyboard, and press RETURN to select an entry.

Optionally you can press TAB to filter the list via an input field.

Uses:

This is ideal for choosing videos, roms, etc.  For example launch the
given video file:

   $ find . -name '*.avi' -print | sysbox choose-stdin -exec 'xine "{}"'

See also 'sysbox help choose-file'.`
}

// Execute is invoked if the user specifies `choose-stdin` as the subcommand.
func (cs *chooseSTDINCommand) Execute(args []string) int {

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
		// Save this away
		//
		cs.stdin = append(cs.stdin, line)

		//
		// Loop again
		//
		line, err = scanner.ReadString(byte('\n'))
	}

	//
	// Launch the UI
	//
	chooser := chooseui.New(cs.stdin)
	choice := chooser.Choose()

	//
	// Did something get chosen?  If not terminate
	//
	if choice == "" {

		return 1
	}

	//
	// Are we executing?
	//
	if cs.exec != "" {

		//
		// Split into command and arguments
		//
		pieces := strings.Fields(cs.exec)

		//
		// Expand the args - this is horrid
		//
		// Is a hack to ensure that things work if we
		// have a selected filename with spaces inside it.
		//
		toRun := []string{}

		for _, piece := range pieces {
			piece = strings.ReplaceAll(piece, "{}", choice)
			toRun = append(toRun, piece)
		}

		//
		// Run it.
		//
		cmd := exec.Command(toRun[0], toRun[1:]...)
		out, errr := cmd.CombinedOutput()
		if errr != nil {
			fmt.Printf("Error running '%s': %s\n", cs.exec, errr.Error())
			return 1
		}

		//
		// Show the output
		//
		fmt.Printf("%s", out)

		//
		// Otherwise we're done
		//
		return 0

	}

	fmt.Printf("%s\n", choice)
	return 0
}
