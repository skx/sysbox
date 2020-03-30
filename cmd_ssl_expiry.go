package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/skx/subcommands"
)

// Structure for our options and state.
type SSLExpiryCommand struct {

	// We embed the NoFlags option, because we accept no command-line flags.
	subcommands.NoFlags
}

// Info returns the name of this subcommand.
func (s *SSLExpiryCommand) Info() (string, string) {
	return "ssl-expiry",
		`Report how long until an SSL certificate expires.

Details:

This sub-command shows the number of hours/days until the SSL
certificate presented upon a remote host expires.  The value
displayed is the minimum expiration time of the certificate and
any bundled-chains served with it.

Examples:

Report on an SSL certificate:

     $ gobox ssl-expiry https://example.com/
     $ gobox ssl-expiry example.com

Report on an SMTP-certificate:

     $ gobox ssl-expiry smtp.gmail.com:465
`

}

// Execute runs our sub-command.
func (s *SSLExpiryCommand) Execute(args []string) int {

	//
	// Ensure we have an argument
	//
	if len(args) < 1 {
		fmt.Printf("You must specify the host(s) to test.\n")
		return 1
	}

	// For each argument
	for _, arg := range args {

		fmt.Printf("%s\n", arg)

		hours, err := s.SSLExpiration(arg)
		if err != nil {
			fmt.Printf("\tERROR:%s\n", err.Error())
		} else {
			fmt.Printf("\t%d hours (%d days)\n", hours, hours/24)
		}
	}

	return 0
}

// SSLExpiration returns the number of hours remaining for a given
// SSL certificate chain.
func (s *SSLExpiryCommand) SSLExpiration(host string) (int64, error) {

	// Expiry time, in hours
	var hours int64
	hours = -1

	//
	// If the string matches http[s]://, then strip it off
	//
	re, err := regexp.Compile(`^https?:\/\/([^\/]+)`)
	if err != nil {
		return 0, err
	}
	res := re.FindAllStringSubmatch(host, -1)
	for _, v := range res {
		host = v[1]
	}

	//
	// If no port is specified default to :443
	//
	p := strings.Index(host, ":")
	if p == -1 {
		host += ":443"
	}

	//
	// Connect, with sane timeout
	//
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Second * 2}, "tcp", host, nil)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	timeNow := time.Now()
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for _, cert := range chain {

			// Get the expiration time, in hours.
			expiresIn := int64(cert.NotAfter.Sub(timeNow).Hours())

			// If we've not checked anything this is the benchmark
			if hours == -1 {
				hours = expiresIn
			} else {
				// Otherwise replace our result if the
				// certificate is going to expire more
				// recently than the current "winner".
				if expiresIn < hours {
					hours = expiresIn
				}
			}
		}
	}

	return hours, nil
}
