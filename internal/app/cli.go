package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/pkg/api"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

var apiClient api.Clienter = liquipedia.NewClient()

// RunProgram prints matches as table on terminal
func (a *App) Run(ctx context.Context) error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return err
	}

	defer f.Close()

	prog := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		return err
	}

	return nil
}

type commander interface {
	GetName() string
	GetAction() func(*cli.Context) error
}

func GetSubCommands() []*cli.Command {
	cmds := []commander{
		&listCommand{},
		&detailsCommand{},
	}

	var subCommands []*cli.Command
	for _, cmd := range cmds {
		subCommands = append(subCommands, &cli.Command{
			Name:   cmd.GetName(),
			Action: cmd.GetAction(),
		})
	}

	return subCommands
}

type listCommand struct {
}

func (c *listCommand) GetName() string {
	return "list"
}

func (c *listCommand) GetAction() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		client := liquipedia.NewClient()
		matches, err := client.GetScheduledMatches(ctx.Context)
		if err != nil {
			return err
		}

		for i, m := range matches {
			log.Printf("%d. %s vs %s", i+1, m.Team1().TeamProfileLink, m.Team2().TeamProfileLink)
		}

		return nil
	}
}

type detailsCommand struct {
}

func (c *detailsCommand) GetName() string {
	return "details"
}

func (c *detailsCommand) GetAction() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		client := liquipedia.NewClient()
		team, err := client.GetTeamDetailsPage(ctx.Context, "https://liquipedia.net/dota2/OG")
		if err != nil {
			return err
		}

		bTeam, err := json.MarshalIndent(team, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(bTeam))

		return nil
	}
}
