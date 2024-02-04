package model

import (
	"fmt"
	"strings"
	"time"
)

type Match struct {
	Start                      time.Time   `json:"start"`
	CompetitionType            string      `json:"competitionType"`
	IsConcludedMatch           bool        `json:"isConcludedMatch"`
	Status                     MatchStatus `json:"status"`
	Tournament                 Tournament  `json:"tournament"`
	Tier                       string      `json:"tier,omitempty"`
	CompetitionTypeDescription string      `json:"competitionTypeDescription,omitempty"`
	ContentType                string      `json:"contentType,omitempty"`
	Name                       string      `json:"name,omitempty"`
	StatusDescription          string      `json:"statusDescription,omitempty"`
	StreamingURL               string      `json:"streamingURL,omitempty"`
	ID                         string      `json:"id"`
	UrlsDescriptions           struct {
		Logo string `json:"logo,omitempty"`
	} `json:"urlsDescriptions"`
	Teams         []*Team  `json:"teams"`
	VideoOnDemand []string `json:"videoOnDemand"`
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
	hasStreamingIcon   = "\ufa42"
	hasNoStreamingIcon = "\u25b7"
	hasTeamProfile     = "\u25c6"
	hasNoTeamProfile   = "\u25c7"
)

func (m *Match) Title() string {
	var typeAndScores string
	tmp := defaultTemplate

	switch m.Status {
	case StatusLive:
		prefix := hasNoStreamingIcon
		if m.StreamingURL != "" {
			prefix = hasStreamingIcon
		}

		tmp = fmt.Sprintf("%s %s", tmp, prefix)
		fallthrough
	case StatusFinished:
		typeAndScores = fmt.Sprintf(tmp, m.Team1().Score, m.Team2().Score, m.Status)
	case StatusComing:
		fallthrough
	default:
		typeAndScores = fmt.Sprintf("%s", m.Status)
	}

	return fmt.Sprintf("%s %s", typeAndScores, m.GeneralTitle())
}

func (m *Match) Description() string {
	return fmt.Sprintf("[%s] - %s", m.Start.Format("2006-01-02 15:04"), m.Tournament.Name)
}

// GeneralTitle uses for filtering only
func (m *Match) GeneralTitle() string {
	team1Name := m.Team1().FullName
	if team1Name == "" {
		team1Name = "TBD"
	}

	team2Name := m.Team2().FullName
	if team2Name == "" {
		team2Name = "TBD"
	}

	team1Icon := hasNoTeamProfile + m.Team1().TeamProfileLink
	if m.Team1().TeamProfileLink != "" {
		team1Icon = hasTeamProfile
	}

	team2Icon := hasNoTeamProfile + m.Team2().TeamProfileLink
	if m.Team2().TeamProfileLink != "" {
		team2Icon = hasTeamProfile
	}

	vs := fmt.Sprintf("%s %s vs. %s %s", team1Icon, team1Name, team2Icon, team2Name)
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
