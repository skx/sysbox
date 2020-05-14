package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type urlsCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags

	// Regular expression we find matches with.
	reg *regexp.Regexp
}

// Info returns the name of this subcommand.
func (u *urlsCommand) Info() (string, string) {
	return "urls", `Extract URLs from text.

Details:

This command extracts URLs from STDIN, or the named files, and
prints them.  Only http and https URLs will be extracted, and we
operate with a regular expression so we're a little naive.

Examples:

  $ echo "https://example.com/ test " | sysbox urls
  $ sysbox urls ~/Org/bookmarks.org

Limitations:

Since we're doing a naive job there are limitations, the most obvious
one is that we use a simple regular expression to find URLs.  I've
chosen break URLs when I hit a ')' or ']' character, which means markdown
files can be parsed neatly.  This does mean it is possible valid links
will be truncated.

For example Wikipedia will contain links like this, which will be truncated
incorrectly:

  http://en.wikipedia.org/...(foo)

(i.e The trailing ')' will be removed.)`
}

// Match our regular expression against the given reader
func (u *urlsCommand) process(reader *bufio.Reader) {

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
		matches := u.reg.FindAllStringSubmatch(line, -1)
		for _, v := range matches {
			if len(v) > 0 {
				fmt.Printf("%s\n", v[1])
			}
		}
		line, err = reader.ReadString(byte('\n'))
	}
}

// Execute is invoked if the user specifies `urls` as the subcommand.
func (u *urlsCommand) Execute(args []string) int {

	//
	// Naive pattern for URL matching.
	//
	// NOTE: This stops when we hit characters that are valid
	//       for example ")", "]", ",", "'", "\", etc.
	//
	//       This is helpful for Markdown documents, however it IS
	//       wrong.
	//
	pattern := "(https?://[^\\\\\"'` \n\r\t\\]\\,)]+)"
	u.reg = regexp.MustCompile(pattern)

	//
	// Read from STDIN
	//
	if len(args) == 0 {

		scanner := bufio.NewReader(os.Stdin)

		u.process(scanner)

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
		u.process(reader)
	}

	return 0
}
