package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// Structure for our options and state.
type splayCommand struct {
	max     int
	verbose bool
}

// Arguments adds per-command args to the object.
func (r *splayCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&r.max, "maximum", 300, "The maximum amount of time to sleep for")
	f.BoolVar(&r.verbose, "verbose", false, "Should we be verbose")

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

Give each script a random-delay via adding a call to the splay subcommand.`
}

// Execute is invoked if the user specifies `version` as the subcommand.
func (s *splayCommand) Execute(args []string) int {

	// Ensure we seed our random number generator.
	rand.Seed(time.Now().UnixNano())

	// Get the delay-time.
	delay := rand.Intn(s.max)
	if s.verbose {
		fmt.Printf("Sleeping for for %d seconds, from max splay-time of %d\n", delay, s.max)
	}

	// Sleep
	time.Sleep(time.Duration(delay) * time.Second)
	return 0
}
