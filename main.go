package main

import (
	"log"
	"os"

	"github.com/rami3l/ting/cmd"
)

func main() {
	app := cmd.App()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
