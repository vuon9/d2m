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

func parseUpComingMatchesPage(matches *[]*api.Match) colly.HTMLCallback {
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

		e.ForEach("tr", func(i int, h *colly.HTMLElement) {
			h.ForEach("td.team-left", func(i int, e *colly.HTMLElement) {
				teamLeft := e.ChildAttr("span", "data-highlightingclass")
				if teamLeft != "" {
					team0.FullName = strings.TrimSpace(teamLeft)
				}

				if team0.FullName == "TBD" {
					return
				}

				team0.ShortName = e.ChildText("span.team-template-text")

				// skip parsing scores if the match is not started yet
				t1PotentialRelativeURLs := e.ChildAttrs("a", "href")
				for _, t1PotentialRelativeURL := range t1PotentialRelativeURLs {
					// the sequence of potential relative URLs is not always the same on each match
					// then better to check if the URL is valid or not
					if isValidTeamURL(t1PotentialRelativeURL) {
						team0.TeamProfileLink = secureDomain + t1PotentialRelativeURL
					} else {
						team0.TeamProfileLink = ""
						break
					}
				}
			})

			e.ForEach("td.team-right", func(i int, e *colly.HTMLElement) {
				teamRight := e.ChildAttr("span", "data-highlightingclass")
				if teamRight != "" {
					team1.FullName = strings.TrimSpace(teamRight)
				}

				// skip parsing scores if the match is not started yet
				if team1.FullName == "TBD" {
					return
				}

				team1.ShortName = e.ChildText("span.team-template-text")

				// the sequence of potential relative URLs is not always the same on each match
				// then better to check if the URL is valid or not
				t2PotentialRelativeURLs := e.ChildAttrs("a", "href")
				for _, t2PotentialRelativeURL := range t2PotentialRelativeURLs {
					if isValidTeamURL(t2PotentialRelativeURL) {
						team1.TeamProfileLink = secureDomain + t2PotentialRelativeURL
					} else {
						team1.TeamProfileLink = ""
						break
					}
				}
			})
		})

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

				team0.Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[0]))
				team1.Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[1]))
			}
		}

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

func parseLiveMatchDetailsPage(ctx context.Context, req *http.Request) ([]*api.LiveTeam, error) {
	// TODO: Implement
	return nil, nil
}

func parseTeamProfilePage(team *api.Team) colly.HTMLCallback {
	type playerTableSelector struct {
		tableSelector string
		activeStatus  api.PlayerStatus
	}

	schemas := []playerTableSelector{
		{
			activeStatus:  api.Active,
			tableSelector: "h3:has(span#Active) + div.table-responsive > table.roster-card tr.Player",
		},
		{
			activeStatus:  api.Inactive,
			tableSelector: "h3:has(span#Inactive) + div.table-responsive > table.roster-card tr.Player",
		},
		{
			activeStatus: api.Former,
			// Only take the active former table, because there are many inactive former player tables
			tableSelector: "h3:has(span#Former) + div.active .table-responsive > table.roster-card tr.Player",
		},
		{
			activeStatus:  api.StandIn,
			tableSelector: "h3:has(span#StandIn) + div.table-responsive > table.roster-card tr.Player",
		},
	}

	return func(h *colly.HTMLElement) {
		team.FullName = h.ChildText("h1#firstHeading span")

		for _, pps := range schemas {
			h.ForEach(pps.tableSelector, func(_ int, h *colly.HTMLElement) {
				team.PlayerRoster = append(team.PlayerRoster, &api.Player{
					ID:   h.ChildText("td.ID a"),
					Name: h.ChildText("td.Name"),
					Position: func() api.Position {
						rawP := h.ChildText("td.Position")
						if rawP == "" {
							return api.PosUnknown
						}

						rawP = rawP[len(rawP)-1:]
						p, _ := strconv.ParseInt(rawP, 10, 64)
						return api.Position(p)
					}(),
					JoinDate:       sanitizeDateOfPlayerRosterTable(h, "td.Position + td.Date i"),
					LeaveDate:      sanitizeDateOfPlayerRosterTable(h, "td.Date + td.Date i"),
					ActiveStatus:   pps.activeStatus,
					ProfilePageURL: secureDomain + h.ChildAttr("td.ID a", "href"),
				})
			})
		}
	}
}

func parseTournamentPage(ctx colly.Context, req *http.Request) (*api.Tournament, error) {
	// TODO: Implement
	return nil, nil
}
