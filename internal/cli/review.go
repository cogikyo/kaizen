package cli

import (
	"fmt"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ReviewCmd struct{}

func (c *ReviewCmd) Run() error {
	due, err := db.GetDueReviews()
	if err != nil {
		return err
	}

	if len(due) == 0 {
		ui.Success("no reviews due today")
		return nil
	}

	fmt.Printf("\n%s%d problems due for review%s\n\n", ui.Bold, len(due), ui.Reset)

	for _, p := range due {
		kyuStr := ""
		if p.Kyu != nil {
			kyuStr = fmt.Sprintf(" %s[%d-%s]%s", ui.Dim, *p.Kyu, kyuNames[*p.Kyu-1], ui.Reset)
		}
		fmt.Printf("  %s%s%s\n", p.Path, kyuStr, ui.Reset)
	}

	fmt.Printf("\n%skaizen solve <problem>%s\n\n", ui.Dim, ui.Reset)

	return nil
}
