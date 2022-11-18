package types

type GameName uint8

const (
	Dota2 GameName = iota + 1
	LeagueOfLegends
	CsGO
	Valorant
)
