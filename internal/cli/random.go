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

	fmt.Printf("\n%s%s%s\n", ui.Bold, p.Path, ui.Reset)

	if p.Kyu != nil {
		fmt.Printf("  %skyu:%s %d-%s\n", ui.Dim, ui.Reset, *p.Kyu, kyuNames[*p.Kyu-1])
	}
	if p.Tags != nil {
		fmt.Printf("  %stags:%s %s\n", ui.Dim, ui.Reset, *p.Tags)
	}
	if p.Source != nil {
		fmt.Printf("  %ssource:%s %s\n", ui.Dim, ui.Reset, *p.Source)
	}
	if p.URL != nil {
		fmt.Printf("  %surl:%s %s\n", ui.Dim, ui.Reset, *p.URL)
	}

	attempts, passes, totalTime, _ := db.GetProblemStats(p.Path)
	if attempts > 0 {
		fmt.Printf("  %spracticed:%s %d attempts, %d passed, %s total\n",
			ui.Dim, ui.Reset, attempts, passes, formatDuration(totalTime))
	}

	fmt.Printf("\n%skz solve %s%s\n\n", ui.Dim, p.Name, ui.Reset)

	return nil
}
