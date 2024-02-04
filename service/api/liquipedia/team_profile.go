package liquipedia

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/service/api/models"
)

type teamProfilePageParser struct {
	req          *http.Request
	rootSelector string
	team         *models.Team
}

func NewTeamProfilePageParser() *teamProfilePageParser {
	return &teamProfilePageParser{
		rootSelector: "body",
		team:         new(models.Team),
	}
}

func (p *teamProfilePageParser) RootSelector() string {
	return p.rootSelector
}

func (p *teamProfilePageParser) Result() (*models.Team, error) {
	return p.team, nil
}

func (p *teamProfilePageParser) Parse() colly.HTMLCallback {
	return func(e *colly.HTMLElement) {
		p.team.TeamProfileLink = e.Request.URL.String()
		p.team.FullName = e.ChildText("h1#firstHeading .mw-page-title-main")

		type playerTableSelector struct {
			tableSelector string
			activeStatus  models.PlayerStatus
		}

		schemas := []playerTableSelector{
			{
				activeStatus:  models.Active,
				tableSelector: "h3:has(span#Active_Roster) + div.table-responsive > table.roster-card tr.Player",
			},
			{
				activeStatus:  models.Active,
				tableSelector: "h3:has(span#Active) + div.table-responsive > table.roster-card tr.Player",
			},
		}

		p.team.PlayerRoster = make([]*models.Player, 0)
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

func (p *teamProfilePageParser) parsePlayerRoster(h *colly.HTMLElement, s models.PlayerStatus) *models.Player {
	id := h.ChildText("td.ID")
	position := h.ChildText("td.Position")

	if strings.TrimSpace(id) == "" || strings.TrimSpace(position) == "" {
		return nil
	}

	return &models.Player{
		ID:   id,
		Name: h.ChildText("td.Name"),
		Position: func() models.Position {
			rawP := position
			if rawP == "" {
				return models.PosUnknown
			}

			rawP = rawP[len(rawP)-1:]
			p, _ := strconv.Atoi(rawP)

			return models.Position(p)
		}(),
		JoinDate:       sanitizeDateOfPlayerRosterTable(h, "td.Position + td.Date i"),
		LeaveDate:      sanitizeDateOfPlayerRosterTable(h, "td.Date + td.Date i"),
		ActiveStatus:   s,
		ProfilePageURL: secureDomain + h.ChildAttr("td.ID a", "href"),
	}
}
