package main

import (
	"os"

	"github.com/rami3l/ting/cmd"
)

func main() {
	app := cmd.App()
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
