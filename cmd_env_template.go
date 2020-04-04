package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// Structure for our options and state.
type envTemplateCommand struct {
}

// Arguments adds per-command args to the object.
func (et *envTemplateCommand) Arguments(f *flag.FlagSet) {
}

// Info returns the name of this subcommand.
func (et *envTemplateCommand) Info() (string, string) {
	return "env-template", `Populate a template-file with environmental variables.

Details:

This command is a slight reworking of the standard 'envsubst' command,
which might not be available upon systems by default.

The intention is that you can substitute environmental variables into
simple (golang) template-files.


Examples:

Consider the case where you have a shell with $PATH and $USER available
you might wish to expand those into a file.  The file could contain:

    Hello {{env USER}} your path is {{env "PATH"}}

Expand the contents via:

    $ sysbox env-template path/to/template

Using the standard golang text/template facilities you can use conditionals
and process variables.  For example splitting $PATH into parts:

    // template.in - shows $PATH entries one by one
    {{$path := env "PATH"}}
    {{$arr := split $path ":"}}
    {{range $k, $v := $arr}}
      {{$k}} {{$v}}
    {{end}}

`

}

// Execute is invoked if the user specifies `with-lock` as the subcommand.
func (et *envTemplateCommand) Execute(args []string) int {

	//
	// Ensure we have an argument
	//
	if len(args) < 1 {
		fmt.Printf("You must specify the template to expand\n")
		return 1
	}

	fail := 0

	for _, file := range args {
		err := et.expandFile(file)
		if err != nil {
			fmt.Printf("error processing %s %s\n", file, err.Error())
			fail = 1
		}
	}
	return fail
}

// expandFile does the file expansion
func (et *envTemplateCommand) expandFile(path string) error {

	// Load the file
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	//
	// Define a helper-function that users can call to get
	// the variables they've set.
	//
	funcMap := template.FuncMap{
		"env": func(s string) string {
			return (os.Getenv(s))
		},
		"split": func(in string, delim string) []string {
			return strings.Split(in, delim)
		},
	}

	// Parse the file
	t := template.Must(template.New("t1").Funcs(funcMap).Parse(string(content)))

	// Render
	err = t.Execute(os.Stdout, nil)

	return err
}
