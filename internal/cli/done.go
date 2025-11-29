package cli

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

	fmt.Printf("%srunning tests...%s\n", ui.Dim, ui.Reset)

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
		ui.Success(fmt.Sprintf("%s completed (passed) %s", p.Path, formatDuration(duration)))
	} else {
		fmt.Printf("%s✗%s %s completed (failed) %s\n", ui.Red, ui.Reset, p.Path, formatDuration(duration))
	}

	return nil
}
