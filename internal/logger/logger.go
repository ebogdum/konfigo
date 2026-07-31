// Package logger provides structured logging functionality for the Konfigo application.
// It supports different log levels and conditional logging based on verbosity settings.
package logger

import (
	"log"
	"os"
	"sync/atomic"
)

var (
	verboseFlag atomic.Bool
	quietFlag   atomic.Bool
	logWriter   = log.New(os.Stderr, "", log.LstdFlags)
)

// Init initializes the logger with the desired verbosity settings.
// When debug is true, DEBUG-level messages are printed.
// When quiet is true, INFO-level messages are suppressed.
func Init(debug, quiet bool) {
	verboseFlag.Store(debug)
	quietFlag.Store(quiet)
}

// Log prints a standard log message, unless in quiet mode.
func Log(format string, v ...interface{}) {
	if !quietFlag.Load() {
		logWriter.Printf(format, v...)
	}
}

// Debug prints a verbose/debug message, only if in verbose mode.
func Debug(format string, v ...interface{}) {
	if verboseFlag.Load() {
		logWriter.Printf("DEBUG: "+format, v...)
	}
}

// Warn prints a warning message unless in quiet mode.
func Warn(format string, v ...interface{}) {
	if !quietFlag.Load() {
		logWriter.Printf("WARN: "+format, v...)
	}
}

// Error prints an error message unconditionally.
func Error(format string, v ...interface{}) {
	logWriter.Printf("ERROR: "+format, v...)
}
