package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
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

func torrentBar(t *torrent.Torrent) {
	go func() {
		if t.Info() == nil {
			<-t.GotInfo()
		}
		var lastLine string
		for {
			var completedPieces, partialPieces int
			psrs := t.PieceStateRuns()
			for _, r := range psrs {
				if r.Complete {
					completedPieces += r.Length
				}
				if r.Partial {
					partialPieces += r.Length
				}
			}
			line := fmt.Sprintf(
				"downloading %q: %s/%s, %d/%d pieces completed (%d partial)\n",
				t.Name(),
				humanize.Bytes(uint64(t.BytesCompleted())),
				humanize.Bytes(uint64(t.Length())),
				completedPieces,
				t.NumPieces(),
				partialPieces,
			)
			if line != lastLine {
				lastLine = line
				os.Stdout.WriteString(line)
			}
			time.Sleep(time.Second)
		}
	}()
}

// Execute is invoked if the user specifies `torrent` as the subcommand.
func (s *torrentCommand) Execute(args []string) int {

	//
	// Ensure we have only a single argument.
	//
	if len(args) != 1 {
		fmt.Printf("You must specify a magnet link to download\n")
		return 1
	}

	//
	// Ensure we have a magnet:-link
	//
	if !strings.HasPrefix(args[0], "magnet:") {
		fmt.Printf("Usage: $sysbox torrent magnet:?....\n")
		return 1
	}

	//
	// Create the default client-configuration.
	//
	clientConfig := torrent.NewDefaultClientConfig()

	//
	// Create the client.
	//
	c, err := torrent.NewClient(clientConfig)
	if err != nil {
		fmt.Printf("failed to create torrent client: %s\n", err.Error())
		return 1
	}
	defer c.Close()

	//
	// Add the magnet link
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
	// Spawn a progress-bar.
	//
	torrentBar(t)

	//
	// Await information to be loaded from the torrent.
	//
	<-t.GotInfo()

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
