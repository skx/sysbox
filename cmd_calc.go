package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/skx/subcommands"
	"github.com/skx/sysbox/calc"
)

// Structure for our options and state.
type calcCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (c *calcCommand) Info() (string, string) {
	return "calc", `A simple (floating-point) calculator.

Details:

This command allows you to evaluate simple mathematical operations,
with support for floating-point operations - something the standard
'expr' command does not support.

Example:

   $ sysbox calc 3 + 3
   $ sysbox calc '1 / 3 * 9'

Note here we can join arguments, or accept a quoted string.  The arguments
must be quoted if you use '*' because otherwise the shell's globbing would
cause surprises.

Repl:

If you execute this command with no arguments you'll be dropped into a REPL
environment.  This environment is almost 100% identical to the non-interactive
use, with the exception that you can define variables:

   $ sysbox calc
   calc> let a = 3
   3
   calc> a * 3
   9
   calc> a / 9
   0.3333
   calc> exit

If you prefer you can handle assignments without "let":

   calc> a = 1; b = 2 ; c = 3
   3
   calc> a + b * c
   7
   calc> exit`
}

// Show the result of a calculation
func (c *calcCommand) showResult(out *calc.Token) error {

	if out.Type == calc.ERROR {
		return fmt.Errorf("error: %s", out.Value.(string))
	}
	if out.Type != calc.NUMBER {
		return fmt.Errorf("unexpected output (not a number): %v", out)
	}

	//
	// Show the result as an int, if possible.
	//
	result := out.Value.(float64)
	if float64(int(result)) == result {
		fmt.Printf("%d\n", int(result))
		return nil
	}

	//
	// strip trailing "0"
	//
	// First convert to string, then remove each
	// final zero.
	output := fmt.Sprintf("%f", result)
	for strings.HasSuffix(output, "0") {
		output = strings.TrimSuffix(output, "0")
	}
	fmt.Printf("%s\n", output)
	return nil
}

// Execute is invoked if the user specifies `calc` as the subcommand.
func (c *calcCommand) Execute(args []string) int {

	//
	// Join all arguments, in case we have been given "3", "+", "4".
	//
	input := ""

	for _, arg := range args {
		input += arg
		input += " "
	}

	//
	// Create a new evaluator
	//
	cal := calc.New()

	//
	// If we have no arguments then we're in the repl.
	//
	// Otherwise we process the input.
	//
	if len(input) > 0 {

		//
		// Load the script
		//
		cal.Load(input)

		//
		// Run it.
		//
		out := cal.Run()

		//
		// Show the result.
		//
		err := c.showResult(out)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			return 1
		}

		return 0
	}

	//
	// Repl.
	//
	scanner := bufio.NewScanner(os.Stdin)

	//
	// Show the prompt and read the lines
	//
	fmt.Printf("calc> ")
	for scanner.Scan() {

		//
		// Get the input, and trim it
		//
		input := scanner.Text()
		input = strings.TrimSpace(input)

		//
		// Exit ?
		//
		if strings.HasPrefix(input, "exit") ||
			strings.HasPrefix(input, "quit") {
			return 0
		}

		//
		// Ignore it, unless it is non-empty
		//
		if input != "" {

			//
			// Load the script
			//
			cal.Load(input)

			//
			// Run it.
			//
			out := cal.Run()

			//
			// Show the result.
			//
			err := c.showResult(out)
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}

		}

		// Show the prompt again, before looping back for
		// more input.
		fmt.Printf("calc> ")
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	//
	// All done
	//
	return 0
}
