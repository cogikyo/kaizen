package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"cogikyo/kaizen/internal/db"
	"cogikyo/kaizen/internal/templates"
	"cogikyo/kaizen/internal/ui"
)

type NewCmd struct {
	Name string `arg:"" optional:"" help:"Problem name"`
}

func (c *NewCmd) Run() error {
	name := c.Name
	if name == "" {
		name = ui.Prompt("name")
		if name == "" {
			ui.Error("name required")
			return nil
		}
	}

	section := promptSection()
	if section == "" {
		return nil
	}

	kyu := promptKyu()
	tags := promptTags()
	source, url := promptSource()

	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	dir := filepath.Join(section, slug)

	if db.NameExists(slug) {
		ui.Error(fmt.Sprintf("problem %q already exists", slug))
		return nil
	}

	if _, err := os.Stat(dir); err == nil {
		ui.Error(fmt.Sprintf("%q already exists", dir))
		return nil
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var kyuNum string
	if kyu > 0 {
		kyuNum = fmt.Sprintf("%d", kyu)
	}

	data := map[string]any{
		"Name":    name,
		"Slug":    slug,
		"Section": section,
		"Package": strings.ReplaceAll(slug, "-", ""),
		"Date":    time.Now().Format("2006-01-02"),
		"Kyu":     kyuNum,
		"Tags":    strings.Join(tags, ", "),
		"Source":  source,
		"URL":     url,
	}

	files := []struct {
		tmpl string
		dest string
	}{
		{"solution.go.tmpl", filepath.Join(dir, "solution.go")},
		{"solution_test.go.tmpl", filepath.Join(dir, "solution_test.go")},
		{"README.md.tmpl", filepath.Join(dir, "README.md")},
	}

	for _, f := range files {
		if err := renderTemplate(f.tmpl, f.dest, data); err != nil {
			return err
		}
	}

	p := db.Problem{
		Path:      dir,
		Section:   section,
		Name:      slug,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	if kyu > 0 {
		p.Kyu = &kyu
	}
	if len(tags) > 0 {
		t := strings.Join(tags, ", ")
		p.Tags = &t
	}
	if source != "" {
		p.Source = &source
	}
	if url != "" {
		p.URL = &url
	}
	db.InsertProblem(p)

	fmt.Println()
	ui.Success(dir)
	return nil
}

func promptSection() string {
	sections, _ := db.GetSections()
	if len(sections) == 0 {
		entries, _ := os.ReadDir(".")
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") && e.Name() != "cmd" && e.Name() != "internal" {
				sections = append(sections, e.Name())
			}
		}
	}

	idx, val := ui.PromptSelect("section", sections, true)
	if val == "" {
		ui.Error("section required")
		return ""
	}

	if idx < 0 {
		slug := strings.ToLower(strings.ReplaceAll(val, " ", "_"))
		if _, err := os.Stat(slug); os.IsNotExist(err) {
			cmd := InitCmd{Name: slug}
			cmd.Run()
		}
		return slug
	}

	return val
}

var kyuNames = []string{"elite", "difficult", "hard", "medium", "normal", "easy"}

func promptKyu() int {
	options := make([]string, 6)
	for i := range 6 {
		options[i] = fmt.Sprintf("%d %s", i+1, kyuNames[i])
	}

	ui.Heading("kyu")
	for i, opt := range options {
		ui.ListItem(i+1, opt)
	}
	ui.ActionItem("-", "skip")

	input := ui.PromptInline("select")
	if input == "" || input == "-" {
		return 0
	}

	for i, name := range kyuNames {
		if input == fmt.Sprintf("%d", i+1) || strings.EqualFold(input, name) {
			return i + 1
		}
	}
	return 0
}

func promptTags() []string {
	existing, _ := db.GetTags()
	var options []string
	for t := range existing {
		options = append(options, t)
	}
	return ui.PromptMultiSelect("tags", options, "space/comma separated")
}

func promptSource() (string, string) {
	sources, _ := db.GetSources()
	idx, val := ui.PromptSelect("source", sources, true)
	if val == "" || val == "-" {
		return "", ""
	}

	var source string
	if idx >= 0 {
		source = val
	} else {
		source = strings.ToLower(val)
	}

	url := ui.Prompt("url (optional)")
	return source, url
}

func renderTemplate(name, dest string, data map[string]any) error {
	content, err := templates.FS.ReadFile(name)
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(string(content))
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
