package api

import (
	"context"

	"github.com/vuon9/d2m/pkg/api/fromtestdata"
	"github.com/vuon9/d2m/pkg/api/liquipedia"
	"github.com/vuon9/d2m/pkg/api/model"
)

type Clienter interface {
	GetScheduledMatches(ctx context.Context) ([]*model.Match, error)
	GetTeamDetailsPage(ctx context.Context, url string) (*model.Team, error)
}

type TestDataClienter interface {
	Clienter
	FromFileMap(testDataMap *model.TestDataFileMap)
	ReadFiles() error
}

var _ Clienter = liquipedia.NewClient()
var _ TestDataClienter = fromtestdata.NewClient()
