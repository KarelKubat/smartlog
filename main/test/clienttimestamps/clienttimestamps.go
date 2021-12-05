package main

import (
	"fmt"
	"os"
	"time"

	"smartlog/client/any"
)

func main() {
	var err error

	checkErr := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err.Error())
			os.Exit(1)
		}
	}

	cl1, err := any.New("file://stdout")
	checkErr(err)

	cl2, err := any.New("file://stdout")
	checkErr(err)
	cl2.TimeFormat = time.RFC3339

	cl1.Info("hello from client #1")
	cl2.Info("hello from client #2")
}
