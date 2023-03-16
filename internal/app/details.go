package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api"
)

type detailsModel struct {
	spinner spinner.Model
	match   *api.Match
}

func newDetailsModel(match *api.Match) tea.Model {
	m := &detailsModel{
		match: match,
	}
	m.resetSpinner()

	return m
}

func (m *detailsModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
}

func (m *detailsModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *detailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *detailsModel) View() string {
	// sp := "\n"
	// sp += m.spinner.View()

	teamDetails := "\n"
	teamDetails += fmt.Sprintf("Tournament logo: %s, Page URL: %s\n", m.match.Tournament.Urls.Logo, m.match.Tournament.Urls.Page)

	for _, team := range m.match.Teams {

		teamDetails += fmt.Sprintf("Team: %s, short name: %s, player roster: \n", team.FullName, team.ShortName)
		for _, pl := range team.PlayerRoster {
			teamDetails += fmt.Sprintf("ID: %s\n", pl.ID)
			teamDetails += fmt.Sprintf("Name: %s\n", pl.Name)
			teamDetails += fmt.Sprintf("JoinDate: %s\n", pl.JoinDate)
			teamDetails += fmt.Sprintf("LeaveDate: %s\n", pl.LeaveDate)
			teamDetails += fmt.Sprintf("NewTeam: %s\n", pl.NewTeam)
			teamDetails += fmt.Sprintf("Position: %d\n", pl.Position)
			teamDetails += fmt.Sprintf("ActiveStatus: %d\n", pl.ActiveStatus)
			teamDetails += fmt.Sprintf("IsCaptain: %t\n", pl.IsCaptain)

		}
	}

	// matchTitle := "\n" + m.match.GeneralTitle()
	sp := lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("\n " + string(teamDetails))

	return sp
}

// TODO: Remove??
func newTableModel() table.Model {
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
