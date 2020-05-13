package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// Structure for our options and state.
type httpdCommand struct {
	host string
	port int
	path string
}

// Arguments adds per-command args to the object.
func (h *httpdCommand) Arguments(f *flag.FlagSet) {

	f.StringVar(&h.path, "path", ".", "The directory to use as the HTTP root directory")
	f.StringVar(&h.host, "host", "127.0.0.1", "The host to bind upon (use 0.0.0.0 for remote access)")
	f.IntVar(&h.port, "port", 3000, "The port to listen upon")

}

// Info returns the name of this subcommand.
func (h *httpdCommand) Info() (string, string) {
	return "httpd", `A simple HTTP server.

Details:

This command implements a simple HTTP-server, which defaults to serving
the contents found beneath the current working directory.

By default the content is served to the localhost only, but that can
be changed.

Examples:

$ sysbox httpd
2020/04/01 21:36:27 Serving upon http://127.0.0.1:3000/

$ sysbox httpd -host=0.0.0.0 -port 8080
2020/04/01 21:36:45 Serving upon http://0.0.0.0:8080/`

}

// Execute is invoked if the user specifies `httpd` as the subcommand.
func (h *httpdCommand) Execute(args []string) int {

	//
	// Create a static-file server, based upon the
	// path we're treating as our root-directory.
	//
	fs := http.FileServer(http.Dir(h.path))
	http.Handle("/", fs)

	//
	// Build up the listen address.
	//
	listen := fmt.Sprintf("%s:%d", h.host, h.port)

	//
	// Log our start, and begin serving.
	//
	log.Printf("Serving upon http://%s/\n", listen)
	http.ListenAndServe(listen, logRequest(http.DefaultServeMux))
	return 0
}

// logRequest dumps the request to the console.
//
// Of course we don't know the return-code, but this is good enough
// for most of my use-cases.
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
