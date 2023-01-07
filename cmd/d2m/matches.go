package d2m

import (
	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/d2m"
	"github.com/vuon9/d2m/pkg/api/types"
)

type matchCmdInfo struct {
	cmdName  string
	aliases  []string
	usage    string
	gameName types.GameName
}

var (
	matchesSubCmds = []matchCmdInfo{
		{
			cmdName:  "dota2",
			aliases:  []string{"d2", "dota"},
			usage:    "Dota 2 matches",
			gameName: types.Dota2,
		},
	}
)

// MatchesCmds returns a slice of cli.Command for matches subcommands
func MatchesCmds() []*cli.Command {
	cliActionFn := func(gName types.GameName) cli.ActionFunc {
		return func(c *cli.Context) error {
			return d2m.GetCLIMatches(c.Context, gName)
		}
	}

	cmds := make([]*cli.Command, len(matchesSubCmds))
	for i, game := range matchesSubCmds {
		cmds[i] = &cli.Command{
			Name:     game.cmdName,
			Aliases:  game.aliases,
			Usage:    game.usage,
			HelpName: "d2m matches: " + game.cmdName,
			Action:   cliActionFn(game.gameName),
		}
	}

	return cmds
}
