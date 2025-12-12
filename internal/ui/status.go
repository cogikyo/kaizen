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
	EndDate    time.Time
	Weeks      int
	TableWidth int
	Days       []int
	MaxCount   int
	Months     []string
}

func BuildStatus() *StatusData {
	termWidth := TermWidth()
	availableWidth := termWidth - 4 - 1
	weeks := availableWidth / (cellWidth + 1)
	if weeks > 52 {
		weeks = 52
	}
	if weeks < 4 {
		weeks = 4
	}

	now := time.Now()
	endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local)
	endDate := endOfMonth
	for endDate.Weekday() != time.Saturday {
		endDate = endDate.AddDate(0, 0, 1)
	}

	startDate := endDate.AddDate(0, 0, -(weeks*7 - 1))
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	counts, _ := db.GetSessionCounts(weeks * 7)
	countMap := make(map[string]int)
	maxCount := 0
	for _, c := range counts {
		countMap[c.Date] = c.Count
		if c.Count > maxCount {
			maxCount = c.Count
		}
	}

	days := make([]int, weeks*7)
	for i := range weeks * 7 {
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
		if !firstOfMonth.Before(startDate) && !firstOfMonth.Before(weekSunday) && !firstOfMonth.After(weekSaturday) {
			months[w] = firstOfMonth.Format("Jan")
			continue
		}

		firstOfNext := firstOfMonth.AddDate(0, 1, 0)
		if !firstOfNext.Before(startDate) && !firstOfNext.Before(weekSunday) && !firstOfNext.After(weekSaturday) {
			months[w] = firstOfNext.Format("Jan")
		}
	}

	tableWidth := 4 + (weeks * (cellWidth + 1)) + 1

	return &StatusData{
		StartDate:  startDate,
		EndDate:    endDate,
		Weeks:      weeks,
		TableWidth: tableWidth,
		Days:       days,
		MaxCount:   maxCount,
		Months:     months,
	}
}

func RenderTable(s *StatusData) {
	weeks := s.Weeks
	now := time.Now()
	currentMonth := now.Format("Jan")
	currentWeekday := int(now.Weekday())

	monthLine := "    "
	for w := range weeks {
		if s.Months[w] != "" {
			pad := cellWidth + 1 - len(s.Months[w])
			left := pad / 2
			right := pad - left
			monthStr := s.Months[w]
			if monthStr == currentMonth {
				monthStr = Accent(monthStr)
			} else {
				monthStr = Muted(monthStr)
			}
			monthLine += strings.Repeat(" ", left) + monthStr + strings.Repeat(" ", right)
		} else {
			monthLine += strings.Repeat(" ", cellWidth+1)
		}
	}
	fmt.Println(monthLine)

	fmt.Printf("    %s┌", brightBlack)
	fmt.Print(strings.Repeat(strings.Repeat("─", cellWidth)+"┬", weeks-1))
	fmt.Println(strings.Repeat("─", cellWidth) + "┐" + reset)

	dayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for row := range 7 {
		label := dayLabels[row]
		if row == currentWeekday {
			label = Primary(label)
		} else {
			label = Muted(label)
		}
		fmt.Printf("%s %s│%s", label, brightBlack, reset)
		for w := range weeks {
			idx := w*7 + row
			if idx >= len(s.Days) {
				fmt.Printf("%s%s│%s", strings.Repeat(" ", cellWidth), brightBlack, reset)
				continue
			}
			fmt.Printf("%s%s│%s", cellBlock(s.Days[idx], s.MaxCount), brightBlack, reset)
		}
		fmt.Println()
	}

	fmt.Printf("    %s└", brightBlack)
	fmt.Print(strings.Repeat(strings.Repeat("─", cellWidth)+"┴", weeks-1))
	fmt.Println(strings.Repeat("─", cellWidth) + "┘" + reset)
}

func RenderStatus() {
	s := BuildStatus()
	RenderTable(s)
}

func cellBlock(count, maxCount int) string {
	if count == -1 {
		return strings.Repeat(" ", cellWidth)
	}
	if count == -2 {
		return brightBlack + " ·· " + reset
	}
	if count == 0 {
		return black + " ░░ " + reset
	}
	if maxCount == 0 {
		maxCount = 1
	}

	ratio := float64(count) / float64(maxCount)
	var color string
	switch {
	case ratio >= 1.0:
		color = yellow
	case ratio >= 0.66:
		color = cyan
	case ratio >= 0.33:
		color = blue
	case ratio >= 0.20:
		color = blue + dim
	default:
		color = dim
	}
	return color + " ██ " + reset
}
