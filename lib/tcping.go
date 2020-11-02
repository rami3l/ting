package lib

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// Indicates a connection timeout.
const TimedOut = -1

type Stats []int64

type TcpingClient struct {
	Host     string
	Port     int
	tryCount int
	// Timeout in `ms`.
	timeOut  int64
	outputOn bool
}

func NewTcpingClient(host string, port *int, tryCount *int, timeOut *int64) *TcpingClient {
	port_, tryCount_, timeOut_ := 80, 10, int64(5000)
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
		fmt.Printf("Connecting to `%s`", net.JoinHostPort(host, strconv.Itoa(port)))
	}

	asyncConnect := func(done chan struct{}) {
		err = socket.Connect(host, port)
		done <- struct{}{}
	}

	done := make(chan struct{})
	t0 := time.Now()
	timer := time.NewTimer(time.Duration(c.timeOut) * time.Millisecond)

	go asyncConnect(done)

	select {
	case <-done:
		t := time.Since(t0)
		if err != nil {
			if c.outputOn {
				fmt.Printf(": %s\n", err)
			}
			return
		}
		remoteAddr = (*socket.Conn).RemoteAddr()
		if c.outputOn {
			fmt.Printf(" (%s)", remoteAddr)
		}
		responseTime = t.Nanoseconds()
		if c.outputOn {
			fmt.Printf(": time=%.2fms\n", float32(t)/1e6)
		}
		return

	case <-timer.C:
		responseTime = TimedOut
		if c.outputOn {
			fmt.Printf(": timed out after %.2fms\n", float32(c.timeOut))
		}
		return
	}
}
