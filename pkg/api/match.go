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
	StreamingURL               string     `json:"streamingURL"`
	ID                         string     `json:"id"`
	UrlsDescriptions           struct {
		Logo string `json:"logo"`
	} `json:"urlsDescriptions"`
	Teams            []*Team     `json:"teams"`
	IsConcludedMatch bool        `json:"isConcludedMatch"`
	Status           MatchStatus `json:"status"`
	VideoOnDemand    []string    `json:"videoOnDemand"`
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
	defaultTemplate    = "[%d:%d] %s"
	hasStreamingIcon   = "\u25b6"
	hasNoStreamingIcon = "\u25b7"
	hasTeamProfile     = "\u25c6"
	hasNoTeamProfile   = "\u25c7"
)

func (m *Match) Title() string {
	var typeAndScores string
	tmp := defaultTemplate

	switch m.Status {
	case StatusLive:
		tmp = "[%d:%d] " + hasNoStreamingIcon + " %s"
		if m.StreamingURL == "" {
			tmp = "[%d:%d] " + hasStreamingIcon + " %s"
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
		vs += fmt.Sprintf(" (%s)", m.CompetitionType)
	}

	return vs
}

func (m *Match) FilterValue() string {
	return strings.Join([]string{
		m.Team1().FullName,
		m.Team2().FullName,
		m.Tournament.Name,
	}, " ")
}
