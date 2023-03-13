package api

type PlayerStatus uint8

const (
	Active PlayerStatus = iota
	Inactive
	Former
	StandIn
)

type Position uint8

const (
	Pos1 Position = iota
	Pos2
	Pos3
	Pos4
	Pos5
	Sub
)

type Player struct {
	ID           string       `json:"gameID"`
	Name         string       `json:"name"`
	JoinDate     string       `json:"joinDate"`
	LeaveDate    string       `json:"leaveDate"`
	NewTeam      string       `json:"newTeam"`
	Position     Position     `json:"position"`
	ActiveStatus PlayerStatus `json:"isActive"`
	IsCaptain    bool         `json:"isCaptain"`
}
