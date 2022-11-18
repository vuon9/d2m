package d2m

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func AppCmd() *cli.App {
	return &cli.App{
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
				Subcommands: MatchesCmds(),
			},
		},
	}
}
