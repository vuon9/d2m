package api

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
