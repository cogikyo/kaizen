package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func Prompt(label string) string {
	fmt.Printf("%s%s:%s ", Dim, label, Reset)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func PromptConfirm(label string) bool {
	fmt.Printf("%s%s [y/N]:%s ", Dim, label, Reset)
	line, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(line)) == "y"
}

func PromptSelect(label string, options []string, allowNew bool) (int, string) {
	fmt.Printf("\n%s%s%s\n", Dim, label, Reset)
	for i, opt := range options {
		fmt.Printf("  %s%d%s %s\n", Cyan, i+1, Reset, opt)
	}
	if allowNew {
		fmt.Printf("  %s+%s new\n", Green, Reset)
	}

	input := Prompt("select")
	if input == "" {
		return -1, ""
	}

	if n, err := strconv.Atoi(input); err == nil && n > 0 && n <= len(options) {
		return n - 1, options[n-1]
	}

	return -1, input
}

func PromptMultiSelect(label string, options []string, hint string) []string {
	fmt.Printf("\n%s%s%s %s(%s)%s\n", Dim, label, Reset, Dim, hint, Reset)
	for i, opt := range options {
		fmt.Printf("  %s%d%s %s\n", Cyan, i+1, Reset, opt)
	}
	if len(options) > 0 {
		fmt.Printf("  %s+%s add new\n", Green, Reset)
	}

	input := Prompt("select")
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

func Success(msg string) {
	fmt.Printf("%s✓%s %s\n", Green, Reset, msg)
}

func Error(msg string) {
	fmt.Printf("\n%s ✗ error:%s %s%s%s\n", Bold+Red, Reset, Red+Dim, msg, Reset)
}

func Info(msg string) {
	fmt.Printf("%s%s%s\n", Dim, msg, Reset)
}
