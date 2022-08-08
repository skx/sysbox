package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Structure for our options and state.
type validateXMLCommand struct {

	// comma-separated list of files to exclude, as set by the
	// command-line flag.
	exclude string

	// an array of patterns to exclude, calculated from the
	// exclude setting above.
	excluded []string

	// Should we report on what we're testing.
	verbose bool
}

// identReader is a hack which allows us to ignore character-conversion
// issues, depending on the encoded-characterset of the XML input.
//
// We use this because we care little for the attributes/values, instead
// wanting to check for tag-validity.
func identReader(encoding string, input io.Reader) (io.Reader, error) {
	return input, nil
}

// Arguments adds per-command args to the object.
func (vx *validateXMLCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vx.verbose, "verbose", false, "Should we be verbose")
	f.StringVar(&vx.exclude, "exclude", "", "Comma-separated list of patterns to exclude files from the check")

}

// Info returns the name of this subcommand.
func (vx *validateXMLCommand) Info() (string, string) {
	return "validate-xml", `Validate all XML files for syntax.

Details:

This command allows you to validate XML files, by default searching
recursively beneath the current directory for all files which match
the pattern '*.xml'.

If you prefer you may specify a number of directories or files:

- Any file specified will be checked.
- Any directory specified will be recursively scanned for matching files.
  - Files that do not have a '.xml' suffix will be ignored.

Example:

    $ sysbox validate-xml -verbose file1.xml file2.xml ..
    $ sysbox validate-xml -exclude=foo /dir/1/path /file/1/path ..
`
}

// Validate a single file
func (vx *validateXMLCommand) validateFile(path string) error {

	// Exclude this file? Based on the supplied list though?
	for _, ex := range vx.excluded {
		if strings.Contains(path, ex) {
			if vx.verbose {
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

	// Store the results here.
	var result interface{}

	// Decode into the results, taking care that we
	// wire up some magic to avoid caring about the
	// encoding/character-set issues.
	r := bytes.NewReader(data)
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = identReader
	err = decoder.Decode(&result)

	// Show the error if there was one, but otherwise only show
	// the success if running verbosely.
	if err != nil {
		fmt.Printf("ERROR\t%s - %s\n", path, err.Error())
	} else {
		if vx.verbose {
			fmt.Printf("OK\t%s\n", path)
		}
	}
	return err
}

// Execute is invoked if the user specifies `validate-xml` as the subcommand.
func (vx *validateXMLCommand) Execute(args []string) int {

	// Did we find at least one file with an error?
	failed := false

	// Create our array of excluded patterns if something
	// should be excluded.
	if vx.exclude != "" {
		vx.excluded = strings.Split(vx.exclude, ",")
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
			files, err := FindFiles(arg, []string{".xml"})
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
		err := vx.validateFile(file)
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
