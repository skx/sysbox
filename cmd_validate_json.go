package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Structure for our options and state.
type validateJSONCommand struct {

	// comma-separated list of files to exclude, as set by the
	// command-line flag.
	exclude string

	// an array of patterns to exclude, calculated from the
	// exclude setting above.
	excluded []string

	// Should we report on what we're testing.
	verbose bool
}

// Arguments adds per-command args to the object.
func (vj *validateJSONCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vj.verbose, "verbose", false, "Should we be verbose")
	f.StringVar(&vj.exclude, "exclude", "", "Comma-separated list of patterns to exclude files from the check")

}

// Info returns the name of this subcommand.
func (vj *validateJSONCommand) Info() (string, string) {
	return "validate-json", `Validate all JSON files for syntax.

Details:

This command allows you to validate JSON files, by default searching
recursively beneath the current directory for all files which match
the pattern '*.json'.

If you prefer you may specify a number of directories or files:

- Any file specified will be checked.
- Any directory specified will be recursively scanned for matching files.
  - Files that do not have a '.json' suffix will be ignored.

Example:

    $ sysbox validate-json -verbose file1.json file2.json ..
    $ sysbox validate-json -exclude=foo /dir/1/path /file/1/path ..
`
}

// Validate a single file
func (vj *validateJSONCommand) validateFile(path string) error {

	// Exclude this file? Based on the supplied list though?
	for _, ex := range vj.excluded {
		if strings.Contains(path, ex) {
			if vj.verbose {
				fmt.Printf("SKIPPED\t%s - matched '%s'\n", path, ex)
			}
			return nil
		}
	}

	// Read the file-contents
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Deserialize - receiving an error if that failed.
	var result interface{}
	err = json.Unmarshal(data, &result)

	// Show the error if there was one, but otherwise only show
	// the success if running verbosely.
	if err != nil {
		fmt.Printf("ERROR\t%s - %s\n", path, err.Error())
	} else {
		if vj.verbose {
			fmt.Printf("OK\t%s\n", path)
		}
	}
	return err
}

// Execute is invoked if the user specifies `validate-json` as the subcommand.
func (vj *validateJSONCommand) Execute(args []string) int {

	// Did we find at least one file with an error?
	failed := false

	// Create our array of excluded patterns if something
	// should be excluded.
	if vj.exclude != "" {
		vj.excluded = strings.Split(vj.exclude, ",")
	}

	// Add a fake argument if nothing is present, because we
	// want to process the current directory (recursively) by default.
	if len(args) < 1 {
		args = append(args, ".")
	}

	// We can handle file/directory names as arguments.  If a
	// directory is specified then we process it recursively.
	//
	// We'll start by building up a list of all the files to test,
	// before we begin the process of testing.  We'll make sure
	// our list is unique to cut down on any unnecessary I/O.
	todo := make(map[string]bool)

	// For each argument ..
	for _, arg := range args {

		// Check that it actually exists.
		info, err := os.Stat(arg)
		if os.IsNotExist(err) {
			fmt.Printf("The path does not exist: %s\n", arg)
			continue
		}

		// Error?
		if err != nil {
			fmt.Printf("Failed to stat(%s): %s\n", arg, err.Error())
			continue
		}

		// A directory?
		if info.Mode().IsDir() {

			// Find suitable entries in the directory
			files, err := FindFiles(arg, []string{".json"})
			if err != nil {
				fmt.Printf("Error finding files in %s: %s\n", arg, err.Error())
				continue
			}

			// Then record each one.
			for _, ent := range files {
				todo[ent] = true
			}
		} else {

			// OK the entry we were given is just a file,
			// so we'll save the path away.
			todo[arg] = true
		}
	}

	//
	// Now we have a list of files to process.
	//
	for file := range todo {

		// Run the validation, and note the result
		err := vj.validateFile(file)
		if err != nil {
			failed = true
		}
	}

	// Setup a suitable exit-code
	if failed {
		return 1
	}

	return 0
}
