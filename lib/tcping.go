package lib

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Indicates a connection timeout.
const TimedOut = -1

type TcpingClient struct {
	Host        string
	Port        int
	tryCount    int
	tryInterval time.Duration
	timeout     time.Duration
	outputOn    bool
}

func NewTcpingClient(host string) *TcpingClient {
	return &TcpingClient{
		Host:        host,
		Port:        80,
		tryCount:    5,
		tryInterval: 1 * time.Second,
		timeout:     5 * time.Second,
	}
}

func (c *TcpingClient) SetPort(port int) *TcpingClient {
	c.Port = port
	return c
}

func (c *TcpingClient) SetTryCount(tryCount int) *TcpingClient {
	c.tryCount = tryCount
	return c
}

func (c *TcpingClient) SetTryInterval(tryInterval time.Duration) *TcpingClient {
	c.tryInterval = tryInterval
	return c
}

func (c *TcpingClient) SetTimeout(timeout time.Duration) *TcpingClient {
	c.timeout = timeout
	return c
}

func (c *TcpingClient) EnableOutput() *TcpingClient {
	c.outputOn = true
	return c
}

func (c TcpingClient) HostAndPort() string {
	return JoinHostPort(c.Host, c.Port)
}

func (c TcpingClient) RunOnce() (responseTime time.Duration, remoteAddr net.Addr, err error) {
	socket := NewSocket("tcp")
	if c.outputOn {
		fmt.Printf("Connecting to `%s`", c.HostAndPort())
	}

	asyncConnect := func(done chan struct{}) {
		err = socket.Connect(c.Host, c.Port)
		done <- struct{}{}
	}

	done := make(chan struct{})
	t0 := time.Now()
	timer := time.NewTimer(c.timeout)

	go asyncConnect(done)

	select {
	case <-done:
		responseTime = time.Since(t0)
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
		if c.outputOn {
			fmt.Printf(": time=%.2fms\n", microseconds(responseTime))
		}
		return

	case <-timer.C:
		responseTime = TimedOut
		if c.outputOn {
			fmt.Printf(": timed out after %.2fms\n", microseconds(responseTime))
		}
		return
	}
}

func (c TcpingClient) Run() (s Stats, err error) {
	// Handle SIGINT and SIGTERM
	signalNotifier := make(chan os.Signal, 5)
	signal.Notify(signalNotifier, os.Interrupt, syscall.SIGTERM)

	results := []Result{}

Loop:
	for i := 0; i < c.tryCount; func() { time.Sleep(c.tryInterval); i++ }() {
		select {
		case <-signalNotifier:
			fmt.Println("\r- Ctrl+C")
			break Loop
		default:
		}
		if c.outputOn {
			fmt.Printf("%3d> ", i)
		}
		if responseTime, remoteAddr, err := c.RunOnce(); err != nil {
			return Stats{Results: results}, err
		} else {
			results = append(results, Result{
				ResponseTime: responseTime,
				RemoteAddr:   remoteAddr,
			})
		}
	}

	s = Stats{Results: results}
	if c.outputOn {
		count := s.Count()
		succCount := s.SuccCount()
		failCount := count - succCount
		succRate := float32(succCount) / float32(count)
		fmt.Printf(`
--- %s tcping statistics ---
%d connections, %d succeeded, %d failed, %.2f%% success rate
minimum = %.2fms, maximum = %.2fms, average = %.2fms
`,
			c.HostAndPort(), count, succCount, failCount, succRate*100,
			microseconds(s.MinTime()), microseconds(s.MaxTime()), microseconds(s.AvgTime()),
		)
	}
	return
}
