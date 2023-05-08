package liquipedia

import (
	"context"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/pkg/api"
)

var (
	defaultHeaders http.Header = map[string][]string{
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"},
		"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"Accept-Language": {"en-US,en;q=0.9,vi;q=0.8"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Cache-Control":   {"max-age=0"},
	}

	upComingPageUrl = secureDomain + "/dota2/Liquipedia:Upcoming_and_ongoing_matches"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (cre *Client) GetScheduledMatches(ctx context.Context) ([]*api.Match, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, upComingPageUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header = defaultHeaders

	matches := make([]*api.Match, 0)
	err = crawl(req, "div.matches-list > div:nth-child(2) table.infobox_matches_content > tbody", parseUpComingMatchesPage(&matches))
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Start.After(matches[j].Start) {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	return matches, nil
}

func (cre *Client) GetTeamDetailsPage(ctx context.Context, url string) (*api.Team, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = defaultHeaders

	team := new(api.Team)
	team.TeamProfileLink = url

	err = crawl(req, "body", func(h *colly.HTMLElement) {
		team.FullName = h.ChildText("h1#firstHeading span")
		team.PlayerRoster = parseTeamProfilePage(h)
	})

	if err != nil {
		return nil, err
	}

	return team, nil
}
