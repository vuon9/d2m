package d2m

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api/types"
)

// GetCLIMatches prints matches as table on terminal
func GetCLIMatches(ctx context.Context) error {
	matches, err := GetMatches(ctx, types.Dota2)
	if err != nil {
		return err
	}

	prog := tea.NewProgram(newModel(matches))
	if _, err := prog.Run(); err != nil {
		return err
	}

	return nil
}