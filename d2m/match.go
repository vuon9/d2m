package d2m

import (
	"context"
	"time"

	"github.com/vuon9/d2m/pkg/api/haglund"
	"github.com/vuon9/d2m/pkg/api/types"
)

type MatchesByDate map[time.Time]types.MatchSlice

func GetMatches(ctx context.Context, gameName types.GameName) (types.MatchSlice, error) {
	client, err := haglund.NewHaglundClient()
	if err != nil {
		return nil, err
	}

	matches, err := client.GetScheduledMatches(ctx, gameName)
	if err != nil {
		return nil, err
	}

	// loop thorugh matches and sort them by ascending date
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Start.After(matches[j].Start) {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches, nil
}
