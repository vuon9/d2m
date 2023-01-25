package d2m

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
	"github.com/vuon9/d2m/pkg/api/types"
)

func GetMatches(ctx context.Context, gameName types.GameName) (types.MatchSlice, error) {
	client, err := liquipedia.NewClient()
	if err != nil {
		return nil, err
	}

	matches, err := client.GetScheduledMatches(ctx, gameName)
	if err != nil {
		return nil, err
	}

	// loop through matches and sort them by ascending date
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Start.After(matches[j].Start) {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches, nil
}



type matchFilter uint8
const (
	all matchFilter = iota
	FromToday
	Today
	Tomorrow
	Yesterday
	Live
	Finished
	Coming
)

func (d *delegator) filterMatches(mf matchFilter) []list.Item {
	var newList []list.Item

	for _, originItem := range d.originItems {
		matcher, ok := originItem.(Matchable)
		if !ok {
			continue
		}

		var isEligible bool

		switch mf {
		case FromToday:
			isEligible = !matcher.StartTime().Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour))
		case Today:
			isEligible = matcher.StartTime().Day() == time.Now().Day()
		case Tomorrow:
			isEligible = matcher.StartTime().Day() == time.Now().AddDate(0, 0, 1).Day()
		case Yesterday:
			isEligible = matcher.StartTime().Day() == time.Now().AddDate(0, 0, -1).Day()
		case Live:
			isEligible = matcher.Status() == "Live"
		case Finished:
			isEligible = matcher.Status() == "Finished"
		case Coming:
			isEligible = matcher.Status() == "Coming"
		default:
			continue
		}

		if isEligible {
			newList = append(newList, originItem)
		}
	}

	return newList
}
