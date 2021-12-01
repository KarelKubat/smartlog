package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type ClientType int

const (
	Std ClientType = iota
	File
	Network
)

func (t ClientType) String() string {
	return []string{"std://", "file://", "udp:// or tcp://"}[t]
}

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"

	infoTag   byte = 'I'
	warnTag        = 'W'
	errorTag       = 'E'
	separator      = '|'
	space          = ' '
)

type Client struct {
	// Present in all loggers
	Type       ClientType // std, file, network
	TimeFormat string     // defaults to YYYY-MM-DD HH:MM:SS localtime
	Writer     io.Writer  // writer for Info(f), Warn(f), Error(f)

	// Present in file loggers
	Filename string

	// Present in network loggers
	Conn net.Conn
}

func (c *Client) String() string {
	return fmt.Sprintf("%v...", c.Type)
}

func (c *Client) Info(msg string) error {
	return sendToWriter(c.Writer, c.timeStamp(), infoTag, msg)
}

func (c *Client) Infof(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Info(full)
}

func (c *Client) Warn(msg string) error {
	return sendToWriter(c.Writer, c.timeStamp(), warnTag, msg)
}

func (c *Client) Warnf(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Warn(full)
}

func (c *Client) Error(msg string) error {
	if err := sendToWriter(c.Writer, c.timeStamp(), errorTag, msg); err != nil {
		return err
	}
	os.Exit(1)
	return nil // to satisfy the prototype
}

func (c *Client) Errorf(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Error(full)
}

func (c *Client) Passthru(msg []byte) error {
	_, err := c.Writer.Write(msg)
	return err
}

func (c *Client) timeStamp() []byte {
	now := time.Now()
	var stamp string
	if c.TimeFormat == "" {
		stamp = now.String()
	} else {
		stamp = now.Format(c.TimeFormat)
	}
	return []byte(stamp)
}

func sendToWriter(wr io.Writer, stamp []byte, lev byte, msg string) error {
	if _, err := wr.Write(stamp); err != nil {
		return err
	}
	if _, err := wr.Write([]byte{space, separator, space, lev, space, separator, space}); err != nil {
		return err
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	if _, err := wr.Write([]byte(msg)); err != nil {
		return err
	}
	return nil
}
