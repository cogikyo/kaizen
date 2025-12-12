package browse

import (
	"fmt"
	"strings"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ListCmd struct {
	Section string `arg:"" optional:"" help:"Filter by section"`
	Kyu     int    `short:"k" help:"Filter by kyu level (1-6)"`
	Tag     string `short:"t" help:"Filter by tag"`
}

func (c *ListCmd) Run() error {
	var kyu *int
	if c.Kyu > 0 {
		kyu = &c.Kyu
	}

	problems, err := db.GetProblems(c.Section, kyu, c.Tag)
	if err != nil {
		return err
	}

	if len(problems) == 0 {
		ui.Info("no problems found")
		return nil
	}

	currentSection := ""
	for _, p := range problems {
		if p.Section != currentSection {
			if currentSection != "" {
				fmt.Println()
			}
			ui.Header(p.Section + "/")
			currentSection = p.Section
		}

		kyuStr := ""
		if p.Kyu != nil {
			kyuStr = " " + ui.Muted(fmt.Sprintf("[%d-%s]", *p.Kyu, db.KyuNames[*p.Kyu-1]))
		}

		tagStr := ""
		if tags := p.TagList(); len(tags) > 0 {
			tagStr = " " + ui.Count(strings.Join(tags, ", "))
		}

		fmt.Printf("  %s%s%s\n", p.Name, kyuStr, tagStr)
	}

	fmt.Println()
	return nil
}

