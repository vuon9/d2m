package internal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"github.com/vuon9/d2m/internal/viewmodel"
)

func NewCLI() *cli.App {
	return &cli.App{
		Name: "d2m",
		Action: func(*cli.Context) error {
			// f, err := tea.LogToFile("resources/logs/debug.log", "debug")
			// if err != nil {
			// 	return err
			// }

			// defer f.Close()

			prog := tea.NewProgram(viewmodel.NewMainScreen(), tea.WithAltScreen())
			if _, err := prog.Run(); err != nil {
				return err
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "from-testdata",
				Action: func(*cli.Context) error {
					viewmodel.UseTestAPIClient()

					prog := tea.NewProgram(viewmodel.NewMainScreen(), tea.WithAltScreen())
					if _, err := prog.Run(); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}
}
