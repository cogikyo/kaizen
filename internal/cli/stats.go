package cli

import (
	"fmt"
	"sort"
	"strings"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type StatsCmd struct{}

func (c *StatsCmd) Run() error {
	fmt.Println()

	stats, _ := db.GetStats()

	fmt.Printf("  %sactivity%s\n", ui.Dim, ui.Reset)
	fmt.Printf("  %s%d%s sessions  %s%d%s attempts  %s%d%s problems practiced\n",
		ui.Blue, stats.TotalSessions, ui.Reset,
		ui.Cyan, stats.TotalAttempts, ui.Reset,
		ui.Green, stats.UniqueProblems, ui.Reset)
	fmt.Printf("  %s%d%s passed  %s%d%s failed  %s%.0f%%%s pass rate\n",
		ui.Green, stats.TotalPassed, ui.Reset,
		ui.Red, stats.TotalFailed, ui.Reset,
		ui.Dim, stats.PassRate, ui.Reset)
	if stats.TotalTime > 0 {
		fmt.Printf("  %s total time  %s avg/session\n",
			formatDuration(stats.TotalTime),
			formatDuration(stats.AvgSessionTime))
	}

	fmt.Printf("\n  %sstreak%s\n", ui.Dim, ui.Reset)
	fmt.Printf("  %s%d%s current  %s%d%s longest\n",
		ui.Green, stats.CurrentStreak, ui.Reset,
		ui.Yellow, stats.LongestStreak, ui.Reset)

	problems, _ := db.GetProblems("", nil, "")

	fmt.Printf("\n  %sby kyu%s\n", ui.Dim, ui.Reset)
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
		bar := strings.Repeat("█", barLen)
		fmt.Printf("  %s%d%s %s%-8s%s %s%s%s %s(%d)%s\n",
			ui.Yellow, k, ui.Reset,
			ui.Dim, kyuNames[k-1], ui.Reset,
			ui.Cyan, bar, ui.Reset,
			ui.Dim, cnt, ui.Reset)
	}

	fmt.Printf("\n  %sby section%s\n", ui.Dim, ui.Reset)
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
		bar := strings.Repeat("█", barLen)
		fmt.Printf("  %-12s %s%s%s %s(%d)%s\n", s, ui.Blue, bar, ui.Reset, ui.Dim, cnt, ui.Reset)
	}

	tags, _ := db.GetTags()
	if len(tags) > 0 {
		fmt.Printf("\n  %sby tag%s\n", ui.Dim, ui.Reset)
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
			bar := strings.Repeat("█", barLen)
			fmt.Printf("  %-12s %s%s%s %s(%d)%s\n", tc.tag, ui.Green, bar, ui.Reset, ui.Dim, tc.count, ui.Reset)
		}
	}

	fmt.Println()
	return nil
}
