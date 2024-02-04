package viewmodel

import (
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vuon9/d2m/pkg/api/model"
)

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

func filterMatches(items []*model.Match, mf matchFilter) []list.Item {
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
			isEligible = match.Status == model.StatusLive
		case Finished:
			isEligible = match.Status == model.StatusFinished
		case Coming:
			isEligible = match.Status == model.StatusComing
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
		re, err := regexp.Compile("(?i)" + term)
		if err != nil {
			continue
		}

		if re.MatchString(target) {
			result = append(result, list.Rank{
				Index:          i,
				MatchedIndexes: []int{0},
			})
		}
	}

	return result
}
