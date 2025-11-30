package cli

import (
	"fmt"
	"strings"
	"time"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ProfileCmd struct{}

func (c *ProfileCmd) Run() error {
	stats, _ := db.GetStats()
	status := ui.BuildStatus(ui.DefaultStartMonth())

	fmt.Println()
	printHeader(stats, status)
	fmt.Println()
	ui.RenderTable(status)
	fmt.Println()
	printRecent()
	fmt.Println()
	printNext()
	fmt.Println()

	return nil
}

func printHeader(stats *db.Stats, status *ui.StatusData) {
	w := status.TableWidth

	left := fmt.Sprintf("%s改%s %s%s%s",
		ui.Yellow+ui.Bold, ui.Reset,
		ui.Dim, status.StartDate.Format("06-01-02"), ui.Reset)

	right := fmt.Sprintf("%s%d%s %sday streak%s  %s%d%s %slongest%s",
		ui.Bold+ui.Green, stats.CurrentStreak, ui.Reset,
		ui.Green, ui.Reset,
		ui.Bold+ui.Yellow, stats.LongestStreak, ui.Reset,
		ui.Yellow, ui.Reset)

	leftLen := ui.VisibleLen(left)
	rightLen := ui.VisibleLen(right)
	padding := w - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}
	fmt.Printf("%s%*s%s\n", left, padding, "", right)

	center := fmt.Sprintf("%s%d%s %ssessions%s  %s%d%s %sproblems%s",
		ui.Bold+ui.Blue, stats.TotalSessions, ui.Reset,
		ui.Blue, ui.Reset,
		ui.Bold+ui.Cyan, stats.UniqueProblems, ui.Reset,
		ui.Cyan, ui.Reset)

	centerLen := ui.VisibleLen(center)
	leftPad := (w - centerLen) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", leftPad), center)
}

func printRecent() {
	today, _ := db.GetTodaySessions()
	now := time.Now()

	if len(today) == 0 {
		fmt.Printf("  %s%s%s %s— no practice yet%s\n",
			ui.Dim, now.Format("Monday"), ui.Reset,
			ui.Dim, ui.Reset)
		return
	}

	fmt.Printf("  %s%s%s %s—%s %s%d%s %spracticed%s\n",
		ui.Cyan, now.Format("Monday"), ui.Reset,
		ui.Dim, ui.Reset,
		ui.Bold+ui.Green, len(today), ui.Reset,
		ui.Green, ui.Reset)

	for _, p := range today {
		fmt.Printf("    %s%s %s%s%s\n", ui.Green, ui.Reset, ui.Dim, p, ui.Reset)
	}
}

func printNext() {
	dueCount, _ := db.GetDueCount()
	p, _ := db.GetRandomProblem("", nil, "", true)

	if p != nil {
		if dueCount > 0 {
			ui.Info(fmt.Sprintf("kz solve %s%s%s %s(%d due)%s",
				ui.Yellow, p.Name, ui.Reset+ui.Dim,
				ui.Dim, dueCount, ui.Reset))
		} else {
			ui.Info(fmt.Sprintf("kz solve %s%s%s",
				ui.Yellow, p.Name, ui.Reset))
		}
	} else {
		ui.Info(fmt.Sprintf("kz new %s<problem-name>%s",
			ui.Yellow, ui.Reset))
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
