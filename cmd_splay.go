package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Structure for our options and state.
type splayCommand struct {
	max     int
	verbose bool
}

// Arguments adds per-command args to the object.
func (s *splayCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&s.max, "maximum", 300, "The maximum amount of time to sleep for")
	f.BoolVar(&s.verbose, "verbose", false, "Should we be verbose")

}

// Info returns the name of this subcommand.
func (s *splayCommand) Info() (string, string) {
	return "splay", `Sleep for a random time.

Details:

This command allows you to stagger execution of things via the introduction
of random delays.

The expected use-case is that you have a number of hosts which each wish
to perform a cron-job, but you don't want to overwhelm a central system
by having all those events occur at precisely the same time (which is
likely to happen if you're running with good clocks).

Give each script a random-delay via adding a call to the splay subcommand.

Usage:

We prefer users to specify the splay-time with a parameter, but to allow
natural usage you may specify as the first argument:

   $ sysbox splay --maximum=10 [-verbose]
   $ sysbox splay 10 [-verbose]`
}

// Execute is invoked if the user specifies `splay` as the subcommand.
func (s *splayCommand) Execute(args []string) int {

	// Ensure we seed our random number generator.
	rand.Seed(time.Now().UnixNano())

	// If the user gave an argument then use it.
	//
	// Because people might expect this to work.
	if len(args) > 0 {

		// First argument will be a number
		num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("error converting %s to integer: %s\n", args[0], err.Error())
		}

		// Save it away.
		s.max = num
	}

	// Get the delay-time.
	delay := rand.Intn(s.max)
	if s.verbose {
		fmt.Printf("Sleeping for for %d seconds, from max splay-time of %d\n", delay, s.max)
	}

	// Sleep
	time.Sleep(time.Duration(delay) * time.Second)
	return 0
}
