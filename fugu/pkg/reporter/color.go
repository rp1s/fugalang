package reporter

const (
	ansiReset = "\033[0m"

	ansiBlack           = "\033[0;30m"
	ansiRed             = "\033[0;31m"
	ansiGreen           = "\033[0;32m"
	ansiYellow          = "\033[0;33m"
	ansiBlue            = "\033[0;34m"
	ansiMagenta         = "\033[0;35m"
	ansiCyan            = "\033[0;36m"
	ansiWhite           = "\033[0;37m"
	ansiGray            = "\033[90m"
	ansiPastelYellow    = "\033[38;2;255;245;150m"
	ansi256PastelYellow = "\033[38;5;229m"

	ansiBoldBlack   = "\033[1;30m"
	ansiBoldRed     = "\033[1;31m"
	ansiBoldGreen   = "\033[1;32m"
	ansiBoldYellow  = "\033[1;33m"
	ansiBoldBlue    = "\033[1;34m"
	ansiBoldMagenta = "\033[1;35m"
	ansiBoldCyan    = "\033[1;36m"
	ansiBoldWhite   = "\033[1;37m"

	ansiBold      = "\033[1m"
	ansiUnderline = "\033[4m"

	ansiBgRed   = "\033[41m"
	ansiBgGreen = "\033[42m"
)

func Black(s string) string   { return ansiBlack + s + ansiReset }
func Red(s string) string     { return ansiRed + s + ansiReset }
func Green(s string) string   { return ansiGreen + s + ansiReset }
func Yellow(s string) string  { return ansiYellow + s + ansiReset }
func Blue(s string) string    { return ansiBlue + s + ansiReset }
func Magenta(s string) string { return ansiMagenta + s + ansiReset }
func Cyan(s string) string    { return ansiCyan + s + ansiReset }
func White(s string) string   { return ansiWhite + s + ansiReset }
func Gray(s string) string    { return ansiGray + s + ansiReset }

func BoldBlack(s string) string   { return ansiBoldBlack + s + ansiReset }
func BoldRed(s string) string     { return ansiBoldRed + s + ansiReset }
func BoldGreen(s string) string   { return ansiBoldGreen + s + ansiReset }
func BoldYellow(s string) string  { return ansiBoldYellow + s + ansiReset }
func BoldBlue(s string) string    { return ansiBoldBlue + s + ansiReset }
func BoldMagenta(s string) string { return ansiBoldMagenta + s + ansiReset }
func BoldCyan(s string) string    { return ansiBoldCyan + s + ansiReset }
func BoldWhite(s string) string   { return ansiBoldWhite + s + ansiReset }

func Bold(s string) string      { return ansiBold + s + ansiReset }
func Underline(s string) string { return ansiUnderline + s + ansiReset }

func ErrorLabel(s string) string   { return ansiBoldRed + s + ansiReset }
func WarningLabel(s string) string { return ansiBoldYellow + s + ansiReset }
func NoteLabel(s string) string    { return ansiBoldBlue + s + ansiReset }
func Highlight(s string) string    { return ansiUnderline + ansiBold + s + ansiReset }

func PastelYellow(s string) string    { return ansiPastelYellow + s + ansiReset }
func PastelYellow256(s string) string { return ansi256PastelYellow + s + ansiReset }
