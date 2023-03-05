package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api"
)

type (
	model struct {
		keys            keyMap
		listModel       list.Model
		delegate        list.DefaultDelegate
		items           []list.Item
		isInDetailsMode bool
	}

	MatchItem interface {
		list.Item
		list.DefaultItem
	}
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
	OpenStreamURL    key.Binding
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

type delegateKeyMap struct {
	choose        key.Binding
	openStreamURL key.Binding
}

func (d *delegateKeyMap) FullHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.openStreamURL,
	}
}

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

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 2)

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
		OpenStreamURL:    key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url")),
	}

	exitKeys = map[string]bool{
		"q":      true,
		"ctrl+c": true,
	}

	listKeyMap = &delegateKeyMap{
		choose: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "choose")),
	}
)

func newModel(matches []list.Item) tea.Model {
	return &model{
		listModel: newListView(matches),
		items:     matches,
		keys:      filterKeys,
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
		m.listModel.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch {
		case msg.String() == "enter":
			match, ok := m.listModel.SelectedItem().(*api.Match)
			if ok {
				m.listModel.NewStatusMessage(fmt.Sprintf("Choose match %s", match.GeneralTitle()))
			}

			m.isInDetailsMode = !m.isInDetailsMode
		case key.Matches(msg, m.keys.OpenStreamURL):
			match, ok := m.listModel.SelectedItem().(*api.Match)
			if !ok || match.StreamingURL == "" {
				m.listModel.NewStatusMessage("No stream URL available")
				break
			}

			go func() {
				OpenURL(match.StreamingURL)
			}()

			m.listModel.NewStatusMessage(fmt.Sprintf("Opening stream URL for '%s'", match.GeneralTitle()))
		case key.Matches(msg, m.keys.AllMatches):
			m.listModel.SetItems(filterMatches(m.items, All))
		case key.Matches(msg, m.keys.FromTodayMatches):
			m.listModel.SetItems(filterMatches(m.items, FromToday))
		case key.Matches(msg, m.keys.TomorrowMatches):
			m.listModel.SetItems(filterMatches(m.items, Tomorrow))
		case key.Matches(msg, m.keys.YesterdayMatches):
			m.listModel.SetItems(filterMatches(m.items, Yesterday))
		case key.Matches(msg, m.keys.LiveMatches):
			m.listModel.SetItems(filterMatches(m.items, Live))
		case key.Matches(msg, m.keys.ComingMatches):
			m.listModel.SetItems(filterMatches(m.items, Coming))
		case key.Matches(msg, m.keys.FinishedMatches):
			m.listModel.SetItems(filterMatches(m.items, Finished))
		}
	}

	newListModel, cmd := m.listModel.Update(msg)
	m.listModel = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func newListView(matches []list.Item) list.Model {
	listView := list.New(filterMatches(matches, FromToday), list.NewDefaultDelegate(), 0, 0)
	listView.Title = "D2M - Dota2 Matches Tracker"
	listView.Styles.Title = titleStyle

	return listView
}

func (m *model) View() string {
	var view string
	if m.isInDetailsMode {
		view = m.detailsView()
	} else {
		view = m.listModel.View()
	}

	return appStyle.Render(view)
}

func (m *model) detailsView() string {
	si := m.listModel.SelectedItem().(*api.Match)
	return "this is details view: " + si.GeneralTitle()
}
