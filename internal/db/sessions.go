package db

import (
	"math/rand"
	"time"
)

type Session struct {
	ID              int    `db:"id"`
	ProblemPath     string `db:"problem_path"`
	Date            string `db:"date"`
	Passed          bool   `db:"passed"`
	DurationSeconds int    `db:"duration_seconds"`
	StartedAt       string `db:"started_at"`
	FinishedAt      string `db:"finished_at"`
}

const sessionGap = time.Hour

func RecordSession(problemPath string, passed bool, durationSeconds int) error {
	now := time.Now()
	_, err := conn.Exec(`
		INSERT INTO sessions (problem_path, date, passed, duration_seconds, started_at, finished_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, problemPath, now.Format("2006-01-02"), passed, durationSeconds, now.Format(time.RFC3339), now.Format(time.RFC3339))
	return err
}

type DayCount struct {
	Date  string `db:"date"`
	Count int    `db:"count"`
}

func GetSessionCounts(days int) ([]DayCount, error) {
	start := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	var counts []DayCount
	err := conn.Select(&counts, `
		SELECT date, COUNT(DISTINCT problem_path) as count 
		FROM sessions 
		WHERE date >= ?
		GROUP BY date 
		ORDER BY date
	`, start)
	return counts, err
}

func GetTodaySessions() ([]string, error) {
	today := time.Now().Format("2006-01-02")
	var problems []string
	err := conn.Select(&problems, `
		SELECT DISTINCT problem_path 
		FROM sessions 
		WHERE date = ?
	`, today)
	return problems, err
}

type Stats struct {
	TotalSessions  int
	TotalAttempts  int
	UniqueProblems int
	CurrentStreak  int
	LongestStreak  int
	TotalTime      int
	AvgSessionTime int
	PassRate       float64
	TotalPassed    int
	TotalFailed    int
}

func GetStats() (*Stats, error) {
	s := &Stats{}

	conn.Get(&s.TotalAttempts, "SELECT COUNT(*) FROM sessions")
	conn.Get(&s.UniqueProblems, "SELECT COUNT(DISTINCT problem_path) FROM sessions")
	conn.Get(&s.TotalTime, "SELECT COALESCE(SUM(duration_seconds), 0) FROM sessions")
	conn.Get(&s.TotalPassed, "SELECT COUNT(*) FROM sessions WHERE passed = 1")
	conn.Get(&s.TotalFailed, "SELECT COUNT(*) FROM sessions WHERE passed = 0")

	if s.TotalAttempts > 0 {
		s.PassRate = float64(s.TotalPassed) / float64(s.TotalAttempts) * 100
	}

	s.TotalSessions = countSessions()
	if s.TotalSessions > 0 {
		s.AvgSessionTime = s.TotalTime / s.TotalSessions
	}

	var dates []string
	conn.Select(&dates, "SELECT DISTINCT date FROM sessions ORDER BY date DESC")

	if len(dates) > 0 {
		today := time.Now().Format("2006-01-02")
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		if dates[0] == today || dates[0] == yesterday {
			s.CurrentStreak = 1
			for i := 1; i < len(dates); i++ {
				curr, _ := time.Parse("2006-01-02", dates[i-1])
				prev, _ := time.Parse("2006-01-02", dates[i])
				if curr.Sub(prev).Hours() == 24 {
					s.CurrentStreak++
				} else {
					break
				}
			}
		}

		streak := 1
		for i := 1; i < len(dates); i++ {
			curr, _ := time.Parse("2006-01-02", dates[i-1])
			prev, _ := time.Parse("2006-01-02", dates[i])
			if curr.Sub(prev).Hours() == 24 {
				streak++
				if streak > s.LongestStreak {
					s.LongestStreak = streak
				}
			} else {
				streak = 1
			}
		}
		if streak > s.LongestStreak {
			s.LongestStreak = streak
		}
	}

	return s, nil
}

func countSessions() int {
	var timestamps []string
	conn.Select(&timestamps, "SELECT finished_at FROM sessions ORDER BY finished_at")

	if len(timestamps) == 0 {
		return 0
	}

	sessions := 1
	var lastTime time.Time
	for i, ts := range timestamps {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			continue
		}
		if i == 0 {
			lastTime = t
			continue
		}
		if t.Sub(lastTime) > sessionGap {
			sessions++
		}
		lastTime = t
	}
	return sessions
}

func GetProblemStats(path string) (attempts int, passes int, totalTime int, err error) {
	conn.Get(&attempts, "SELECT COUNT(*) FROM sessions WHERE problem_path = ?", path)
	conn.Get(&passes, "SELECT COUNT(*) FROM sessions WHERE problem_path = ? AND passed = 1", path)
	conn.Get(&totalTime, "SELECT COALESCE(SUM(duration_seconds), 0) FROM sessions WHERE problem_path = ?", path)
	return
}

func ClearSessions() error {
	_, err := conn.Exec("DELETE FROM sessions")
	if err != nil {
		return err
	}
	_, err = conn.Exec("DELETE FROM review_schedule")
	return err
}

func SeedSessions() error {
	problems := []string{
		"data_structures/two-sum",
		"data_structures/valid-parentheses",
		"data_structures/merge-intervals",
		"algorithms/binary-search",
		"algorithms/quick-sort",
		"algorithms/dijkstra",
		"patterns/sliding-window",
		"patterns/two-pointers",
	}

	start := time.Date(2025, 5, 1, 0, 0, 0, 0, time.Local)
	now := time.Now()

	for d := start; !d.After(now); d = d.AddDate(0, 0, 1) {
		if rand.Float32() > 0.4 {
			continue
		}

		count := rand.Intn(3) + 1
		for range count {
			problem := problems[rand.Intn(len(problems))]
			duration := rand.Intn(1800) + 300
			passed := rand.Float32() > 0.15

			hour := rand.Intn(12) + 8
			minute := rand.Intn(60)
			sessionTime := time.Date(d.Year(), d.Month(), d.Day(), hour, minute, 0, 0, time.Local)
			endTime := sessionTime.Add(time.Duration(duration) * time.Second)

			_, err := conn.Exec(`
				INSERT INTO sessions (problem_path, date, passed, duration_seconds, started_at, finished_at)
				VALUES (?, ?, ?, ?, ?, ?)
			`, problem, d.Format("2006-01-02"), passed, duration, sessionTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
