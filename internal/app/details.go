package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type detailsModel struct {
	spinner spinner.Model
	match   string
}

func newDetailsModel(match string) tea.Model {
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
	sp := m.spinner.View()
	sp += lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(" I'm loading")

	return sp
}
