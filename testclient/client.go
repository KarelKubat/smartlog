package main

import (
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Supply 1 argument, e.g. `udp://SERVER:PORT` or `file://this.log`")
		os.Exit(1)
	}

	var err error
	defer func() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}()

	var c *client.Client
	c, err = any.New(os.Args[1])
	if err != nil {
		err = fmt.Errorf(
			"%v\nwhen using a TCP or UDP client: make sure to at least run `nc -ul $PORT` or `nc -tl $PORT` on the receiving host\n", err)
		return
	}

	c.Info("------------- run start -------------")
	for i := 1; i <= 10; i++ {
		if err = c.Infof("informational message %d", i); err != nil {
			return
		}
		if err = c.Warnf("warning message %d", i); err != nil {
			return
		}
	}
	c.Info(lorem)
	if err = c.Info("------------- run end -------------"); err != nil {
		return
	}
}
