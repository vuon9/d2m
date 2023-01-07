package d2m

import (
	"context"
	"time"

	"github.com/vuon9/d2m/pkg/api/haglund"
	"github.com/vuon9/d2m/pkg/api/types"
)

type MatchesByDate map[time.Time]types.MatchSlice

func GetMatches(ctx context.Context, gameName types.GameName) (MatchesByDate, error) {
	client, err := haglund.NewHaglundClient()
	if err != nil {
		return nil, err
	}

	matches, err := client.GetScheduledMatches(ctx, gameName)
	if err != nil {
		return nil, err
	}

	matchesByDate := make(MatchesByDate)
	for _, match := range matches {
		matchDate := match.Start.Local().Truncate(24 * time.Hour)
		matchesByDate[matchDate] = append(matchesByDate[matchDate], match)
	}

	return matchesByDate, nil
}
