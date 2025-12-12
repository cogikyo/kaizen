package dev

import (
	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type SeedCmd struct{}

func (c *SeedCmd) Run() error {
	if !ui.PromptConfirm("Generate fake practice data?") {
		return nil
	}

	if err := db.SeedSessions(); err != nil {
		return err
	}

	ui.Success("seeded practice data")
	return nil
}

