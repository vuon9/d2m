package fromtestdata

import (
	"context"
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	"github.com/vuon9/d2m/pkg/api/model"
)

type fetchedData struct {
	scheduledMatches []*model.Match
	teamDetails      *model.Team
}

// Testdata client is a client that reads test data from files
// instead of fetching from the internet. It's used for testing with fixed data
type testdata struct {
	fileMap     *model.TestDataFileMap
	fetchedData *fetchedData
}

func NewClient() *testdata {
	return &testdata{
		fileMap:     new(model.TestDataFileMap),
		fetchedData: nil,
	}
}

func (c *testdata) GetScheduledMatches(ctx context.Context) ([]*model.Match, error) {
	err := c.ReadFiles()
	if err != nil {
		return nil, err
	}

	return c.fetchedData.scheduledMatches, nil
}

func (c *testdata) GetTeamDetailsPage(ctx context.Context, url string) (*model.Team, error) {
	err := c.ReadFiles()
	if err != nil {
		return nil, err
	}

	return c.fetchedData.teamDetails, nil
}

func (c *testdata) FromFileMap(m *model.TestDataFileMap) {
	c.fileMap = m
}

func (c *testdata) ReadFiles() error {
	if c.fetchedData != nil {
		return nil
	}

	c.fetchedData = new(fetchedData)

	rawMatchList, err := os.ReadFile(c.fileMap.MatchList)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawMatchList, &c.fetchedData.scheduledMatches)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal match list")
	}

	rawTeamDetails, err := os.ReadFile(c.fileMap.Team)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawTeamDetails, &c.fetchedData.teamDetails)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal team details")
	}

	return nil
}
