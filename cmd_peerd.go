package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
)

var (
	mutex *sync.Mutex
)

// Update our list of peers.
func (p *peerdCommand) writePeers(members []*memberlist.Node) {
	mutex.Lock()
	defer mutex.Unlock()

	//
	// A temporary structure to contain our peers
	//
	type Object struct {
		IPs     []string
		Names   []string
		Members map[string]string
	}

	//
	// Create an instance of the object.
	//
	obj := &Object{Members: make(map[string]string)}

	//
	// Populate it
	//
	for _, member := range members {
		obj.IPs = append(obj.IPs, member.Addr.String())
		obj.Names = append(obj.Names, member.Name)
		obj.Members[member.Name] = member.Addr.String()
	}

	//
	// Convert to JSON
	//
	out, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("error marshalling peers to json %s\n", err.Error())
		os.Exit(1)
	}

	//
	// Write to disk
	//
	err = os.WriteFile("/var/tmp/peerd.json", out, 0644)
	if err != nil {
		fmt.Printf("error writing JSON to file %s\n", err.Error())
		os.Exit(1)
	}

}

// eventDelegate is used to report upon changes to our peer-list
type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Println("joined: " + node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Println("left: " + node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Println("updated: " + node.String())
}

// Structure for our options and state.
type peerdCommand struct {

	// Our public IP
	ip string
}

// Arguments adds per-command args to the object.
func (p *peerdCommand) Arguments(f *flag.FlagSet) {
	f.StringVar(&p.ip, "ip", "", "Our public-facing IP address")

}

// Info returns the name of this subcommand.
func (p *peerdCommand) Info() (string, string) {
	return "peerd", `Keep track of peer hosts.

Details:

This command works as a daemon, keeping in constant contact with a set
of peers.  Peers that are known and "up" are tracked and stored in the
JSON file '/var/tmp/peerd.json'.

Usage:

Launch the daemon on one host, with the public IP specified:

    peerd -ip=1.2.3.4

Now launch on the a second host, giving the IP of at least one peer:

    peerd -ip=11.22.33.44 1.2.3.4

Both hosts will know about the other, and will update their local state
file if the other host goes away, or new hosts join.

Firewalling:

The communication happens over port 7946.`
}

// Execute is invoked if the user specifies `peerd` as the subcommand.
func (p *peerdCommand) Execute(args []string) int {

	// Create a mutex
	mutex = &sync.Mutex{}

	// Create config
	config := memberlist.DefaultWANConfig()

	// Setup our external IP
	config.AdvertiseAddr = p.ip

	// Here we log node join/exit
	config.Events = &eventDelegate{}

	// Create the config
	list, err := memberlist.Create(config)
	if err != nil {
		fmt.Printf("Failed to create memberlist: " + err.Error())
		return 1
	}

	// If we have a peer, join it.
	for _, peer := range args {
		_, err := list.Join([]string{peer})
		if err != nil {
			fmt.Printf("Failed to join cluster via peer %s - %s", peer, err.Error())
			return 1
		}
	}

	// Now we'll update every few seconds
	for {

		peers := list.Members()

		// Update our peer-list
		p.writePeers(peers)
		time.Sleep(5 * time.Second)
	}

}
