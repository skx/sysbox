package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

// Structure for our options and state.
type commentsCommand struct {

	// Single-line
	c bool

	// Multi-line
	cpp bool

	// Shell, single-line
	shell bool

	// Pretty-print the comments?
	pretty bool
}

// Arguments adds per-command args to the object.
func (cc *commentsCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&cc.c, "c", true, "Output C-style comments, prefixed with '//'.")
	f.BoolVar(&cc.cpp, "cpp", true, "Output C++-style comments, between '/*' and '*/'")
	f.BoolVar(&cc.shell, "shell", false, "Output shell-style comments, prefixed with '#'")
	f.BoolVar(&cc.pretty, "pretty", false, "Reformat comments for readability")

}

// Info returns the name of this subcommand.
func (cc *commentsCommand) Info() (string, string) {
	return "comments", `Output the comments contained in the given file.

Details:

This naive command outputs the comments which are included in the specified
filename(s).

This is useful if you wish to run spell-checkers, etc.`
}

// showComment writes the comment to the console, after optionally tidying
func (cc *commentsCommand) showComment(comment string) {
	if cc.pretty {
		// Remove newlines
		comment = strings.Replace(comment, "\n", "", -1)

		// Remove " * "
		comment = strings.Replace(comment, " * ", " ", -1)

		// Collapse adjacent spaces
		comment = strings.Join(strings.Fields(comment), " ")

		// Skip empty comments
		if comment == "//" || comment == "#" || comment == "/* */" {
			return
		}
	}

	// Remove trailing newline, so we can safely add one
	comment = strings.TrimSuffix(comment, "\n")
	fmt.Printf("%s\n", comment)
}

// dumpComments dumps the comments from the given file.
func (cc *commentsCommand) dumpComments(filename string) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", filename, err.Error())
		return
	}

	// Offset of the file-contents we're looking at
	offset := 0

	// Are we inside a single/multiline comment at the moment?
	insideShell := false
	insideSingle := false
	insideMultiline := false

	// Current comment
	comment := ""

	// Walk the contents
	for offset < (len(content) - 1) {

		// Get the current character, and the next character
		c := content[offset]
		cn := content[offset+1]

		// If we're currently inside a comment add on the character
		if insideShell || insideSingle || insideMultiline {
			comment += string(c)
		}

		// If we're inside a single-line comment then
		// look for the end
		if insideSingle && c == byte('\n') {
			insideSingle = false
			offset++

			// Show the comment, after tidying.
			cc.showComment(comment)

			comment = ""
			continue
		}

		if insideShell && c == byte('\n') {
			insideShell = false
			offset++

			// Show the comment, after tidying.
			cc.showComment(comment)

			comment = ""
			continue
		}

		// If we're inside a multiline-line comment then
		// look for the end.
		if insideMultiline && c == byte('*') && cn == byte('/') {
			insideMultiline = false
			offset++

			// Show the comment, after tidying.
			comment += "*/"
			cc.showComment(comment)

			comment = ""
			continue
		}

		//
		// OK we're not inside a comment.
		//
		// Has a comment just opened up?
		//
		if cc.c && c == byte('/') && cn == byte('/') {
			insideSingle = true
			comment = "/"
		}

		if cc.shell && c == byte('#') {
			insideShell = true
			comment = "#"
		}

		if cc.cpp && c == byte('/') && cn == byte('*') {
			insideMultiline = true
			comment = "/"
		}

		offset++
	}
}

// Execute is invoked if the user specifies `comments` as the subcommand.
func (cc *commentsCommand) Execute(args []string) int {

	if len(args) <= 0 {
		fmt.Printf("Usage: comments file1 [file2] ..[argN]\n")
		return 1
	}

	for _, file := range args {
		cc.dumpComments(file)
	}
	return 0
}
