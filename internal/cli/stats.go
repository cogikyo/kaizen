package cli

import (
	"fmt"
	"sort"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type StatsCmd struct{}

func (c *StatsCmd) Run() error {
	fmt.Println()

	stats, _ := db.GetStats()

	ui.Section("activity")
	fmt.Printf("  %s sessions  %s attempts  %s problems practiced\n",
		ui.Colorize(fmt.Sprintf("%d", stats.TotalSessions), ui.Blue),
		ui.Colorize(fmt.Sprintf("%d", stats.TotalAttempts), ui.Cyan),
		ui.Colorize(fmt.Sprintf("%d", stats.UniqueProblems), ui.Green))
	fmt.Printf("  %s passed  %s failed  %s pass rate\n",
		ui.Colorize(fmt.Sprintf("%d", stats.TotalPassed), ui.Green),
		ui.Colorize(fmt.Sprintf("%d", stats.TotalFailed), ui.Red),
		ui.Colorize(fmt.Sprintf("%.0f%%", stats.PassRate), ui.Dim))
	if stats.TotalTime > 0 {
		fmt.Printf("  %s total time  %s avg/session\n",
			formatDuration(stats.TotalTime),
			formatDuration(stats.AvgSessionTime))
	}

	fmt.Println()
	ui.Section("streak")
	fmt.Printf("  %s current  %s longest\n",
		ui.Colorize(fmt.Sprintf("%d", stats.CurrentStreak), ui.Green),
		ui.Colorize(fmt.Sprintf("%d", stats.LongestStreak), ui.Yellow))

	problems, _ := db.GetProblems("", nil, "")

	fmt.Println()
	ui.Section("by kyu")
	kyuCounts := make(map[int]int)
	maxKyu := 0
	for _, p := range problems {
		if p.Kyu != nil {
			kyuCounts[*p.Kyu]++
			if kyuCounts[*p.Kyu] > maxKyu {
				maxKyu = kyuCounts[*p.Kyu]
			}
		} else {
			kyuCounts[0]++
		}
	}
	if maxKyu == 0 {
		maxKyu = 1
	}
	for k := 1; k <= 6; k++ {
		cnt := kyuCounts[k]
		barLen := (cnt * 20) / maxKyu
		if cnt > 0 && barLen < 1 {
			barLen = 1
		}
		fmt.Printf("  %s %s%-8s%s %s %s\n",
			ui.Colorize(fmt.Sprintf("%d", k), ui.Yellow),
			ui.Dim, kyuNames[k-1], ui.Reset,
			ui.Bar(barLen, ui.Cyan),
			ui.Colorize(fmt.Sprintf("(%d)", cnt), ui.Dim))
	}

	fmt.Println()
	ui.Section("by section")
	sectionCounts := make(map[string]int)
	maxSection := 0
	for _, p := range problems {
		sectionCounts[p.Section]++
		if sectionCounts[p.Section] > maxSection {
			maxSection = sectionCounts[p.Section]
		}
	}
	if maxSection == 0 {
		maxSection = 1
	}
	for s, cnt := range sectionCounts {
		barLen := (cnt * 20) / maxSection
		if cnt > 0 && barLen < 1 {
			barLen = 1
		}
		fmt.Printf("  %-12s %s %s\n", s, ui.Bar(barLen, ui.Blue), ui.Colorize(fmt.Sprintf("(%d)", cnt), ui.Dim))
	}

	tags, _ := db.GetTags()
	if len(tags) > 0 {
		fmt.Println()
		ui.Section("by tag")
		type tagCount struct {
			tag   string
			count int
		}
		var sorted []tagCount
		maxTag := 0
		for t, cnt := range tags {
			sorted = append(sorted, tagCount{t, cnt})
			if cnt > maxTag {
				maxTag = cnt
			}
		}
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].count > sorted[j].count
		})
		if maxTag == 0 {
			maxTag = 1
		}
		for _, tc := range sorted {
			barLen := (tc.count * 20) / maxTag
			if tc.count > 0 && barLen < 1 {
				barLen = 1
			}
			fmt.Printf("  %-12s %s %s\n", tc.tag, ui.Bar(barLen, ui.Green), ui.Colorize(fmt.Sprintf("(%d)", tc.count), ui.Dim))
		}
	}

	fmt.Println()
	return nil
}
