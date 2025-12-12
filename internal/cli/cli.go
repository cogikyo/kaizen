package cli

import (
	"cogikyo/kaizen/internal/cmds/browse"
	"cogikyo/kaizen/internal/cmds/dev"
	"cogikyo/kaizen/internal/cmds/info"
	"cogikyo/kaizen/internal/cmds/practice"
	"cogikyo/kaizen/internal/cmds/setup"
	"cogikyo/kaizen/internal/db"
)

type CLI struct {
	Init     setup.InitCmd      `cmd:"" aliases:"i" help:"Initialize kaizen repo or create section"`
	New      setup.NewCmd       `cmd:"" aliases:"n" help:"Create new problem (interactive)"`
	Solve    practice.SolveCmd  `cmd:"" aliases:"s" help:"Open problem in editor"`
	Done     practice.DoneCmd   `cmd:"" aliases:"d" help:"Mark problem complete, run tests"`
	Random   browse.RandomCmd   `cmd:"" aliases:"r" help:"Pick random problem to practice"`
	Review   browse.ReviewCmd   `cmd:"" help:"Show problems due for review"`
	Profile  info.ProfileCmd    `cmd:"" aliases:"p" default:"1" help:"Show status and suggest next"`
	History  info.HistoryCmd    `cmd:"" aliases:"h" help:"Show practice history (optional: yy-mm start)"`
	Stats    info.StatsCmd      `cmd:"" help:"Show detailed statistics"`
	List     browse.ListCmd     `cmd:"" aliases:"ls" help:"List problems"`
	Test     practice.TestCmd   `cmd:"" aliases:"t" help:"Run tests"`
	Sections browse.SectionsCmd `cmd:"" help:"List sections"`
	Tags     browse.TagsCmd     `cmd:"" help:"List all tags"`
	Seed     dev.SeedCmd        `cmd:"" hidden:"" help:"Generate fake practice data"`
	Reset    dev.ResetCmd       `cmd:"" hidden:"" help:"Clear all session data"`
}

func OpenDB() error {
	if err := db.Open(); err != nil {
		return err
	}
	return db.SyncProblems()
}
