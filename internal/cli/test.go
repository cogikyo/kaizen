package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type TestCmd struct {
	Target string `arg:"" optional:"" help:"Problem name or section"`
}

func (c *TestCmd) Run() error {
	var path string
	switch {
	case c.Target == "":
		path = "./..."
	case strings.Contains(c.Target, "/"):
		path = "./" + c.Target
	default:
		if p, err := db.GetProblem(c.Target); err == nil {
			path = "./" + p.Path
		} else {
			path = "./" + c.Target + "/..."
		}
	}

	fmt.Printf("%srunning%s %s\n\n", ui.Dim, ui.Reset, path)

	cmd := exec.Command("go", "test", "-v", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
