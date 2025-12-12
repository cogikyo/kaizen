package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func TermWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 40 {
		return 80
	}
	return w
}

func VisibleLen(s string) int {
	inEscape := false
	length := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

func Justified(width int, groups ...string) {
	if len(groups) == 0 {
		fmt.Println()
		return
	}
	if len(groups) == 1 {
		fmt.Println(groups[0])
		return
	}

	totalLen := 0
	for _, g := range groups {
		totalLen += VisibleLen(g)
	}

	gaps := len(groups) - 1
	totalPadding := max(width-totalLen, gaps)
	perGap := totalPadding / gaps
	extra := totalPadding % gaps

	var out string
	for i, g := range groups {
		out += g
		if i < gaps {
			pad := perGap
			if i < extra {
				pad++
			}
			for range pad {
				out += " "
			}
		}
	}
	fmt.Println(out)
}
