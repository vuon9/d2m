package api

import "time"

type Link struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

type Team struct {
	ID                     string    `json:"id"`
	ShortName              string    `json:"shortName"`
	FullName               string    `json:"fullName"`
	TeamProfileLink        string    `json:"teamProfileLink"`
	Score                  int       `json:"score"`
	MatchResult            string    `json:"matchResult"`
	MatchResultDescription string    `json:"matchResultDescription"`
	LogoPrimaryColorRgb    string    `json:"logoPrimaryColorRgb"`
	LogoPrimaryColorHsl    string    `json:"logoPrimaryColorHsl"`
	Location               string    `json:"location"`
	Region                 string    `json:"region"`
	Manager                string    `json:"manager"`
	TeamCaptain            string    `json:"teamCaptain"`
	AppoxTotalWinnings     float32   `json:"appoxTotalWinnings"`
	Links                  []*Link   `json:"links"`
	CreatedAt              time.Time `json:"createdAt"`
	Urls                   struct {
		Logo   string `json:"logo"`
		Search string `json:"search"`
	} `json:"urls"`
	UrlsDescriptions struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
	Players       []*Player `json:"players"`
	RecentMatches []*Match  `json:"recentMatches"`
}
