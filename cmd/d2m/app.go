package d2m

import (
	"time"

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
			today := time.Now().Truncate(24 * time.Hour)
			return app.GetCLIMatches(c.Context, d2m.WithDate(today))
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
			{
				Name:     "today",
				Aliases:  []string{"t"},
				Usage:    "Matches today",
				HelpName: "d2m today",
				Action: func(c *cli.Context) error {
					today := time.Now().Truncate(24 * time.Hour)
					return app.GetCLIMatches(c.Context, d2m.WithDate(today))
				},
			},
			{
				Name:     "tomorrow",
				Aliases:  []string{"m"},
				Usage:    "Matches tomorrow",
				HelpName: "d2m tomorrow",
				Action: func(c *cli.Context) error {
					tomorrow := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
					return app.GetCLIMatches(c.Context, d2m.WithDate(tomorrow))
				},
			},
			{
				Name:     "yesterday",
				Aliases:  []string{"y"},
				Usage:    "Matches yesterday",
				HelpName: "d2m yesterday",
				Action: func(c *cli.Context) error {
					yesterday := time.Now().Truncate(24 * time.Hour).Add(-24 * time.Hour)
					return app.GetCLIMatches(c.Context, d2m.WithDate(yesterday))
				},
			},
		},
	}
}
