package main

import (
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/skx/subcommands"
)

// Structure for our options and state.
type torrentCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (s *torrentCommand) Info() (string, string) {
	return "torrent", `Download a torrent-file.

Details:

This is a simple bittorrent client, which allows downloading the torrent
files located on the command-line.
Example:

    $ sysbox torrent magnet:?xt=urn:btih:ZOCMZQIPFFW7OLLMIC5HUB6BPCSDEOQU`
}

// Execute is invoked if the user specifies `torrent` as the subcommand.
func (s *torrentCommand) Execute(args []string) int {

	//
	// Ensure we have only a single argument.
	//
	if len(args) != 1 {
		fmt.Printf("You must specify a single magnet link")
		return 1
	}

	//
	// Create a client.
	//
	c, err := torrent.NewClient(nil)
	if err != nil {
		fmt.Printf("failed to create torrent client: %s\n", err.Error())
		return 1
	}
	defer c.Close()

	//
	// Add each magnet link
	//
	t, err := c.AddMagnet(args[0])
	if err != nil {
		fmt.Printf("failed to add magnet-link %s: %s\n", args[0], err.Error())
		return 1
	}

	//
	// Record our start-time.
	//
	start := time.Now()

	//
	// Await information to be loaded from the torrent.
	//
	<-t.GotInfo()

	//
	// Get the torrent content-list.
	//
	files := t.Files()

	//
	// Show header.
	//
	if len(files) == 0 {
		fmt.Printf("torrent contains no files\n")
		return 1
	} else if len(files) == 1 {
		fmt.Printf("torrent contains the following file:\n")
	} else {
		fmt.Printf("torrent contains the following %d files:\n", len(files))
	}

	//
	// Show files.
	//
	for _, file := range files {
		fmt.Printf("\t%s\n", file.DisplayPath())
	}

	//
	// Await completion of download.
	//
	t.DownloadAll()
	c.WaitAll()

	//
	// Completed.
	//
	fmt.Printf("download complete; %s\n", time.Since(start))

	return 0
}
