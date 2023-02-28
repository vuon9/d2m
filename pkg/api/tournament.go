package api

import "time"

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
