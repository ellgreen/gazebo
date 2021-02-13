package debug

import (
	"fmt"
	"os"
)

var debugging bool

// Enable turns on debugging
func Enable() {
	debugging = true
}

// Disable turns off debugging
func Disable() {
	debugging = false
}

// Enabled returns whether debugging is turned on
func Enabled() bool {
	return debugging
}

// Printf prints formatted output to os.Stderr
func Printf(format string, args ...interface{}) {
	if Enabled() {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
