// Package log is an attempt to provide a drop-in replacement for https://pkg.go.dev/log.
// Some methods of Golang's `log` package are provided as wrapped replacements. Feel free to add more :)
package log

import (
	"fmt"
	"strings"

	"github.com/KarelKubat/smartlog/client"
)

func Fatal(v ...interface{}) {
	var m strings.Builder
	for _, p := range v {
		fmt.Fprintf(&m, fmt.Sprintf("%v", p))
	}
	client.Fatal(m.String())
}

func Fatalf(format string, v ...interface{}) {
	client.Fatalf(format, v...)
}

func Print(v ...interface{}) {
	var m strings.Builder
	for _, p := range v {
		fmt.Fprintf(&m, fmt.Sprintf("%v", p))
	}
	client.Info(m.String())
}

func Printf(format string, v ...interface{}) {
	client.Infof(format, v...)
}
