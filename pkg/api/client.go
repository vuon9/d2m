package api

import "context"

type Clienter interface {
	GetScheduledMatches(ctx context.Context) ([]*Match, error)
}