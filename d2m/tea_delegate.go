package d2m

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api"
)

type ExtendedBubbleItem interface {
	list.Item
	GeneralTitle() string
	Title() string
	Description () string
}

var (
	// Safety check to make sure our type implements the interface.
	_ ExtendedBubbleItem = (*api.Match)(nil)
)

type delegateKeyMap struct {
	choose key.Binding
	all key.Binding
	fromToday key.Binding
	today key.Binding
	tomorrow key.Binding
	yesterday key.Binding
	live key.Binding
	finished key.Binding
	coming key.Binding
}

func (m delegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.choose},
		// {m.all, m.fromToday},
		// {m.today, m.tomorrow, m.yesterday},
		// {m.live, m.finished, m.coming},
	}
}

func (m delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.choose,
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
		all:       filterKeys[All],
		fromToday: filterKeys[FromToday],
		today:     filterKeys[Today],
		tomorrow:  filterKeys[Tomorrow],
		yesterday: filterKeys[Yesterday],
		live:      filterKeys[Live],
		finished:  filterKeys[Finished],
		coming:    filterKeys[Coming],
	}
}

type ListItem interface{
	Title() string
	Description() string
}

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	de := list.NewDefaultDelegate()

	de.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(ExtendedBubbleItem); ok {
			title = i.GeneralTitle()
		} else {
			return nil
		}

		switch msg := msg.(type) { //nolint:gocritic
		case tea.KeyMsg:
			switch { //nolint:gocritic
				case key.Matches(msg, keys.choose):
					return m.NewStatusMessage(fmt.Sprintf("Current match is '%s'", title))
			}
		}

		return nil
	}

	help := []key.Binding{
		keys.choose,
		keys.live,
		keys.fromToday,
		keys.coming,
		keys.finished,
		keys.all,
		keys.today,
		keys.tomorrow,
		keys.yesterday,
	}

	de.ShortHelpFunc = func() []key.Binding {
		return help
	}

	de.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return de
}
