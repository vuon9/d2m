package viewmodels

import (
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vuon9/d2m/service/api/models"
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

func filterMatches(items []*models.Match, mf matchFilter) []list.Item {
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
			isEligible = match.Status == models.StatusLive
		case Finished:
			isEligible = match.Status == models.StatusFinished
		case Coming:
			isEligible = match.Status == models.StatusComing
		default:
			continue
		}

		if isEligible {
			filteredItems = append(filteredItems, match)
		}
	}

	return filteredItems
}

func RegexFilter(term string, targets []string) []list.Rank {
	var result []list.Rank

	for i, target := range targets {
		regexp, err := regexp.Compile(term)
		if err != nil {
			continue
		}

		if regexp.MatchString(target) {
			result = append(result, list.Rank{
				Index:          i,
				MatchedIndexes: []int{0},
			})
		}
	}

	return result
}