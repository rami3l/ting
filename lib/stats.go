package lib

import (
	"net"
	"time"
)

// Result is a record of a single run of the Tcping client.
type Result struct {
	RemoteAddr   net.Addr
	Error        error
	ResponseTime time.Duration
}

// Stats is the collection of Tcping client run records (Results).
type Stats struct {
	Results []Result
}

// Count returns the total number of records.
func (s Stats) Count() (c int) {
	return len(s.Results)
}

// SuccCount returns the number of success connections.
func (s Stats) SuccCount() (sc int) {
	for _, r := range s.Results {
		if r.Error == nil {
			sc++
		}
	}
	return
}

// FailCount returns the number of failed connections.
func (s Stats) FailCount() (fc int) {
	for _, r := range s.Results {
		if r.Error != nil {
			fc++
		}
	}
	return
}

// MaxTime returns the maximum connection time.
func (s Stats) MaxTime() (maxt time.Duration) {
	if s.Count() <= 0 {
		return
	}
	maxt = s.Results[0].ResponseTime
	for _, r := range s.Results[1:] {
		if t := r.ResponseTime; t > maxt {
			maxt = t
		}
	}
	return
}

// MinTime returns the minimum connection time.
func (s Stats) MinTime() (mint time.Duration) {
	if s.Count() <= 0 {
		return
	}
	mint = s.Results[0].ResponseTime
	for _, r := range s.Results[1:] {
		if t := r.ResponseTime; t < mint {
			mint = t
		}
	}
	return
}

// AvgTime returns the average connection time.
func (s Stats) AvgTime() (avgt time.Duration) {
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
