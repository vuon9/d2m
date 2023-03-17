package liquipedia

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
)

var secureDomain = "https://liquipedia.net"

var allowedDomains = []string{
	"liquipedia.net",
	"www.liquipedia.net",
}

func crawl(req *http.Request, rootSelector string, onHTMLFunc colly.HTMLCallback) error {
	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		// colly.CacheDir("./_cache"),
	)

	c.OnRequest(func(r *colly.Request) {
		for k, v := range req.Header.Clone() {
			r.Headers.Set(k, v[0])
		}
	})

	c.OnHTML(rootSelector, onHTMLFunc)

	err := c.Visit(req.URL.String())
	if err != nil {
		return err
	}

	return nil
}

func isValidTeamURL(potentialURL string) bool {
	return !strings.Contains(potentialURL, "/dota2/index.php?title=")
}

// buildStreamPageLink builds a link to the stream page on liquipedia
func buildStreamPageLink(channelName string) string {
	return fmt.Sprintf("%s/dota2/Special:Stream/twitch/%s", secureDomain, channelName)
}
