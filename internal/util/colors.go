package util

import "fmt"

type AnsiColor string

// General Colors
const (
	AnsiColorReset   AnsiColor = "\033[0m"
	AnsiColorRed     AnsiColor = "\033[31m"
	AnsiColorGreen   AnsiColor = "\033[32m"
	AnsiColorYellow  AnsiColor = "\033[33m"
	AnsiColorBlue    AnsiColor = "\033[34m"
	AnsiColorMagenta AnsiColor = "\033[35m"
	AnsiColorCyan    AnsiColor = "\033[36m"
	AnsiColorGray    AnsiColor = "\033[37m"
	AnsiColorWhite   AnsiColor = "\033[97m"
)

// Bright Colors
const (
	AnsiBrightRed     AnsiColor = "\033[91m"
	AnsiBrightGreen   AnsiColor = "\033[92m"
	AnsiBrightYellow  AnsiColor = "\033[93m"
	AnsiBrightBlue    AnsiColor = "\033[94m"
	AnsiBrightMagenta AnsiColor = "\033[95m"
	AnsiBrightCyan    AnsiColor = "\033[96m"
)

// Background Colors
const (
	AnsiBackgroundBlack   AnsiColor = "\033[40m"
	AnsiBackgroundRed     AnsiColor = "\033[41m"
	AnsiBackgroundGreen   AnsiColor = "\033[42m"
	AnsiBackgroundYellow  AnsiColor = "\033[43m"
	AnsiBackgroundBlue    AnsiColor = "\033[44m"
	AnsiBackgroundMagenta AnsiColor = "\033[45m"
	AnsiBackgroundCyan    AnsiColor = "\033[46m"
	AnsiBackgroundWhite   AnsiColor = "\033[47m"
)

// Text Styles
const (
	AnsiStyleBold          AnsiColor = "\033[1m"
	AnsiStyleFaint         AnsiColor = "\033[2m"
	AnsiStyleItalic        AnsiColor = "\033[3m"
	AnsiStyleUnderline     AnsiColor = "\033[4m"
	AnsiStyleBlink         AnsiColor = "\033[5m"
	AnsiStyleBlinkRapid    AnsiColor = "\033[6m"
	AnsiStyleReverse       AnsiColor = "\033[7m"
	AnsiStyleInvisible     AnsiColor = "\033[8m"
	AnsiStyleStrikethrough AnsiColor = "\033[9m"
)

func PrintInColor(color AnsiColor, input ...any) {
	fmt.Print(string(color))
	fmt.Print(input...)
	fmt.Print(string(AnsiColorReset))
}

func PrintLnInColor(color AnsiColor, input ...any) {
	fmt.Print(string(color))
	fmt.Println(input...)
	fmt.Print(string(AnsiColorReset))
}
