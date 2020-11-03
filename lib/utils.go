package lib

import (
	"net"
	"strconv"
)

func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
