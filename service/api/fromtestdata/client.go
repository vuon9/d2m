package fromtestdata

import (
	"context"

	"github.com/vuon9/d2m/service/api/models"
)

type testdata struct {
}

func NewClient() *testdata {
	return &testdata{}
}

func (c *testdata) GetScheduledMatches(ctx context.Context) ([]*models.Match, error) {
	return nil, nil
}

func (c *testdata) GetTeamDetailsPage(ctx context.Context, url string) (*models.Team, error) {
	return nil, nil
}

func (c *testdata) FromTestData(m *models.TestDataMap) error {
	return nil
}
