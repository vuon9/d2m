package d2m

import (
	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/d2m"
	"github.com/vuon9/d2m/pkg/api/types"
)

func AppCmd() *cli.App {
	return &cli.App{
		Name:  "d2m",
		Usage: "Dota2 matches schedule on terminal",
		Action: func(c *cli.Context) error {
			return d2m.GetCLIMatches(c.Context, types.Dota2)
		},
	}
}
