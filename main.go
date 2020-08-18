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
	}
	app.Action = func(ctx *cli.Context) error {
		srv, err := NewServer(ctx.String("listen"))
		if err != nil {
			return err
		}

		return srv.Run()
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
