package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
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
