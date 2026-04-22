//go:build windows
// +build windows

package core

import (
	"os"
)

// SetupSignalHandler is a no-op on Windows as SIGWINCH is not available
func SetupSignalHandler(sigChan chan os.Signal) {
	// Windows doesn't support SIGWINCH, so we do nothing
}