package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api"
)

type MatchItem interface {
	list.Item
	list.DefaultItem
}

// MatchItem is an interface that all match items must implement.
var _ MatchItem = (*api.Match)(nil)

type delegateKeyMap struct {
	choose        key.Binding
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

type itemKeyMap struct {
	ChooseMatch   key.Binding
	OpenStreamURL key.Binding
}

func (m itemKeyMap) FullHelp() []key.Binding {
	return []key.Binding{
		m.ChooseMatch,
		m.OpenStreamURL,
	}
}

var itemKeys = itemKeyMap{
	ChooseMatch:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
	OpenStreamURL: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url")),
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose:        key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
		openStreamURL: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url")),
	}
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
	help := append(filterKeys.FullHelp(), itemKeys.FullHelp()...)

	de.ShortHelpFunc = func() []key.Binding {
		return help
	}

	de.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return de
}
