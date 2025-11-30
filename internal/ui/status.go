package ui

import (
	"fmt"
	"strings"
	"time"

	"cogikyo/kaizen/internal/db"
)

const cellWidth = 4

type StatusData struct {
	StartDate  time.Time
	Weeks      int
	TableWidth int
	Days       []int
	MaxCount   int
	Months     []string
}

func BuildStatus(startMonth time.Time) *StatusData {
	counts, _ := db.GetSessionCounts(400)

	countMap := make(map[string]int)
	maxCount := 0
	for _, c := range counts {
		countMap[c.Date] = c.Count
		if c.Count > maxCount {
			maxCount = c.Count
		}
	}

	now := time.Now()
	logicalStart := time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	startDate := logicalStart
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	const weeks = 29
	totalDays := weeks * 7

	days := make([]int, totalDays)
	for i := range totalDays {
		d := startDate.AddDate(0, 0, i)
		if d.After(now) {
			days[i] = -2
			continue
		}
		days[i] = countMap[d.Format("2006-01-02")]
	}

	months := make([]string, weeks)
	for w := range weeks {
		weekSunday := startDate.AddDate(0, 0, w*7)
		weekSaturday := weekSunday.AddDate(0, 0, 6)

		firstOfMonth := time.Date(weekSunday.Year(), weekSunday.Month(), 1, 0, 0, 0, 0, time.Local)
		if !firstOfMonth.Before(logicalStart) && !firstOfMonth.Before(weekSunday) && !firstOfMonth.After(weekSaturday) {
			months[w] = firstOfMonth.Format("Jan")
			continue
		}

		firstOfNext := firstOfMonth.AddDate(0, 1, 0)
		if !firstOfNext.Before(logicalStart) && !firstOfNext.Before(weekSunday) && !firstOfNext.After(weekSaturday) {
			months[w] = firstOfNext.Format("Jan")
		}
	}

	tableWidth := 4 + (weeks * (cellWidth + 1)) + 1

	return &StatusData{
		StartDate:  startDate,
		Weeks:      weeks,
		TableWidth: tableWidth,
		Days:       days,
		MaxCount:   maxCount,
		Months:     months,
	}
}

func RenderTable(s *StatusData) {
	weeks := s.Weeks

	fmt.Printf("    %s┌", Grey)
	fmt.Print(strings.Repeat(strings.Repeat("─", cellWidth)+"┬", weeks-1))
	fmt.Println(strings.Repeat("─", cellWidth) + "┐" + Reset)

	monthLine := "    " + Grey + "│" + Reset
	for w := range weeks {
		if s.Months[w] != "" {
			pad := cellWidth - len(s.Months[w])
			left := pad / 2
			right := pad - left
			monthLine += fmt.Sprintf("%s%s%s%s%s", strings.Repeat(" ", left), Dim, s.Months[w], Reset, strings.Repeat(" ", right))
		} else {
			monthLine += strings.Repeat(" ", cellWidth)
		}
		monthLine += Grey + "│" + Reset
	}
	fmt.Println(monthLine)

	fmt.Printf("    %s├", Grey)
	fmt.Print(strings.Repeat(strings.Repeat("─", cellWidth)+"┼", weeks-1))
	fmt.Println(strings.Repeat("─", cellWidth) + "┤" + Reset)

	dayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for row := range 7 {
		fmt.Printf("%s%s%s %s│%s", Dim, dayLabels[row], Reset, Grey, Reset)
		for w := range weeks {
			idx := w*7 + row
			if idx >= len(s.Days) {
				fmt.Printf("%s%s│%s", strings.Repeat(" ", cellWidth), Grey, Reset)
				continue
			}
			fmt.Printf("%s%s│%s", cellBlock(s.Days[idx], s.MaxCount), Grey, Reset)
		}
		fmt.Println()
	}

	fmt.Printf("    %s└", Grey)
	fmt.Print(strings.Repeat(strings.Repeat("─", cellWidth)+"┴", weeks-1))
	fmt.Println(strings.Repeat("─", cellWidth) + "┘" + Reset)
}

func RenderStatus(startMonth time.Time) {
	s := BuildStatus(startMonth)
	RenderTable(s)
}

func DefaultStartMonth() time.Time {
	now := time.Now()
	start := now.AddDate(0, 0, -200)
	return time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
}

func cellBlock(count, maxCount int) string {
	if count == -1 {
		return strings.Repeat(" ", cellWidth)
	}
	if count == -2 {
		return Grey + " ·· " + Reset
	}
	if count == 0 {
		return Black + " ░░ " + Reset
	}
	if maxCount == 0 {
		maxCount = 1
	}

	ratio := float64(count) / float64(maxCount)
	var color string
	switch {
	case ratio >= 1.0:
		color = Yellow
	case ratio >= 0.66:
		color = Cyan
	case ratio >= 0.33:
		color = Blue
	case ratio >= 0.20:
		color = Blue + Dim
	default:
		color = Dim
	}
	return color + " ██ " + Reset
}
