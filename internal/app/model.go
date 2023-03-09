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
	KeyAllMatches       key.Binding
	KeyFromTodayMatches key.Binding
	KeyTodayMatches     key.Binding
	KeyTomorrowMatches  key.Binding
	KeyYesterdayMatches key.Binding
	KeyLiveMatches      key.Binding
	KeyFinishedMatches  key.Binding
	KeyComingMatches    key.Binding
	KeyOpenStreamURL    key.Binding
}

// TODO: Missing help menu for these keys
func (km keyMap) FullHelp() []key.Binding {
	return []key.Binding{
		km.KeyAllMatches,
		km.KeyFromTodayMatches,
		km.KeyTodayMatches,
		km.KeyTomorrowMatches,
		km.KeyYesterdayMatches,
		km.KeyLiveMatches,
		km.KeyFinishedMatches,
		km.KeyComingMatches,
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

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 2)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	KeyAllMatches       = key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "all"))
	KeyFromTodayMatches = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "from today"))
	KeyTodayMatches     = key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "today"))
	KeyTomorrowMatches  = key.NewBinding(key.WithKeys("m"), key.WithHelp("m", "tomorrow"))
	KeyYesterdayMatches = key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "yesterday"))
	KeyLiveMatches      = key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "live"))
	KeyFinishedMatches  = key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "finished"))
	KeyComingMatches    = key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "coming"))
	KeyOpenStreamURL    = key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open stream url"))

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

type matchFilterKeys map[matchFilter]key.Binding

func (m matchFilterKeys) FullHelp() []key.Binding {
	var keys []key.Binding
	for _, k := range m {
		keys = append(keys, k)
	}

	return keys
}

func IsFilterKey(msg tea.KeyMsg) bool {
	for _, k := range filterKeys {
		if key.Matches(msg, k) {
			return true
		}
	}

	return false
}

var filterKeys = matchFilterKeys{
	All:       KeyAllMatches,
	FromToday: KeyFromTodayMatches,
	Today:     KeyTodayMatches,
	Tomorrow:  KeyTomorrowMatches,
	Yesterday: KeyYesterdayMatches,
	Live:      KeyLiveMatches,
	Finished:  KeyFinishedMatches,
	Coming:    KeyComingMatches,
}

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

// DoFilterSuccessful is used to filter matches by key
func (m *model) DoFilterSuccessful(msg tea.KeyMsg) bool {
	switch {
	case key.Matches(msg, KeyAllMatches):
		m.listModel.SetItems(filterMatches(m.items, All))
	case key.Matches(msg, KeyFromTodayMatches):
		m.listModel.SetItems(filterMatches(m.items, FromToday))
	case key.Matches(msg, KeyTodayMatches):
		m.listModel.SetItems(filterMatches(m.items, Today))
	case key.Matches(msg, KeyTomorrowMatches):
		m.listModel.SetItems(filterMatches(m.items, Tomorrow))
	case key.Matches(msg, KeyYesterdayMatches):
		m.listModel.SetItems(filterMatches(m.items, Yesterday))
	case key.Matches(msg, KeyLiveMatches):
		m.listModel.SetItems(filterMatches(m.items, Live))
	case key.Matches(msg, KeyFinishedMatches):
		m.listModel.SetItems(filterMatches(m.items, Finished))
	case key.Matches(msg, KeyComingMatches):
		m.listModel.SetItems(filterMatches(m.items, Coming))
	default:
		return false
	}

	return true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) { //nolint:gocritic
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.listModel.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		// If the list is filtering, we want to skip all the keys in this state
		if m.listModel.FilterState() == list.Filtering {
			break
		}

		// common keys, work in all states
		switch m.appState {
		case showListMatch:
			switch { //nolint:gocritic
			case exitKeys[msg.String()]:
				return m, tea.Quit
			case key.Matches(msg, KeyOpenStreamURL):
				m.openStreamingURL()
				return m, nil
			case msg.String() == "enter":
				match, ok := m.listModel.SelectedItem().(*api.Match)
				if ok {
					m.listModel.NewStatusMessage(fmt.Sprintf("Choose match %s", match.GeneralTitle()))
					m.appState = showDetailsMatch
				}

				return m, nil
			case m.DoFilterSuccessful(msg):
				return m, nil
			}
		case showDetailsMatch:
			switch {
			case exitKeys[msg.String()]:
				m.appState = showListMatch
				return m, nil
			default:
				var cmd tea.Cmd
				m.detailsModel, cmd = m.detailsModel.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	var cmd tea.Cmd
	m.listModel, cmd = m.listModel.Update(msg)
	cmds = append(cmds, cmd)

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

func (m model) View() string {
	view := m.listModel.View()
	if m.appState == showDetailsMatch {
		view = m.detailsModel.View()
	}

	return appStyle.Render(view)
}

func newListView(matches []list.Item) list.Model {
	listView := list.New(filterMatches(matches, FromToday), list.NewDefaultDelegate(), 0, 0)
	listView.Title = "D2M - Dota2 Matches Tracker"
	listView.Styles.Title = titleStyle
	listView.AdditionalFullHelpKeys = filterKeys.FullHelp

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
