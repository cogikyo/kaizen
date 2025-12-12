package practice

import (
	"fmt"
	"os/exec"
	"time"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type DoneCmd struct {
	Problem string `arg:"" optional:"" help:"Problem name"`
}

func (c *DoneCmd) Run() error {
	problem := c.Problem
	if problem == "" {
		problem = promptProblem()
		if problem == "" {
			return nil
		}
	}

	p, err := db.GetProblem(problem)
	if err != nil {
		ui.Error(fmt.Sprintf("problem %q not found", problem))
		return nil
	}

	return runDone(p)
}

func runDone(p *db.Problem) error {
	ui.Hint("running tests...")

	start := time.Now()
	cmd := exec.Command("go", "test", "-v", "./"+p.Path)
	output, err := cmd.CombinedOutput()
	duration := int(time.Since(start).Seconds())

	passed := err == nil

	fmt.Println(string(output))

	if err := db.RecordSession(p.Path, passed, duration); err != nil {
		return err
	}

	if err := db.UpdateReview(p.Path, passed); err != nil {
		return err
	}

	if passed {
		ui.Success(fmt.Sprintf("%s completed (passed) %s", p.Path, ui.FormatDuration(duration)))
	} else {
		ui.Warn(fmt.Sprintf("%s completed (failed) %s", p.Path, ui.FormatDuration(duration)))
	}

	return nil
}
