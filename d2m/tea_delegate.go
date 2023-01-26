package d2m

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api"
)

type delegateKeyMap struct {
	choose key.Binding
	openStreamURL key.Binding
}

func (m delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.choose},
		{m.openStreamURL},
	}
}

func (m delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.choose,
		m.openStreamURL,
	}
}

var (
	itemKeys = keyMaps{
		ChooseMatch: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
		OpenStreamURL: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url")),
	}
)

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
		openStreamURL: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url")),
	}
}

type ListItem interface{
	Title() string
	Description() string
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	de := list.NewDefaultDelegate()

	de.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		match, ok := m.SelectedItem().(*api.Match)
		if !ok {
			return m.NewStatusMessage("Invalid match")
		}

		title := match.GeneralTitle()

		switch msg := msg.(type) { //nolint:gocritic
		case tea.KeyMsg:
			switch { //nolint:gocritic
				case key.Matches(msg, keys.choose):
					return m.NewStatusMessage(fmt.Sprintf("Current match is '%s'", title))
				case key.Matches(msg, keys.openStreamURL):
					if match.StreamingURL == "" {
						return m.NewStatusMessage(fmt.Sprintf("Match '%s' is not live", title))
					}

					go func() {
						OpenURL(match.StreamingURL)
					}()
					return m.NewStatusMessage(fmt.Sprintf("Opening stream URL for '%s'", title))
			}
		}

		return nil
	}

	// Merge filterKeys and itemKeys into []key.Binding slice as help vars
	help := make([]key.Binding, 0)
	for _, filterKey := range filterKeys {
		help = append(help, filterKey)
	}

	for _, item := range itemKeys {
		help = append(help, item)
	}

	de.ShortHelpFunc = func() []key.Binding {
		return help
	}

	de.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return de
}
