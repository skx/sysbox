package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type cppCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags

	// regular expression to find include-statements
	include *regexp.Regexp

	// regular expression to find exec-statements
	exec *regexp.Regexp
}

// Info returns the name of this subcommand.
func (c *cppCommand) Info() (string, string) {
	return "cpp", `Trivial CPP-like preprocessor.

Details:

This command is a minimal implementation of something that looks a little
like the standard C preprocessor, cpp.

We only support two directives at the moment:

* #include <file>
  * You can use #include "file/path" if you prefer
* #execute command argument1 .. argument2
  * The command is executed via the shell, so you can pipe, etc.

Example:

Given the following file:

    before
    #include "/etc/passwd"
    #execute /bin/ls -l | grep " 2 "
    after

You can expand that via either of these two commands:

    $ sysbox cpp file.in
    $ cat file.in | sysbox cpp
`
}

// Process the contents of the given reader.
func (c *cppCommand) process(reader *bufio.Reader) {

	//
	// Read line by line.
	//
	// Usually we'd use bufio.Scanner, however that can
	// report problems with lines that are too long:
	//
	//   Error: bufio.Scanner: token too long
	//
	// Instead we use the bufio.ReadString method to avoid it.
	//
	line, err := reader.ReadString(byte('\n'))
	for err == nil {

		//
		// Do we have an #include-line?
		//
		matches := c.include.FindAllStringSubmatch(line, -1)
		for _, v := range matches {
			if len(v) > 0 {

				file := v[1]

				//
				// Now we should have a file to include
				//
				dat, derr := os.ReadFile(file)
				if derr != nil {
					fmt.Printf("error including: %s - %s\n", file, derr.Error())
					return
				}

				line = string(dat)
			}
		}

		//
		// Do we have "#execute" ?
		//
		matches = c.exec.FindAllStringSubmatch(line, -1)
		for _, v := range matches {
			if len(v) < 1 {
				continue
			}

			cmd := exec.Command("/bin/bash", "-c", v[1])
			out, derrr := cmd.CombinedOutput()
			if derrr != nil {
				fmt.Printf("Error running '%v': %s\n", cmd, derrr.Error())
				return
			}

			line = string(out)
		}

		fmt.Printf("%s", line)

		// Loop again
		line, err = reader.ReadString(byte('\n'))
	}
}

// Execute is invoked if the user specifies `cpp` as the subcommand.
func (c *cppCommand) Execute(args []string) int {

	//
	// Setup our regular expressions.
	//
	c.include = regexp.MustCompile("^#\\s*include\\s+[\"<]([^\">]+)[\">]")
	c.exec = regexp.MustCompile("^#\\s*execute\\s+([^\n\r]+)[\r\n]$")

	//
	// Read from STDIN
	//
	if len(args) == 0 {

		scanner := bufio.NewReader(os.Stdin)

		c.process(scanner)

		return 0
	}

	//
	// Otherwise each named file
	//
	for _, file := range args {

		handle, err := os.Open(file)
		if err != nil {
			fmt.Printf("error opening %s : %s\n", file, err.Error())
			return 1
		}

		reader := bufio.NewReader(handle)
		c.process(reader)
	}

	return 0
}
