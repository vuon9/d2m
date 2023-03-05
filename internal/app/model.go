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
		listModel    list.Model
		detailsModel table.Model
		items        []list.Item
		appState     appState
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

	listKeys = keyMap{
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

type appState uint8

const (
	showListMatch appState = iota
	showDetailsMatch
)

func newModel(matches []list.Item) tea.Model {
	return &model{
		listModel:    newListView(matches),
		detailsModel: newDetailsView(),
		items:        matches,
		appState:     showListMatch,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

type matchFilterKey struct {
	Key    key.Binding
	Filter matchFilter
}

type matchFilterKeys []matchFilterKey

var filterKeys = matchFilterKeys{
	{
		Key:    listKeys.AllMatches,
		Filter: All,
	},
	{
		Key:    listKeys.FromTodayMatches,
		Filter: FromToday,
	},
	{
		Key:    listKeys.TodayMatches,
		Filter: Today,
	},
	{
		Key:    listKeys.TomorrowMatches,
		Filter: Tomorrow,
	},
	{
		Key:    listKeys.YesterdayMatches,
		Filter: Yesterday,
	},
	{
		Key:    listKeys.LiveMatches,
		Filter: Live,
	},
	{
		Key:    listKeys.FinishedMatches,
		Filter: Finished,
	},
	{
		Key:    listKeys.ComingMatches,
		Filter: Coming,
	},
}

func (fk matchFilterKeys) Match(msg tea.KeyMsg) (found bool, filter matchFilter) {
	for _, k := range fk {
		if key.Matches(msg, k.Key) {
			return true, k.Filter
		}
	}

	return false, All
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.listModel.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		// common keys, work in all states
		switch { //nolint:gocritic
		case key.Matches(msg, listKeys.OpenStreamURL):
			m.openStreamingURL()
		}

		// keys by view states
		switch m.appState {
		case showListMatch:
			var cmd tea.Cmd
			m.listModel, cmd = m.listModel.Update(msg)

			switch {
			case msg.String() == "enter":
				match, ok := m.listModel.SelectedItem().(*api.Match)
				if ok {
					m.listModel.NewStatusMessage(fmt.Sprintf("Choose match %s", match.GeneralTitle()))
				}

				m.appState = showDetailsMatch
			default:
				found, mf := filterKeys.Match(msg)
				if found {
					m.listModel.SetItems(filterMatches(m.items, mf))
				}
			}

			cmds = append(cmds, cmd)

		case showDetailsMatch:
			var cmd tea.Cmd
			m.detailsModel, cmd = m.detailsModel.Update(msg)

			if msg.String() == "esc" {
				m.appState = showListMatch
			}

			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) openStreamingURL() {
	match, ok := m.listModel.SelectedItem().(*api.Match)
	if !ok || match.StreamingURL == "" {
		m.listModel.NewStatusMessage("No stream URL available")
		return
	}

	go func() {
		OpenURL(match.StreamingURL)
	}()

	m.listModel.NewStatusMessage(fmt.Sprintf("Opening stream URL for '%s'", match.GeneralTitle()))
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

func (m model) View() string {
	view := m.listModel.View()
	if m.appState == showDetailsMatch {
		view = m.detailsModel.View()
	}

	return appStyle.Render(view)
}
