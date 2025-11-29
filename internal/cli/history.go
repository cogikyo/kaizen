package cli

import (
	"fmt"
	"time"

	"cogikyo/kaizen/internal/ui"
)

type HistoryCmd struct {
	Year int `arg:"" optional:"" help:"Year to show (defaults to current)"`
}

func (c *HistoryCmd) Run() error {
	year := c.Year
	if year == 0 {
		year = time.Now().Year()
	}

	fmt.Printf("\n%s%d%s\n\n", ui.Bold, year, ui.Reset)
	ui.RenderHistoryYear(year)
	fmt.Println()
	return nil
}
