package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"smartlog/client"
	"smartlog/client/any"
)

// Test msg containing newlines and empty lines
const lorem = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse 
cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident, 
sunt in culpa qui officia deserunt mollit anim id est laborum.

`

const usage = `
Supply 1 argument: where to send to, e.g. udp://SERVER:PORT or file://this.log
Use flag -n <nr> to control the number of sent messages.
`

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}()

	nFlag := flag.Int("n", 100, "number of messages to send")
	flag.Parse()
	if flag.NArg() != 1 {
		err = errors.New(usage)
		return
	}

	var c *client.Client
	c, err = any.New(flag.Arg(0))
	if err != nil {
		err = fmt.Errorf(
			"%v\nwhen using a TCP or UDP client: make sure to at least run `nc -ul $PORT` or `nc -tl $PORT` on the receiving host\n", err)
		return
	}

	c.Info("------------- run start -------------")
	nMessages := 1
	for nMessages <= *nFlag {
		if err = c.Infof("informational message %d", nMessages); err != nil {
			return
		}
		nMessages++
		if err = c.Warnf("warning message %d", nMessages); err != nil {
			return
		}
		nMessages++
		if err = c.Info(fmt.Sprintf("%v: %v", nMessages, lorem)); err != nil {
			return
		}
		nMessages++
	}
	if err = c.Info("------------- run end -------------"); err != nil {
		return
	}
}
