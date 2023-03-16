package liquipedia

import (
	"context"
	"net/http"

	"github.com/vuon9/d2m/pkg/api"
)

var (
	upComingPageUrl = secureDomain + "/dota2/Liquipedia:Upcoming_and_ongoing_matches"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (cre *Client) GetScheduledMatches(ctx context.Context) ([]*api.Match, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, upComingPageUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9,vi;q=0.8")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cache-Control", "max-age=0")

	return parseUpComingPage(ctx, req)

}
