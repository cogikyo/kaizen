package db

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Problem struct {
	Path      string  `db:"path"`
	Section   string  `db:"section"`
	Name      string  `db:"name"`
	Kyu       *int    `db:"kyu"`
	Tags      *string `db:"tags"`
	Source    *string `db:"source"`
	URL       *string `db:"url"`
	CreatedAt string  `db:"created_at"`
}

func (p Problem) TagList() []string {
	if p.Tags == nil || *p.Tags == "" {
		return nil
	}
	parts := strings.Split(*p.Tags, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func SyncProblems() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	var found []Problem
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == "cmd" || e.Name() == "internal" {
			continue
		}

		section := e.Name()
		subs, err := os.ReadDir(section)
		if err != nil {
			continue
		}

		for _, s := range subs {
			if !s.IsDir() {
				continue
			}

			p := Problem{
				Path:      filepath.Join(section, s.Name()),
				Section:   section,
				Name:      s.Name(),
				CreatedAt: time.Now().Format(time.RFC3339),
			}

			readme := filepath.Join(p.Path, "README.md")
			if meta := parseReadmeMeta(readme); meta != nil {
				if kyu, ok := meta["kyu"]; ok && kyu != "" {
					k := parseKyu(kyu)
					p.Kyu = &k
				}
				if tags, ok := meta["tags"]; ok && tags != "" {
					p.Tags = &tags
				}
				if source, ok := meta["source"]; ok && source != "" {
					p.Source = &source
				}
				if url, ok := meta["url"]; ok && url != "" {
					p.URL = &url
				}
			}

			found = append(found, p)
		}
	}

	for _, p := range found {
		_, err := conn.Exec(`
			INSERT INTO problems (path, section, name, kyu, tags, source, url, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(path) DO UPDATE SET
				kyu = excluded.kyu,
				tags = excluded.tags,
				source = excluded.source,
				url = excluded.url
		`, p.Path, p.Section, p.Name, p.Kyu, p.Tags, p.Source, p.URL, p.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseReadmeMeta(path string) map[string]string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() || scanner.Text() != "---" {
		return nil
	}

	meta := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			break
		}
		if idx := strings.Index(line, ":"); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			val = strings.Trim(val, "[]\"'")
			meta[key] = val
		}
	}
	return meta
}

var kyuFromString = map[string]int{
	"1": 1, "elite": 1,
	"2": 2, "difficult": 2,
	"3": 3, "hard": 3,
	"4": 4, "medium": 4,
	"5": 5, "normal": 5,
	"6": 6, "easy": 6,
}

func parseKyu(s string) int {
	return kyuFromString[strings.ToLower(s)]
}

func GetProblems(section string, kyu *int, tag string) ([]Problem, error) {
	query := "SELECT * FROM problems WHERE 1=1"
	args := []any{}

	if section != "" {
		query += " AND section = ?"
		args = append(args, section)
	}
	if kyu != nil {
		query += " AND kyu = ?"
		args = append(args, *kyu)
	}
	if tag != "" {
		query += " AND tags LIKE ?"
		args = append(args, "%"+tag+"%")
	}

	query += " ORDER BY section, name"

	var problems []Problem
	err := conn.Select(&problems, query, args...)
	return problems, err
}

func GetProblem(pathOrName string) (*Problem, error) {
	var p Problem
	err := conn.Get(&p, "SELECT * FROM problems WHERE path = ? OR name = ?", pathOrName, pathOrName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func NameExists(name string) bool {
	var count int
	conn.Get(&count, "SELECT COUNT(*) FROM problems WHERE name = ?", name)
	return count > 0
}

func GetSections() ([]string, error) {
	var sections []string
	err := conn.Select(&sections, "SELECT DISTINCT section FROM problems ORDER BY section")
	return sections, err
}

func GetTags() (map[string]int, error) {
	problems, err := GetProblems("", nil, "")
	if err != nil {
		return nil, err
	}

	tags := make(map[string]int)
	for _, p := range problems {
		for _, t := range p.TagList() {
			tags[t]++
		}
	}
	return tags, nil
}

func GetSources() ([]string, error) {
	var sources []string
	err := conn.Select(&sources, "SELECT DISTINCT source FROM problems WHERE source IS NOT NULL AND source != '' ORDER BY source")
	return sources, err
}

func InsertProblem(p Problem) error {
	_, err := conn.Exec(`
		INSERT INTO problems (path, section, name, kyu, tags, source, url, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, p.Path, p.Section, p.Name, p.Kyu, p.Tags, p.Source, p.URL, p.CreatedAt)
	return err
}
