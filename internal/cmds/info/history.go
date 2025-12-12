package info

import (
	"fmt"

	"cogikyo/kaizen/internal/ui"
)

type HistoryCmd struct{}

func (c *HistoryCmd) Run() error {
	fmt.Println()
	ui.RenderStatus()
	fmt.Println()
	return nil
}

