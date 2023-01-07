package liquipedia

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/pkg/api/types"
)

var allowedDomains = []string{
	"liquipedia.net",
	"www.liquipedia.net",
}

func ScarppingHTML(ctx context.Context, req *http.Request) (types.MatchSlice, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		// colly.CacheDir("./_cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		for k, v := range req.Header.Clone() {
			r.Headers.Set(k, v[0])
		}
	})

	uniqueMatches := map[string]*types.Match{}
	c.OnHTML("table.infobox_matches_content > tbody", func(e *colly.HTMLElement) {
		match := types.Match{
			Teams: []*types.Team{
				{},
				{},
			},
		}

		teamLeft := e.ChildAttr("tr > td.team-left span", "data-highlightingclass")
		if teamLeft != "" {
			match.Teams[0].FullName = strings.TrimSpace(teamLeft)
		}

		teamRight := e.ChildAttr("tr > td.team-right span", "data-highlightingclass")
		if teamRight != "" {
			match.Teams[1].FullName = strings.TrimSpace(teamRight)
		}

		versus := e.ChildText("tr > td.versus")
		if versus != "" {
			match.CompetitionType = strings.ReplaceAll(versus, "(", " (")

			// Skip parsing scores if the match is not started yet
			if strings.Contains(versus, "vs") {
				match.Status = "Coming"
			} else if strings.Contains(versus, "Bo") {
				match.Status = "Live"
			} else {
				match.Status = "Finished"
			}

			if match.Status == "Finished" || match.Status == "Live" {
				rawScores := strings.Split(versus, ":")
				score0, _ := strconv.ParseInt(strings.TrimSpace(rawScores[0]), 10, 64)
				score1, _ := strconv.ParseInt(strings.TrimSpace(rawScores[1]), 10, 64)

				match.Teams[0].Score = int(score0)
				match.Teams[1].Score = int(score1)
			}
		}

		e.ForEach("tr > td.match-filler", func(_ int, el *colly.HTMLElement) {
			match.Tournament.Name = el.ChildText("div:nth-child(1) > div:nth-child(1) a")
			match.Tournament.Urls.Logo = el.ChildAttr("div:nth-child(1) span.league-icon-small-image img", "src")

			// Get start time of match
			dataStartTimestamp := el.ChildAttr("span > span.timer-object", "data-timestamp")
			startTimestamp, _ := strconv.ParseInt(dataStartTimestamp, 10, 64)
			match.Start = time.Unix(startTimestamp, 0)
		})

		// Store with has to avoid duplicate matches
		h := md5.New()
		_, _ = io.WriteString(h, match.Teams[0].FullName+match.Teams[1].FullName+match.Start.String())
		uniqueMatches[fmt.Sprintf("%x", h.Sum(nil))] = &match
	})

	err := c.Visit(req.URL.String())
	if err != nil {
		return nil, err
	}

	matches := types.MatchSlice{}
	for k, m := range uniqueMatches {
		matches = append(matches, m)
		uniqueMatches[k] = m
	}

	return matches, nil
}

func checkBoMatch(boType string, score0, score1 int64) bool {
	checkBoMaps := map[string][]string{
		"(Bo2)": {"2:0", "0:2", "1:1"},
		"(Bo3)": {"1:2", "2:1", "2:0"},
		"(Bo5)": {"1:3", "3:1", "2:3", "3:2", "3:0", "0:3"},
	}

	pScore := fmt.Sprintf("%d:%d", score0, score1)
	for _, score := range checkBoMaps[boType] {
		if pScore == score {
			return true
		}
	}

	return false
}
