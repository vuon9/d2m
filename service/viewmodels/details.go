package viewmodels

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vuon9/d2m/service/api/liquipedia"
	"github.com/vuon9/d2m/service/api/models"
)

type matchDetail struct {
	spinner     spinner.Model
	match       *models.Match
	fetchIsDone bool
	lastErr     error
}

func newDetailsModel(match *models.Match) tea.Model {
	m := &matchDetail{
		match: match,
	}
	m.resetSpinner()

	return m
}

func fetchTeams(teams []*models.Team) func() tea.Msg {
	return func() tea.Msg {
		urls := []string{}

		for _, team := range teams {
			if team.TeamProfileLink == "" {
				continue
			}

			urls = append(urls, team.TeamProfileLink)
		}

		wg := sync.WaitGroup{}
		teams := make([]*models.Team, 0)

		for _, url := range urls {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				team, err := liquipedia.NewClient().GetTeamDetailsPage(context.TODO(), url)
				if err != nil {
					team = new(models.Team)
					team.LastError = errors.New(fmt.Sprintf("Error while fetching team details: %s, url: %s", err.Error(), url))
				}

				teams = append(teams, team)
			}(url)
		}

		wg.Wait()

		return teams
	}
}

func (m *matchDetail) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
}

func (m *matchDetail) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
		fetchTeams(m.match.Teams),
	)
}

func (m *matchDetail) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			return m, nil
		}
	case []*models.Team:
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

func (m *matchDetail) View() string {
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

	for i, t := range m.match.Teams {
		activePlayers := []table.Row{}
		teamDetails += fmt.Sprintf("Team %d: %s\n", i+1, t.FullName)
		for j, p := range t.PlayerRoster {
			activePlayers = append(activePlayers, table.Row{
				fmt.Sprintf("%d", j+1),
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
