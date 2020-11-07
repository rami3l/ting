package lib

import (
	"net"
	"strconv"
	"time"
)

func microseconds(t time.Duration) float32 {
	return float32(t.Nanoseconds()) / 1e6
}

func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
