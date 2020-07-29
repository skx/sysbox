package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

// Comment is a structure to hold a language/syntax for comments.
//
// A comment is denoted as the content between a start-marker, and an
// end-marker.  For single-line comments we define the end-marker as
// being a newline.
type Comment struct {

	// The text which denotes the start of a comment.
	//
	// For C++ this might be `/*`, for a shell-script it might be `#`.
	start string

	// The text which denotes the end of a comment.
	//
	// For C++ this might be `*/`, for a shell-script it might be `\n`.
	end string
}

// Structure for our options and state.
type commentsCommand struct {

	// The styles of comments to be enabled, as set by the command-line.
	style string

	// Pretty-print the comments?
	pretty bool

	// The comments we're matching
	patterns []Comment
}

// Arguments adds per-command args to the object.
func (cc *commentsCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&cc.style, "style", "c,cpp", "A comma-separated list of the comment-styles to use")
	f.BoolVar(&cc.pretty, "pretty", false, "Reformat comments for readability")

}

// Info returns the name of this subcommand.
func (cc *commentsCommand) Info() (string, string) {
	return "comments", `Output the comments contained in the given file.

Details:

This naive command outputs the comments which are included in the specified
filename(s). This is useful if you wish to run spell-checkers, etc.

There is support for outputting single-line and multi-line comments for C,
C++, Lua, and Golang.  Additional options are welcome.  By default C, and
C++ are enabled.  To only use Lua comments you could run:

    $ sysbox comments --style=lua *.lua`
}

// showComment writes the comment to the console, after optionally tidying
func (cc *commentsCommand) showComment(comment string) {
	if cc.pretty {
		// Remove newlines
		comment = strings.Replace(comment, "\n", " ", -1)

		// Remove " * "
		comment = strings.Replace(comment, " * ", " ", -1)

		// Collapse adjacent spaces
		comment = strings.Join(strings.Fields(comment), " ")

		// Skip empty comments; i.e. just literal matches of
		// the opening pattern.
		for _, pattern := range cc.patterns {
			if comment == pattern.start {
				return
			}
		}
	}

	// Remove trailing newline, so we can safely add one
	comment = strings.TrimSuffix(comment, "\n")
	fmt.Printf("%s\n", comment)
}

// dumpComments dumps the comments from the given file.
func (cc *commentsCommand) dumpComments(filename string) {

	// Read the content
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error reading %s: %s\n", filename, err.Error())
		return
	}

	// Convert our internal patterns to a series of regular expressions.
	var r []*regexp.Regexp
	for _, pattern := range cc.patterns {
		reg := "(?s)" + regexp.QuoteMeta(pattern.start) + "(.*?)" + regexp.QuoteMeta(pattern.end)
		r = append(r, regexp.MustCompile(reg))
	}

	// Now for each regexp do the matching over the whole input.
	for _, re := range r {
		out := re.FindAllSubmatch(content, -1)
		for _, match := range out {
			cc.showComment(string(match[0]))
		}
	}

}

// Execute is invoked if the user specifies `comments` as the subcommand.
func (cc *commentsCommand) Execute(args []string) int {

	// Map of known patterns, by name
	known := make(map[string][]Comment)

	// Populate with the patterns.
	known["ada"] = []Comment{Comment{start: "--", end: "\n"}}
	known["apl"] = []Comment{Comment{start: "⍝", end: "\n"}}
	known["applescript"] = []Comment{Comment{start: "(*", end: "*)"},
		Comment{start: "--", end: "\n"}}
	known["asm"] = []Comment{Comment{start: ";", end: "\n"}}
	known["basic"] = []Comment{Comment{start: "REM", end: "\n"}}
	known["c"] = []Comment{Comment{start: "//", end: "\n"}}
	known["coldfusion"] = []Comment{Comment{start: "<!---", end: "--->"}}
	known["cpp"] = []Comment{Comment{start: "/*", end: "*/"}}
	known["go"] = []Comment{Comment{start: "/*", end: "*/"},
		Comment{start: "//", end: "\n"},
	}
	known["html"] = []Comment{Comment{start: "<!--", end: "-->"}}
	known["haskell"] = []Comment{Comment{start: "{-", end: "-}"},
		Comment{start: "--", end: "\n"}}
	known["lisp"] = []Comment{Comment{start: ";", end: "\n"}}
	known["java"] = []Comment{Comment{start: "/*", end: "*/"},
		Comment{start: "//", end: "\n"}}
	known["javascript"] = []Comment{Comment{start: "/*", end: "*/"},
		Comment{start: "//", end: "\n"}}
	known["lua"] = []Comment{Comment{start: "--[[", end: "--]]"},
		Comment{start: "-- ", end: "\n"}}
	known["matlab"] = []Comment{Comment{start: "%{", end: "%}"},
		Comment{start: "% ", end: "\n"}}
	known["pascal"] = []Comment{Comment{start: "(*", end: "*)"}}
	known["perl"] = []Comment{Comment{start: "#", end: "\n"}}
	known["php"] = []Comment{Comment{start: "/*", end: "*/"},
		Comment{start: "//", end: "\n"},
		Comment{start: "#", end: "\n"},
	}
	known["python"] = []Comment{Comment{start: "#", end: "\n"}}
	known["ruby"] = []Comment{Comment{start: "#", end: "\n"}}
	known["shell"] = []Comment{Comment{start: "#", end: "\n"}}
	known["swift"] = []Comment{Comment{start: "/*", end: "*/"},
		Comment{start: "//", end: "\n"}}
	known["sql"] = []Comment{Comment{start: "--", end: "\n"}}
	known["xml"] = []Comment{Comment{start: "<!--", end: "-->"}}

	// Ensure we have at least one filename specified.
	if len(args) <= 0 {
		fmt.Printf("Usage: comments file1 [file2] ..[argN]\n")
		return 1
	}

	// Load the patterns the user selected.
	for _, kind := range strings.Split(cc.style, ",") {

		// Lookup the choice
		pat, ok := known[kind]

		// Not found?  That's an error
		if !ok {
			fmt.Printf("Unknown style %s, valid options include:\n", kind)

			keys := make([]string, 0, len(known))
			for k := range known {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				fmt.Printf("\t%s\n", k)
			}
			return 1
		}

		// Otherwise add it to the list.
		cc.patterns = append(cc.patterns, pat...)
	}

	// Now process the input file(s)
	for _, file := range args {
		cc.dumpComments(file)
	}

	// All done.
	return 0
}
