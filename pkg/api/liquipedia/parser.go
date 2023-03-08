package liquipedia

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/samber/lo"
	"github.com/vuon9/d2m/pkg/api"
)

var rootDomain = "https://liquipedia.net"

var allowedDomains = []string{
	"liquipedia.net",
	"www.liquipedia.net",
}

func isValidTeamURL(potentialURL string) bool {
	return !strings.Contains(potentialURL, "/dota2/index.php?title=")
}

func parseUpComingPage(ctx context.Context, req *http.Request) ([]*api.Match, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		// colly.CacheDir("./_cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		for k, v := range req.Header.Clone() {
			r.Headers.Set(k, v[0])
		}
	})

	uniqueMatches := map[string]*api.Match{}
	c.OnHTML("table.infobox_matches_content > tbody", func(e *colly.HTMLElement) {
		match := api.Match{
			Teams: []*api.Team{
				{},
				{},
			},
		}

		teamLeft := e.ChildAttr("tr > td.team-left span", "data-highlightingclass")
		if teamLeft != "" {
			match.Teams[0].FullName = strings.TrimSpace(teamLeft)
			refLink := e.ChildAttr("tr > td.team-left a", "href")
			match.Teams[0].TeamProfileLink = rootDomain + refLink
		}

		// skip parsing scores if the match is not started yet
		t1PotentialRelativeURLs := e.ChildAttrs("tr > td.team-left a", "href")
		for _, t1PotentialRelativeURL := range t1PotentialRelativeURLs {
			// the sequence of potential relative URLs is not always the same on each match
			// then better to check if the URL is valid or not
			if isValidTeamURL(t1PotentialRelativeURL) {
				match.Teams[0].TeamProfileLink = rootDomain + t1PotentialRelativeURL
			} else {
				match.Teams[0].TeamProfileLink = ""
				break
			}
		}

		teamRight := e.ChildAttr("tr > td.team-right span", "data-highlightingclass")
		if teamRight != "" {
			match.Teams[1].FullName = strings.TrimSpace(teamRight)
		}

		// skip parsing scores if the match is not started yet
		if match.Teams[1].FullName != "TBD" {
			// the sequence of potential relative URLs is not always the same on each match
			// then better to check if the URL is valid or not
			t2PotentialRelativeURLs := e.ChildAttrs("tr > td.team-right a", "href")
			for _, t2PotentialRelativeURL := range t2PotentialRelativeURLs {
				if isValidTeamURL(t2PotentialRelativeURL) {
					match.Teams[1].TeamProfileLink = rootDomain + t2PotentialRelativeURL
				} else {
					match.Teams[1].TeamProfileLink = ""
					break
				}
			}
		}

		versus := e.ChildText("tr > td.versus")
		if versus != "" {
			re := regexp.MustCompile(`([A-Z]\w{2})`)
			match.CompetitionType = string(re.Find([]byte(versus)))

			// Skip parsing scores if the match is not started yet
			switch {
			case strings.Contains(versus, "vs"):
				match.Status = api.StatusComing
			case strings.Contains(versus, "Bo"):
				match.Status = api.StatusLive
			default:
				match.Status = api.StatusFinished
			}

			if lo.Contains([]api.MatchStatus{api.StatusFinished, api.StatusLive}, match.Status) {
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

			// Get twitch channel name
			twitchChannelName := el.ChildAttr("span > span.timer-object", "data-stream-twitch")
			if twitchChannelName != "" {
				match.StreamingURL = buildStreamPageLink(twitchChannelName)
				match.HasStreamingURL = true
			}
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

	matches := make([]*api.Match, 0)
	for k, m := range uniqueMatches {
		matches = append(matches, m)
		uniqueMatches[k] = m
	}

	return matches, nil
}

// buildStreamPageLink builds a link to the stream page on liquipedia
func buildStreamPageLink(channelName string) string {
	return fmt.Sprintf("https://liquipedia.net/dota2/Special:Stream/twitch/%s", channelName)
}
