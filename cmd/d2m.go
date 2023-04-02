package cmd

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	iapp "github.com/vuon9/d2m/internal/app"
)

func Execute() {
	app := &cli.App{
		Name: "d2m",
		Action: func(*cli.Context) error {
			prog := iapp.NewApp()
			return prog.Run(context.Background())
		},
		Commands: []*cli.Command{
			{
				Name:        "test",
				Subcommands: iapp.GetSubCommands(),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
