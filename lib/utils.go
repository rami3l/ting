package lib

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func SprintDuration(
	format string,
	t time.Duration,
	unit time.Duration,
) string {
	// fmt.Sprintf("%v", time.Microsecond) => "1µs"
	unitName := fmt.Sprintf("%v", unit)[1:]
	unitCount := float32(t.Nanoseconds()) / float32(unit)
	return fmt.Sprintf(format, unitCount) + unitName
}

func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}
