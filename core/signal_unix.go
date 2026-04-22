//go:build !windows
// +build !windows

package core

import (
	"os"
	"os/signal"
	"syscall"
)

// SetupSignalHandler sets up signal handling for SIGWINCH (terminal resize)
// This function is not available on Windows
func SetupSignalHandler(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGWINCH)
}