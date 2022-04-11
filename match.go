package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Match struct {
	Name                       string `json:"name"`
	Status                     string `json:"status"`
	StatusDescription          string `json:"statusDescription"`
	CompetitionType            string `json:"competitionType"`
	CompetitionTypeDescription string `json:"competitionTypeDescription"`
	ContentType                string `json:"contentType"`
	Tier                       string `json:"tier"`
	Tournament                 struct {
		Name                string    `json:"name"`
		PrizePool           string    `json:"prizePool"`
		LogoPrimaryColorRgb string    `json:"logoPrimaryColorRgb"`
		LogoPrimaryColorHsl string    `json:"logoPrimaryColorHsl"`
		Start               time.Time `json:"start"`
		End                 time.Time `json:"end"`
		ID                  string    `json:"id"`
		Urls                struct {
			Logo         string `json:"logo"`
			BannerImage  string `json:"bannerImage"`
			DefaultImage string `json:"defaultImage"`
			SquareImage  string `json:"squareImage"`
			Thumbnail    string `json:"thumbnail"`
			Default      string `json:"default"`
			Search       string `json:"search"`
		} `json:"urls"`
		UrlsDescriptions struct {
			Logo    string `json:"logo"`
			Default string `json:"default"`
		} `json:"urlsDescriptions"`
	} `json:"tournament"`
	Teams            []*Team   `json:"teams"`
	Start            time.Time `json:"start"`
	ID               string    `json:"id"`
	VideoGameID      string    `json:"videoGameId"`
	UrlsDescriptions struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
}

type Team struct {
	ShortName              string `json:"shortName"`
	FullName               string `json:"fullName"`
	Score                  int    `json:"score"`
	MatchResult            string `json:"matchResult"`
	MatchResultDescription string `json:"matchResultDescription"`
	LogoPrimaryColorRgb    string `json:"logoPrimaryColorRgb"`
	LogoPrimaryColorHsl    string `json:"logoPrimaryColorHsl"`
	ID                     string `json:"id"`
	Urls                   struct {
		Logo   string `json:"logo"`
		Search string `json:"search"`
	} `json:"urls"`
	UrlsDescriptions struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
}

var (
	ErrMatchUnauthorized = errors.New("unauthorized")
	scheduleMatchesURL   = "https://esportshub.azure-api.net/schedule/matches"
)

type GameName string

var Dota2 GameName = "dota2"

var videoGameIds = map[GameName]string{
	Dota2: "51b8bf37-fede-45d5-3943-fef79b0fa628",
}

type ScheduleMatches struct {
	Matches []*Match `json:"matches"`
}

func getScheduledMatches(cre *MatchAPICredentials, videoGameId string) (*ScheduleMatches, error) {
	params := url.Values{}
	params.Add("referenceDateTime", time.Now().Format("2006-01-02T15:04:05.000Z")) // 2006-01-02T15:04:05Z07:00
	params.Add("direction", "Forward")
	params.Add("videoGameIds", videoGameId)
	params.Add("limit", "30")
	params.Add("withObjects", "teams,tournaments")

	req, err := http.NewRequest(http.MethodGet, scheduleMatchesURL, nil)
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
