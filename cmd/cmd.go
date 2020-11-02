package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Main() {
	fmt.Printf("Hello, ting!")
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
