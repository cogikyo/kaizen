package info

import (
	"fmt"
	"time"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ProfileCmd struct{}

func (c *ProfileCmd) Run() error {
	stats, _ := db.GetStats()
	status := ui.BuildStatus()

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
	w := status.TableWidth - 4
	activity := ui.Primary("sessions", stats.TotalSessions) + "  " + ui.Count("problems", stats.UniqueProblems)
	dateRange := ui.InlineTitle(status.StartDate.Format("06-01-02") + " — " + status.EndDate.Format("06-01-02"))
	streak := ui.Positive("day streak", stats.CurrentStreak) + "  " + ui.Accent("longest", stats.LongestStreak)
	fmt.Print("    ")
	ui.Justified(w, activity, dateRange, streak)
}

func printRecent() {
	today, _ := db.GetTodaySessions()
	now := time.Now()

	if len(today) == 0 {
		fmt.Printf("  %s — %s\n",
			ui.Muted(now.Format("Monday")),
			ui.Muted("no practice yet"))
		return
	}

	fmt.Printf("  %s %s %s\n",
		ui.Count(now.Format("Monday")),
		ui.Muted("—"),
		ui.Positive("practiced", len(today)))

	for _, p := range today {
		fmt.Printf("    %s %s\n", ui.Positive(""), ui.Muted(p))
	}
}

func printNext() {
	dueCount, _ := db.GetDueCount()
	p, _ := db.GetRandomProblem("", nil, "", true)

	if p != nil {
		if dueCount > 0 {
			ui.Info(fmt.Sprintf("kz solve %s %s",
				ui.Accent(p.Name),
				ui.Muted(fmt.Sprintf("(%d due)", dueCount))))
		} else {
			ui.Info("kz solve " + ui.Accent(p.Name))
		}
	} else {
		ui.Info("kz new " + ui.Accent("<problem-name>"))
	}
}

