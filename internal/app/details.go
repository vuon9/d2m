package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type detailsModel struct {
	spinner spinner.Model
	match   string
}

func newDetailsModel(match string) detailsModel {
	return detailsModel{
		spinner: spinner.NewModel(),
		match:   match,
	}
}

func (m detailsModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m detailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m detailsModel) View() string {
	sp := m.spinner.View()
	sp += " I'm loading"

	return sp
}
