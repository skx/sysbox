package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Structure for our options and state.
type validateJSONCommand struct {

	// comma-separated list of files to exclude
	exclude string

	// Should we report on what we're testing.
	verbose bool
}

// Arguments adds per-command args to the object.
func (vj *validateJSONCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vj.verbose, "verbose", false, "Should we be verbose")
	f.StringVar(&vj.exclude, "exclude", "", "Comma-separated list of files to exclude")

}

// Info returns the name of this subcommand.
func (vj *validateJSONCommand) Info() (string, string) {
	return "validate-json", `Validate all JSON files for syntax.

Details:

This command finds all files which match the pattern '*.json', and
attempts to load them, validating syntax.

By default the filesystem is walked from the current working directory,
but if you prefer you may specify a starting-directory name as the single
argument to the sub-command.`
}

// validateJSON finds and tests all files beneath the named directory.
func (vj *validateJSONCommand) validateJSON(path string) bool {

	//
	// Did we see a failure?
	//
	fail := false

	//
	// Find all files
	//
	files, err := FindFiles(path, []string{".json"})

	//
	// Failure?
	//
	if err != nil {
		fmt.Printf("Error looking for files: %s\n", err.Error())
		os.Exit(1)
	}

	//
	// Split excluded files, if any
	//
	var excluded []string
	if vj.exclude != "" {
		excluded = strings.Split(vj.exclude, ",")
	}

	//
	// Now we walk the list of files we're going to process,
	// and we process each one.
	//
	for _, file := range files {

		//
		// We default to not excluding files.
		//
		exclude := false

		//
		// Exclude this file?
		//
		for _, ex := range excluded {

			if strings.Contains(file, ex) {
				exclude = true
			}
		}

		if vj.verbose {
			if exclude {
				fmt.Printf("Excluded: %s\n", file)
			} else {
				fmt.Printf("Testing: %s\n", file)
			}
		}
		if exclude {
			continue
		}

		err := vj.validateFile(file)
		if err != nil {
			fmt.Printf("%s : %s\n", file, err.Error())
			fail = true
		}

	}

	return fail
}

// Validate a single file
func (vj *validateJSONCommand) validateFile(path string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}

// Execute is invoked if the user specifies `validate-json` as the subcommand.
func (vj *validateJSONCommand) Execute(args []string) int {

	path := "."

	if len(args) > 0 {
		path = args[0]

	}

	if vj.validateJSON(path) {
		return 1
	}
	return 0
}
