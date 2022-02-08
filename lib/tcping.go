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

// TcpingClient is a ping-like speed test client, but works under TCP.
type TcpingClient struct {
	Host        string
	Port        int
	tryCount    int
	tryInterval time.Duration
	timeout     time.Duration
	outputOn    bool
}

// NewTcpingClient initializes a TcpingClient in default settings.
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

// EnableOutput turns on the output of TcpingClient to stdout.
func (c *TcpingClient) EnableOutput() *TcpingClient {
	c.outputOn = true
	return c
}

// HostAndPort returns the "host:port" pair of this client.
func (c TcpingClient) HostAndPort() string {
	return JoinHostPort(c.Host, c.Port)
}

// RunOnce makes a single tcping test.
func (c TcpingClient) RunOnce() (responseTime time.Duration, remoteAddr net.Addr, err error) {
	socket := NewSocket("tcp")
	if c.outputOn {
		fmt.Printf("Connecting to `%s`", c.HostAndPort())
	}

	done := make(chan time.Duration)
	t0 := time.Now()
	timer := time.NewTimer(c.timeout)
	responseTime = TimedOut

	go func() {
		err = socket.Connect(c.Host, c.Port)
		done <- time.Since(t0)
	}()

	select {
	// Connection finished (or returned an error) before timeout.
	case t := <-done:
		if err != nil {
			if c.outputOn {
				fmt.Printf(": %s\n", err)
			}
			return
		}
		responseTime = t
		remoteAddr = (*socket.Conn).RemoteAddr()
		if c.outputOn {
			fmt.Printf(
				" (%s): time=%s\n",
				remoteAddr,
				SprintDuration("%.2f", responseTime, time.Millisecond),
			)
		}

	// Connection timed out.
	case <-timer.C:
		if c.outputOn {
			fmt.Printf(
				": timed out after %s\n",
				SprintDuration("%.2f", c.timeout, time.Millisecond),
			)
		}
	}

	return
}

// Run makes several consequent tcping tests and analyzes the overall result.
func (c TcpingClient) Run() (s Stats) {
	// Handle SIGINT and SIGTERM
	signalNotifier := make(chan os.Signal, 5)
	signal.Notify(signalNotifier, os.Interrupt, syscall.SIGTERM)

	results := []Result{}

Loop:
	for i := 0; i < c.tryCount; func() { time.Sleep(c.tryInterval); i++ }() {
		var (
			responseTime time.Duration
			remoteAddr   net.Addr
			err          error
		)

		// Notifier of a finished tcping test.
		done := make(chan struct{})

		go func() {
			// Show the number of tries.
			if c.outputOn {
				fmt.Printf("%3d> ", i)
			}
			// We discard all errors here ON PURPOSE:
			// errors should not stop the looping.
			responseTime, remoteAddr, err = c.RunOnce()
			done <- struct{}{}
		}()

		// If we have received a signal, we need to break the loop early.
		select {
		case <-signalNotifier:
			fmt.Println("\r  <Ctrl+C>")
			break Loop
		case <-done:
			results = append(results, Result{
				ResponseTime: responseTime,
				RemoteAddr:   remoteAddr,
				Error:        err,
			})
		}
	}

	// Analyze and print the final result.
	s = Stats{Results: results}
	if c.outputOn {
		count := s.Count()
		succCount := s.SuccCount()
		failCount := count - succCount
		succRate := float32(succCount) / float32(count)
		minTime := SprintDuration("%.2f", s.MinTime(), time.Millisecond)
		maxTime := SprintDuration("%.2f", s.MaxTime(), time.Millisecond)
		avgTime := SprintDuration("%.2f", s.AvgTime(), time.Millisecond)

		fmt.Printf(`
--- %s tcping statistics ---
%d connections, %d succeeded, %d failed, %.2f%% success rate
minimum = %s, maximum = %s, average = %s
`,
			c.HostAndPort(),
			count, succCount, failCount, succRate*100,
			minTime, maxTime, avgTime,
		)
	}
	return
}
