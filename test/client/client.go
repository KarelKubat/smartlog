package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

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
Usage: client [FLAGS] DEST
DEST defines where to send to, e.g. udp://SERVER:PORT or file://this.log
FLAGS may be:
  -n NR : send at least NR messages
  -v    : display what is being sent
  -t    : show timing
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
	vFlag := flag.Bool("v", false, "log locally what is being sent")
	tFlag := flag.Bool("t", true, "show timing")
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

	nMessages := 0
	sendf := func(f func(string, ...interface{}) error, msg string, args ...interface{}) error {
		err := c.Infof(msg, args...)
		if *vFlag {
			client.Infof("sent: "+msg, args...)
		}
		nMessages++
		return err
	}

	start := time.Now()
	sendf(c.Infof, "------------- run start -------------")
	for nMessages <= *nFlag {
		if err = sendf(c.Infof, "informational message %d", nMessages); err != nil {
			return
		}
		if err = sendf(c.Warnf, "warning message %d", nMessages); err != nil {
			return
		}
		if err = sendf(c.Infof, "%v: %v", nMessages, lorem); err != nil {
			return
		}
	}
	if err = sendf(c.Infof, "------------- run end -------------"); err != nil {
		return
	}
	if *tFlag {
		client.Infof("sent %v messages in %v", nMessages, time.Now().Sub(start))
	}
}
