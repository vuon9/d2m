package liquipedia

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/pkg/api/model"
)

type teamProfilePageParser struct {
	req          *http.Request
	rootSelector string
	team         *model.Team
}

func NewTeamProfilePageParser() *teamProfilePageParser {
	return &teamProfilePageParser{
		rootSelector: "body",
		team:         new(model.Team),
	}
}

func (p *teamProfilePageParser) RootSelector() string {
	return p.rootSelector
}

func (p *teamProfilePageParser) Result() (*model.Team, error) {
	return p.team, nil
}

func (p *teamProfilePageParser) Parse() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		p.team.TeamProfileLink = e.Request.URL.String()
		p.team.FullName = e.ChildText("h1#firstHeading .mw-page-title-main")

		type playerTableSelector struct {
			tableSelector string
			activeStatus  model.PlayerStatus
		}

		schemas := []playerTableSelector{
			{
				activeStatus:  model.Active,
				tableSelector: "h3:has(span#Active_Roster) + div.table-responsive > table.roster-card tr.Player",
			},
			{
				activeStatus:  model.Active,
				tableSelector: "h3:has(span#Active) + div.table-responsive > table.roster-card tr.Player",
			},
		}

		p.team.PlayerRoster = make([]*model.Player, 0)
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

func (p *teamProfilePageParser) parsePlayerRoster(h *colly.HTMLElement, s model.PlayerStatus) *model.Player {
	id := h.ChildText("td.ID")
	position := h.ChildText("td.Position")

	if strings.TrimSpace(id) == "" || strings.TrimSpace(position) == "" {
		return nil
	}

	return &model.Player{
		ID:   id,
		Name: h.ChildText("td.Name"),
		Position: func() model.Position {
			rawP := position
			if rawP == "" {
				return model.PosUnknown
			}

			rawP = rawP[len(rawP)-1:]
			p, _ := strconv.Atoi(rawP)

			return model.Position(p)
		}(),
		JoinDate:       sanitizeDateOfPlayerRosterTable(h, "td.Position + td.Date i"),
		LeaveDate:      sanitizeDateOfPlayerRosterTable(h, "td.Date + td.Date i"),
		ActiveStatus:   s,
		ProfilePageURL: secureDomain + h.ChildAttr("td.ID a", "href"),
	}
}
