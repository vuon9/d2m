package api

import (
	"context"

	"github.com/vuon9/d2m/service/api/fromtestdata"
	"github.com/vuon9/d2m/service/api/liquipedia"
	"github.com/vuon9/d2m/service/api/models"
)

type Clienter interface {
	GetScheduledMatches(ctx context.Context) ([]*models.Match, error)
	GetTeamDetailsPage(ctx context.Context, url string) (*models.Team, error)
}

type TestDataClienter interface {
	Clienter
	FromTestData(testDataMap *models.TestDataMap) error
}

var _ Clienter = liquipedia.NewClient()
var _ TestDataClienter = fromtestdata.NewClient()
