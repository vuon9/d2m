package command

import (
	"fmt"

	"github.com/vuon9/d2m/internal/esporthub"
)

type GameName int

const (
	GameDota2 GameName = iota
)

var mapToEsportHubGameName = map[GameName]esporthub.GameName{
	GameDota2: esporthub.Dota2,
}

func GetMatches(g GameName) {
	cre := esporthub.FetchCredentials()

	gameName := mapToEsportHubGameName[g]
	scheduleMatches, err := esporthub.GetScheduledMatches(cre, gameName)
	if err != nil {
		panic("couldn't get matches")
	}

	groupByDates := make(map[string][]*esporthub.Match)
	dates := make([]string, 0)

	for _, m := range scheduleMatches.Matches {
		k := m.Start.Local().Format("Mon, 02 Jan 2006")
		if _, ok := groupByDates[k]; !ok {
			dates = append(dates, k)
		}
		groupByDates[k] = append(groupByDates[k], m)
	}

	for _, date := range dates {
		fmt.Printf("\n***** [ %s ] *****\n", date)
		for _, match := range groupByDates[date] {

			if len(match.Teams) < 2 {
				continue
			}

			matchTime := match.Start.Local().Format("15:04")

			var status string
			switch match.Status {
			case "Resolved":
				status = "[Finish]"
			case "Unresolved":
				status = fmt.Sprintf("[Coming - %s]", matchTime)
			case "Live":
				status = fmt.Sprintf("[Live - %s]", matchTime)
			}

			fmt.Printf("* %s %s vs. %s\n--------------------------\n",
				status,
				match.Teams[0].FullName,
				match.Teams[1].FullName,
			)
		}
	}
}
