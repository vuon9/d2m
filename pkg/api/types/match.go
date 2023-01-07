package types

import (
	"time"
)

type MatchSlice []*Match

type Match struct {
	Name                       string     `json:"name"`
	Status                     string     `json:"status"`
	StatusDescription          string     `json:"statusDescription"`
	CompetitionType            string     `json:"competitionType"`
	CompetitionTypeDescription string     `json:"competitionTypeDescription"`
	ContentType                string     `json:"contentType"`
	Tier                       string     `json:"tier"`
	Tournament                 Tournament `json:"tournament"`
	Teams                      []*Team    `json:"teams"`
	Start                      time.Time  `json:"start"`
	ID                         string     `json:"id"`
	VideoGameID                string     `json:"videoGameId"`
	UrlsDescriptions           struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
}

type Tournament struct {
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
}

func (m *Match) Team1() *Team {
	if len(m.Teams) > 0 {
		return m.Teams[0]
	}

	return &Team{
		FullName: "TBD",
	}
}

func (m *Match) Team2() *Team {
	if len(m.Teams) > 1 {
		return m.Teams[1]
	}

	return &Team{
		FullName: "TBD",
	}
}

// TimebasedFriendlyStatus returns a friendly status based on the start time of the match
// it's inaccuracy because it doesn't take into account the status of the match
func (m *Match) TimebasedFriendlyStatus() string {
	now := time.Now()
	threeHours := 3 * time.Hour

	if m.Start.Before(now) {
		if now.Sub(m.Start) > threeHours {
			return "Finish"
		}

		return "Live"
	}

	return "Coming"
}

func (m *Match) FriendlyStatus() string {
	if m.Status == "" {
		return m.TimebasedFriendlyStatus()
	}

	return m.Status
}
