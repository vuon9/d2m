package api

import "time"

type Link struct {
	Link string `json:"link"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type League struct {
	IconURL string `json:"iconURL"`
	Name    string `json:"name"`
	Link    *Link  `json:"link"`
}

type LiveTeam struct {
	ID        string    `json:"id"`
	ShortName string    `json:"shortName"`
	FullName  string    `json:"fullName"`
	Players   []*Player `json:"rosters"`
}

type Team struct {
	ID                     string    `json:"id"`
	ShortName              string    `json:"shortName"`
	FullName               string    `json:"fullName"`
	Overview               string    `json:"overview"`
	History                string    `json:"history"`
	TeamProfileLink        string    `json:"teamProfileLink"`
	PlayerRoster           []*Player `json:"players"`
	Score                  int       `json:"score"`
	MatchResult            string    `json:"matchResult"`
	MatchResultDescription string    `json:"matchResultDescription"`
	LogoPrimaryColorRgb    string    `json:"logoPrimaryColorRgb"`
	LogoPrimaryColorHsl    string    `json:"logoPrimaryColorHsl"`
	Location               string    `json:"location"`
	Region                 string    `json:"region"`
	Coach                  string    `json:"coach"`
	Manager                string    `json:"manager"`
	TeamCaptain            string    `json:"teamCaptain"`
	AppoxTotalWinnings     float32   `json:"appoxTotalWinnings"`
	Links                  []*Link   `json:"links"`
	CreatedAt              time.Time `json:"createdAt"`
	Archivements           []*League `json:"archivements"`
	RecentMatches          []*Match  `json:"recentMatches"`
	Articles               []*Link   `json:"articles"`
}
