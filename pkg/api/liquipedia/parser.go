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

		e.ForEach("tr", func(_ int, h *colly.HTMLElement) {
			h.ForEach("td.team-left", parseTeam(team0))
			h.ForEach("td.team-right", parseTeam(team1))
		})

		parseMatchStateAndScores(e, match)

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

func parseTeam(team *api.Team) func(i int, e *colly.HTMLElement) {
	return func(_ int, e *colly.HTMLElement) {
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
}

func parseMatchStateAndScores(e *colly.HTMLElement, match *api.Match) {
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

			match.Team1().Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[0]))
			match.Team2().Score, _ = strconv.Atoi(strings.TrimSpace(rawScores[1]))
		}
	}
}

func parseLiveMatchDetailsPage(ctx context.Context, req *http.Request) ([]*api.LiveTeam, error) {
	// TODO: Implement
	return nil, nil
}

type playerTableSelector struct {
	tableSelector string
	activeStatus  api.PlayerStatus
}

func parseTeamProfilePage(team *api.Team) colly.HTMLCallback {
	// TODO: I need to answer to myself, why I need to use a slice of playerTableSelector instead of just one selector?
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
			h.ForEach(pps.tableSelector, parsePlayerRoster(pps.activeStatus, team))
		}
	}
}

func parsePlayerRoster(playerStatus api.PlayerStatus, team *api.Team) func (_ int, h *colly.HTMLElement) {
	return func (_ int, h *colly.HTMLElement) {
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
			ActiveStatus:   playerStatus,
			ProfilePageURL: secureDomain + h.ChildAttr("td.ID a", "href"),
		})
	}
}

func parseTournamentPage(ctx colly.Context, req *http.Request) (*api.Tournament, error) {
	// TODO: Implement
	return nil, nil
}
