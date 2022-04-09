package main

type MatchAPICredentials struct {
	ClientID           string
	HubSubscriptionKey string
}

type Config struct {
	MatchCredentials MatchAPICredentials `tom:"match_credentials"`
}

func retrieveConfig() *Config {
	return nil
}

func retrieveMatchAPICredentials() *MatchAPICredentials {
	return nil
}

func writeConfig() (bool, error) {
	return true, nil
}
