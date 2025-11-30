package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	Grey          = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

func TermWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 40 {
		return 80
	}
	return w
}

func Success(msg string) {
	fmt.Printf("%s%s %s%s%s\n", BrightGreen, Reset, Green, msg, Reset)
}

func Error(msg string) {
	fmt.Printf("%s%s %s%s%s\n", BrightRed, Reset, Red, msg, Reset)
}

func Warn(msg string) {
	fmt.Printf("%s%s %s%s%s\n", BrightYellow, Reset, Yellow, msg, Reset)
}

func Ask(msg string) {
	fmt.Printf("%s%s %s%s%s ", BrightCyan, Reset, Cyan, msg, Reset)
}

func Heading(msg string) {
	fmt.Printf("%s%s␁%s%s%s\n", BrightBlue, Reset, Blue, msg, Reset)
}

func Info(msg string) {
	fmt.Printf("%s%s %s%s%s\n", BrightBlue, Reset, Dim, msg, Reset)
}

func Stat(value any, label string, color string) string {
	return fmt.Sprintf("%s%s%v%s%s%s", Bold+color, "", value, Reset+color, " "+label, Reset)
}

func Justified(left, right string, width int) {
	leftLen := VisibleLen(left)
	rightLen := VisibleLen(right)
	padding := width - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}
	fmt.Printf("%s%*s%s\n", left, padding, "", right)
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
