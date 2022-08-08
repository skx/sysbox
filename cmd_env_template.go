package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type envTemplateCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (et *envTemplateCommand) Info() (string, string) {
	return "env-template", `Populate a template-file with environmental variables.

Details:

This command is a slight reworking of the standard 'envsubst' command,
which might not be available upon systems by default, along with extra
support for file-inclusion (which supports the inclusion of other files,
along with extra behavior such as 'grep' and inserting regions of files
between matches of a start/end pair of regular expressions).

The basic use-case of this sub-command is to allow substituting
environmental variables into simple (golang) template-files.

However there are extra facilities, as noted above.

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


Inclusion Examples:

The basic case of including a file could be handled like so:

   Before
   {{include "/etc/passwd"}}
   After

You can also include only lines matching a particular regular
expression:

   {{grep "/etc/passwd" "^(root|nobody):"}}

Or lines between a pair of marker (regular expressions):

   {{between "/etc/passwd" "^root" "^bin"}}

If the input-file contains a '|' prefix it will instead read the output
of running the named command - so you shouldn't process user-submitted
templates, as that is a potential security-risk.

NOTE: Using 'between' includes the lines that match too, not just the region
between them.  If you regard this as a bug please file an issue.

`

}

// Execute is invoked if the user specifies `env-template` as the subcommand.
func (et *envTemplateCommand) Execute(args []string) int {

	//
	// Ensure we have an argument
	//
	if len(args) < 1 {
		fmt.Printf("You must specify the template to expand.\n")
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
	var err error
	var content []byte
	content, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	//
	// Define a helper-function that are available within the
	// templates we process.
	//
	funcMap := template.FuncMap{
		"between": between,
		"env":     env,
		"grep":    grep,
		"include": include,
		"split":   split,
	}

	// Parse the file
	t := template.Must(template.New("t1").Funcs(funcMap).Parse(string(content)))

	// Render
	err = t.Execute(os.Stdout, nil)

	return err
}
