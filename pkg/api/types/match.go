package types

import (
	"time"
)

type MatchStatus string

func (m MatchStatus) String() string {
	return string(m)
}

const (
	MatchStatusComing   MatchStatus = "Coming"
	MatchStatusLive     MatchStatus = "Live"
	MatchStatusFinished MatchStatus = "Finished"
)

type MatchSlice []*Match

type Match struct {
	Start                      time.Time  `json:"start"`
	Tournament                 Tournament `json:"tournament"`
	Name                       string     `json:"name"`
	CompetitionType            string     `json:"competitionType"`
	CompetitionTypeDescription string     `json:"competitionTypeDescription"`
	ContentType                string     `json:"contentType"`
	Tier                       string     `json:"tier"`
	StatusDescription          string     `json:"statusDescription"`
	VideoGameID                string     `json:"videoGameId"`
	Status                     string     `json:"status"`
	ID                         string     `json:"id"`
	UrlsDescriptions           struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
	Teams []*Team `json:"teams"`
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
func (m *Match) TimebasedFriendlyStatus() MatchStatus {
	now := time.Now()
	threeHours := 3 * time.Hour

	if m.Start.Before(now) {
		if now.Sub(m.Start) > threeHours {
			return MatchStatusFinished
		}

		return MatchStatusLive
	}

	return MatchStatusComing
}

func (m *Match) FriendlyStatus() MatchStatus {
	if m.Status == "" {
		return m.TimebasedFriendlyStatus()
	}

	return MatchStatus(m.Status)
}
