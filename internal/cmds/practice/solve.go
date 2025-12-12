package practice

import (
	"fmt"
	"os"
	"os/exec"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type SolveCmd struct {
	Problem string `arg:"" optional:"" help:"Problem name"`
	NoDone  bool   `short:"n" help:"Skip running tests on exit"`
}

func (c *SolveCmd) Run() error {
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

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	ui.Hint(fmt.Sprintf("opening %s in %s", p.Path, editor))
	fmt.Println()

	cmd := exec.Command(editor, p.Path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	if !c.NoDone {
		fmt.Println()
		return runDone(p)
	}

	return nil
}

func promptProblem() string {
	problems, _ := db.GetProblems("", nil, "")
	if len(problems) == 0 {
		ui.Info("no problems yet")
		return ""
	}

	options := make([]string, len(problems))
	for i, p := range problems {
		options[i] = p.Path
	}

	_, val := ui.PromptSelect("problem", options, false)
	return val
}
