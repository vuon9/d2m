package d2m

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api"
)

type (
	model struct {
		listModel    list.Model
		detailsModel tea.Model
		items        []*api.Match
		spinner      spinner.Model
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

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	defaultBodyStyle = lipgloss.NewStyle().MarginLeft(2)

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
	return helpOptions
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

var helpOptions = []key.Binding{
	KeyAllMatches,
	KeyFromTodayMatches,
	KeyTodayMatches,
	KeyTomorrowMatches,
	KeyYesterdayMatches,
	KeyLiveMatches,
	KeyFinishedMatches,
	KeyComingMatches,
}

func newModel() tea.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot

	return &model{
		spinner:   sp,
		listModel: newListView(),
		appState:  showListMatch,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		getMatches,
	)
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

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Commons handling
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.listModel.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		// If the list is filtering, we want to skip all the keys in this state
		if m.listModel.FilterState() == list.Filtering {
			break
		}

		switch {
		case exitKeys[msg.String()]:
			if m.appState == showDetailsMatch {
				m.appState = showListMatch
				return m, nil
			}

			return m, tea.Quit
		case key.Matches(msg, KeyOpenStreamURL):
			m.openStreamingURL()
		}
	}

	// Handling with specific app state
	switch m.appState {
	case showDetailsMatch:
		var cmd tea.Cmd
		m.detailsModel, cmd = m.detailsModel.Update(msg)
		cmds = append(cmds, cmd)
	case showListMatch:
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		case []*api.Match:
			m.items = msg
			m.listModel.SetItems(filterMatches(msg, FromToday))
		case tea.KeyMsg:
			// If the list is filtering, we want to skip all the keys in this state
			if m.listModel.FilterState() == list.Filtering {
				break
			}

			switch {
			case msg.String() == "enter":
				match, ok := m.listModel.SelectedItem().(*api.Match)
				hasAnyLinks := false
				for i, t := range match.Teams {
					if t.TeamProfileLink != "" {
						hasAnyLinks = true
						break
					}

					if i == 1 {
						break
					}
				}

				if ok && hasAnyLinks {
					m.listModel.NewStatusMessage(fmt.Sprintf("Choose match %s", match.GeneralTitle()))
					m.appState = showDetailsMatch
					m.detailsModel = newDetailsModel(match)
					return m, m.detailsModel.Init()
				} else {
					m.listModel.NewStatusMessage("No team info available")
				}

				return m, nil
			case m.DoFilterSuccessful(msg):
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.listModel, cmd = m.listModel.Update(msg)
		cmds = append(cmds, cmd)
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

func (m model) View() string {
	title := titleStyle.Render(m.listModel.Title) + "\n\n"
	view := title

	if m.appState == showListMatch {
		if m.items != nil {
			view = m.listModel.View()
		} else {
			view += m.spinner.View() + " Fetching matches"
			view = defaultBodyStyle.Render(view)
		}
	}

	if m.appState == showDetailsMatch {
		view = title + m.detailsModel.View()
	}

	return appStyle.Render(view)
}

func getMatches() tea.Msg {
	matches, err := apiClient.GetScheduledMatches(context.TODO())
	if err != nil {
		return err
	}

	return matches
}

func newListView() list.Model {
	listView := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listView.AdditionalFullHelpKeys = filterKeys.FullHelp
	listView.Filter = RegexFilter

	listView.Title = "D2M - Dota2 Matches Tracker"
	listView.Styles.Title = titleStyle

	return listView
}
