package lib

import (
	"fmt"
	"net"
	"time"
)

// Indicates a connection timeout.
const TimedOut = -1

type Stats []int64

type TcpingClient struct {
	Host     string
	Port     int
	tryCount int
	timeOut  int
	outputOn bool
}

func NewTcpingClient(host string, port *int, tryCount *int, timeOut *int) *TcpingClient {
	port_, tryCount_, timeOut_ := 80, 10, 5000
	if port != nil {
		port_ = *port
	}
	if tryCount != nil {
		tryCount_ = *tryCount
	}
	if timeOut != nil {
		timeOut_ = *timeOut
	}
	return &TcpingClient{
		Host:     host,
		Port:     port_,
		tryCount: tryCount_,
		timeOut:  timeOut_,
	}
}

func (c *TcpingClient) EnableOutput() *TcpingClient {
	c.outputOn = true
	return c
}

func (c TcpingClient) RunOnce() (responseTime int64, remoteAddr net.Addr, err error) {
	socket := NewSocket("tcp")
	host, port := c.Host, c.Port
	if c.outputOn {
		fmt.Printf("Connecting to %s", net.JoinHostPort(host, string(port)))
	}

	asyncConnect := func(done chan struct{}) {
		err = socket.Connect(host, port)
		done <- struct{}{}
	}

	done := make(chan struct{})
	t0 := time.Now()
	timer := time.NewTimer(time.Duration(c.timeOut))

	go asyncConnect(done)

	select {
	case <-done:
		t := time.Since(t0)
		remoteAddr = (*socket.Conn).RemoteAddr()
		if c.outputOn {
			fmt.Printf(" (%s)", remoteAddr)
		}
		if err != nil {
			if c.outputOn {
				fmt.Printf(": %s", err)
			}
			return
		}
		responseTime = t.Nanoseconds()
		if c.outputOn {
			fmt.Printf(": time=%.2fms", float32(t)/1e3)
		}
		return

	case <-timer.C:
		responseTime = TimedOut
		remoteAddr = (*socket.Conn).RemoteAddr()
		if c.outputOn {
			fmt.Printf(" (%s): timed out", remoteAddr)
		}
		return
	}
}
