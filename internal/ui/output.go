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

func Stat(value any, label string, color string) string {
	return fmt.Sprintf("%s%s%v%s%s%s", Bold+color, "", value, Reset+color, " "+label, Reset)
}

func Justified(left, right string, width int) {
	leftLen := VisibleLen(left)
	rightLen := VisibleLen(right)
	padding := max(width-leftLen-rightLen, 1)
	fmt.Printf("%s%*s%s\n", left, padding, "", right)
}
