package liquipedia

import (
	"net/http"
	"strconv"
	"sync"

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
	team := new(api.Team)

	return func(e *colly.HTMLElement) {
		team.TeamProfileLink = e.Request.URL.String()
		team.FullName = e.ChildText("h1#firstHeading span")

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
				activeStatus:  api.Former,
				tableSelector: "h3:has(span#Former) + div.active .table-responsive > table.roster-card tr.Player",
			},
			{
				activeStatus:  api.StandIn,
				tableSelector: "h3:has(span#StandIn) + div.table-responsive > table.roster-card tr.Player",
			},
		}

		players := make([]*api.Player, 0)

		wg := sync.WaitGroup{}
		wg.Add(len(schemas))
		for _, schema := range schemas {
			go func(s playerTableSelector, pl *[]*api.Player) {
				e.ForEach(s.tableSelector, func(_ int, h *colly.HTMLElement) {
					*pl = append(*pl, p.parsePlayerRoster(h, s.activeStatus))
				})
				wg.Done()
			}(schema, &players)
		}

		wg.Wait()

		team.PlayerRoster = players

		p.team = team
	}
}

func (p *teamProfilePageParser) parsePlayerRoster(h *colly.HTMLElement, s api.PlayerStatus) *api.Player {
	return &api.Player{
		ID:   h.ChildText("td.ID a"),
		Name: h.ChildText("td.Name"),
		Position: func() api.Position {
			rawP := h.ChildText("td.Position")
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
