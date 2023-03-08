package api

import (
	"fmt"
	"strings"
	"time"
)

type Match struct {
	Start                      time.Time  `json:"start"`
	Tournament                 Tournament `json:"tournament"`
	Tier                       string     `json:"tier"`
	CompetitionType            string     `json:"competitionType"`
	CompetitionTypeDescription string     `json:"competitionTypeDescription"`
	ContentType                string     `json:"contentType"`
	Name                       string     `json:"name"`
	StatusDescription          string     `json:"statusDescription"`
	HasStreamingURL            bool       `json:"hasStreamingURL"`
	StreamingURL               string     `json:"streamingURL"`
	ID                         string     `json:"id"`
	UrlsDescriptions           struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
	Teams  []*Team     `json:"teams"`
	Status MatchStatus `json:"status"`
}

var defaultTeam = &Team{
	FullName: "TBD",
}

func (m MatchStatus) String() string {
	return matchStatuses[m]
}

func (m *Match) Team1() *Team {
	if len(m.Teams) > 0 {
		return m.Teams[0]
	}

	return defaultTeam
}

func (m *Match) Team2() *Team {
	if len(m.Teams) > 1 {
		return m.Teams[1]
	}

	return defaultTeam
}

var (
	defaultTemplate  = "[%d:%d] - %s"
	hasStreamingIcon = "\u003e\u003e"
	hasTeamProfile   = "\u2b1f"
	hasNoTeamProfile = "\u25a2"
)

func (m *Match) Title() string {
	var typeAndScores string
	tmp := defaultTemplate

	switch m.Status {
	case StatusLive:
		if m.HasStreamingURL {
			tmp = "[%d:%d] - " + hasStreamingIcon + " %s"
		}

		fallthrough
	case StatusFinished:
		typeAndScores = fmt.Sprintf(tmp, m.Team1().Score, m.Team2().Score, m.Status)
	case StatusComing:
		fallthrough
	default:
		typeAndScores = fmt.Sprintf("%s", m.Status)
	}

	return fmt.Sprintf("%s - %s", typeAndScores, m.GeneralTitle())
}

func (m *Match) Description() string {
	return fmt.Sprintf("[%s] - %s", m.Start.Format("2006-01-02 15:04"), m.Tournament.Name)
}

// GeneralTitle uses for filtering only
func (m *Match) GeneralTitle() string {
	team1Name := hasTeamProfile
	if m.Team1().TeamProfileLink == "" {
		team1Name = hasNoTeamProfile
	}

	team2Name := hasTeamProfile
	if m.Team2().TeamProfileLink == "" {
		team2Name = hasNoTeamProfile
	}

	team1Name += " " + m.Team1().FullName
	team2Name += " " + m.Team2().FullName

	vs := fmt.Sprintf("%s vs. %s", team1Name, team2Name)
	if m.CompetitionType != "" {
		vs = fmt.Sprintf("%s (%s)", vs, m.CompetitionType)
	}

	return vs
}

func (m *Match) FilterValue() string {
	return strings.Join([]string{
		m.Status.String(),
		m.GeneralTitle(),
		m.Tournament.Name,
		m.Description(),
	}, " ")
}
