package debug

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

// os independent line break
var lineBreak string

// Standard output (modified during unit testing)
var out io.Writer = os.Stdout

// Init os independent lineBreak from runtime.GOOS
func init() {
	switch runtime.GOOS {
	case "windows":
		lineBreak = "\r\n"
	default:
		lineBreak = "\n"
	}
}

// Dump variable with type for debug
func Dump(variable any) {
	_, _ = fmt.Fprintf(out, "%T: %+v%s", variable, variable, lineBreak)
}
