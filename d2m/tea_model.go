package d2m

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	filterKeys = map[matchFilter]key.Binding{
		All: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "all")),
		FromToday: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "from today")),
		Today: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "today")),
		Tomorrow: key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "tomorrow")),
		Yesterday: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yesterday")),
		Live: key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "live")),
		Finished: key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "finished")),
		Coming: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "coming")),
	}

	exitKeys = map[string]bool{
		"q": true,
		"esc": true,
		"ctrl+c": true,
	}
)

type model struct {
	delegate list.DefaultDelegate
	keys     map[matchFilter]key.Binding
	list     list.Model
	items    []list.Item
}

func newModel(matches []list.Item) tea.Model {
	delegate := newItemDelegate(newDelegateKeyMap())
	items := filterMatches(matches, FromToday)

	matchList := list.New(items, delegate, 0, 0)
	matchList.Title = "D2M - Dota2 Matches"
	matchList.Styles.Title = titleStyle

	return &model{
		list: matchList,
		delegate: delegate,
		items: matches,
		keys: filterKeys,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch {
		case exitKeys[msg.String()]:
			return m, tea.Quit
		case key.Matches(msg, m.keys[All]):
			m.list.SetItems(filterMatches(m.items, All))
		case key.Matches(msg, m.keys[FromToday]):
			m.list.SetItems(filterMatches(m.items, FromToday))
		case key.Matches(msg, m.keys[Tomorrow]):
			m.list.SetItems(filterMatches(m.items, Tomorrow))
		case key.Matches(msg, m.keys[Yesterday]):
			m.list.SetItems(filterMatches(m.items, Yesterday))
		case key.Matches(msg, m.keys[Live]):
			m.list.SetItems(filterMatches(m.items, Live))
		case key.Matches(msg, m.keys[Coming]):
			m.list.SetItems(filterMatches(m.items, Coming))
		case key.Matches(msg, m.keys[Finished]):
			m.list.SetItems(filterMatches(m.items, Finished))
		}
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return appStyle.Render(m.list.View())
}
