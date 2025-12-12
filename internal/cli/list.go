package cli

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
			ui.Println(ui.Bold, p.Section+"/")
			currentSection = p.Section
		}

		kyuStr := ""
		if p.Kyu != nil {
			kyuStr = " " + ui.Colorize(fmt.Sprintf("[%d-%s]", *p.Kyu, kyuNames[*p.Kyu-1]), ui.Dim)
		}

		tagStr := ""
		if tags := p.TagList(); len(tags) > 0 {
			tagStr = " " + ui.Colorize(strings.Join(tags, ", "), ui.Cyan)
		}

		fmt.Printf("  %s%s%s\n", p.Name, kyuStr, tagStr)
	}

	fmt.Println()
	return nil
}

type SectionsCmd struct{}

func (c *SectionsCmd) Run() error {
	sections, err := db.GetSections()
	if err != nil {
		return err
	}

	if len(sections) == 0 {
		ui.Info("no sections yet")
		return nil
	}

	for _, s := range sections {
		ui.ListItem(0, s)
	}
	return nil
}

type TagsCmd struct{}

func (c *TagsCmd) Run() error {
	tags, err := db.GetTags()
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		ui.Info("no tags yet")
		return nil
	}

	for t, count := range tags {
		fmt.Printf("  %s %s\n", t, ui.Colorize(fmt.Sprintf("(%d)", count), ui.Dim))
	}
	return nil
}
