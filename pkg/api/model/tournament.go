package model

import "time"

type Tournament struct {
	Name                string    `json:"name"`
	PrizePool           string    `json:"prizePool"`
	LogoPrimaryColorRgb string    `json:"logoPrimaryColorRgb,omitempty"`
	LogoPrimaryColorHsl string    `json:"logoPrimaryColorHsl,omitempty"`
	Start               time.Time `json:"start,omitempty"`
	End                 time.Time `json:"end,omitempty"`
	ID                  string    `json:"id"`
	Urls                struct {
		Logo         string `json:"logo"`
		Page         string `json:"pageUrl"`
		BannerImage  string `json:"bannerImage,omitempty"`
		DefaultImage string `json:"defaultImage,omitempty"`
		SquareImage  string `json:"squareImage,omitempty"`
		Thumbnail    string `json:"thumbnail,omitempty"`
		Default      string `json:"default,omitempty"`
		Search       string `json:"search,omitempty"`
	} `json:"urls"`
}
