package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"smartlog/client"
	"smartlog/linebuf"
	"smartlog/msg"
	"smartlog/uri"
)

const (
	chSize    = 1024              // # of messages that may be buffered while fanning out
	dropInfo  = chSize * 75 / 100 // drop Info(f) when 75% full
	dropDebug = chSize * 50 / 100 // drop Debug(f) when 50% full
)

var (
	RestartWait     = time.Second / 10 // waittime between listener restarts
	RestartAttempts = 10               // # of restart attempts
)

type Server struct {
	URI         *uri.URI         // URI this was constructed from
	clients     []*client.Client // clients to fan out to
	bufCh       chan []byte      // msg channel for fanout to clients
	tcpListener net.Listener     // in the case of a TCP server
	udpConn     *net.UDPConn     // in the case of a UDP server
	closed      bool             // true upon server.Close()
}

func New(u string) (*Server, error) {
	// Parse URI, we support: tcp://mush:port and udp://mush:port
	ur, err := uri.New(u)
	if err != nil {
		return nil, err
	}

	s := &Server{
		URI:   ur,
		bufCh: make(chan []byte, chSize),
	}

	// Set the connection
	switch ur.Scheme {
	case uri.TCP:
		if err := s.tcpStartListener(); err != nil {
			return nil, err
		}
	case uri.UDP:
		if err := s.udpStartListener(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%v: only udp:// or tcp:// servers are supported", s)
	}

	return s, nil
}

func (s *Server) String() string {
	return fmt.Sprintf("%v", s.URI)
}

func (s *Server) AddClient(c *client.Client) {
	s.clients = append(s.clients, c)
}

func (s *Server) Serve() error {
	go func() {
		s.fanout()
	}()

	switch s.URI.Scheme {
	case uri.TCP:
		if err := s.tcpServe(); err != nil {
			return fmt.Errorf("%v: TCP server stopped: %v", s, err)
		}
	case uri.UDP:
		if err := s.udpServe(); err != nil {
			return fmt.Errorf("%v: UDP server stopped: %v", s, err)
		}
	default:
		errors.New("internal foobar, unhandled case in server.Serve")
	}
	return nil
}

func (s *Server) Close() error {
	s.closed = true

	var err error
	switch s.URI.Scheme {
	case uri.TCP:
		err = s.tcpListener.Close()
	case uri.UDP:
		err = s.udpConn.Close()
	}
	return err
}

func (s *Server) udpStartListener() error {
	var err error
	var addr *net.UDPAddr

	for i := 0; i < RestartAttempts; i++ {
		time.Sleep(RestartWait * time.Duration(i))
		addr, err = net.ResolveUDPAddr(s.URI.Scheme.String(), strings.Join(s.URI.Parts, ":"))
		if err != nil {
			return fmt.Errorf("%v: failed to resolve address: %v", s, err)
		}
		s.udpConn, err = net.ListenUDP(s.URI.Scheme.String(), addr)
		if err == nil {
			return nil
		}
		if err != nil {
			return fmt.Errorf("%v: failed to start UDP listener: %v", s, err)
		}
	}
	return fmt.Errorf("%v: failed to start UDP listener: %v", s, err)
}

func (s *Server) udpServe() error {
	line := linebuf.New()

	// Don't return unless the server gets closed.
	for {
		for {
			buf := make([]byte, 1024)
			n, addr, err := s.udpConn.ReadFromUDP(buf)
			if err != nil {
				if s.closed {
					return nil
				}
				client.Warnf("%v: failed to handle UDP connection from %v: %v", s, addr, err)
				if err := s.udpStartListener(); err != nil {
					return err
				}
			}
			if n == 0 {
				continue
			}
			line.Add(buf, n)
			for line.Complete() {
				s.bufCh <- line.Statement()
			}
		}
	}
}

func (s *Server) tcpStartListener() error {
	var err error
	for i := 0; i < RestartAttempts; i++ {
		time.Sleep(RestartWait * time.Duration(i))
		s.tcpListener, err = net.Listen(s.URI.Scheme.String(), strings.Join(s.URI.Parts, ":"))
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("%v: failed to start TCP listener: %v", s, err)
}

func (s *Server) tcpServe() error {
	// Don't return unless the connection gets closed.
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			if s.closed {
				return nil
			}
			client.Warnf("%v: failed to accept TCP connection: %v", s, err)
			continue // restart listener
		}
		go s.handleTCPConnection(conn)
	}
}

func (s *Server) handleTCPConnection(conn net.Conn) {
	line := linebuf.New()
	var err error
	defer func() {
		for line.Complete() {
			s.bufCh <- line.Statement()
		}
		if err != nil {
			fmt.Println(err)
		}
		if err != nil && err.Error() != "EOF" {
			client.Warnf("%v: failed to handle TCP connection from %v: %v", s, conn.RemoteAddr(), err)
		}
		conn.Close()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n > 0 {
			line.Add(buf, n)
			for line.Complete() {
				s.bufCh <- line.Statement()
			}
		}
		if err != nil {
			return
		}
	}
}

func (s *Server) fanout() {
	// fanout() never returns, but forever consumes messages from the bufCh
	var dropped bool

	for {
		buf := <-s.bufCh

		// The threshold to drop debug messages is lowest. If that is overrun then we need to reparse the message,
		// see what type it is and maybe drop it.
		chLen := len(s.bufCh)
		if chLen > dropDebug {
			var shouldDrop bool
			t := msg.TypeFromBytes(buf)
			switch t {
			case msg.Debug:
				shouldDrop = true
			case msg.Info:
				if chLen > dropInfo {
					shouldDrop = true
				}
			}
			if shouldDrop {
				if !dropped {
					dropped = true
					client.Warnf("%v: dropping debug/info message(s), %v already buffered, limit %v", s, chLen, dropInfo)
					continue
				}
			}
		}

		dropped = false
		var wg sync.WaitGroup
		for _, c := range s.clients {
			wg.Add(1)
			go func(c *client.Client, buf []byte) {
				if err := c.Passthru(buf); err != nil {
					client.Warnf("%v: failed to fanout to client %v: %v, buf %v",
						s, c, err, strings.TrimRight(string(buf), "\n"))
				}
				wg.Done()
			}(c, buf)
		}
		wg.Wait()
	}
}
