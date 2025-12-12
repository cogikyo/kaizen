package ui

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func Prompt(label string) string {
	Ask(label + ":")
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func PromptInline(label string) string {
	InlineLabel(label)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func PromptConfirm(label string) bool {
	Ask(label + " [y/N]:")
	line, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(line)) == "y"
}

func PromptSelect(label string, options []string, allowNew bool) (int, string) {
	Label(label)
	for i, opt := range options {
		ListItem(i+1, opt)
	}
	if allowNew {
		ActionItem("+", "new")
	}

	input := PromptInline("select")
	if input == "" {
		return -1, ""
	}

	if n, err := strconv.Atoi(input); err == nil && n > 0 && n <= len(options) {
		return n - 1, options[n-1]
	}

	return -1, input
}

func PromptMultiSelect(label string, options []string, hint string) []string {
	LabelHint(label, hint)
	for i, opt := range options {
		ListItem(i+1, opt)
	}
	if len(options) > 0 {
		ActionItem("+", "add new")
	}

	input := PromptInline("select")
	if input == "" {
		return nil
	}

	var result []string
	seen := make(map[string]bool)

	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == ' '
	})

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		if n, err := strconv.Atoi(p); err == nil && n > 0 && n <= len(options) {
			val := options[n-1]
			if !seen[val] {
				result = append(result, val)
				seen[val] = true
			}
		} else {
			val := strings.ToLower(p)
			if !seen[val] {
				result = append(result, val)
				seen[val] = true
			}
		}
	}

	return result
}
