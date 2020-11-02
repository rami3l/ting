package lib

import (
	"errors"
	"net"
	"strconv"
)

type Socket struct {
	Conn     *net.Conn
	ConnType string
}

func NewSocket(connType string) *Socket {
	return &Socket{
		Conn:     nil,
		ConnType: connType,
	}
}

func (s *Socket) Connect(host string, port int) (err error) {
	conn, err := net.Dial(s.ConnType, net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return
	}
	s.Conn = &conn
	return
}

func (s *Socket) Close() (err error) {
	if s.Conn == nil {
		return errors.New("connection not established")
	}
	return (*s.Conn).Close()
}
