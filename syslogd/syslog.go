package syslogd

import (
	"log"
	"net"

	"github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"
)

type Format int

const (
	RFC3164 Format = iota // RFC3164: http://www.ietf.org/rfc/rfc3164.txt
	RFC5424               // RFC5424: http://www.ietf.org/rfc/rfc5424.txt
)

type Server struct {
	connections []net.Conn
	format      Format
	channel     chan syslogparser.LogParts
}

// NewServer build a new server
func NewServer() *Server {
	return &Server{}
}

// SetFormat the syslog format
func (s *Server) SetFormat(format Format) {
	s.format = format
}

// ListenUDP add a server to listen on an UDP addr
func (s *Server) ListenUDP(addr string) error {
	log.Printf("started UDP listener")
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	s.connections = append(s.connections, connection)
	return nil
}

// ListenUnixgram add a server to listen on a unix socket
func (s *Server) ListenUnixgram(addr string) error {
	unixAddr, err := net.ResolveUnixAddr("unixgram", addr)
	if err != nil {
		return err
	}

	connection, err := net.ListenUnixgram("unixgram", unixAddr)
	if err != nil {
		return err
	}

	s.connections = append(s.connections, connection)
	return nil
}

// Start reading syslog records
func (s *Server) Start(channel chan syslogparser.LogParts) {
	s.channel = channel
	for _, conn := range s.connections {
		s.readDatagrams(conn)
	}
}

type readFunc func([]byte) (int, net.Addr, error)

func (s *Server) readDatagrams(conn net.Conn) {
	log.Printf("scanning connection %s", conn.LocalAddr())

	switch c := conn.(type) {
	case *net.UDPConn:
		go s.read(func(buf []byte) (int, net.Addr, error) {
			return c.ReadFromUDP(buf)
		})
	case *net.UnixConn:
		go s.read(func(buf []byte) (int, net.Addr, error) {
			return c.ReadFromUnix(buf)
		})
	}

}

func (s *Server) read(read readFunc) {

	buf := make([]byte, 4096)

	for {
		n, addr, err := read(buf)

		if err != nil {
			break
		}

		parts, err := s.parse(buf[:n])

		parts["source"] = addr.String()

		s.channel <- parts

	}

}

func (s *Server) parse(buf []byte) (syslogparser.LogParts, error) {

	var p syslogparser.LogParser

	switch s.format {
	case RFC3164:
		p = rfc3164.NewParser(buf)
	case RFC5424:
		p = rfc5424.NewParser(buf)
	default:
		p = rfc3164.NewParser(buf)
	}

	err := p.Parse()

	if err != nil {
		return nil, err
	}

	parts := p.Dump()

	return parts, nil
}
