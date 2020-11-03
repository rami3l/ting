package cmd

import (
	"fmt"

	"github.com/rami3l/ting/lib"
)

func Main() {
	fmt.Printf("Hello, ting!\n")
	testTcpingClient()
}

func testTcpingClient() {
	client := lib.NewTcpingClient("google.com", nil, nil, nil).EnableOutput()
	client.Run()
}
