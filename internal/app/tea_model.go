package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyName string

const (
	// Filter keys
	AllMatches       = keyName("all")
	FromTodayMatches = keyName("from_today")
	TodayMatches     = keyName("today")
	TomorrowMatches  = keyName("tomorrow")
	YesterdayMatches = keyName("yesterday")
	LiveMatches      = keyName("live")
	FinishedMatches  = keyName("finished")
	ComingMatches    = keyName("coming")

	// Delegate keys
	ChooseMatch   = keyName("choose")
	OpenStreamURL = keyName("open_stream_url")
)

type keyMap struct {
	AllMatches       key.Binding
	FromTodayMatches key.Binding
	TodayMatches     key.Binding
	TomorrowMatches  key.Binding
	YesterdayMatches key.Binding
	LiveMatches      key.Binding
	FinishedMatches  key.Binding
	ComingMatches    key.Binding
}

func (km keyMap) FullHelp() []key.Binding {
	return []key.Binding{
		km.AllMatches,
		km.FromTodayMatches,
		km.TodayMatches,
		km.TomorrowMatches,
		km.YesterdayMatches,
		km.LiveMatches,
		km.FinishedMatches,
		km.ComingMatches,
	}
}

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	filterKeys = keyMap{
		AllMatches:       key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "all")),
		FromTodayMatches: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "from today")),
		TodayMatches:     key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "today")),
		TomorrowMatches:  key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "tomorrow")),
		YesterdayMatches: key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yesterday")),
		LiveMatches:      key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "live")),
		FinishedMatches:  key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "finished")),
		ComingMatches:    key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "coming")),
	}

	exitKeys = map[string]bool{
		"q":      true,
		"ctrl+c": true,
	}
)

type model struct {
	delegate list.DefaultDelegate
	keys     keyMap
	list     list.Model
	items    []list.Item
}

func newModel(matches []list.Item) tea.Model {
	delegate := newItemDelegate(newDelegateKeyMap())
	items := filterMatches(matches, FromToday)

	matchList := list.New(items, delegate, 0, 0)
	matchList.Title = "D2M - Dota2 Matches Tracker"
	matchList.Styles.Title = titleStyle

	return &model{
		list:     matchList,
		delegate: delegate,
		items:    matches,
		keys:     filterKeys,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok {
		if exitKeys[msg.String()] {
			return m, tea.Quit
		}
	}

	switch msg := msg.(type) { //nolint:gocritic
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.AllMatches):
			m.list.SetItems(filterMatches(m.items, All))
		case key.Matches(msg, m.keys.FromTodayMatches):
			m.list.SetItems(filterMatches(m.items, FromToday))
		case key.Matches(msg, m.keys.TomorrowMatches):
			m.list.SetItems(filterMatches(m.items, Tomorrow))
		case key.Matches(msg, m.keys.YesterdayMatches):
			m.list.SetItems(filterMatches(m.items, Yesterday))
		case key.Matches(msg, m.keys.LiveMatches):
			m.list.SetItems(filterMatches(m.items, Live))
		case key.Matches(msg, m.keys.ComingMatches):
			m.list.SetItems(filterMatches(m.items, Coming))
		case key.Matches(msg, m.keys.FinishedMatches):
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
