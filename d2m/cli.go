package d2m

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

	items := make([]list.Item, 0)
	for _, match := range matches {
		items = append(items, match)
	}

	prog := tea.NewProgram(newModel(items))
	if _, err := prog.Run(); err != nil {
		return err
	}

	return nil
}