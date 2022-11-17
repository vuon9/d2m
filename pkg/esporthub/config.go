package esporthub

type GameName uint8

const (
	Dota2 GameName = iota + 1
	LeagueOfLegends
	CsGO
	Valorant
)

var videoGameMaps = map[GameName]string{
	Dota2:           "51b8bf37-fede-45d5-3943-fef79b0fa628",
	LeagueOfLegends: "dcd754dc-9f53-0f8b-3bfc-9c401f16138b",
	CsGO:            "8f345aa7-ff29-2efd-48fd-8230fd8795aa",
	Valorant:        "72d5fb42-ec96-44e0-ae93-2e1cfb2c1836",
}
