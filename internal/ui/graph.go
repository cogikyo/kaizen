package ui

import (
	"fmt"
	"time"

	"cogikyo/kaizen/internal/db"
)

func RenderHistory() {
	RenderHistoryYear(time.Now().Year())

	if db.HasDataBefore(time.Now().Year()) {
		fmt.Printf("\n  %srun 'kz history %d' to see last year%s\n", Dim, time.Now().Year()-1, Reset)
	}
}

func RenderHistoryYear(year int) {
	counts, err := db.GetSessionCountsForYear(year)
	if err != nil {
		return
	}

	countMap := make(map[string]int)
	maxCount := 0
	for _, c := range counts {
		countMap[c.Date] = c.Count
		if c.Count > maxCount {
			maxCount = c.Count
		}
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	prevYearCounts, _ := db.GetSessionCountsForYear(year - 1)
	for _, c := range prevYearCounts {
		countMap[c.Date] = c.Count
	}

	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.Local)
	now := time.Now()

	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
	days := make([]int, totalDays)
	for i := range totalDays {
		d := startDate.AddDate(0, 0, i)
		if d.After(now) {
			days[i] = -2
			continue
		}
		days[i] = countMap[d.Format("2006-01-02")]
	}

	weeks := (totalDays + 6) / 7

	months := []string{"    "}
	currentMonth := -1
	for w := range weeks {
		d := startDate.AddDate(0, 0, w*7)
		m := int(d.Month())
		if m != currentMonth && d.Year() == year {
			currentMonth = m
			months = append(months, d.Format("Jan")+" ")
		} else {
			months = append(months, "    ")
		}
	}

	monthLine := ""
	for i := 0; i < min(weeks+1, len(months)); i++ {
		monthLine += months[i]
	}
	fmt.Println(monthLine)

	dayLabels := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for row := range 7 {
		fmt.Printf("%s ", dayLabels[row])
		for w := range weeks {
			idx := w*7 + row
			if idx >= len(days) {
				break
			}
			fmt.Print(cellColor(days[idx], maxCount))
		}
		fmt.Println(Reset)
	}
}

func cellColor(count, maxCount int) string {
	if count == -1 {
		return "    "
	}
	if count == -2 {
		return Black + "··· " + Reset
	}
	if count == 0 {
		return Black + "░░░ " + Reset
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
	return color + "███ " + Reset
}
