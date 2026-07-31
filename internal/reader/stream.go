package reader

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// ReadStdin reads content from standard input.
// It validates that stdin is actually a pipe and not a terminal.
//
// The same maxFileSize ceiling that applies to file inputs applies here: stdin
// is an unbounded stream, so reading it without a limit would let a large or
// endless pipe exhaust memory even though the equivalent file would be rejected.
func ReadStdin() ([]byte, error) {
	// Check if stdin is a terminal (not a pipe)
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat stdin: %w", err)
	}

	if (info.Mode() & os.ModeCharDevice) != 0 {
		return nil, errors.New("stdin is a terminal, not a pipe")
	}

	// Read one byte past the limit so an oversize stream is detected without
	// buffering all of it.
	content, err := io.ReadAll(io.LimitReader(os.Stdin, maxFileSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %w", err)
	}
	if len(content) > maxFileSize {
		return nil, fmt.Errorf("stdin input exceeds limit %d bytes", maxFileSize)
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
