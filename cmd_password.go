package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// Structure for our options and state.
type passwordCommand struct {

	// The length of the password to generate
	length int

	// Specials?
	specials bool

	// Digits?
	digits bool
}

// Arguments adds per-command args to the object.
func (p *passwordCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&p.length, "length", 15, "The length of the password to generate")
	f.BoolVar(&p.specials, "specials", true, "Should we use special characters?")
	f.BoolVar(&p.digits, "digits", true, "Should we use digits?")
}

// Info returns the name of this subcommand.
func (p *passwordCommand) Info() (string, string) {
	return "make-password", `Generate a random password.

Details:

This command generates a simple random password, by default being 12
characters long.  You can tweak the alphabet used via the command-line
flags if necessary.`
}

// Execute is invoked if the user specifies `make-password` as the subcommand.
func (p *passwordCommand) Execute(args []string) int {
	rand.Seed(time.Now().UnixNano())

	// Alphabets we use for generation
	digits := "0123456789"
	specials := "~=&+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz"

	// Extend our alphabet if we should
	if p.digits {
		all = all + digits
	}
	if p.specials {
		all = all + specials
	}

	// Make a buffer and fill it with all characters
	buf := make([]byte, p.length)
	for i := 0; i < p.length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	// Add a digit if we should.
	//
	// We might already have them present, because our `all`
	// alphabet was used already.  But this ensures we have at
	// least one digit present.
	if p.digits {
		buf[0] = digits[rand.Intn(len(digits))]
	}

	// Add a special-character if we should.
	//
	// We might already have them present, because our `all`
	// alphabet was used already.  But this ensures we have at
	// least one special-character present.
	if p.specials {
		buf[1] = specials[rand.Intn(len(specials))]
	}

	// Shuffle and output
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	fmt.Printf("%s\n", buf)

	return 0
}
