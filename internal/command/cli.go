package command

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/vuon9/d2m/pkg/esporthub"
)

// GetCLIMatches prints matches as table on terminal
func GetCLIMatches(ctx context.Context, gameName esporthub.GameName) error {
	matchesByDate, err := GetMatchesByDate(ctx, gameName)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "Status", "Team 1", "Team 2", "Score"})
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
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for date, matches := range matchesByDate {
		for _, match := range matches {
			tableRow := []string{
				date.Format("2006-01-02 15:04"),
				match.Status(),
				match.Team1().FullName,
				match.Team2().FullName,
				fmt.Sprintf("%d - %d", match.Team1().Score, match.Team2().Score),
			}
			table.Append(tableRow)
		}
	}

	table.Render()

	return nil
}
