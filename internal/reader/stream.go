package reader

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ReadStdin reads content from standard input.
// It validates that stdin is actually a pipe and not a terminal.
func ReadStdin() ([]byte, error) {
	// Check if stdin is a terminal (not a pipe)
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}
	
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("stdin is a terminal, not a pipe")
	}
	
	// Read all content from stdin
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}
	
	return content, nil
}

// IsStdinAvailable checks if there is data available on stdin without reading it.
func IsStdinAvailable() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	
	// Check if stdin is not a terminal (i.e., it's a pipe or redirect)
	return (info.Mode() & os.ModeCharDevice) == 0
}

// ValidateStdinFormat ensures that when reading from stdin, a format is specified.
func ValidateStdinFormat(formatOverride string) error {
	if formatOverride == "" {
		return errors.New("reading from stdin requires an input format flag (-sj, -sy, -st, or -se)")
	}
	return nil
}
