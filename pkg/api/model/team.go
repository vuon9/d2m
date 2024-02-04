package model

import "time"

type Link struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type League struct {
	Link    *Link  `json:"link"`
	IconURL string `json:"iconURL"`
	Name    string `json:"name"`
}

type LiveTeam struct {
	ID        string    `json:"id"`
	ShortName string    `json:"shortName"`
	FullName  string    `json:"fullName"`
	Players   []*Player `json:"rosters"`
}

type Team struct {
	CreatedAt              time.Time `json:"createdAt,omitempty"`
	LastError              error     `json:"lastError,omitempty"`
	Location               string    `json:"location,omitempty"`
	MatchResultDescription string    `json:"matchResultDescription,omitempty"`
	Region                 string    `json:"region,omitempty"`
	Overview               string    `json:"overview,omitempty"`
	History                string    `json:"history,omitempty"`
	ID                     string    `json:"id"`
	MatchResult            string    `json:"matchResult,omitempty"`
	TeamProfileLink        string    `json:"teamProfileLink"`
	Coach                  string    `json:"coach,omitempty"`
	LogoPrimaryColorHsl    string    `json:"logoPrimaryColorHsl,omitempty"`
	ShortName              string    `json:"shortName"`
	FullName               string    `json:"fullName"`
	LogoPrimaryColorRgb    string    `json:"logoPrimaryColorRgb,omitempty"`
	Manager                string    `json:"manager,omitempty"`
	TeamCaptain            string    `json:"teamCaptain,omitempty"`
	Links                  []*Link   `json:"links,omitempty"`
	Archivements           []*League `json:"archivements,omitempty"`
	RecentMatches          []*Match  `json:"recentMatches,omitempty"`
	Articles               []*Link   `json:"articles,omitempty"`
	PlayerRoster           []*Player `json:"players,omitempty"`
	Score                  int       `json:"score"`
	AppoxTotalWinnings     float32   `json:"appoxTotalWinnings,omitempty"`
}
