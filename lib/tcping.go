package lib

import (
	"fmt"
	"net"
	"time"
)

// Indicates a connection timeout.
const TimedOut = -1

func microseconds(t time.Duration) float32 {
	return float32(t.Nanoseconds()) / 1e6
}

type Result struct {
	ResponseTime time.Duration
	RemoteAddr   net.Addr
}

type Stats struct {
	Results []Result
}

func (s Stats) Count() (c int) {
	return len(s.Results)
}

func (s Stats) SuccCount() (sc int) {
	for _, r := range s.Results {
		if r.ResponseTime.Nanoseconds() > 0 {
			sc++
		}
	}
	return
}

func (s Stats) FailCount() (fc int) {
	for _, r := range s.Results {
		if r.ResponseTime.Nanoseconds() <= 0 {
			fc++
		}
	}
	return
}

func (s Stats) MaxTime() (mt time.Duration) {
	if s.Count() <= 0 {
		return
	}
	mt = s.Results[0].ResponseTime
	for _, r := range s.Results[1:] {
		if t := r.ResponseTime; t > mt {
			mt = t
		}
	}
	return
}

func (s Stats) MinTime() (mt time.Duration) {
	if s.Count() <= 0 {
		return
	}
	mt = s.Results[0].ResponseTime
	for _, r := range s.Results[1:] {
		if t := r.ResponseTime; t < mt {
			mt = t
		}
	}
	return
}

func (s Stats) AvgTime() (at time.Duration) {
	if s.Count() <= 0 {
		return
	}
	var st time.Duration
	for _, r := range s.Results {
		st += r.ResponseTime
	}
	avg := st.Nanoseconds() / int64(s.Count())
	return time.Duration(avg) * time.Nanosecond
}

type TcpingClient struct {
	Host     string
	Port     int
	tryCount int
	timeOut  time.Duration
	outputOn bool
}

func NewTcpingClient(host string, port *int, tryCount *int, timeOut *time.Duration) *TcpingClient {
	port_, tryCount_, timeOut_ := 80, 10, 5*time.Second
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
	timer := time.NewTimer(c.timeOut)

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
	results := []Result{}
	for i := 0; i < c.tryCount; func() { time.Sleep(1 * time.Second); i++ }() {
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
		/*
			--- api.github.com[:80] tcping statistics ---
			10 connections, 10 succeeded, 0 failed, 100.00% success rate
			minimum = 233.51ms, maximum = 251.77ms, average = 243.40ms
		*/
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
