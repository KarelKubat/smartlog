package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"smartlog/uri"
)

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
	TimeFormat string    // defaults to YYYY-MM-DD HH:MM:SS localtime
	Writer     io.Writer // writer for Info(f), Warn(f), Error(f)
	URI        *uri.URI  // URI from which the client was constructed

	// Present in network loggers
	Conn net.Conn
}

func (c *Client) String() string {
	return fmt.Sprintf("%v", c.URI)
}

func (c *Client) Info(msg string) error {
	return c.sendToWriter(infoTag, msg)
}

func (c *Client) Infof(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Info(full)
}

func (c *Client) Warn(msg string) error {
	return c.sendToWriter(warnTag, msg)
}

func (c *Client) Warnf(format string, args ...interface{}) error {
	full := fmt.Sprintf(format, args...)
	return c.Warn(full)
}

func (c *Client) Error(msg string) error {
	if err := c.sendToWriter(errorTag, msg); err != nil {
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

func (c *Client) sendToWriter(lev byte, msg string) error {
	prefix := append(c.timeStamp(), space, separator, space, lev, space, separator, space)

	for _, line := range strings.Split(msg, "\n") {
		if line == "" {
			continue
		}
		out := append(append(prefix, []byte(line)...), '\n')
		if _, err := c.Writer.Write(out); err != nil {
			return fmt.Errorf("write failure to %v: %v", c, err)
		}
	}
	return nil
}
