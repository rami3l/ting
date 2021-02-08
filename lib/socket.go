package lib

import (
	"errors"
	"net"
)

// Socket is a wrapper of net.Conn.
type Socket struct {
	Conn     *net.Conn
	ConnType string
}

// NewSocket initiates a Socket in default settings.
func NewSocket(connType string) *Socket {
	return &Socket{
		Conn:     nil,
		ConnType: connType,
	}
}

// Connect tries to open the socket connection.
func (s *Socket) Connect(host string, port int) (err error) {
	conn, err := net.Dial(s.ConnType, JoinHostPort(host, port))
	s.Conn = &conn
	return
}

// Close tries to close the socket connection.
func (s *Socket) Close() (err error) {
	if s.Conn == nil {
		return errors.New("connection not established")
	}
	return (*s.Conn).Close()
}
