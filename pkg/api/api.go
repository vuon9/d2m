package api

import (
	"context"
	"errors"

	"github.com/vuon9/d2m/pkg/api/types"
)

var (
	ErrMatchUnauthorized = errors.New("unauthorized")
)

type Apier interface {
	FetchScheduledMatches(ctx context.Context, gameName types.GameName) ([]*types.Match, error)
}
