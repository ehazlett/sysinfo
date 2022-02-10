package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Usage = "system info over http"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "listen",
			Aliases: []string{"l"},
			Usage:   "listen address",
			Value:   ":8080",
		},
		&cli.StringFlag{
			Name:    "trace-endpoint",
			Aliases: []string{"t"},
			Usage:   "send opentelemetry traces to this endpoint",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "app-version",
			EnvVars: []string{"APP_VERSION"},
			Hidden:  true,
			Value:   "latest",
		},
	}
	app.Action = func(clix *cli.Context) error {
		srv, err := NewServer(clix.String("listen"), clix.String("app-version"), clix.String("trace-endpoint"))
		if err != nil {
			return err
		}

		return srv.Run()
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
