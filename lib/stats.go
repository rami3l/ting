package lib

import (
	"net"
	"time"
)

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
		if r.ResponseTime > 0 {
			sc++
		}
	}
	return
}

func (s Stats) FailCount() (fc int) {
	for _, r := range s.Results {
		if r.ResponseTime <= 0 {
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
