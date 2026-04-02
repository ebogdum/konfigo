// Package logger provides structured logging functionality for the Konfigo application.
// It supports different log levels and conditional logging based on verbosity settings.
package logger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
)

// LogLevel represents the severity of a log message.
type LogLevel int

const (
	// DEBUG level for detailed debugging information
	DEBUG LogLevel = iota
	// INFO level for general informational messages
	INFO
	// WARN level for warning messages that don't stop execution
	WARN
	// ERROR level for error messages
	ERROR
)

// String returns the string representation of a log level.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

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

// LogAtLevel logs a message at the specified level.
func LogAtLevel(level LogLevel, format string, v ...interface{}) {
	switch level {
	case DEBUG:
		Debug(format, v...)
	case INFO:
		Log(format, v...)
	case WARN:
		Warn(format, v...)
	case ERROR:
		Error(format, v...)
	}
}

// IsVerbose returns true if verbose logging is enabled.
func IsVerbose() bool {
	return verboseFlag.Load()
}

// IsQuiet returns true if quiet mode is enabled.
func IsQuiet() bool {
	return quietFlag.Load()
}

// SetOutput sets the output destination for the logger.
func SetOutput(output *os.File) {
	logWriter.SetOutput(output)
}

// Logf is a convenience function for formatted logging.
func Logf(level LogLevel, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	LogAtLevel(level, "%s", message)
}
