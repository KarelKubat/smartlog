package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"smartlog/client"
	"smartlog/client/any"
	"smartlog/server"
)

var (
	fThreads  = flag.Int("threads", 1, "nr of parallel client runs")
	fDuration = flag.Duration("duration", time.Second, "load duration")
	fServer   = flag.String("server", "udp://:2025", "server to start")
	fClients  = flag.String("clients", "none://x", "comma-separated list of server clients")
)

func main() {
	checkErr(parseCmdLine())

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

	// Start indicated threads.
	msgsPerThread := make([]int, *fThreads)
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < *fThreads; i++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			// Load until we're out of time.
			msgNr := 0
			for {
				if time.Now().Sub(start) > *fDuration {
					break
				}
				msgNr++
				checkErr(client.Warnf("hello there, this is message # %v", msgNr))
			}
			msgsPerThread[slot] = msgNr
		}(i)
	}
	wg.Wait()

	total := 0
	for i := 0; i < *fThreads; i++ {
		fmt.Printf("thread %v sent %v messages\n", i, msgsPerThread[i])
		total += msgsPerThread[i]
	}
	fmt.Printf("average: %v in %v\n", total / *fThreads, *fDuration)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func parseCmdLine() error {
	flag.Parse()
	if len(flag.Args()) != 1 || flag.Arg(0) != "go" {
		return errors.New("Usage: `load [FLAGS] go`\nRun `load -h` to see what flags you may set.")
	}
	if *fThreads < 1 {
		return errors.New("--threads NR: must be a positive number")
	}
	return nil
}
