package esporthub

import (
	"errors"
	"time"
)

var (
	ErrMatchUnauthorized = errors.New("unauthorized")
)

type EsportHubClient struct {
	ClientID           string
	HubSubscriptionKey string
}

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

type ScheduleMatches struct {
	Matches []*Match `json:"matches"`
}
