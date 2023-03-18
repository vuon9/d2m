package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	iapp "github.com/vuon9/d2m/internal/app"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
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
				Name: "test",
				Subcommands: []*cli.Command{
					{
						Name: "list",
						Action: func(cCtx *cli.Context) error {
							client := liquipedia.NewClient()
							matches, err := client.GetScheduledMatches(cCtx.Context)
							if err != nil {
								return err
							}

							for i, m := range matches {
								log.Printf("%d. %s vs %s", i+1, m.Team1().TeamProfileLink, m.Team2().TeamProfileLink)
							}

							return nil
						},
					},
					{
						Name: "details",
						Action: func(cCtx *cli.Context) error {
							client := liquipedia.NewClient()
							team, err := client.GetTeamDetailsPage(cCtx.Context, "https://liquipedia.net/dota2/OG")
							if err != nil {
								return err
							}

							bTeam, err := json.MarshalIndent(team, "", "  ")
							if err != nil {
								return err
							}

							fmt.Println(string(bTeam))

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
