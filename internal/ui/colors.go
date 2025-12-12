package ui

import (
	"fmt"
	"strings"
)

const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"

	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	brightBlack   = "\033[90m"
	brightRed     = "\033[91m"
	brightGreen   = "\033[92m"
	brightYellow  = "\033[93m"
	brightBlue    = "\033[94m"
	brightMagenta = "\033[95m"
	brightCyan    = "\033[96m"
	brightWhite   = "\033[97m"
)

func Muted(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightBlack + fmt.Sprint(prefix[0]) + reset + dim + " " + fmt.Sprint(msg) + reset
	}
	return dim + fmt.Sprint(msg) + reset
}

func Primary(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightBlue + fmt.Sprint(prefix[0]) + reset + blue + " " + fmt.Sprint(msg) + reset
	}
	return blue + fmt.Sprint(msg) + reset
}

func Accent(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightYellow + fmt.Sprint(prefix[0]) + reset + yellow + " " + fmt.Sprint(msg) + reset
	}
	return yellow + fmt.Sprint(msg) + reset
}

func Positive(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightGreen + fmt.Sprint(prefix[0]) + reset + green + " " + fmt.Sprint(msg) + reset
	}
	return green + fmt.Sprint(msg) + reset
}

func Negative(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightRed + fmt.Sprint(prefix[0]) + reset + red + " " + fmt.Sprint(msg) + reset
	}
	return red + fmt.Sprint(msg) + reset
}

func Count(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightCyan + fmt.Sprint(prefix[0]) + reset + cyan + " " + fmt.Sprint(msg) + reset
	}
	return cyan + fmt.Sprint(msg) + reset
}

func Subtle(msg any, prefix ...any) string {
	if len(prefix) > 0 {
		return bold + brightBlack + fmt.Sprint(prefix[0]) + reset + brightBlack + " " + fmt.Sprint(msg) + reset
	}
	return brightBlack + fmt.Sprint(msg) + reset
}

func styled(icon, msg, iconColor, textColor string, subMessage ...string) {
	subMessageStr := ""
	if len(subMessage) > 0 {
		subMessageStr = " " + strings.Join(subMessage, " ")
	}
	fmt.Printf("%s%s%s%s %s%s%s%s\n", iconColor, bold, icon, reset, textColor, msg, subMessageStr, reset)
}

func styledInline(icon, msg, iconColor, textColor string) {
	fmt.Printf("%s%s%s%s %s%s%s ", iconColor, bold, icon, reset, textColor, msg, reset)
}

func Success(msg string, tag ...bool) {
	if len(tag) > 0 && tag[0] {
		styled("✓", "[SUCCESS] "+msg, brightGreen, green)
	} else {
		styled("✓", msg, brightGreen, green)
	}
}

func Error(msg string, tag ...bool) {
	if len(tag) > 0 && tag[0] {
		styled("✗", "[ERROR] "+msg, brightRed, red)
	} else {
		styled("✗", msg, brightRed, red)
	}
}

func Warn(msg string, tag ...bool) {
	if len(tag) > 0 && tag[0] {
		styled("!", "[WARN] "+msg, brightYellow, yellow)
	} else {
		styled("!", msg, brightYellow, yellow)
	}
}

func Info(msg string, tag ...bool) {
	if len(tag) > 0 && tag[0] {
		styled("i", "[INFO] "+msg, brightBlue, blue)
	} else {
		styled("i", msg, brightBlue, blue)
	}
}

func Title(msg string) { styled("◆◆", strings.ToUpper(msg), brightWhite, bold+white, "◆◆") }
func Hint(msg string)  { styled("*", msg, brightBlack, brightBlack) }

func InlineTitle(msg string) string {
	return dim + "◆ " + reset + bold + brightWhite + strings.ToUpper(msg) + reset + dim + " ◆" + reset
}

func InlineInfo(msg any) string {
	return bold + brightBlue + "i" + reset + " " + blue + fmt.Sprint(msg) + reset
}
func Debug(msg string) { styled("›", msg, brightMagenta, magenta) }
func Ask(msg string)   { styledInline("?", msg, brightCyan, cyan) }

func Header(msg string) {
	fmt.Printf("%s%s%s\n", bold, msg, reset)
}

func SubHeader(msg string) {
	fmt.Printf("  %s%s%s\n", dim, msg, reset)
}

func ListItem(n int, label string) {
	if n > 0 {
		fmt.Printf("  %s%d%s %s\n", dim, n, reset, label)
	} else {
		fmt.Printf("  %s\n", label)
	}
}

func ActionItem(icon, label string) {
	fmt.Printf("  %s%s%s %s\n", green, icon, reset, label)
}

func Label(text string) {
	fmt.Printf("\n%s?%s %s%s%s\n", brightCyan, reset, cyan, text, reset)
}

func LabelHint(text, hint string) {
	fmt.Printf("\n%s?%s %s%s%s %s(%s)%s\n", brightCyan, reset, cyan, text, reset, dim, hint, reset)
}

func InlineLabel(text string) {
	fmt.Printf("  %s%s:%s ", dim, text, reset)
}

func Section(title string) {
	fmt.Printf("  %s%s%s\n", dim, title, reset)
}

func Heading(title string) {
	fmt.Printf("\n%s%s%s\n", dim, title, reset)
}

func Field(label, value string) {
	fmt.Printf("  %s%s:%s %s\n", dim, label, reset, value)
}

func bar(length int, color string) string {
	if length < 0 {
		length = 0
	}
	return color + strings.Repeat("█", length) + reset
}

func BarPrimary(length int) string  { return bar(length, blue) }
func BarPositive(length int) string { return bar(length, green) }
func BarCount(length int) string    { return bar(length, cyan) }
func BarAccent(length int) string   { return bar(length, yellow) }
