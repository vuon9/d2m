package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/pkg/api"
)

type detailsModel struct {
	spinner     spinner.Model
	match       *api.Match
	fetchIsDone bool
	lastErr     error
}

func newDetailsModel(match *api.Match) tea.Model {
	m := &detailsModel{
		match: match,
	}
	m.resetSpinner()

	return m
}

func fetchTeams(match *api.Match) func() tea.Msg {
	return func() tea.Msg {
		urls := []string{}

		for _, team := range match.Teams {
			if team.TeamProfileLink == "" {
				continue
			}

			urls = append(urls, team.TeamProfileLink)
		}

		wg := sync.WaitGroup{}
		teams := make([]*api.Team, 0)

		var lastErr error

		for _, url := range urls {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()
				team, err := apiClient.GetTeamDetailsPage(context.TODO(), url)
				if err != nil {
					lastErr = errors.New(fmt.Sprintf("Error while fetching team details: %s, url: %s", err.Error(), url))
					return
				}

				teams = append(teams, team)
			}(url)
		}

		wg.Wait()

		if lastErr != nil {
			return lastErr
		}

		return teams
	}
}

func (m *detailsModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
}

func (m *detailsModel) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		fetchTeams(m.match),
	)
}

func (m *detailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			return m, nil
		}
	case []*api.Team:
		m.fetchIsDone = true
		m.match.Teams = msg

		return m, nil
	case error:
		m.fetchIsDone = true
		m.lastErr = msg

		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *detailsModel) View() string {
	if !m.fetchIsDone {
		return m.spinner.View() + " Fetching teams"
	}

	if m.lastErr != nil {
		return m.lastErr.Error()
	}

	playerTableCols := []table.Column{
		{
			Title: "No.",
			Width: 3,
		},
		{
			Title: "In-game ID",
			Width: 15,
		},
		{
			Title: "Name",
			Width: 20,
		},
		{
			Title: "Position",
			Width: 10,
		},
		{
			Title: "Join Date",
			Width: 10,
		},
	}

	teamDetails := headerStyle.Render("Player Roster - Active") + "\n\n"

	activePlayers := []table.Row{}
	for i, t := range m.match.Teams {
		teamDetails += fmt.Sprintf("Team %d: %s\n", i+1, t.FullName)
		for _, p := range t.PlayerRoster {
			activePlayers = append(activePlayers, table.Row{
				fmt.Sprintf("%d", i+1),
				p.ID,
				p.Name,
				p.Position.String(),
				p.JoinDate,
			})
		}

		t := table.New(
			table.WithColumns(playerTableCols),
			table.WithRows(activePlayers),
			table.WithHeight(9),
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
		t.SetStyles(s)

		teamDetails += t.View() + "\n"
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("\n " + teamDetails)
}
