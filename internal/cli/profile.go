package cli

import (
	"fmt"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ProfileCmd struct{}

func (c *ProfileCmd) Run() error {
	fmt.Println()
	printHeader()
	fmt.Println()

	stats, err := db.GetStats()
	if err != nil {
		return err
	}

	fmt.Printf("  %s%d%s sessions%s  %s%d%s problems%s  ",
		ui.Bold+ui.Blue, stats.TotalSessions, ui.Reset+ui.Blue, ui.Reset,
		ui.Bold+ui.Cyan, stats.UniqueProblems, ui.Reset+ui.Cyan, ui.Reset)
	fmt.Printf("%s%d%s day streak%s  %s%d%s longest%s\n",
		ui.Bold+ui.Green, stats.CurrentStreak, ui.Reset+ui.Green, ui.Reset,
		ui.Bold+ui.Yellow, stats.LongestStreak, ui.Reset+ui.Yellow, ui.Reset)

	dueCount, _ := db.GetDueCount()
	if dueCount > 0 {
		fmt.Printf("  %s%d%s reviews due\n", ui.Yellow, dueCount, ui.Reset)
	}

	fmt.Println()
	ui.RenderHistory()

	today, _ := db.GetTodaySessions()
	if len(today) > 0 {
		fmt.Printf("\n%stoday%s\n", ui.Dim, ui.Reset)
		for _, p := range today {
			fmt.Printf("  %s✓%s %s\n", ui.Green, ui.Reset, p)
		}
	}

	fmt.Println()
	return nil
}

func printHeader() {
	p, _ := db.GetRandomProblem("", nil, "", true)
	if p != nil {
		fmt.Printf("%skaizen%s %skz solve %s%s\n", ui.Yellow+ui.Bold, ui.Reset, ui.Dim, p.Name, ui.Reset)
	} else {
		fmt.Printf("%skaizen%s %sno problem found, use 'kz new' to create one%s\n", ui.Yellow+ui.Bold, ui.Reset, ui.Dim, ui.Reset)
	}
}

func formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%dm", seconds/60)
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	return fmt.Sprintf("%dh%dm", h, m)
}
