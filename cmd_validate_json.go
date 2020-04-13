package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Structure for our options and state.
type validateJSONCommand struct {
	verbose bool
}

// Arguments adds per-command args to the object.
func (vj *validateJSONCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vj.verbose, "verbose", false, "Should we be verbose")

}

// Info returns the name of this subcommand.
func (jv *validateJSONCommand) Info() (string, string) {
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
	// Files we found.
	//
	var fileList []string

	//
	// Did we see a failure?
	//
	fail := false

	//
	// Find all files
	//
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".json") && !f.IsDir() {
			fileList = append(fileList, path)
		}
		return err
	})

	if err != nil {
		fmt.Printf("Error looking for files: %s\n", err.Error())
		os.Exit(1)
	}

	//
	// Now we walk the list of files we're going to process,
	// and we process each one.
	//
	for _, file := range fileList {

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
