package d2m

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api/types"
)

type Matcher struct {
	Keyword    string
	FilterType string
}

func NewMatcher() *Matcher {
	return &Matcher{}
}

type MatcherOption func(*Matcher)

func WithMatchStatus(status types.MatchStatus) MatcherOption {
	return func(matcher *Matcher) {
		matcher.Keyword = status.String()
		matcher.FilterType = "status"
	}
}

func WithDate(date time.Time) MatcherOption {
	return func(matcher *Matcher) {
		matcher.Keyword = date.Truncate(24 * time.Hour).Format("2006-01-02")
		matcher.FilterType = "date"
	}
}

// GetCLIMatches prints matches as table on terminal
func (m *Matcher) GetCLIMatches(ctx context.Context, options ...MatcherOption) error {
	for _, mo := range options {
		mo(m)
	}

	prev5Hours := time.Now().Add(-24 * time.Hour)
	matches, err := GetMatches(ctx, types.Dota2)
	if err != nil {
		return err
	}

	finalMatches := make(types.MatchSlice, 0)

	for _, match := range matches {
		matchStatus := match.FriendlyStatus()
		matchStartByDate := match.Start.Truncate(24 * time.Hour)
		matchStartByDateIsToday := matchStartByDate.Equal(time.Now().Truncate(24 * time.Hour))

		if m.FilterType == "status" {
			if m.Keyword == "today" && !matchStartByDateIsToday {
				continue
			} else if matchStatus != types.MatchStatus(m.Keyword) {
				continue
			}
		}

		if m.FilterType == "date" && matchStartByDate.Format("2006-01-02") != m.Keyword {
			continue
		}

		if match.Start.Before(prev5Hours) {
			continue
		}

		finalMatches = append(finalMatches, match)
	}

	prog := tea.NewProgram(newModel(matches))
	if _, err := prog.Run(); err != nil {
		return err
	}

	return nil
}

type delegateKeyMap struct {
	choose key.Binding
}

func (m delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.choose},
	}
}

func (m delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.choose,
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
	}
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(*item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) { //nolint:gocritic
		case tea.KeyMsg:
			switch { //nolint:gocritic
				case key.Matches(msg, keys.choose):
					return m.NewStatusMessage("You chose " + title)
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}
