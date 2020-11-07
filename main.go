package main

import (
	"log"

	"github.com/rami3l/ting/cmd"
)

func main() {
	app := cmd.App()
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
