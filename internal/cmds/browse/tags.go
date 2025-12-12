package browse

import (
	"fmt"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

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
		fmt.Printf("  %s %s\n", t, ui.Muted(fmt.Sprintf("(%d)", count)))
	}
	return nil
}

