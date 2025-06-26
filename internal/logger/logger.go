// Package logger provides structured logging functionality for the Konfigo application.
// It supports different log levels and conditional logging based on verbosity settings.
package logger

import (
	"fmt"
	"log"
	"os"
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
	isVerbose = false
	isQuiet   = false
	logger    = log.New(os.Stderr, "", log.LstdFlags)
)

// Init initializes the logger with the desired verbosity settings.
func Init(verbose, quiet bool) {
	isVerbose = verbose
	isQuiet = quiet
}

// Log prints a standard log message, unless in quiet mode.
func Log(format string, v ...interface{}) {
	if !isQuiet {
		logger.Printf(format, v...)
	}
}

// Debug prints a verbose/debug message, only if in verbose mode.
func Debug(format string, v ...interface{}) {
	if isVerbose {
		logger.Printf("DEBUG: "+format, v...)
	}
}

// Warn prints a warning message, only if in verbose mode.
func Warn(format string, v ...interface{}) {
	if isVerbose {
		logger.Printf("WARN: "+format, v...)
	}
}

// Error prints an error message unconditionally.
func Error(format string, v ...interface{}) {
	logger.Printf("ERROR: "+format, v...)
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
	return isVerbose
}

// IsQuiet returns true if quiet mode is enabled.
func IsQuiet() bool {
	return isQuiet
}

// SetOutput sets the output destination for the logger.
func SetOutput(output *os.File) {
	logger.SetOutput(output)
}

// Logf is a convenience function for formatted logging.
func Logf(level LogLevel, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	LogAtLevel(level, "%s", message)
}
