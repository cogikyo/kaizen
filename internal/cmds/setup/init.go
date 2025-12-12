package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/ui"
)

type InitCmd struct {
	Name string `arg:"" optional:"" help:"Section name (omit to initialize repo)"`
}

func (c *InitCmd) Run() error {
	if c.Name == "" {
		return initRepo()
	}
	return initSection(c.Name)
}

func initRepo() error {
	if db.Exists() {
		ui.Info("already initialized")
		return nil
	}

	if err := db.Create(); err != nil {
		return err
	}

	ui.Success("initialized .kaizen/")
	return nil
}

func initSection(name string) error {
	if !db.Exists() {
		ui.Warn("not a kaizen directory (run 'kz init' first)")
		return nil
	}

	slug := strings.ToLower(strings.ReplaceAll(name, " ", "_"))

	if _, err := os.Stat(slug); err == nil {
		ui.Error(fmt.Sprintf("section %q already exists", slug))
		return nil
	}

	if err := os.MkdirAll(slug, 0755); err != nil {
		return err
	}

	readme := filepath.Join(slug, "README.md")
	content := fmt.Sprintf("# %s\n", strings.ReplaceAll(slug, "_", " "))
	if err := os.WriteFile(readme, []byte(content), 0644); err != nil {
		return err
	}

	ui.Success(slug + "/")
	return nil
}

