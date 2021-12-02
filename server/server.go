package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"smartlog/client"
	"smartlog/linebuf"
	"smartlog/uri"
)

const (
	chSize = 1024 // # of messages that may be buffered while fanning out
)

type Server struct {
	URI         *uri.URI         // URI this was constructed from
	clients     []*client.Client // clients to fan out to
	bufCh       chan []byte      // msg channel for fanout to clients
	tcpListener net.Listener     // in the case of a TCP server
	udpConn     *net.UDPConn     // in the case of a UDP server
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
		s.tcpListener, err = net.Listen(fmt.Sprintf("%v", ur.Scheme), strings.Join(ur.Parts, ":"))
		if err != nil {
			return nil, fmt.Errorf("%v: failed to start TCP listener: %v", s, err)
		}
	case uri.UDP:
		addr, err := net.ResolveUDPAddr(fmt.Sprintf("%v", ur.Scheme), strings.Join(ur.Parts, ":"))
		if err != nil {
			return nil, fmt.Errorf("%v: failed to resolve address: %v", s, err)
		}
		s.udpConn, err = net.ListenUDP(fmt.Sprintf("%v", ur.Scheme), addr)
		if err != nil {
			return nil, fmt.Errorf("%v: failed to start UDP listener: %v", s, err)
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
		s.tcpServe()
	case uri.UDP:
		s.udpServe()
	default:
		errors.New("internal foobar, unhandled case in server.Serve")
	}
	return fmt.Errorf("%v: server stopped", s)
}

func (s *Server) Close() error {
	var err error
	switch s.URI.Scheme {
	case uri.TCP:
		err = s.tcpListener.Close()
	case uri.UDP:
		err = s.udpConn.Close()
	}
	return err
}

func (s *Server) udpServe() {
	line := linebuf.New()

	// Serve doesn't return
	for {
		for {
			buf := make([]byte, 1024)
			n, addr, err := s.udpConn.ReadFromUDP(buf)
			if err != nil {
				client.Warnf("%v: error while handling UDP connection from %v: %v", s, addr, err)
				continue
			}
			if n == 0 {
				return
			}
			line.Add(buf, n)
			for line.Complete() {
				s.bufCh <- line.Statement()
			}
		}
	}
}

func (s *Server) tcpServe() {
	for {
		conn, err := s.tcpListener.Accept()
		if err != nil {
			client.Warnf("%v: error accepting TCP connection: %v", s, err)
			continue
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
			client.Warnf("%v: error while handling TCP connection from %v: %v", s, conn.RemoteAddr(), err)
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
	for {
		buf := <-s.bufCh
		var wg sync.WaitGroup
		for _, c := range s.clients {
			wg.Add(1)
			go func(c *client.Client, buf []byte) {
				if err := c.Passthru(buf); err != nil {
					client.Warnf("%v: error while fanning out to client %v: %v, buf %v", s, c, err, string(buf))
				}
				wg.Done()
			}(c, buf)
		}
		wg.Wait()
	}
}
