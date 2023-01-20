package d2m

import (
	cli "github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/d2m"
	"github.com/vuon9/d2m/pkg/api/types"
)

func AppCmd() *cli.App {
	app := d2m.NewMatcher()

	return &cli.App{
		Name:  "d2m",
		Usage: "Dota2 matches schedule on terminal",
		Action: func(c *cli.Context) error {
			return app.GetCLIMatches(c.Context)
		},
		Commands: []*cli.Command{
			{
				Name:     "live",
				Aliases:  []string{"l"},
				Usage:    "Live matches",
				HelpName: "d2m live",
				Action: func(c *cli.Context) error {
					return app.GetCLIMatches(c.Context, d2m.WithMatchStatus(types.MatchStatusLive))
				},
			},
			{
				Name:     "coming",
				Aliases:  []string{"u"},
				Usage:    "Upcoming matches",
				HelpName: "d2m coming",
				Action: func(c *cli.Context) error {
					return app.GetCLIMatches(c.Context, d2m.WithMatchStatus(types.MatchStatusComing))
				},
			},
			{
				Name:     "finished",
				Aliases:  []string{"f"},
				Usage:    "Finished matches",
				HelpName: "d2m finished",
				Action: func(c *cli.Context) error {
					return app.GetCLIMatches(c.Context, d2m.WithMatchStatus(types.MatchStatusFinished))
				},
			},
		},
	}
}
