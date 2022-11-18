package esporthub

import (
	"errors"
	"regexp"
	"strings"
)

// parseCredentials parses the credentials from the given html content.
func parseCredentials(str string) (*EsportHubClient, error) {
	re := regexp.MustCompile(`(?:clientId|hubSubscriptionKey)+:"+([a-zA-Z0-9]+)"`)

	var cre EsportHubClient
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		if strings.HasPrefix(match[0], "clientId") {
			cre.ClientID = match[1]
		}

		if strings.HasPrefix(match[0], "hubSubscriptionKey") {
			cre.HubSubscriptionKey = match[1]
		}
	}

	if (cre == EsportHubClient{}) {
		return nil, errors.New("couldn't find credentials for esporthub client")
	}

	return &cre, nil
}
