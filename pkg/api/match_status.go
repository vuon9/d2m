package api

type MatchStatus uint8

const (
	StatusComing MatchStatus = iota
	StatusLive
	StatusFinished
)

var matchStatuses = map[MatchStatus]string{
	StatusComing:   "Coming",
	StatusLive:     "Live",
	StatusFinished: "Finished",
}
