package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api"
)

type (
	model struct {
		keys         keyMap
		listModel    list.Model
		detailsModel table.Model
		items        []list.Item
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

var (
	showListMatch    = true
	showDetailsMatch = false
)

func newModel(matches []list.Item) tea.Model {
	return &model{
		listModel:    newListView(matches),
		detailsModel: newDetailsView(),
		items:        matches,
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

			showDetailsMatch = !showDetailsMatch
			showListMatch = !showListMatch
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
		case key.Matches(msg, filterKeys.AllMatches):
			m.listModel.SetItems(filterMatches(m.items, All))
		case key.Matches(msg, filterKeys.FromTodayMatches):
			m.listModel.SetItems(filterMatches(m.items, FromToday))
		case key.Matches(msg, filterKeys.TomorrowMatches):
			m.listModel.SetItems(filterMatches(m.items, Tomorrow))
		case key.Matches(msg, filterKeys.YesterdayMatches):
			m.listModel.SetItems(filterMatches(m.items, Yesterday))
		case key.Matches(msg, filterKeys.LiveMatches):
			m.listModel.SetItems(filterMatches(m.items, Live))
		case key.Matches(msg, filterKeys.ComingMatches):
			m.listModel.SetItems(filterMatches(m.items, Coming))
		case key.Matches(msg, filterKeys.FinishedMatches):
			m.listModel.SetItems(filterMatches(m.items, Finished))
		}
	}

	var cmd tea.Cmd
	if showListMatch {
		m.listModel, cmd = m.listModel.Update(msg)
	}

	if showDetailsMatch {
		m.detailsModel, cmd = m.detailsModel.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func newListView(matches []list.Item) list.Model {
	listView := list.New(filterMatches(matches, FromToday), list.NewDefaultDelegate(), 0, 0)
	listView.Title = "D2M - Dota2 Matches Tracker"
	listView.Styles.Title = titleStyle

	return listView
}

func newDetailsView() table.Model {
	columns := []table.Column{
		{Title: "Player", Width: 10},
		{Title: "Hero", Width: 10},
		{Title: "Team", Width: 10},
	}

	rows := []table.Row{
		{"player1", "hero1", "Liquid"},
		{"player2", "hero2", "OG"},
	}

	tableView := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	tableView.SetStyles(s)

	return tableView
}

func (m *model) View() string {
	view := m.listModel.View()
	if showDetailsMatch {
		view = m.detailsModel.View()
	}

	return appStyle.Render(view)
}
