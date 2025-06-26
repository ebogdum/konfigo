package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Input/Output errors
	ErrorTypeFileRead       ErrorType = "FILE_READ"
	ErrorTypeFileWrite      ErrorType = "FILE_WRITE"
	ErrorTypeStdinRead      ErrorType = "STDIN_READ"
	
	// Parsing errors
	ErrorTypeParsing        ErrorType = "PARSING"
	ErrorTypeFormatDetect   ErrorType = "FORMAT_DETECT"
	ErrorTypeInvalidFormat  ErrorType = "INVALID_FORMAT"
	
	// Schema errors
	ErrorTypeSchemaLoad     ErrorType = "SCHEMA_LOAD"
	ErrorTypeSchemaProcess  ErrorType = "SCHEMA_PROCESS"
	ErrorTypeValidation     ErrorType = "VALIDATION"
	
	// Variable errors
	ErrorTypeVarResolution  ErrorType = "VAR_RESOLUTION"
	ErrorTypeVarSubstitute  ErrorType = "VAR_SUBSTITUTE"
	
	// Configuration errors
	ErrorTypeConfigMerge    ErrorType = "CONFIG_MERGE"
	ErrorTypeImmutableField ErrorType = "IMMUTABLE_FIELD"
	
	// CLI errors
	ErrorTypeCLIFlag        ErrorType = "CLI_FLAG"
	ErrorTypeCLIValidation  ErrorType = "CLI_VALIDATION"
	
	// Internal errors
	ErrorTypeInternal       ErrorType = "INTERNAL"
	ErrorTypeDeepCopy       ErrorType = "DEEP_COPY"
)

// KonfigoError represents a structured error with context
type KonfigoError struct {
	Type        ErrorType
	Message     string
	Path        string    // Configuration path (e.g., "database.host")
	FilePath    string    // File path where error occurred
	Line        int       // Line number (if applicable)
	Column      int       // Column number (if applicable)
	Cause       error     // Underlying error
	StackTrace  []string  // Stack trace
	Context     map[string]interface{} // Additional context
}

// Error implements the error interface
func (e *KonfigoError) Error() string {
	var parts []string
	
	// Add error type
	parts = append(parts, fmt.Sprintf("[%s]", e.Type))
	
	// Add file context if available
	if e.FilePath != "" {
		if e.Line > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", e.FilePath, e.Line))
		} else {
			parts = append(parts, e.FilePath)
		}
	}
	
	// Add config path if available
	if e.Path != "" {
		parts = append(parts, fmt.Sprintf("path:%s", e.Path))
	}
	
	// Add main message
	parts = append(parts, e.Message)
	
	// Join all parts
	result := strings.Join(parts, " ")
	
	// Add cause if present
	if e.Cause != nil {
		result += fmt.Sprintf(" (caused by: %v)", e.Cause)
	}
	
	return result
}

