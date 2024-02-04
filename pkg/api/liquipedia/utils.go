package liquipedia

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vuon9/d2m/pkg/api/model"
)

var secureDomain = "https://liquipedia.net"

var allowedDomains = []string{
	"liquipedia.net",
	"www.liquipedia.net",
}

type CrawData interface {
	[]*model.Match | *model.Team
}

func crawl[T CrawData](req *http.Request, parser PageParser[T]) (T, error) {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		// colly.CacheDir("./_cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		for k, v := range req.Header.Clone() {
			r.Headers.Set(k, v[0])
		}
	})

	var em T
	c.OnHTML(parser.RootSelector(), parser.Parse())

	err := c.Visit(req.URL.String())
	if err != nil {
		return em, err
	}

	return parser.Result()
}

func isValidTeamURL(potentialURL string) bool {
	return !strings.Contains(potentialURL, "/dota2/index.php?title=")
}

// buildStreamPageLink builds a link to the stream page on liquipedia
func buildStreamPageLink(channelName string) string {
	return fmt.Sprintf("%s/dota2/Special:Stream/twitch/%s", secureDomain, channelName)
}

// Remove ref link's text [1] [11] by separating the string by "[" and taking the first element
func sanitizeDateOfPlayerRosterTable(h *colly.HTMLElement, selector string) string {
	return strings.Split(h.ChildText(selector), "[")[0]
}
