package d2m

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type delegateKeyMap struct {
	choose key.Binding
	all key.Binding
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
		{m.all},
		{m.today, m.tomorrow, m.yesterday},
		{m.live, m.finished, m.coming},
	}
}

func (m delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		m.choose,
		m.all,
		m.today,
		m.tomorrow,
		m.yesterday,
		m.live,
		m.finished,
		m.coming,
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
		all: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "all")),
		today: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "today")),
		tomorrow: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "tomorrow")),
		yesterday: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yesterday")),
		live: key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "live")),
		finished: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "finished")),
		coming: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "coming")),
	}
}

type delegator struct {
	originItems []list.Item
}

func (d *delegator) newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	de := list.NewDefaultDelegate()

	de.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
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
				case key.Matches(msg, keys.all):
					return m.SetItems(d.originItems)
				case key.Matches(msg, keys.today):
					return m.SetItems(d.filterMatches(today))
				case key.Matches(msg, keys.tomorrow):
					return m.SetItems(d.filterMatches(tomorrow))
				case key.Matches(msg, keys.yesterday):
					return m.SetItems(d.filterMatches(yesterday))
				case key.Matches(msg, keys.live):
					return m.SetItems(d.filterMatches(live))
				case key.Matches(msg, keys.finished):
					return m.SetItems(d.filterMatches(finished))
				case key.Matches(msg, keys.coming):
					return m.SetItems(d.filterMatches(coming))
			}
		}

		return nil
	}

	help := []key.Binding{
		keys.choose,
		keys.all,
		keys.today,
		keys.tomorrow,
		keys.yesterday,
		keys.live,
		keys.finished,
		keys.coming,
	}

	de.ShortHelpFunc = func() []key.Binding {
		return help
	}

	de.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return de
}