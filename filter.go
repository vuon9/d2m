package d2m

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vuon9/d2m/pkg/api"
)

type matchFilter uint8

const (
	All matchFilter = iota
	FromToday
	Today
	Tomorrow
	Yesterday
	Live
	Finished
	Coming
)

func filterMatches(items []*api.Match, mf matchFilter) []list.Item {
	var filteredItems []list.Item

	for _, match := range items {
		var isEligible bool

		switch mf {
		case All:
			isEligible = true
		case FromToday:
			t := time.Now()
			isEligible = match.Start.After(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
		case Today:
			isEligible = match.Start.Day() == time.Now().Day()
		case Tomorrow:
			isEligible = match.Start.Day() == time.Now().AddDate(0, 0, 1).Day()
		case Yesterday:
			isEligible = match.Start.Day() == time.Now().AddDate(0, 0, -1).Day()
		case Live:
			isEligible = match.Status == api.StatusLive
		case Finished:
			isEligible = match.Status == api.StatusFinished
		case Coming:
			isEligible = match.Status == api.StatusComing
		default:
			continue
		}

		if isEligible {
			filteredItems = append(filteredItems, match)
		}
	}

	return filteredItems
}
