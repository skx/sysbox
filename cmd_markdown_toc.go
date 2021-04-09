package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Structure for our options and state.
type markdownTOCCommand struct {

	// Maximum level to include.
	max int
}

// tocItem holds state for a single entry
type tocItem struct {
	// Depth of the entry
	depth int

	// The content of the header (text).
	content string
}

// String converts the tocItem to a string
func (t tocItem) String() string {

	// Characters dropped from anchors
	droppedChars := []string{
		"\"", "'", "`", ".",
		"!", ",", "~", "&",
		"%", "^", "*", "#",
		"/", "\\",
		"@", "|",
		"(", ")",
		"{", "}",
		"[", "]",
	}

	// link is lowercase
	link := strings.ToLower(t.content)

	// Remove the characters
	for _, c := range droppedChars {
		link = strings.Replace(link, c, "", -1)
	}

	// Replace everything else with "-"
	link = strings.Replace(link, " ", "-", -1)
	link = "#" + link

	return fmt.Sprintf("%v* [%v](%v) \n",
		strings.Repeat(" ", 2*(t.depth-1)),
		t.content,
		link)
}

// Arguments adds per-command args to the object.
func (m *markdownTOCCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&m.max, "max", 100, "The maximum nesting level to generate.")

}

// Info returns the name of this subcommand.
func (m *markdownTOCCommand) Info() (string, string) {
	return "markdown-toc", `Create a table-of-contents for a markdown file.

Details:

This command allows you to generate a (github-themed) table of contents
for a given markdown file.


Usage:

   $ sysbox markdown-toc README.md
   $ sysbox markdown-toc < README.md`
}

// process handles the generation of the TOC from the given reader
func (m *markdownTOCCommand) process(reader *bufio.Reader) error {

	fileScanner := bufio.NewScanner(reader)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		headerCount := m.countHashes(line)

		if headerCount >= 1 && headerCount < m.max {

			// Create an item for this header
			item := tocItem{
				depth:   headerCount,
				content: line[headerCount+1:],
			}

			// Print it
			fmt.Print(item.String())
		}
	}

	if err := fileScanner.Err(); err != nil {
		return err
	}

	return nil
}

// counts hashes at the beginning of a line
func (m *markdownTOCCommand) countHashes(s string) int {
	for i, c := range s {
		if c != '#' {
			return i
		}
	}
	return len(s)
}

// Execute is invoked if the user specifies `markdown-toc` as the subcommand.
func (m *markdownTOCCommand) Execute(args []string) int {

	var err error

	// No file?  Use STDIN
	if len(args) == 0 {

		scanner := bufio.NewReader(os.Stdin)
		err = m.process(scanner)

		if err != nil {
			fmt.Printf("error processing STDIN - %s\n", err.Error())
			return 1
		}
		return 0
	}

	// Otherwise each named file
	for _, file := range args {

		handle, err2 := os.Open(file)
		if err2 != nil {
			fmt.Printf("error opening %s: %s\n", file, err2.Error())
			return 1
		}

		reader := bufio.NewReader(handle)
		err = m.process(reader)
		if err != nil {
			fmt.Printf("error processing %s: %s\n", file, err.Error())
			return 1
		}
	}

	return 0
}
