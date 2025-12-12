package ui

import (
	"fmt"
	"strings"
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

	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
)

func Colorize(text, color string) string {
	return color + text + Reset
}

func StripANSI(s string) string {
	var b strings.Builder
	inEscape := false
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
		b.WriteRune(r)
	}
	return b.String()
}

func styled(icon, msg, iconColor, textColor string, subMessage ...string) {
	subMessageStr := ""
	if len(subMessage) > 0 {
		subMessageStr = " " + strings.Join(subMessage, " ")
	}
	fmt.Printf("%s%s%s%s %s%s%s%s\n", iconColor, Bold, icon, Reset, textColor, msg, subMessageStr, Reset)
}

func styledInline(icon, msg, iconColor, textColor string) {
	fmt.Printf("%s%s%s%s %s%s%s ", iconColor, Bold, icon, Reset, textColor, msg, Reset)
}

func Success(msg string) { styled("✓", msg, BrightGreen, Green) }
func Error(msg string)   { styled("✗", msg, BrightRed, Red) }
func Warn(msg string)    { styled("!", msg, BrightYellow, Yellow) }
func Info(msg string)    { styled("i", msg, BrightBlue, Blue) }
func Title(msg string)   { styled("◆◆", strings.ToUpper(msg), BrightWhite, Bold+White, "◆◆") }
func Hint(msg string)    { styled("*", msg, BrightBlack, BrightBlack) }
func Debug(msg string)   { styled("›", msg, BrightMagenta, Magenta) }
func Ask(msg string)     { styledInline("?", msg, BrightCyan, Cyan) }

func Styled(icon, msg, iconColor, textColor string) {
	styled(icon, msg, iconColor, textColor)
}

func Print(color, msg string) {
	fmt.Print(color + msg + Reset)
}

func Println(color, msg string) {
	fmt.Println(color + msg + Reset)
}

func Printf(color, format string, args ...any) {
	fmt.Print(color + fmt.Sprintf(format, args...) + Reset)
}

func ListItem(n int, label string) {
	if n > 0 {
		fmt.Printf("  %s%d%s %s\n", Dim, n, Reset, label)
	} else {
		fmt.Printf("  %s\n", label)
	}
}

func ActionItem(icon, label string) {
	fmt.Printf("  %s%s%s %s\n", Green, icon, Reset, label)
}

func Label(text string) {
	fmt.Printf("\n%s?%s %s%s%s\n", BrightCyan, Reset, Cyan, text, Reset)
}

func LabelHint(text, hint string) {
	fmt.Printf("\n%s?%s %s%s%s %s(%s)%s\n", BrightCyan, Reset, Cyan, text, Reset, Dim, hint, Reset)
}

func InlineLabel(text string) {
	fmt.Printf("  %s%s:%s ", Dim, text, Reset)
}

func Section(title string) {
	fmt.Printf("  %s%s%s\n", Dim, title, Reset)
}

func Heading(title string) {
	fmt.Printf("\n%s%s%s\n", Dim, title, Reset)
}

func Field(label, value string) {
	fmt.Printf("  %s%s:%s %s\n", Dim, label, Reset, value)
}

func Bar(length int, color string) string {
	if length < 0 {
		length = 0
	}
	return color + strings.Repeat("█", length) + Reset
}
