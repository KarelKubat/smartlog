package server

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"smartlog/client"
	"smartlog/linebuf"
	"smartlog/uri"
)

const (
	chSize = 1024 // # of messages that may be buffered while fanning out
)

type Server struct {
	serverType  Type             // tcp or udp
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

	// Port must be valid
	port, err := uri.Port(u, ur.Parts[1])
	if err != nil {
		return nil, err
	}

	s := &Server{
		bufCh: make(chan []byte, chSize),
	}

	// Set the connection
	switch ur.Scheme {
	case uri.TCP:
		s.serverType = tcp
		s.tcpListener, err = net.Listen("tcp", fmt.Sprintf(":%v", port))
	case uri.UDP:
		s.serverType = udp
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%v", port))
		if err != nil {
			return nil, err
		}
		s.udpConn, err = net.ListenUDP("udp", addr)
		if err != nil {
			return nil, err
		}
	default:
		return nil, uri.Errorf(u, "only udp:// or tcp:// servers are supported (not %v)", ur.Scheme)
	}

	return s, nil
}

func (s *Server) AddClient(c *client.Client) {
	s.clients = append(s.clients, c)
}

func (s *Server) Serve() error {
	go func() {
		s.fanout()
	}()

	switch s.serverType {
	case tcp:
		s.tcpServe()
	case udp:
		s.udpServe()
	default:
		errors.New("internal foobar, unhandled case in server.Serve")
	}
	return errors.New("server stopped")
}

func (s *Server) udpServe() {
	line := linebuf.New()

	// Serve doesn't return
	for {
		for {
			buf := make([]byte, 1024)
			n, addr, err := s.udpConn.ReadFromUDP(buf)
			if err != nil {
				client.Warnf("error while handling UDP connection from %v: %v", addr, err)
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
			client.Warnf("error accepting TCP connection: %v", err)
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
		if err != nil && err.Error() != "EOF" {
			client.Warnf("error while handling TCP connection from %v: %v", conn.RemoteAddr(), err)
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
			go func(buf []byte) {
				if err := c.Passthru(buf); err != nil {
					client.Warnf("error while fanning out to client %v: %v, buf %v", c, err, string(buf))
				}
				wg.Done()
			}(buf)
		}
		wg.Wait()
	}
}
