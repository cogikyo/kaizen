package db

import (
	"math/rand"
	"time"
)

type ReviewSchedule struct {
	ProblemPath  string  `db:"problem_path"`
	NextReview   string  `db:"next_review"`
	IntervalDays int     `db:"interval_days"`
	EaseFactor   float64 `db:"ease_factor"`
}

func UpdateReview(problemPath string, passed bool) error {
	var schedule ReviewSchedule
	err := conn.Get(&schedule, "SELECT * FROM review_schedule WHERE problem_path = ?", problemPath)

	if err != nil {
		schedule = ReviewSchedule{
			ProblemPath:  problemPath,
			IntervalDays: 1,
			EaseFactor:   2.5,
		}
	}

	if passed {
		if schedule.IntervalDays == 1 {
			schedule.IntervalDays = 2
		} else {
			schedule.IntervalDays = int(float64(schedule.IntervalDays) * schedule.EaseFactor)
		}
		if schedule.IntervalDays > 60 {
			schedule.IntervalDays = 60
		}
		schedule.EaseFactor += 0.1
		if schedule.EaseFactor > 3.0 {
			schedule.EaseFactor = 3.0
		}
	} else {
		schedule.IntervalDays = 1
		schedule.EaseFactor -= 0.2
		if schedule.EaseFactor < 1.3 {
			schedule.EaseFactor = 1.3
		}
	}

	schedule.NextReview = time.Now().AddDate(0, 0, schedule.IntervalDays).Format("2006-01-02")

	_, err = conn.Exec(`
		INSERT INTO review_schedule (problem_path, next_review, interval_days, ease_factor)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(problem_path) DO UPDATE SET
			next_review = excluded.next_review,
			interval_days = excluded.interval_days,
			ease_factor = excluded.ease_factor
	`, schedule.ProblemPath, schedule.NextReview, schedule.IntervalDays, schedule.EaseFactor)

	return err
}

func GetDueReviews() ([]Problem, error) {
	today := time.Now().Format("2006-01-02")

	var problems []Problem
	err := conn.Select(&problems, `
		SELECT p.* FROM problems p
		JOIN review_schedule r ON p.path = r.problem_path
		WHERE r.next_review <= ?
		ORDER BY r.next_review, r.interval_days
	`, today)
	return problems, err
}

func GetDueCount() (int, error) {
	today := time.Now().Format("2006-01-02")
	var count int
	err := conn.Get(&count, "SELECT COUNT(*) FROM review_schedule WHERE next_review <= ?", today)
	return count, err
}

func GetRandomProblem(section string, kyu *int, tag string, preferDue bool) (*Problem, error) {
	if preferDue {
		due, err := GetDueReviews()
		if err == nil && len(due) > 0 {
			filtered := due
			if section != "" || kyu != nil || tag != "" {
				filtered = nil
				for _, p := range due {
					if section != "" && p.Section != section {
						continue
					}
					if kyu != nil && (p.Kyu == nil || *p.Kyu != *kyu) {
						continue
					}
					if tag != "" {
						hasTag := false
						for _, t := range p.TagList() {
							if t == tag {
								hasTag = true
								break
							}
						}
						if !hasTag {
							continue
						}
					}
					filtered = append(filtered, p)
				}
			}
			if len(filtered) > 0 {
				return &filtered[rand.Intn(len(filtered))], nil
			}
		}
	}

	problems, err := GetProblems(section, kyu, tag)
	if err != nil {
		return nil, err
	}
	if len(problems) == 0 {
		return nil, nil
	}

	return &problems[rand.Intn(len(problems))], nil
}
