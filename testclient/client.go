package main

import (
	"fmt"
	"os"

	"smartlog/client"
	"smartlog/client/any"
)

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
			"%v\nmake sure to at least run `nc -ul $PORT` or `nc -tl $PORT` on the receiving host\n", err)
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
	if err = c.Info("------------- run end -------------"); err != nil {
		return
	}
}
