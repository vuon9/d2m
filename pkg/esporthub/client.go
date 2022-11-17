package esporthub

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly"
)

const (
	scheduleMatchesURL = "https://esportshub.azure-api.net/schedule/matches"
	homePageURL        = "https://www.msn.com/en-us/esports/calendar/dota2/matches?ocid=winp2oct"
)

func NewEsportHubClient() (*EsportHubClient, error) {
	var scriptContent string

	c := colly.NewCollector()

	c.OnHTML("div[id=esportshub]", func(e *colly.HTMLElement) {
		scriptContent = e.DOM.Next().Text()
	})

	if err := c.Visit(homePageURL); err != nil {
		return nil, err
	}

	return parseCredentials(scriptContent)
}

func (cre *EsportHubClient) GetScheduledMatches(ctx context.Context, gameName GameName) (*ScheduleMatches, error) {
	params := url.Values{}

	now := time.Now()
	startedToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	params.Add("referenceDateTime", startedToday.Format("2006-01-02T15:04:05.000Z")) // 2006-01-02T15:04:05Z07:00
	params.Add("direction", "Forward")
	params.Add("videoGameIds", videoGameMaps[gameName])
	params.Add("limit", "30")
	params.Add("withObjects", "teams,tournaments")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, scheduleMatchesURL, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Add("Ocp-Apim-Subscription-Key", cre.HubSubscriptionKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrMatchUnauthorized
	}

	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var scheduleMatches ScheduleMatches
		err = json.Unmarshal(body, &scheduleMatches)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse JSON: %s", err.Error())
		}

		return &scheduleMatches, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
}