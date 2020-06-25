package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/skx/sysbox/templatedcmd"
)

// Structure for our options and state.
type execSTDINCommand struct {

	// testing the command
	dryRun bool

	// parallel job count
	parallel int

	// verbose flag
	verbose bool

	// field separator
	split string
}

// Command holds a command we're going to execute in a worker-process.
//
// (Command in this sense is a system-binary / external process.)
type Command struct {

	// args holds the command + args to execute.
	args []string
}

// Arguments adds per-command args to the object.
func (es *execSTDINCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&es.dryRun, "dry-run", false, "Don't run the command.")
	f.BoolVar(&es.verbose, "verbose", false, "Be verbose.")
	f.IntVar(&es.parallel, "parallel", 1, "How many jobs to run in parallel.")
	f.StringVar(&es.split, "split", "", "Split on a different character.")

}

// worker reads a command to execute from the channel, and executes it.
//
// The result is pushed back, but ignored.
func (es *execSTDINCommand) worker(id int, jobs <-chan Command, results chan<- int) {
	for j := range jobs {

		// Run the command, and get the output?
		cmd := exec.Command(j.args[0], j.args[1:]...)
		out, errr := cmd.CombinedOutput()

		// error?
		if errr != nil {
			fmt.Printf("Error running '%s': %s\n", strings.Join(j.args, " "), errr.Error())
		} else {

			// Show the output
			fmt.Printf("%s", out)
		}

		// Send a result to our output channel.
		results <- 1
	}
}

// Info returns the name of this subcommand.
func (es *execSTDINCommand) Info() (string, string) {
	return "exec-stdin", `Execute a command for each line of STDIN.

Details:

This command reads lines from STDIN, and executes the specified command with
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

  $ cat /etc/passwd | sysbox exec-stdin -split=: groups {1}

If you wish you can run the commands in parallel, using the -parallel flag
to denote how many simultaneous executions are permitted.

The only other flag is '-verbose', to show the command that would be
executed and '-dry-run' to avoid running anything.`
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
	// Ensure we have a command.
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
	// The jobs we're going to add.
	//
	// We save these away so that we can allow parallel execution later.
	//
	var toRun []Command

	//
	// Read a line
	//
	line, err := scanner.ReadString(byte('\n'))
	for err == nil && line != "" {

		//
		// Create the command to execute
		//
		run := templatedcmd.Expand(cmd, line, es.split)

		//
		// Show command if being verbose
		//
		if es.verbose || es.dryRun {
			fmt.Printf("%s\n", strings.Join(run, " "))
		}

		//
		// If we're not in "pretend"-mode then we'll save the
		// constructed command away.
		//
		if !es.dryRun {
			toRun = append(toRun, Command{args: run})
		}

		//
		// Loop again
		//
		line, err = scanner.ReadString(byte('\n'))
	}

	//
	// We've built up all the commands we're going to run now.
	//
	// Get the number, and create suitable channels.
	//
	num := len(toRun)
	jobs := make(chan Command, num)
	results := make(chan int, num)

	//
	// Launch the appropriate number of parallel workers.
	//
	for w := 1; w <= es.parallel; w++ {
		go es.worker(w, jobs, results)
	}

	//
	// Add all the pending jobs.
	//
	for _, j := range toRun {
		jobs <- j
	}
	close(jobs)

	//
	// Await all the results.
	//
	for a := 1; a <= num; a++ {
		<-results
	}
	return 0
}
