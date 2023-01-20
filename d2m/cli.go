package d2m

import (
	"context"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/vuon9/d2m/pkg/api/types"
)

type Matcher struct {
	status types.MatchStatus
}

func NewMatcher() *Matcher {
	return &Matcher{
		status: types.MatchStatusDefault,
	}
}

type MatcherOption func(*Matcher)

func WithMatchStatus(status types.MatchStatus) MatcherOption {
	return func(matcher *Matcher) {
		matcher.status = status
	}
}

// GetCLIMatches prints matches as table on terminal
func (m *Matcher) GetCLIMatches(ctx context.Context, options ...MatcherOption) error {
	matches, err := GetMatches(ctx, types.Dota2)
	if err != nil {
		return err
	}

	for _, mo := range options {
		mo(m)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Tournament", "Team 1", "vs.", "Team 2", "Status"})
	table.SetBorder(true)
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetCenterSeparator(" ")
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
		tablewriter.ALIGN_LEFT,
	})

	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(false)
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)

	prev5Hours := time.Now().Add(-24 * time.Hour)

	for _, match := range matches {
		frStatus := match.FriendlyStatus()
		if frStatus != m.status {
			continue
		}

		if match.Start.Before(prev5Hours) {
			continue
		}

		tableRow := []string{
			match.Start.Format("2006-01-02 15:04"),
			match.Tournament.Name,
			match.Team1().FullName,
			match.CompetitionType,
			match.Team2().FullName,
			frStatus.String(),
		}
		table.Append(tableRow)
	}

	table.Render()

	// print match as beautiful JSON
	// matchJSON, err := json.MarshalIndent(match, "", "  ")
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(string(matchJSON))

	return nil
}
