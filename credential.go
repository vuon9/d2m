package main

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var homePageURL = "https://www.msn.com/en-us/esports/calendar/dota2/matches?ocid=winp2oct"

func fetchCredentials() *MatchAPICredentials {
	var scriptContent string

	c := colly.NewCollector()
	c.OnHTML("div[id=esportshub]", func(e *colly.HTMLElement) {
		scriptContent = e.DOM.Next().Text()
	})

	c.Visit(homePageURL)

	return parseCredentials(scriptContent)
}

func parseCredentials(str string) *MatchAPICredentials {
	re := regexp.MustCompile(`(?:clientId|hubSubscriptionKey)+:"+([a-zA-Z0-9]+)"`)

	var cre MatchAPICredentials
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		if strings.HasPrefix(match[0], "clientId") {
			cre.ClientID = match[1]
		}

		if strings.HasPrefix(match[0], "hubSubscriptionKey") {
			cre.HubSubscriptionKey = match[1]
		}
	}

	if (cre == MatchAPICredentials{}) {
		return nil
	}

	return &cre
}
