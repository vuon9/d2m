package haglund

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vuon9/d2m/pkg/api/types"
)

const (
	scheduleMatchesURL = "https://dota.haglund.dev/v1/matches"
)

type HaglundClient struct {
}

type Match struct {
	Hash       string    `json:"hash"`
	MatchType  string    `json:"matchType"`
	StreamURL  string    `json:"streamUrl"`
	StartsAt   time.Time `json:"startsAt"`
	LeagueName string    `json:"leagueName"`
	LeagueURL  string    `json:"leagueUrl"`
	Teams      []Team    `json:"teams"`
}

type Team struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Response []Match

func NewHaglundClient() (*HaglundClient, error) {
	return &HaglundClient{}, nil
}

func (cre *HaglundClient) GetScheduledMatches(ctx context.Context, gameName types.GameName) (types.MatchSlice, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, scheduleMatchesURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// parse json into types.MatchSlice
	var matches types.MatchSlice
	if res.StatusCode != http.StatusOK {

		return nil, errors.New("It's not OK")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse JSON: %s", err.Error())
	}

	for _, match := range resp {
		// only get matches which has start time within +/- 1 day
		if match.StartsAt.Before(time.Now().Add(-24*time.Hour)) || match.StartsAt.After(time.Now().Add(24*time.Hour)) {
			continue
		}

		m := types.Match{
			Start:                      match.StartsAt,
			CompetitionType:            match.MatchType,
			CompetitionTypeDescription: match.LeagueName,
			Tournament: types.Tournament{
				Name: match.LeagueName,
			},
			Teams: make([]*types.Team, 0),
		}

		for _, team := range match.Teams {
			t := &types.Team{
				FullName: team.Name,
			}
			m.Teams = append(m.Teams, t)
		}

		matches = append(matches, &m)
	}

	return matches, nil
}
