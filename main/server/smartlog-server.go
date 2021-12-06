package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"smartlog/client/any"

	"smartlog/server"
)

const (
	version = "0.01"
	usage   = `
This is the Smartlog Server, catcher and forwarder of Smartlog-client generated
messages; version ` + version + `

Usage: smartlog-server [FLAGS] SERVERADDRESS CLIENT [CLIENT...]
Where:

  SERVERADDRESS defines what the server listens to and must be in the form:
    udp://HOSTNAME:PORT : (leave out the HOSTNAME to listen to all IPs), or
    tcp://HOSTNAME:PORT : (again, the HOSTNAME can be left out)

  CLIENTS defines where received messages are fanned out to. At least one must
  be given. Use one or more of:
    file://stdout      : dumps to stdout
    file://FILENAME     : appends to FILENAME
    tcp://HOSTNAME:PORT : forwards to a next hop over TCP
    udp://HOSTNAME:PORT : forwards to a next hop over UDP
    none://WHATEVER     : discards, useful for testing

  FLAGS may be:
`
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Supported flag(s)
	flagS := flag.Duration("s", 0, "stop server after stated duration, 0 = serve forever")

	// Parse options, show usage when that fails.
	flag.Usage = usageFunc
	flag.Parse()
	// We need at least 2 positional arguments, or show usage and stop.
	if flag.NArg() < 2 {
		usageFunc()
		return errors.New("(not enough arguments)")
	}

	// Start serving
	srv, err := server.New(flag.Arg(0))
	if err != nil {
		return err
	}
	if *flagS > 0 {
		go func() {
			time.Sleep(*flagS)
			srv.Close()
		}()
	}

	// Add clients from the commandline
	for _, uri := range flag.Args()[1:] {
		cl, err := any.New(uri)
		if err != nil {
			return err
		}
		srv.AddClient(cl)
	}
	return srv.Serve()
}

func usageFunc() {
	fmt.Fprintf(os.Stderr, usage)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stdout)
}