// Unwrap returns the underlying cause
func (e *KonfigoError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a specific type
func (e *KonfigoError) Is(target error) bool {
	if targetErr, ok := target.(*KonfigoError); ok {
		return e.Type == targetErr.Type
	}
	return false
}

// WithContext adds additional context to the error
func (e *KonfigoError) WithContext(key string, value interface{}) *KonfigoError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// GetContext retrieves context value
func (e *KonfigoError) GetContext(key string) (interface{}, bool) {
	if e.Context == nil {
		return nil, false
	}
	value, exists := e.Context[key]
	return value, exists
}

// NewError creates a new KonfigoError
func NewError(errorType ErrorType, message string) *KonfigoError {
	return &KonfigoError{
		Type:       errorType,
		Message:    message,
		StackTrace: captureStackTrace(2), // Skip this function and the caller
	}
}

// NewErrorf creates a new KonfigoError with formatted message
func NewErrorf(errorType ErrorType, format string, args ...interface{}) *KonfigoError {
	return &KonfigoError{
		Type:       errorType,
		Message:    fmt.Sprintf(format, args...),
		StackTrace: captureStackTrace(2),
	}
}

// WrapError wraps an existing error with KonfigoError context
func WrapError(errorType ErrorType, message string, cause error) *KonfigoError {
	return &KonfigoError{
		Type:       errorType,
		Message:    message,
		Cause:      cause,
		StackTrace: captureStackTrace(2),
	}
}

// WrapErrorf wraps an existing error with formatted message
func WrapErrorf(errorType ErrorType, cause error, format string, args ...interface{}) *KonfigoError {
	return &KonfigoError{
		Type:       errorType,
		Message:    fmt.Sprintf(format, args...),
		Cause:      cause,
		StackTrace: captureStackTrace(2),
	}
}

// FileError creates a file-related error
func FileError(filePath string, cause error, message string) *KonfigoError {
	return &KonfigoError{
		Type:       ErrorTypeFileRead,
		Message:    message,
		FilePath:   filePath,
		Cause:      cause,
		StackTrace: captureStackTrace(2),
	}
}

// ValidationError creates a validation error
func ValidationError(path string, message string) *KonfigoError {
	return &KonfigoError{
		Type:       ErrorTypeValidation,
		Message:    message,
		Path:       path,
		StackTrace: captureStackTrace(2),
	}
}

// ParsingError creates a parsing error with line/column information
func ParsingError(filePath string, line, column int, message string, cause error) *KonfigoError {
	return &KonfigoError{
		Type:       ErrorTypeParsing,
		Message:    message,
		FilePath:   filePath,
		Line:       line,
		Column:     column,
		Cause:      cause,
		StackTrace: captureStackTrace(2),
	}
}

// ConfigError creates a configuration-related error
func ConfigError(path string, message string) *KonfigoError {
	return &KonfigoError{
		Type:       ErrorTypeConfigMerge,
		Message:    message,
		Path:       path,
		StackTrace: captureStackTrace(2),
	}
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) []string {
	var trace []string
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip+1, pc)
	frames := runtime.CallersFrames(pc[:n])
	
	for {
		frame, more := frames.Next()
		if !strings.Contains(frame.File, "runtime/") {
			trace = append(trace, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
		}
		if !more {
			break
		}
	}
	
	return trace
}

// IsType checks if an error is of a specific KonfigoError type
func IsType(err error, errorType ErrorType) bool {
	if kErr, ok := err.(*KonfigoError); ok {
		return kErr.Type == errorType
	}
	return false
}

// GetType returns the error type of a KonfigoError, or empty string for other errors
func GetType(err error) ErrorType {
	if kErr, ok := err.(*KonfigoError); ok {
		return kErr.Type
	}
	return ""
}

// FormatUserFriendly formats an error for user-friendly display
func FormatUserFriendly(err error) string {
	kErr, ok := err.(*KonfigoError)
	if !ok {
		return err.Error()
	}
	
	switch kErr.Type {
	case ErrorTypeFileRead:
		return fmt.Sprintf("Cannot read file '%s': %s", kErr.FilePath, kErr.Message)
	case ErrorTypeFileWrite:
		return fmt.Sprintf("Cannot write file '%s': %s", kErr.FilePath, kErr.Message)
	case ErrorTypeParsing:
		if kErr.Line > 0 {
			return fmt.Sprintf("Parsing error in '%s' at line %d: %s", kErr.FilePath, kErr.Line, kErr.Message)
		}
		return fmt.Sprintf("Parsing error in '%s': %s", kErr.FilePath, kErr.Message)
	case ErrorTypeValidation:
		return fmt.Sprintf("Validation failed for '%s': %s", kErr.Path, kErr.Message)
	case ErrorTypeVarResolution:
		return fmt.Sprintf("Variable resolution failed: %s", kErr.Message)
	case ErrorTypeCLIFlag:
		return fmt.Sprintf("Invalid command line option: %s", kErr.Message)
	default:
		return kErr.Error()
	}
}
