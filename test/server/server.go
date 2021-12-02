package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"smartlog/client"
	anyclient "smartlog/client/any"

	"smartlog/server"
)

const usage = `
Usage: server [FLAGS] SERVERADDRESS CLIENTS
SERVERADDRESS defines what the server listens to and must be in the form:
  udp://HOSTNAME:PORT (leave out the HOSTNAME to listen to all available addresses), or
  tcp://HOSTNAME:PORT (again, the HOSTNAME can be left out)
CLIENTS defines where received messages are fanned out to:
  udp://HOSTNAME:PORT or tcp://HOSTNAME:PORT forwards over the network
  file://stdout dumps to stdout, file://FILENAME overwrites FILENAME and dumps in there
  none://WHATEVER discards, useful for testing
FLAGS may be:
  -s DURATION: stops the server after the stated duration, useful for testing
`

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()

	flagS := flag.Duration("s", 0, "stop server after stated duration, 0 = serve forever")
	flag.Parse()
	if flag.NArg() < 2 {
		err = errors.New(usage)
		return
	}

	var srv *server.Server
	srv, err = server.New(flag.Arg(0))
	if err != nil {
		return
	}
	if *flagS > 0 {
		go func() {
			time.Sleep(*flagS)
			srv.Close()
		}()
	}

	for _, uri := range flag.Args()[1:] {
		var cl *client.Client
		cl, err = anyclient.New(uri)
		if err != nil {
			return
		}
		srv.AddClient(cl)
	}
	err = srv.Serve()
}
