package liquipedia

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/samber/lo"
	"github.com/vuon9/d2m/pkg/api"
)

type UpComingMatchesPageParser struct {
}

func (p *UpComingMatchesPageParser) Parse(anyMatches any) colly.HTMLCallback {
	matches := anyMatches.(*[]*api.Match)
	matchHash := make(map[string]struct{})

	return func(e *colly.HTMLElement) {
		team0 := new(api.Team)
		team1 := new(api.Team)

		match := &api.Match{
			Teams: []*api.Team{
				team0,
				team1,
			},
		}

		e.ForEach("tr", func(_ int, h *colly.HTMLElement) {
			h.ForEach("td.team-left", func(_ int, h *colly.HTMLElement) {
				p.parseTeam(h, team0)
			})

			h.ForEach("td.team-right", func(_ int, h *colly.HTMLElement) {
				p.parseTeam(h, team1)
			})
		})

		p.parseMatchStateAndScores(e, match)

		e.ForEach("tr > td.match-filler", func(_ int, el *colly.HTMLElement) {
			match.Tournament.Name = el.ChildText("div:nth-child(1) > div:nth-child(1) a")
			el.ForEach("div:nth-child(1) span.league-icon-small-image", func(_ int, h *colly.HTMLElement) {
				match.Tournament.Urls.Page = secureDomain + h.ChildAttr("a", "href")
				match.Tournament.Urls.Logo = h.ChildAttr("img", "src")
			})

			el.ForEach("span > span.timer-object", func(_ int, h *colly.HTMLElement) {
				// Get start time of match
				dataStartTimestamp := h.Attr("data-timestamp")
				startTimestamp, _ := strconv.ParseInt(dataStartTimestamp, 10, 64)
				match.Start = time.Unix(startTimestamp, 0)

				// Get twitch channel name
				twitchChannelName := h.Attr("data-stream-twitch")
				if twitchChannelName != "" {
					match.StreamingURL = buildStreamPageLink(twitchChannelName)
				}
			})
		})

		match.IsConcludedMatch = team0.FullName != "TBD" && team1.FullName != "TBD"

		// Generate a hash to use for checking duplicate content which probably parsed in the previous iteration
		// A limitation that when a match has TBDs for both team and happen in the same time, then it would not show correctly
		h := md5.New()
		_, _ = io.WriteString(h, team0.FullName+team1.FullName+match.CompetitionType+match.Start.String())
		hashMatchID := fmt.Sprintf("%x", h.Sum(nil))

		// Only add new item if the hash is new
		if _, found := matchHash[hashMatchID]; !found {
			matchHash[hashMatchID] = struct{}{}
			*matches = append(*matches, match)
		}
	}
}

func (p *UpComingMatchesPageParser) parseTeam(e *colly.HTMLElement, team *api.Team) {
	teamName := e.ChildText("span.team-template-text")
	if teamName != "" {
		team.ShortName = teamName
	}

	teamFullName := e.ChildAttr("span", "data-highlightingclass")
	if teamFullName != "" {
		team.FullName = strings.TrimSpace(teamFullName)
	}

	if team.FullName == "TBD" {
		return
	}

	potentialRelURLs := e.ChildAttrs("a", "href")
	for _, t1PotentialRelativeURL := range potentialRelURLs {
		// the sequence of potential relative URLs is not always the same on each match
		// then better to check if the URL is valid or not
		if team.TeamProfileLink == "" && isValidTeamURL(t1PotentialRelativeURL) {
			team.TeamProfileLink = secureDomain + t1PotentialRelativeURL
			break
		}
	}
}

func (p *UpComingMatchesPageParser) parseMatchStateAndScores(e *colly.HTMLElement, match *api.Match) {
	versus := e.ChildText("tr > td.versus")
	if versus == "" {
		return
	}

	match.CompetitionType = string(regexp.MustCompile(`([A-Z]\w{2})`).Find([]byte(versus)))

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

		match.Team1().Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[0]))
		match.Team2().Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[1]))
	}
}
