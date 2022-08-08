package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
)

// Structure for our options and state.
type fingerdCommand struct {
	port int
}

// Arguments adds per-command args to the object.
func (fc *fingerdCommand) Arguments(f *flag.FlagSet) {
	f.IntVar(&fc.port, "port", 79, "The port to listen upon")
}

// Info returns the name of this subcommand.
func (fc *fingerdCommand) Info() (string, string) {
	return "fingerd", `A small finger daemon.

Details:

This command provides a simple finger server, which allows remote users
to finger your local users.

The file ~/.plan will be served to any remote clients who inspect your
users.

Examples:

   $ sysbox fingerd &
   $ echo "I like cakes" > ~/.plan
   $ finger $USER@localhost

Security:

To allow this to be started as a non-root user you'll want to
run something like:

   $ sudo setcap cap_net_bind_service=+ep /path/to/sysbox

This is better than dropping privileges and starting as root
as a result of the lack of reliability of the latter.  See
https://github.com/golang/go/issues/1435 for details

The alternative would be to bind to :7979 and use iptables
to redirect access from :79 -> 127.0.0.1:7979.

Something like this for external access:

   # iptables -t nat -A PREROUTING -p tcp -m tcp --dport 79 -j REDIRECT --to-ports 7979

And finally for localhost access:

   # iptables -t nat -A OUTPUT -o lo -p tcp --dport 79 -j REDIRECT --to-port 7979
`
}

// Execute is invoked if the user specifies `fingerd` as the subcommand.
func (fc *fingerdCommand) Execute(args []string) int {

	// Listen
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(fc.port))
	if err != nil {
		fmt.Printf("failed to bind to port %d:n%s\n",
			fc.port, err.Error())
		return 1
	}

	// Accept
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go fc.handleConnection(conn)
	}
}

func (fc *fingerdCommand) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	usr, _, _ := reader.ReadLine()

	info, err := fc.getUserInfo(string(usr))
	if err != nil {
		conn.Write([]byte(err.Error()))
	} else {
		conn.Write(info)
	}
}

func (fc *fingerdCommand) getUserInfo(usr string) ([]byte, error) {
	u, e := user.Lookup(usr)
	if e != nil {
		return nil, e
	}
	data, err := os.ReadFile(u.HomeDir + "/.plan")
	if err != nil {
		return data, errors.New("user doesn't have a .plan file")
	}
	return data, nil
}
