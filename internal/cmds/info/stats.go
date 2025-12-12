package info

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
	fmt.Printf("  %s  %s  %s\n",
		ui.Primary("sessions", stats.TotalSessions),
		ui.Count("attempts", stats.TotalAttempts),
		ui.Positive("problems practiced", stats.UniqueProblems))
	fmt.Printf("  %s  %s  %s\n",
		ui.Positive("passed", stats.TotalPassed),
		ui.Negative("failed", stats.TotalFailed),
		ui.Muted(fmt.Sprintf("%.0f%% pass rate", stats.PassRate)))
	if stats.TotalTime > 0 {
		fmt.Printf("  %s total time  %s avg/session\n",
			ui.FormatDuration(stats.TotalTime),
			ui.FormatDuration(stats.AvgSessionTime))
	}

	fmt.Println()
	ui.Section("streak")
	fmt.Printf("  %s  %s\n",
		ui.Positive("current", stats.CurrentStreak),
		ui.Accent("longest", stats.LongestStreak))

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
		fmt.Printf("  %s %s %s %s\n",
			ui.Accent(k),
			ui.Muted(fmt.Sprintf("%-8s", db.KyuNames[k-1])),
			ui.BarCount(barLen),
			ui.Muted(fmt.Sprintf("(%d)", cnt)))
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
		fmt.Printf("  %-12s %s %s\n", s, ui.BarPrimary(barLen), ui.Muted(fmt.Sprintf("(%d)", cnt)))
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
			fmt.Printf("  %-12s %s %s\n", tc.tag, ui.BarPositive(barLen), ui.Muted(fmt.Sprintf("(%d)", tc.count)))
		}
	}

	fmt.Println()
	return nil
}

