package api

type Player struct {
	GameID   string `json:"gameID"`
	Name     string `json:"name"`
	JoinDate string `json:"joinDate"`
	Position uint8  `json:"position"`
}
