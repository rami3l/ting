package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/rami3l/ting/lib"
	"github.com/urfave/cli/v2"
)

func App() (app *cli.App) {
	app = &cli.App{
		Name:  "ting",
		Usage: "Ping over a tcp connection",
	}

	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "TCP port `INT`",
			Value:   80,
		},

		&cli.Float64Flag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "Wait `DEC` second(s) between pings",
			Value:   1,
		},

		&cli.IntFlag{
			Name:    "count",
			Aliases: []string{"n"},
			Usage:   "Try `INT` times",
			Value:   5,
		},

		&cli.Float64Flag{
			Name:    "timeout",
			Aliases: []string{"w"},
			Usage:   "Wait `DEC` second(s) for a response",
			Value:   5,
		},
	}

	app.Action = func(c *cli.Context) (err error) {
		if !c.Args().Present() {
			fmt.Printf("Error: Missing HOST")
			return
		}

		for _, host := range c.Args().Slice() {
			port := c.Int("port")
			tryCount := c.Int("count")
			tryInterval, err := time.ParseDuration(fmt.Sprintf("%fs", c.Float64("interval")))
			if err != nil {
				return err
			}
			timeout, err := time.ParseDuration(fmt.Sprintf("%fs", c.Float64("timeout")))
			if err != nil {
				return err
			}
			client := lib.NewTcpingClient(host).
				SetPort(port).
				SetTryCount(tryCount).
				SetTryInterval(tryInterval).
				SetTimeout(timeout).
				EnableOutput()
			if _, err := client.Run(); err != nil {
				return err
			}
		}

		return
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	return
}
