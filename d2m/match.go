package d2m

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vuon9/d2m/pkg/api"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
)

type Tracker interface {
	GetMatches(ctx context.Context) ([]*api.Match, error)
}

type tracker struct {
	client api.Clienter
}

var (
	_ Tracker = (*tracker)(nil)
)

func NewTracker() *tracker {
	return &tracker{
		client: liquipedia.NewClient(),
	}
}

func (d *tracker) GetMatches(ctx context.Context) ([]*api.Match, error) {
	matches, err := d.client.GetScheduledMatches(ctx)
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
	All matchFilter = iota
	FromToday
	Today
	Tomorrow
	Yesterday
	Live
	Finished
	Coming
)

func filterMatches(items []list.Item, mf matchFilter) []list.Item {
	var filteredItems []list.Item

	for _, match := range items {
		match, ok := match.(*api.Match)
		if !ok {
			continue
		}

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
			isEligible = match.Status == "Live"
		case Finished:
			isEligible = match.Status == "Finished"
		case Coming:
			isEligible = match.Status == "Coming"
		default:
			continue
		}

		if isEligible {
			filteredItems = append(filteredItems, match)
		}
	}

	return filteredItems
}
