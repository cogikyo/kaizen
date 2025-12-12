package dev

import (
	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type ResetCmd struct{}

func (c *ResetCmd) Run() error {
	if !ui.PromptConfirm("Clear all session data? This cannot be undone.") {
		return nil
	}

	if err := db.ClearSessions(); err != nil {
		return err
	}

	ui.Success("cleared all session data")
	return nil
}

