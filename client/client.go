package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/KarelKubat/smartlog/msg"
	"github.com/KarelKubat/smartlog/uri"
)

var (
	RestartWait     = time.Second / 10 // waittime between listener restarts
	RestartAttempts = 10               // # of restart attempts
)

type Client struct {
	// May be set by client code
	TimeFormat     string // defaults to YYYY-MM-DD HH:MM:SS localtime
	DebugThreshold uint8  // defaults to 0

	// Set by implementations
	Writer     io.Writer // writer for Info(f), Warn(f), Error(f)
	URI        *uri.URI  // URI from which the client was constructed
	Conn       net.Conn  // Only in network loggers
	IsTrueFile bool      // Only in file loggers
	Buffer     [][]byte  // only in HTTP loggers
}

func (c *Client) String() string {
	return fmt.Sprintf("%v", c.URI)
}

func (c *Client) Debug(lev uint8, message string) error {
	if lev > c.DebugThreshold {
		return nil
	}
	return c.sendToWriter(msg.Debug, message)
}

func (c *Client) Debugf(lev uint8, format string, args ...interface{}) error {
	return c.Debug(lev, fmt.Sprintf(format, args...))
}

func (c *Client) Info(message string) error {
	return c.sendToWriter(msg.Info, message)
}

func (c *Client) Infof(format string, args ...interface{}) error {
	return c.Info(fmt.Sprintf(format, args...))
}

func (c *Client) Warn(message string) error {
	return c.sendToWriter(msg.Warn, message)
}

func (c *Client) Warnf(format string, args ...interface{}) error {
	return c.Warn(fmt.Sprintf(format, args...))
}

func (c *Client) Fatal(message string) error {
	if err := c.sendToWriter(msg.Fatal, message); err != nil {
		return err
	}
	os.Exit(1)
	return nil // to satisfy the prototype
}

func (c *Client) Fatalf(format string, args ...interface{}) error {
	return c.Fatal(fmt.Sprintf(format, args...))
}

// Called by the server to pass messages already containing a timestamp etc. to clients.
func (c *Client) Passthru(buf []byte) error {
	if c.URI.Scheme == uri.None {
		return nil
	}
	return c.write(buf)
}

// Called by file:// clients.
func (c *Client) OpenFile() error {
	if c.URI.Scheme != uri.File || c.URI.Parts[0] == "stdout" {
		return fmt.Errorf("%v: attempt to open a file but this URI is not file-based", c)
	}
	var err error
	c.Writer, err = os.OpenFile(c.URI.Parts[0], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("%v: failed to create file: %v", c, err)
	}
	c.IsTrueFile = true
	return nil
}

// Called by network clients (tcp:// or udp://).
func (c *Client) Connect() error {
	if c.URI.Scheme != uri.TCP && c.URI.Scheme != uri.UDP {
		return fmt.Errorf("internal foobar, client.Connect isn't meant for %v", c)
	}
	var err error
	for i := 0; i < RestartAttempts; i++ {
		time.Sleep(RestartWait * time.Duration(i))
		c.Conn, err = net.Dial(c.URI.Scheme.String(), strings.Join(c.URI.Parts, ":"))
		if err == nil {
			c.Writer = c.Conn
			return nil
		}
	}
	return fmt.Errorf("%v: failed to (re)connect: %v", c, err)
}

func (c *Client) sendToWriter(lev msg.MsgType, message string) error {
	if c.URI.Scheme == uri.None {
		return nil
	}

	for _, buf := range msg.BytesFromMessage(&msg.Message{
		Type:       lev,
		TimeFormat: c.TimeFormat,
		Message:    message,
	}) {
		if err := c.write(buf); err != nil {
			return err
		}
	}

	// If the file disappears, reopen it
	if c.IsTrueFile {
		_, err := os.Stat(c.URI.Parts[0])
		if err != nil {
			return c.OpenFile()
		}
	}

	return nil
}

func (c *Client) write(buf []byte) error {
	nWritten := 0
	for nWritten < len(buf) {
		n, err := c.Writer.Write(buf)
		if err != nil {
			// Write errors on clients try to reconnect, if the error is due a broken pipe.
			// Otherwise the error goes to the caller for handling.
			if c.URI.Scheme != uri.TCP && c.URI.Scheme != uri.UDP || !strings.Contains(err.Error(), "broken pipe") {
				return fmt.Errorf("%v: write failure: %v", c, err)
			}
			Warnf("%v: write failure on network client: %v", c, err)
			if err := c.Connect(); err != nil {
				return err
			}
			n, err = c.Writer.Write(buf)
			if err != nil {
				return fmt.Errorf("%v: write failure despite reconnecting: %v", c, err)
			}
		}
		nWritten += n
	}
	return nil
}
