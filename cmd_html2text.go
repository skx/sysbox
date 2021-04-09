package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skx/subcommands"
	"golang.org/x/net/html"
)

// Structure for our options and state.
type html2TextCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (h2t *html2TextCommand) Info() (string, string) {
	return "html2text", `HTML to text conversion.

This command converts the contents of STDIN, or the named files,
from HTML to text, and prints them to the console.

Examples:

  $ curl --silent https://steve.fi/ | sysbox html2text | less
  $ sysbox html2text /usr/share/doc/gdisk/gdisk.html |less


`
}

func (h2t *html2TextCommand) process(reader *bufio.Reader) {

	domDocTest := html.NewTokenizer(reader)
	previousStartTokenTest := domDocTest.Token()
loopDomTest:

	for {
		tt := domDocTest.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDomTest // End of the document,  done
		case tt == html.StartTagToken:
			previousStartTokenTest = domDocTest.Token()
		case tt == html.TextToken:
			if previousStartTokenTest.Data == "script" ||
				previousStartTokenTest.Data == "style" {
				continue
			}
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDocTest.Text())))
			if len(TxtContent) > 0 {
				fmt.Printf("%s\n", TxtContent)
			}
		}
	}
}

// Execute is invoked if the user specifies `html2text` as the subcommand.
func (h2t *html2TextCommand) Execute(args []string) int {

	//
	// Read from STDIN
	//
	if len(args) == 0 {

		scanner := bufio.NewReader(os.Stdin)

		h2t.process(scanner)

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
		h2t.process(reader)
	}
	return 0
}
