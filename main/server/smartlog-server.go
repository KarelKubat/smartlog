package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"smartlog/client/any"

	"smartlog/server"
)

const usage = `
Usage: smartlog-server [FLAGS] SERVERADDRESS CLIENT(S)
SERVERADDRESS defines what the server listens to and must be in the form:
  udp://HOSTNAME:PORT (leave out the HOSTNAME to listen to all available addresses), or
  tcp://HOSTNAME:PORT (again, the HOSTNAME can be left out)
CLIENTS defines where received messages are fanned out to. At least one must be given.
  udp://HOSTNAME:PORT or tcp://HOSTNAME:PORT forwards over the network
  file://stdout dumps to stdout, file://FILENAME appends to FILENAME
  none://WHATEVER discards, useful for testing
FLAGS may be:
`

func main() {
	// Supported flag(s)
	flagS := flag.Duration("s", 0, "stop server after stated duration, 0 = serve forever")

	// usageFunc() shows how to invoke the server
	usageFunc := func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Usage = usageFunc
	flag.Parse()
	if flag.NArg() < 2 {
		usageFunc()
	}

	// Start serving
	srv, err := server.New(flag.Arg(0))
	checkErr(err)
	if *flagS > 0 {
		go func() {
			time.Sleep(*flagS)
			srv.Close()
		}()
	}

	// Add clients from the commandline
	for _, uri := range flag.Args()[1:] {
		cl, err := any.New(uri)
		checkErr(err)
		srv.AddClient(cl)
	}
	checkErr(srv.Serve())
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
