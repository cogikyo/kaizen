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

	fmt.Println()
	ui.Println(ui.Bold, fmt.Sprintf("%d problems due for review", len(due)))
	fmt.Println()

	for _, p := range due {
		kyuStr := ""
		if p.Kyu != nil {
			kyuStr = fmt.Sprintf(" %s[%d-%s]%s", ui.Dim, *p.Kyu, kyuNames[*p.Kyu-1], ui.Reset)
		}
		fmt.Printf("  %s%s\n", p.Path, kyuStr)
	}

	fmt.Println()
	ui.Hint("kz solve <problem>")
	fmt.Println()

	return nil
}
