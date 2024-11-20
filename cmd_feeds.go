package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/skx/subcommands"
	"golang.org/x/net/html"
)

// Structure for our options and state.
type feedsCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// ErrNoFeeds is used if no feeds are found in a remote URL
var ErrNoFeeds = errors.New("NO-FEED")

// Info returns the name of this subcommand.
func (t *feedsCommand) Info() (string, string) {
	return "feeds", `Extract RSS feeds from remote URLS.

Details:

This command fetches the contents of the specified URL, much like
the 'http-get' command would, and extracts any specified RSS feed
from the contents of that remote URL.

Examples:

  $ sysbox feeds https://blog.steve.fi/`
}

func (t *feedsCommand) FindFeeds(base string) ([]string, error) {

	ret := []string{}

	if !strings.HasPrefix(base, "http") {
		base = "https://" + base
	}

	// Make the request
	response, err := http.Get(base)
	if err != nil {
		return ret, err
	}

	// Get the body.
	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				if len(ret) > 0 {
					return ret, nil
				}
				return ret, ErrNoFeeds
			}
			return ret, fmt.Errorf("%s", z.Err())
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "link" {
				isRSS := false
				u := ""
				for _, attr := range t.Attr {
					if attr.Key == "type" && (attr.Val == "application/rss+xml" || attr.Val == "application/atom+xml") {
						isRSS = true
					}

					if attr.Key == "href" {
						u = attr.Val
					}
				}
				if isRSS {
					if !strings.HasPrefix(u, "http") {
						u, _ = url.JoinPath(base, u)
					}
					ret = append(ret, u)
				}
			}
		}
	}

	// Nothing found?
	if len(ret) == 0 {
		return ret, ErrNoFeeds
	}
	return ret, nil
}

// Execute is invoked if the user specifies `feeds` as the subcommand.
func (t *feedsCommand) Execute(args []string) int {

	// Ensure we have only a single URL
	if len(args) != 1 {
		fmt.Printf("Usage: feeds URL\n")
		return 1
	}

	// The URL
	url := args[0]

	// We'll default to https if the protocol isn't specified.
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	out, err := t.FindFeeds(url)
	if err != nil {
		if err == ErrNoFeeds {
			fmt.Printf("No Feeds found in %s\n", url)
		} else {
			fmt.Printf("Error processing %s: %s\n", url, err)
			return 1
		}
	} else {
		for _, x := range out {
			fmt.Printf("%s\n", x)
		}
	}

	return 0
}
