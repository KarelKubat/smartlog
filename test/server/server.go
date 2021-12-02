package main

import (
	"fmt"
	"os"

	"smartlog/client"
	anyclient "smartlog/client/any"

	"smartlog/server"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Fprintln(os.Stderr, "Supply one server URI and one or more client URIs, e.g. `udp://:2021 file://stdout`")
		os.Exit(1)
	}

	var err error
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}()

	var srv *server.Server
	srv, err = server.New(os.Args[1])
	if err != nil {
		return
	}
	for _, uri := range os.Args[2:] {
		var cl *client.Client
		cl, err = anyclient.New(uri)
		if err != nil {
			return
		}
		srv.AddClient(cl)
	}
	err = srv.Serve()
}
