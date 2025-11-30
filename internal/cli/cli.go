package cli

import (
	"cogikyo/kaizen/internal/db"
)

type Cmd interface {
	Run() error
}

type CLI struct {
	Init     InitCmd     `cmd:"" aliases:"i" help:"Initialize kaizen repo or create section"`
	New      NewCmd      `cmd:"" aliases:"n" help:"Create new problem (interactive)"`
	Solve    SolveCmd    `cmd:"" aliases:"s" help:"Open problem in editor"`
	Done     DoneCmd     `cmd:"" aliases:"d" help:"Mark problem complete, run tests"`
	Random   RandomCmd   `cmd:"" aliases:"r" help:"Pick random problem to practice"`
	Review   ReviewCmd   `cmd:"" help:"Show problems due for review"`
	Profile  ProfileCmd  `cmd:"" aliases:"p" default:"1" help:"Show status and suggest next"`
	History  HistoryCmd  `cmd:"" aliases:"h" help:"Show practice history (optional: yy-mm start)"`
	Stats    StatsCmd    `cmd:"" help:"Show detailed statistics"`
	List     ListCmd     `cmd:"" aliases:"ls" help:"List problems"`
	Test     TestCmd     `cmd:"" aliases:"t" help:"Run tests"`
	Sections SectionsCmd `cmd:"" help:"List sections"`
	Tags     TagsCmd     `cmd:"" help:"List all tags"`
	Seed     SeedCmd     `cmd:"" hidden:"" help:"Generate fake practice data"`
	Reset    ResetCmd    `cmd:"" hidden:"" help:"Clear all session data"`
}

func OpenDB() error {
	if err := db.Open(); err != nil {
		return err
	}
	return db.SyncProblems()
}
