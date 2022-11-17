package command

import (
	"context"
	"fmt"
	"time"

	"github.com/vuon9/d2m/pkg/esporthub"
)

type Match struct {
	esporthub.Match
}

func (m *Match) Team1() *esporthub.Team {
	if len(m.Teams) > 0 {
		return m.Teams[0]
	}

	return &esporthub.Team{
		FullName: "TBD",
	}
}

func (m *Match) Team2() *esporthub.Team {
	if len(m.Teams) > 1 {
		return m.Teams[1]
	}

	return &esporthub.Team{
		FullName: "TBD",
	}
}

func (m *Match) Status() string {
	switch m.Match.Status {
	case "Resolved":
		return "Finish"
	case "Unresolved":
		return "Coming"
	case "Live":
		return "Live"
	default:
		return fmt.Sprintf("[UN] %s", m.Match.Status)
	}
}

type MatchesByDate map[time.Time][]*Match

func GetMatchesByDate(ctx context.Context, gameName esporthub.GameName) (MatchesByDate, error) {
	client, err := esporthub.NewEsportHubClient()
	if err != nil {
		return nil, err
	}

	scheduleMatches, err := client.GetScheduledMatches(ctx, gameName)
	if err != nil {
		return nil, err
	}

	matchesByDate := make(MatchesByDate)
	for _, match := range scheduleMatches.Matches {
		matchDate := match.Start.Local().Truncate(24 * time.Hour)
		match := &Match{
			Match: *match,
		}
		matchesByDate[matchDate] = append(matchesByDate[matchDate], match)
	}

	return matchesByDate, nil
}
