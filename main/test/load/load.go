package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"smartlog/client"
	"smartlog/client/any"
	"smartlog/server"
)

func main() {
	//	fThreads := flag.Int("threads", 1, "nr of parallel client runs")
	fDuration := flag.Duration("duration", time.Second, "load duration")
	fServer := flag.String("server", "udp://:2025", "server to start")
	fClients := flag.String("clients", "none://x", "comma-separated list of server clients")
	flag.Parse()
	if len(flag.Args()) != 1 || flag.Arg(0) != "go" {
		checkErr(errors.New("Usage: `load [FLAGS] go`\nRun `load -h` to see what flags you may set."))
	}

	// Initialize the server, add its clients, start serving, close it when we're done.
	srv, err := server.New(*fServer)
	checkErr(err)
	for _, uri := range strings.Split(*fClients, ",") {
		cl, err := any.New(uri)
		checkErr(err)
		srv.AddClient(cl)
	}
	go func(s *server.Server) {
		checkErr(srv.Serve())
	}(srv)
	defer func(s *server.Server) {
		s.Close()
	}(srv)

	// Repoint the default client to the server.
	client.DefaultClient, err = any.New(*fServer)
	checkErr(err)

	// Load until we're out of time.
	start := time.Now()
	msgNr := 0
	for {
		if time.Now().Sub(start) > *fDuration {
			break
		}
		msgNr++
		err = client.Warnf("hello there, this is message # %v", msgNr)
	}
	fmt.Printf("%v messages sent in %v\n", msgNr, time.Now().Sub(start))
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
