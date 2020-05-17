// Package templatedcmd allows expanding command-lines via a simple
// template-expansion process.
//
// For example the user might wish to run a command with an argument
// like so:
//
//     command {}
//
// But we also support expanding the input into fields, and selecting
// only a single one, as per:
//
//     $ echo "one two" | echo {1}
//     # -> "one"
//
// All arguments are available via "{}" and "{N}" will refer to the
// Nth field of the given input.
//
package templatedcmd

import (
	"regexp"
	"strconv"
	"strings"
)

// Expand performs the expansion of the given input, via the supplied
// template.  As we allow input to be referred to as an array of fields
// we also let the user specify a split-string here.
//
// By default the input is split on whitespace, but you may supply another
// string instead.
func Expand(template string, input string, split string) []string {

	//
	// Regular expression for looking for ${1}, "${2}", "${3}", etc.
	//
	reg := regexp.MustCompile("({[0-9]+})")

	//
	// Trim our input of leading/trailing spaces.
	//
	input = strings.TrimSpace(input)

	//
	// Default to splitting the input on white-space.
	//
	fields := strings.Fields(input)
	if split != "" {
		fields = strings.Split(input, split)
	}

	//
	// The return-value is an array of strings
	//
	cmd := []string{}

	//
	// We'll operate upon a temporary copy of our template,
	// split into fields.
	//
	cmdTmp := strings.Fields(template)

	//
	// For each piece of the template-string look for
	// "{}", and  "{N}", expand appropriately.
	//
	for _, piece := range cmdTmp {

		//
		// Do we have a "{N}" ?
		//
		matches := reg.FindAllStringSubmatch(piece, -1)

		//
		// If so for each match, perform the expansion
		//
		for _, v := range matches {

			//
			// Copy the match and remove the {}
			//
			// So we just have "1", "3", etc.
			//
			match := v[1]
			match = strings.ReplaceAll(match, "{", "")
			match = strings.ReplaceAll(match, "}", "")

			//
			// Convert the string to a number, and if that
			// worked we'll replace it with the appropriately
			// numbered field.
			//
			num, err := strconv.Atoi(match)
			if err == nil {

				//
				// If the field matches then we can replace it
				//
				if num >= 1 && num <= len(fields) {
					piece = strings.ReplaceAll(piece, v[1], fields[num-1])
				} else {
					//
					// Otherwise it's a field that doesn't
					// exist.  So it's replaced with ''.
					//
					piece = strings.ReplaceAll(piece, v[1], "")
				}
			}
		}

		//
		// Now replace "{}" with the complete argument
		//
		piece = strings.ReplaceAll(piece, "{}", input)

		// And append
		cmd = append(cmd, piece)
	}

	//
	// Now we should have an array of expanded strings.
	//
	return cmd
}
