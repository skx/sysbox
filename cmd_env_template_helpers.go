// helper functions for template-expansion.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// between finds the text between two regular expressions, from either
// a file or the output of a given command.
func between(in string, begin string, end string) string {

	var content []byte
	var err error

	// Read the named file/command-output here.
	if strings.HasPrefix(in, "|") {
		content, err = runCommand(strings.TrimPrefix(in, "|"))
	} else {
		content, err = os.ReadFile(in)
	}

	if err != nil {
		return fmt.Sprintf("error reading %s: %s", in, err.Error())
	}

	// temporary holder
	res := []string{}

	// found the open?
	var found bool

	// for each line
	for _, line := range strings.Split(string(content), "\n") {

		// in the section we care about?
		var matched bool
		matched, err = regexp.MatchString(begin, line)
		if err != nil {
			return fmt.Sprintf("error matching %s: %s", begin, err.Error())
		}

		// if we matched add the line
		if matched || found {
			res = append(res, line)
		}

		// if we matched, or we're in a match
		// then skip
		if matched {
			found = true
			continue
		}

		// are we closing a match?
		if found {
			matched, err = regexp.MatchString(end, line)
			if err != nil {
				return fmt.Sprintf("error matching %s: %s", end, err.Error())
			}

			if matched {
				found = false
			}
		}
	}
	return strings.Join(res, "\n")
}

// env returns the contents of an environmental variable.
func env(s string) string {
	return (os.Getenv(s))
}

// grep allows filtering the contents of a file, or output of a command,
// via a regular expression.
func grep(in string, pattern string) string {
	var content []byte
	var err error

	// Read the named file/command-output here.
	if strings.HasPrefix(in, "|") {
		content, err = runCommand(strings.TrimPrefix(in, "|"))
	} else {
		content, err = os.ReadFile(in)
	}

	if err != nil {
		return fmt.Sprintf("error reading %s: %s", in, err.Error())
	}

	var matched bool
	res := []string{}
	for _, line := range strings.Split(string(content), "\n") {
		matched, err = regexp.MatchString(pattern, line)
		if err != nil {
			return fmt.Sprintf("error matching %s: %s", pattern, err.Error())
		}
		if matched {
			res = append(res, line)
		}
	}
	return strings.Join(res, "\n")

}

// include inserts the contents of a file, or output of a command
func include(in string) string {

	var content []byte
	var err error

	// Read the named file/command-output here.
	if strings.HasPrefix(in, "|") {
		content, err = runCommand(strings.TrimPrefix(in, "|"))
	} else {
		content, err = os.ReadFile(in)
	}

	if err != nil {
		return fmt.Sprintf("error reading %s: %s", in, err.Error())
	}
	return (string(content))
}

// runCommand returns the output of running the given command
func runCommand(command string) ([]byte, error) {

	// Build up the thing to run, using a shell so that
	// we can handle pipes/redirection.
	toRun := []string{"/bin/bash", "-c", command}

	// Run the command
	cmd := exec.Command(toRun[0], toRun[1:]...)

	// Get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("error running command '%s' %s", command, err.Error())
	}

	// Strip trailing newline.
	return output, nil
}

// split converts a string to an array.
func split(in string, delim string) []string {
	return strings.Split(in, delim)
}
