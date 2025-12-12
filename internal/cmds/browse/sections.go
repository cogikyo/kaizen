package browse

import (
	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

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

