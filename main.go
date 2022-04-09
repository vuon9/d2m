package main

import (
	"fmt"
)

func main() {
	cre := fetchCredentials()

	scheduleMatches, err := getScheduledMatches(cre)
	if err != nil {
		panic("couldn't get matches")
	}

	for _, match := range scheduleMatches.Matches {
		if len(match.Teams) < 2 {
			continue
		}

		var status string
		switch match.Status {
		case "Resolved":
			status = "[Finish]"
		case "Unresolved":
			status = "[Coming]"
		case "Live":
			status = "[Live]  "
		}

		fmt.Printf("* %s %s vs. %s\n  at %s\n--------------------------\n",
			status,
			match.Teams[0].FullName,
			match.Teams[1].FullName,
			match.Start.Local().Format("15:04 2006-01-02"),
		)
	}
}
