package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

// Structure for our options and state.
type ipsCommand struct {

	// show only IPv4 addresses?
	ipv4 bool

	// show only IPv6 addresses?
	ipv6 bool

	// show local addresses?
	local bool

	// show global/remote addresses?
	remote bool

	// Cached store of network/netmask to IP-range - IPv4
	ip4Ranges map[string]*net.IPNet

	// Cached store of network/netmask to IP-range - IPv6
	ip6Ranges map[string]*net.IPNet
}

// Arguments adds per-command args to the object.
func (i *ipsCommand) Arguments(f *flag.FlagSet) {
	f.BoolVar(&i.ipv4, "4", true, "Should we show IPv4 addresses?")
	f.BoolVar(&i.ipv6, "6", true, "Should we show IPv6 addresses?")
	f.BoolVar(&i.local, "local", true, "Should we show local addresses?")
	f.BoolVar(&i.remote, "remote", true, "Should we show global addresses?")

}

// Info returns the name of this subcommand.
func (i *ipsCommand) Info() (string, string) {
	return "ips", `Show IP address information.

Details:

This command allows you to see local/global IP addresses assigned to
the current host.

By default all IP addresses will be shown, but you can disable protocols
and types of addresses you do not wish to see.

Examples:

$ sysbox ips -4=false
::1
fe80::feaa:14ff:fe32:688
fe80::78e5:95b6:1659:b407

$ sysbox ips -local=false -4=false
2a01:4f9:c010:27d8::1
`
}

// isLocal is a helper to test if an address is "local" or "remote".
func (i *ipsCommand) isLocal(address *net.IPNet) bool {

	localIP4 := []string{
		"10.0.0.0/8",         // RFC1918
		"100.64.0.0/10",      // RFC 6598
		"127.0.0.0/8",        // IPv4 loopback
		"169.254.0.0/16",     // RFC3927 link-local
		"172.16.0.0/12",      // RFC1918
		"192.0.0.0/24",       // RFC 5736
		"192.0.2.0/24",       // RFC 5737
		"192.168.0.0/16",     // RFC1918
		"192.18.0.0/15",      // RFC 2544
		"192.88.99.0/24",     // RFC 3068
		"198.51.100.0/24",    //
		"203.0.113.0/24",     //
		"224.0.0.0/4",        // RFC 3171
		"255.255.255.255/32", // RFC 919 Section 7
	}
	localIP6 := []string{
		"::/128",        // RFC 4291: Unspecified Address
		"100::/64",      // RFC 6666: Discard Address Block
		"2001:2::/48",   // RFC 5180: Benchmarking
		"2001::/23",     // RFC 2928: IETF Protocol Assignments
		"2001::/32",     // RFC 4380: TEREDO
		"2001:db8::/32", // RFC 3849: Documentation
		"::1/128",       // RFC 4291: Loopback Address
		"fc00::/7",      // RFC 4193: Unique-Local
		"fe80::/10",     // RFC 4291: Section 2.5.6 Link-Scoped Unicast
		"ff00::/8",      // RFC 4291: Section 2.7
	}

	// Create our maps
	if i.ip4Ranges == nil {
		i.ip4Ranges = make(map[string]*net.IPNet)
		i.ip6Ranges = make(map[string]*net.IPNet)

		// Join our ranges.
		tmp := localIP4
		tmp = append(tmp, localIP6...)

		// For each network-range.
		for _, entry := range tmp {

			// Parse
			_, block, _ := net.ParseCIDR(entry)

			// Record in the protocol-specific range
			if strings.Contains(entry, ":") {
				i.ip6Ranges[entry] = block
			} else {
				i.ip4Ranges[entry] = block
			}
		}
	}

	// The map we're testing from
	testMap := i.ip4Ranges

	// Are we testing an IPv6 address?
	if strings.Contains(address.String(), ":") {
		testMap = i.ip6Ranges
	}

	// Loop over the appropriate map and test for inclusion
	for _, block := range testMap {
		if block.Contains(address.IP) {
			return true
		}
	}

	// Not found.
	return false
}

// Execute is invoked if the user specifies `ips` as the subcommand.
func (i *ipsCommand) Execute(args []string) int {

	// Get addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("Error finding IPs:%s\n", err.Error())
		return 1
	}

	// For each one
	for _, address := range addrs {

		// cast ..
		ipnet, ok := address.(*net.IPNet)
		if !ok {
			fmt.Printf("Failed to convert %v to IP\n", address)
			return 1
		}

		// If we're not showing locals, then skip if this is.
		if !i.local && i.isLocal(ipnet) {
			continue
		}

		// If we're not showing globals, then skip if this is
		if !i.remote && !i.isLocal(ipnet) {
			continue
		}

		res := ipnet.IP.String()

		// If we're not showing IPv4 and the address is that
		// then skip it
		if !i.ipv4 && !strings.Contains(res, ":") {
			continue
		}

		// If we're not showing IPv6 and the address is that then
		// skip it
		if !i.ipv6 && strings.Contains(res, ":") {
			continue
		}

		fmt.Println(res)
	}
	return 0
}
