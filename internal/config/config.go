package config

import (
	"github.com/vuon9/d2m/internal/esporthub"
)

type Config struct {
	MatchCredentials esporthub.MatchAPICredentials `tom:"match_credentials"`
}
