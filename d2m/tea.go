package d2m

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api/types"
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
)

type model struct {
	list list.Model
}

type item struct {
	title string
	description string
}

func (i *item) Title() string {
	return i.title
}

func (i *item) FilterValue() string {
	return i.title + i.description
}

func (i *item) Description() string {
	return i.description
}

func newModel(matches types.MatchSlice) tea.Model {
	items := make([]list.Item, len(matches), len(matches))
	for i, match := range matches {
		items[i] = &item{
			title: match.Team1().FullName + " vs. " + match.Team2().FullName,
			description: "[" + match.Start.Format("2006-01-02 15:04") + "] - " + match.Tournament.Name,
		}
	}


	matchList := list.New(items, newItemDelegate(newDelegateKeyMap()), 0, 0)
	matchList.Title = "D2M - Dota2 Matches"
	matchList.Styles.Title = titleStyle
	matchList.AdditionalFullHelpKeys = func () []key.Binding {
		return []key.Binding{}
	}

	return &model{
		list: matchList,
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
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
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
