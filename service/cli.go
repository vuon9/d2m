package service

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/service/viewmodels"
)

func NewCLIApp() *cli.App {
	return &cli.App{
		Name: "d2m",
		Action: func(*cli.Context) error {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				return err
			}

			defer f.Close()

			prog := tea.NewProgram(viewmodels.NewMatchList(), tea.WithAltScreen())
			if _, err := prog.Run(); err != nil {
				return err
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "from-testdata",
				// Subcommands: GetSubCommands(),
			},
		},
	}
}
