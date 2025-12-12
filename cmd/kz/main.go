package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"cogikyo/kaizen/internal/cli"
	"cogikyo/kaizen/internal/db"
)

func main() {
	var c cli.CLI
	ctx := kong.Parse(&c,
		kong.Name("kz"),
		kong.Description("kata practice tracker"),
		kong.UsageOnError(),
	)

	cmd := ctx.Command()
	isInit := strings.HasPrefix(cmd, "init")

	if !isInit {
		if err := cli.OpenDB(); err != nil {
			if errors.Is(err, db.ErrNotKaizen) {
				fmt.Println("not a kaizen directory (run 'kz init' to create)")
				os.Exit(1)
			}
			ctx.Errorf("failed to open database: %v", err)
			os.Exit(1)
		}
		defer db.Close()
	}

	if err := ctx.Run(); err != nil {
		ctx.Errorf("%v", err)
		os.Exit(1)
	}
}
