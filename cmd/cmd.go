package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rami3l/ting/lib"
)

func Main() {
	fmt.Printf("Hello, ting!\n")
	testTcpingClient()
}

// SetupInterruptHandler creates a listener on a new goroutine notifying the program
// if it receives an interrupt from the OS.
func SetupInterruptHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		// DeleteFiles()
		os.Exit(0)
	}()
}

func testTcpingClient() {
	client := lib.NewTcpingClient("google.com", nil, nil, nil).EnableOutput()
	client.Run()
}
