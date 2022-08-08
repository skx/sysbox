package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
)

// Structure for our options and state.
type httpGetCommand struct {

	// Show headers?
	headers bool

	// Show body?
	body bool
}

// Arguments adds per-command args to the object.
func (hg *httpGetCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&hg.body, "body", true, "Show the response body.")
	f.BoolVar(&hg.headers, "headers", false, "Show the response headers.")

}

// Info returns the name of this subcommand.
func (hg *httpGetCommand) Info() (string, string) {
	return "http-get", `Download and display the contents of a remote URL.

Details:

This command is very much curl-lite, allowing you to fetch the contents of
a remote URL, with no configuration options of any kind.

While it is unusual to find hosts without curl or wget installed it does
happen, this command will bridge the gap a little.

Examples:

  $ sysbox http-get https://steve.fi/`
}

// Execute is invoked if the user specifies `http-get` as the subcommand.
func (hg *httpGetCommand) Execute(args []string) int {

	// Ensure we have only a single URL
	if len(args) != 1 {
		fmt.Printf("Usage: http-get URL\n")
		return 1
	}

	// Make the request
	response, err := http.Get(args[0])
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return 1
	}

	// Get the body.
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return 1
	}

	// Show header?
	if hg.headers {

		// Keep a list of the headers here for sort/display
		headers := []string{}

		// Copy the headers
		for header := range response.Header {
			headers = append(headers, header)
		}

		// Sort them
		sort.Strings(headers)

		// Output them
		for _, header := range headers {
			fmt.Printf("%s: %s\n", header, response.Header.Get(header))
		}
	}

	// If showing header and body separate them both
	if hg.headers && hg.body {
		fmt.Printf("\n")
	}

	// Show body?
	if hg.body {
		fmt.Printf("%s\n", string(contents))
	}

	return 0
}
