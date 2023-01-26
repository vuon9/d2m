package api

import (
	"fmt"
	"time"

	"github.com/samber/lo"
)

type MatchStatus string

func (m MatchStatus) String() string {
	return string(m)
}

const (
	StatusComing   MatchStatus = "Coming"
	StatusLive     MatchStatus = "Live"
	StatusFinished MatchStatus = "Finished"
)

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
	Status                     MatchStatus     `json:"status"`
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

func (m *Match) Title() string {
	vsOrScores := fmt.Sprintf("%s", m.Status)
	if lo.Contains([]MatchStatus{StatusLive, StatusFinished}, m.Status) {
		vsOrScores = fmt.Sprintf("[%d:%d] - %s", m.Team1().Score, m.Team2().Score, m.Status)
	}

	return fmt.Sprintf("%s - %s",
		vsOrScores,
		m.GeneralTitle(),
	)
}

func (m *Match) GeneralTitle() string {
	return fmt.Sprintf("%s vs. %s", m.Team1().FullName, m.Team2().FullName)
}

func (m *Match) Description() string {
	return fmt.Sprintf("[%s] - %s", m.Start.Format("2006-01-02"),  m.Tournament.Name)
}

func (m *Match) FilterValue() string {
	return m.GeneralTitle() + " " + m.Description()
}

