package api

import "context"

type Clienter interface {
	GetScheduledMatches(ctx context.Context) ([]*Match, error)
	GetTeamDetailsPage(ctx context.Context, url string) (*Team, error)
}
