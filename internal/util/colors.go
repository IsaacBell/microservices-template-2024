package util

import "fmt"

type AnsiColor string

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

func PrintLnInColor(color AnsiColor, input ...any) {
	fmt.Print(string(color))
	fmt.Println(input...)
	fmt.Print(string(AnsiColorReset))
}