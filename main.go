package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/internal/command"
	"github.com/vuon9/d2m/pkg/esporthub"
)

func main() {

	app := &cli.App{
		Name:  "d2m",
		Usage: "Esport matches schedule on terminal",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello, welcome to d2m!")
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			{
				Name:        "matches",
				Aliases:     []string{"m"},
				Usage:       "Options for matches",
				HelpName:    "d2m matches",
				Subcommands: matchesSubCommands(),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func matchesSubCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:    "dota2",
			Usage:   "Get matches of Dota 2",
			Aliases: []string{"d2", "dota"},
			Action: func(c *cli.Context) error {
				return command.GetCLIMatches(c.Context, esporthub.Dota2)
			},
		},
		{
			Name:  "csgo",
			Usage: "Get matches of CS:GO",
			Action: func(c *cli.Context) error {
				return command.GetCLIMatches(c.Context, esporthub.CsGO)
			},
		},
		{
			Name:    "leagueoflegends",
			Aliases: []string{"lol", "league"},
			Usage:   "Get matches of League of Legends",
			Action: func(c *cli.Context) error {
				return command.GetCLIMatches(c.Context, esporthub.LeagueOfLegends)
			},
		},
	}
}
