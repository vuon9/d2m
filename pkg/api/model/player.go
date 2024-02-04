package model

type PlayerStatus uint8

func (ps PlayerStatus) String() string {
	return [...]string{"Active", "Inactive", "Former", "StandIn"}[ps]
}

const (
	Unknown PlayerStatus = iota
	Active
	Inactive
	Former
	StandIn
)

type Position int

func (p Position) String() string {
	return [...]string{"Unknown", "1", "2", "3", "4", "5"}[p]
}

const (
	PosUnknown Position = iota
	Pos1
	Pos2
	Pos3
	Pos4
	Pos5
)

type Player struct {
	ID             string       `json:"gameID"`
	Name           string       `json:"name"`
	JoinDate       string       `json:"joinDate"`
	LeaveDate      string       `json:"leaveDate"`
	NewTeam        string       `json:"newTeam"`
	Position       Position     `json:"position"`
	ActiveStatus   PlayerStatus `json:"isActive"`
	IsCaptain      bool         `json:"isCaptain"`
	ProfilePageURL string       `json:"profilePageURL"`
}
