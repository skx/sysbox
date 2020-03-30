package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type collapseCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (c *collapseCommand) Info() (string, string) {
	return "collapse", `Remove whitespace from input.

Details:

This command reads input and removes all leading and trailing whitespace
from it.  Empty lines are also discarded.`
}

// Execute is invoked if the user specifies `collapse` as the subcommand.
func (c *collapseCommand) Execute(args []string) int {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return 1
	}

	return 0
}
