package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Structure for our options and state.
type validateYAMLCommand struct {

	// Should we be verbose in what we're testing.
	verbose bool
}

// Arguments adds per-command args to the object.
func (vy *validateYAMLCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vy.verbose, "verbose", false, "Should we be verbose")

}

// Info returns the name of this subcommand.
func (vy *validateYAMLCommand) Info() (string, string) {
	return "validate-yaml", `Validate all YAML files for syntax.

Details:

This command finds all files which match the pattern '*.yml', and
'*.yaml' and attempts to load them, validating syntax in the process.

By default the filesystem is walked from the current working directory,
but if you prefer you may specify a starting-directory name as the single
argument to the sub-command.`
}

// validateYAML finds and tests all files beneath the named directory.
func (vy *validateYAMLCommand) validateYAML(path string) bool {

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

		if strings.HasSuffix(path, ".yaml") && !f.IsDir() {
			fileList = append(fileList, path)
		}
		if strings.HasSuffix(path, ".yml") && !f.IsDir() {
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

		if vy.verbose {
			fmt.Printf("Testing: %s\n", file)
		}

		err := vy.validateFile(file)
		if err != nil {
			fmt.Printf("%s : %s\n", file, err.Error())
			fail = true
		}

	}

	return fail
}

// Validate a single file
func (vj *validateYAMLCommand) validateFile(path string) error {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var result interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	return nil
}

// Execute is invoked if the user specifies `validate-yaml` as the subcommand.
func (vj *validateYAMLCommand) Execute(args []string) int {

	path := "."

	if len(args) > 0 {
		path = args[0]

	}

	if vj.validateYAML(path) {
		return 1
	}
	return 0
}
