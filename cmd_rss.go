package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

// Structure for our options and state.
type rssCommand struct {
	// format contains the format-string to display for entries.
	// We recognize "$link" and "$title".  More might be added in the future.
	format string
}

// Arguments adds per-command args to the object.
func (r *rssCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&r.format, "format", "$link", "Specify the format-string to display for entries")
}

// Info returns the name of this subcommand.
func (r *rssCommand) Info() (string, string) {
	return "rss", `Show details from an RSS feed.

Details:

This command fetches the specified URLs as RSS feeds, and shows
their contents in a simple fashion.  By default only the entry URLs
are shown, but a format-string may be used to specify the output.

For example to show the link and title of entries:

    $ sysbox rss -format='$link $title' http://,..

Suggestions for additional fields/details to be displayed are
welcome via issue-reports.

Format String:

Currently the following values are supported:

* $content The content of the entry.
* $date The published date of the entry.
* $guid The GUID of the entry.
* $length The length of the entry.
* $link The link to the entry.
* $title The title of the entry.

Usage:

   $ sysbox rss url1 url2 .. urlN

Note:

Care must be taken to escape, or quote, the '$' character which
is used in the format-string.
`
}

// Process each specified feed.
func (r *rssCommand) processFeed(url string) error {

	// Create the parser with defaults
	fp := gofeed.NewParser()

	// Parse the feed
	feed, err := fp.ParseURL(url)
	if err != nil {
		return err
	}

	// For each entry
	for _, ent := range feed.Items {

		// Get a piece of text, using our format-string
		txt := os.Expand(
			r.format,
			func(s string) string {
				switch s {
				case "content":
					return ent.Content
				case "date":
					return ent.Published
				case "guid":
					return ent.GUID
				case "length":
					return fmt.Sprintf("%d", len(ent.Content))
				case "link":
					return ent.Link
				case "title":
					return ent.Title
				default:
					return s
				}
			},
		)

		// Now show it
		fmt.Println(txt)
	}

	// All good.
	return nil
}

// Execute is invoked if the user specifies `rss` as the subcommand.
func (r *rssCommand) Execute(args []string) int {

	for _, u := range args {

		err := r.processFeed(u)
		if err != nil {
			fmt.Printf("Failed to process %s: %s\n", u, err)
		}
	}

	return 0
}
