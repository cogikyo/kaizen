package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cogikyo/kaizen/internal/ui"
)

type HistoryCmd struct {
	Start string `arg:"" optional:"" help:"Start month (yy-mm format)"`
}

func (c *HistoryCmd) Run() error {
	startMonth := ui.DefaultStartMonth()

	if c.Start != "" {
		parts := strings.Split(c.Start, "-")
		if len(parts) != 2 {
			ui.Error("invalid format, use yy-mm")
			return nil
		}
		year, err1 := strconv.Atoi(parts[0])
		month, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || month < 1 || month > 12 {
			ui.Error("invalid format, use yy-mm")
			return nil
		}
		year += 2000
		startMonth = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	}

	fmt.Println()
	ui.RenderStatus(startMonth)
	fmt.Println()
	return nil
}
