package app

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type app struct {
	tracker Tracker
}

type Apper interface {
	Run(ctx context.Context) error
}

func NewApp() Apper {
	return &app{
		tracker: NewTracker(),
	}
}

// RunProgram prints matches as table on terminal
func (a *app) Run(ctx context.Context) error {
	matches, err := a.tracker.GetMatches(ctx)
	if err != nil {
		return err
	}

	// for _, match := range matches {
	// 	fmt.Println(match.Team1().FullName, match.Team2().FullName, match.Start)
	// 	fmt.Println(match.Team1().TeamProfileLink, match.Team2().TeamProfileLink)
	// 	fmt.Println()
	// }

	items := make([]list.Item, 0)
	for _, match := range matches {
		items = append(items, match)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return err
	}

	defer f.Close()

	prog := tea.NewProgram(newModel(items), tea.WithAltScreen())
	if _, err := prog.Run(); err != nil {
		return err
	}

	return nil
}
