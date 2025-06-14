package logger

import "log"

var (
	isVerbose = false
	isQuiet   = false
)

// Init initializes the logger with the desired verbosity settings.
func Init(verbose, quiet bool) {
	isVerbose = verbose
	isQuiet = quiet
}

// Log prints a standard log message, unless in quiet mode.
func Log(format string, v ...interface{}) {
	if !isQuiet {
		log.Printf(format, v...)
	}
}

// Debug prints a verbose/debug message, only if in verbose mode.
func Debug(format string, v ...interface{}) {
	if isVerbose {
		log.Printf("DEBUG: "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if isVerbose {
		log.Printf("WARN: "+format, v...)
	}
}
