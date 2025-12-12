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

	left := ui.Colorize("改", ui.Yellow+ui.Bold) + " " + ui.Colorize(status.StartDate.Format("06-01-02"), ui.Dim)

	right := ui.Colorize(fmt.Sprintf("%d", stats.CurrentStreak), ui.Bold+ui.Green) + " " +
		ui.Colorize("day streak", ui.Green) + "  " +
		ui.Colorize(fmt.Sprintf("%d", stats.LongestStreak), ui.Bold+ui.Yellow) + " " +
		ui.Colorize("longest", ui.Yellow)

	ui.Justified(left, right, w)

	center := ui.Colorize(fmt.Sprintf("%d", stats.TotalSessions), ui.Bold+ui.Blue) + " " +
		ui.Colorize("sessions", ui.Blue) + "  " +
		ui.Colorize(fmt.Sprintf("%d", stats.UniqueProblems), ui.Bold+ui.Cyan) + " " +
		ui.Colorize("problems", ui.Cyan)

	centerLen := ui.VisibleLen(center)
	leftPad := (w - centerLen) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", leftPad), center)
}

func printRecent() {
	today, _ := db.GetTodaySessions()
	now := time.Now()

	if len(today) == 0 {
		fmt.Printf("  %s — %s\n",
			ui.Colorize(now.Format("Monday"), ui.Dim),
			ui.Colorize("no practice yet", ui.Dim))
		return
	}

	fmt.Printf("  %s %s %s %s\n",
		ui.Colorize(now.Format("Monday"), ui.Cyan),
		ui.Colorize("—", ui.Dim),
		ui.Colorize(fmt.Sprintf("%d", len(today)), ui.Bold+ui.Green),
		ui.Colorize("practiced", ui.Green))

	for _, p := range today {
		fmt.Printf("    %s %s\n", ui.Colorize("", ui.Green), ui.Colorize(p, ui.Dim))
	}
}

func printNext() {
	dueCount, _ := db.GetDueCount()
	p, _ := db.GetRandomProblem("", nil, "", true)

	if p != nil {
		if dueCount > 0 {
			ui.Info(fmt.Sprintf("kz solve %s %s",
				ui.Colorize(p.Name, ui.Yellow),
				ui.Colorize(fmt.Sprintf("(%d due)", dueCount), ui.Dim)))
		} else {
			ui.Info("kz solve " + ui.Colorize(p.Name, ui.Yellow))
		}
	} else {
		ui.Info("kz new " + ui.Colorize("<problem-name>", ui.Yellow))
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
