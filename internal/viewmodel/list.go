package viewmodel

import (
	"github.com/charmbracelet/bubbles/list"
)

func newListView() list.Model {
	listView := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listView.AdditionalFullHelpKeys = filterKeys.FullHelp
	listView.Filter = RegexFilter

	listView.Title = "D2M - Dota2 Matches Tracker"
	listView.Styles.Title = titleStyle

	return listView
}
