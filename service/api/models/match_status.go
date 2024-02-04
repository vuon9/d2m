package models

type MatchStatus uint8

const (
	StatusUnknown MatchStatus = iota
	StatusComing
	StatusLive
	StatusFinished
)

var matchStatuses = map[MatchStatus]string{
	StatusUnknown:  "Unknown",
	StatusComing:   "Coming",
	StatusLive:     "Live",
	StatusFinished: "Finished",
}
