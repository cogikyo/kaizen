package cli

import (
	"fmt"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type RandomCmd struct {
	Section string `short:"s" help:"Filter by section"`
	Kyu     int    `short:"k" help:"Filter by kyu level (1-6)"`
	Tag     string `short:"t" help:"Filter by tag"`
	NoDue   bool   `help:"Don't prioritize due reviews"`
}

func (c *RandomCmd) Run() error {
	var kyu *int
	if c.Kyu > 0 {
		kyu = &c.Kyu
	}

	p, err := db.GetRandomProblem(c.Section, kyu, c.Tag, !c.NoDue)
	if err != nil {
		return err
	}

	if p == nil {
		ui.Info("no problems match filters")
		return nil
	}

	fmt.Println()
	ui.Println(ui.Bold, p.Path)

	if p.Kyu != nil {
		ui.Field("kyu", fmt.Sprintf("%d-%s", *p.Kyu, kyuNames[*p.Kyu-1]))
	}
	if p.Tags != nil {
		ui.Field("tags", *p.Tags)
	}
	if p.Source != nil {
		ui.Field("source", *p.Source)
	}
	if p.URL != nil {
		ui.Field("url", *p.URL)
	}

	attempts, passes, totalTime, _ := db.GetProblemStats(p.Path)
	if attempts > 0 {
		ui.Field("practiced", fmt.Sprintf("%d attempts, %d passed, %s total", attempts, passes, formatDuration(totalTime)))
	}

	ui.Hint("kz solve " + p.Name)
	fmt.Println()

	return nil
}
