package cmd

import (
	"fmt"
	"time"

	"github.com/rami3l/ting/lib"
	"github.com/spf13/cobra"
)

func App() (app *cobra.Command) {
	app = &cobra.Command{
		Use:   "ting [hosts...]",
		Short: "Ping over a tcp connection",
	}

	app.Args = cobra.MinimumNArgs(1)

	port := app.Flags().IntP("port", "p", 80, "Numeric TCP port")
	tryCount := app.Flags().IntP("count", "n", 5, "Number of tries")
	tryInterval := app.Flags().Float32P("interval", "i", 1, "Interval between pings, in seconds")
	timeout := app.Flags().Float32P("timeout", "w", 5, "Maximum time to wait for a response, in seconds")

	app.RunE = func(_ *cobra.Command, args []string) (err error) {
		for _, host := range args {
			tryInterval, err := time.ParseDuration(fmt.Sprintf("%fs", *tryInterval))
			if err != nil {
				return err
			}
			timeout, err := time.ParseDuration(fmt.Sprintf("%fs", *timeout))
			if err != nil {
				return err
			}
			client := lib.NewTcpingClient(host).
				SetPort(*port).
				SetTryCount(*tryCount).
				SetTryInterval(tryInterval).
				SetTimeout(timeout).
				EnableOutput()
			if _, err := client.Run(); err != nil {
				return err
			}
		}

		return
	}

	return
}
