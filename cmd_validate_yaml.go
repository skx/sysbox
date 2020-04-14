package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Structure for our options and state.
type validateYAMLCommand struct {

	// comma-separated list of files to exclude
	exclude string

	// Should we be verbose in what we're testing.
	verbose bool
}

// Arguments adds per-command args to the object.
func (vy *validateYAMLCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&vy.verbose, "verbose", false, "Should we be verbose")
	f.StringVar(&vy.exclude, "exclude", "", "Comma-separated list of files to exclude")
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
	// Did we see a failure?
	//
	fail := false

	//
	// Find all files
	//
	files, err := FindFiles(path, []string{".yaml", ".yml"})

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
	if vy.exclude != "" {
		excluded = strings.Split(vy.exclude, ",")
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

		if vy.verbose {
			if exclude {
				fmt.Printf("Excluded: %s\n", file)
			} else {
				fmt.Printf("Testing: %s\n", file)
			}
		}
		if exclude {
			continue
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
func (vy *validateYAMLCommand) validateFile(path string) error {

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
func (vy *validateYAMLCommand) Execute(args []string) int {

	path := "."

	if len(args) > 0 {
		path = args[0]

	}

	if vy.validateYAML(path) {
		return 1
	}
	return 0
}
