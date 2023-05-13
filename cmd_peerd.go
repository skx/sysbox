package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	// Write to disk the current members of the mesh.
	//
	err = os.WriteFile(p.stateFile, out, 0644)
	if err != nil {
		fmt.Printf("error writing JSON to file %s\n", err.Error())
		os.Exit(1)
	}

}

// eventDelegate is used to report upon changes to our peer-list
type eventDelegate struct {
	// up is a command to run when a peer joins.
	up string

	// down is a command to run when a peer leaves.
	down string
}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	fmt.Println("joined: " + node.String())
	if ed.up != "" {
		ed.RunCommand(ed.up, node)
	}
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	fmt.Println("left: " + node.String())
	if ed.down != "" {
		ed.RunCommand(ed.down, node)
	}
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	fmt.Println("updated: " + node.String())
}

// RunCommand is called to run a command, replacing "${IP}" and ${NAME}
// appropriately.
func (ed *eventDelegate) RunCommand(cmd string, node *memberlist.Node) {

	// Helper to expand $IP, ${IP}, etc
	mapper := func(placeholderName string) string {
		switch placeholderName {
		case "IP", "ip":
			return node.Addr.String()
		case "NAME", "name":
			return node.Name
		}

		return ""
	}

	// Expand the command
	cmd = os.Expand(cmd, mapper)

	// Assume Unix
	shell := "/bin/sh -c"

	switch runtime.GOOS {
	case "windows":
		shell = "cmd /c"
	}

	// Build up the thing to run
	sh := strings.Split(shell, " ")
	sh = append(sh, cmd)

	// Now run
	run := exec.Command(sh[0], sh[1:]...)
	_, err := run.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running '%v': %s\n", sh, err.Error())
	}

}

// Structure for our options and state.
type peerdCommand struct {

	// Our public IP
	ip string

	// The path to our state file
	stateFile string

	// runPeerUp is the (shell) command to run when a node joins.
	runPeerUp string

	// runPeerDown is the (shell) command to run when a node leaves.
	runPeerDown string
}

// Arguments adds per-command args to the object.
func (p *peerdCommand) Arguments(f *flag.FlagSet) {

	f.StringVar(&p.ip, "ip", "", "Our public-facing IP address")
	f.StringVar(&p.stateFile, "state", "/var/tmp/peerd.json", "The file within which to store peer members")
	f.StringVar(&p.runPeerUp, "run-up", "", "The command to run when a node joins the group")
	f.StringVar(&p.runPeerDown, "run-down", "", "The command to run when a node joins the group")
}

// Info returns the name of this subcommand.
func (p *peerdCommand) Info() (string, string) {
	return "peerd", `Keep track of peer hosts.

Details:

This command works as a daemon, keeping in constant contact with a set
of peers.  Peers that are known and "up" are tracked and stored in the
JSON file '/var/tmp/peerd.json'.  You may specify an alternative location
via the -state flag

Usage:

Launch the daemon on one host, with the public IP specified:

    peerd -ip=1.2.3.4

Now launch on the a second host, giving the IP of at least one peer:

    peerd -ip=11.22.33.44 1.2.3.4

Both hosts will know about the other, and will update their local state
file if the other host goes away, or new hosts join.


Notifications:

You can configure a (shell) command to execute when peers join, or leave,
the peer network via '-run-up' and '-run-down'.  Within those commands
$IP and $NAME will be expanded to contain the appropriate detail of the
remote peer which has joined/left the group.


Firewalling considerations:

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
	config.Events = &eventDelegate{
		up:   p.runPeerUp,
		down: p.runPeerDown,
	}

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
