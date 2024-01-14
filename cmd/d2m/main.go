package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m"
)

func main() {
	app := &cli.App{
		Name: "d2m",
		Action: func(*cli.Context) error {
			prog := d2m.NewApp()
			return prog.Run(context.Background())
		},
		Commands: []*cli.Command{
			{
				Name:        "test",
				Subcommands: d2m.GetSubCommands(),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
