package liquipedia

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/pkg/api"
)

type teamProfilePageParser struct {
	req          *http.Request
	rootSelector string
	team         *api.Team
}

func NewTeamProfilePageParser() *teamProfilePageParser {
	return &teamProfilePageParser{
		rootSelector: "body",
		team:         new(api.Team),
	}
}

func (p *teamProfilePageParser) RootSelector() string {
	return p.rootSelector
}

func (p *teamProfilePageParser) Result() (*api.Team, error) {
	return p.team, nil
}

func (p *teamProfilePageParser) Parse() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		p.team.TeamProfileLink = e.Request.URL.String()
		p.team.FullName = e.ChildText("h1#firstHeading .mw-page-title-main")

		type playerTableSelector struct {
			tableSelector string
			activeStatus  api.PlayerStatus
		}

		schemas := []playerTableSelector{
			{
				activeStatus:  api.Active,
				tableSelector: "h3:has(span#Active_Roster) + div.table-responsive > table.roster-card tr.Player",
			},
			{
				activeStatus:  api.Active,
				tableSelector: "h3:has(span#Active) + div.table-responsive > table.roster-card tr.Player",
			},
		}

		p.team.PlayerRoster = make([]*api.Player, 0)
		for _, schema := range schemas {
			e.ForEachWithBreak(schema.tableSelector, func(_ int, h *colly.HTMLElement) bool {
				player := p.parsePlayerRoster(h, schema.activeStatus)
				if player == nil {
					return false
				}

				p.team.PlayerRoster = append(p.team.PlayerRoster, player)
				return true
			})
		}
	}
}

func (p *teamProfilePageParser) parsePlayerRoster(h *colly.HTMLElement, s api.PlayerStatus) *api.Player {
	id := h.ChildText("td.ID")
	position := h.ChildText("td.Position")

	if strings.TrimSpace(id) == "" || strings.TrimSpace(position) == "" {
		return nil
	}

	return &api.Player{
		ID:   id,
		Name: h.ChildText("td.Name"),
		Position: func() api.Position {
			rawP := position
			if rawP == "" {
				return api.PosUnknown
			}

			rawP = rawP[len(rawP)-1:]
			p, _ := strconv.Atoi(rawP)

			return api.Position(p)
		}(),
		JoinDate:       sanitizeDateOfPlayerRosterTable(h, "td.Position + td.Date i"),
		LeaveDate:      sanitizeDateOfPlayerRosterTable(h, "td.Date + td.Date i"),
		ActiveStatus:   s,
		ProfilePageURL: secureDomain + h.ChildAttr("td.ID a", "href"),
	}
}
