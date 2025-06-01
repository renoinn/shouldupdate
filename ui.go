package main

import (
	"fmt"
	"os"
)

// Gruvbox-inspired 256-Color ANSI Codes
const (
	colorReset     = "\033[0m"
	colorFgDefault = "\033[38;5;250m" // light cream/grey (ebdbb2)
	colorBgDefault = "\033[48;5;235m" // dark grey (282828) - Note: Background typically set by terminal

	// Foreground colors based on Gruvbox Dark palette
	colorRedFg     = "\033[38;5;167m" // Red (fb4934)
	colorGreenFg   = "\033[38;5;142m" // Green (b8bb26)
	colorYellowFg  = "\033[38;5;214m" // Yellow (fabd2f)
	colorBlueFg    = "\033[38;5;109m" // Blue (83a598)
	colorMagentaFg = "\033[38;5;175m" // Magenta (d3869b)
	colorCyanFg    = "\033[38;5;108m" // Cyan (8ec07c)
	colorOrangeFg  = "\033[38;5;208m" // Orange (fe8019)

	// Text attributes
	colorBold = "\033[1m"
)

// Simpler 16-color ANSI versions (Bright) - can be used as fallbacks or for simplicity
const (
	ansiReset   = "\033[0m" // Same as colorReset
	ansiRed     = "\033[91m"
	ansiGreen   = "\033[92m"
	ansiYellow  = "\033[93m"
	ansiBlue    = "\033[94m"
	ansiMagenta = "\033[95m"
	ansiCyan    = "\033[96m"
	ansiWhite   = "\033[97m" // For default bright text
	ansiBold    = "\033[1m"  // Same as colorBold
)

// Colorize wraps text with a given ANSI color code. It does NOT add a reset code.
// The reset is expected to be handled by the calling print function at the end of the full line.
func Colorize(text string, colorCode string) string {
	return colorCode + text + colorFgDefault // Return to default fg after this specific color
}

// PrintError prints a formatted error message to os.Stderr in red.
// Prefix: "Error: "
func PrintError(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", colorRedFg, message, colorReset)
}

// PrintSuccess prints a formatted success message to os.Stdout in green.
// Prefix: "Success: "
func PrintSuccess(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("%sSuccess: %s%s\n", colorGreenFg, message, colorReset)
}

// PrintInfo prints a formatted informational message to os.Stdout in yellow.
// Prefix: "Info: "
func PrintInfo(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Printf("%sInfo: %s%s\n", colorYellowFg, message, colorReset)
}

// PrintMessage prints a formatted message to os.Stdout.
// It applies colorFgDefault to the base message and ensures a reset at the end.
// Arguments can be pre-colorized using Colorize.
func PrintMessage(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	// Ensure the base of the message is in default fg, then reset everything.
	fmt.Printf("%s%s%s\n", colorFgDefault, message, colorReset)
}

// PrintHeader prints a formatted header message to os.Stdout, bold and in blue.
// Format: "== <message> =="
func PrintHeader(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	// Header is bold blue, then reset. The message itself is part of this.
	fmt.Printf("%s%s== %s ==%s\n", colorBold, colorBlueFg, message, colorReset)
}

// PrintUsageMessage prints a general usage line, typically to Stderr. Uses default foreground color.
func PrintUsageMessage(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	// Usage messages often go to Stderr for consistency with tool output conventions
	fmt.Fprintf(os.Stderr, "%s%s%s\n", colorFgDefault, message, colorReset)
}
