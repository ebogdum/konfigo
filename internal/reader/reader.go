// Package reader provides file and stream reading capabilities for configuration data.
//
// This package handles reading configuration data from various sources including:
// - Local files with automatic format detection
// - Standard input (stdin) with explicit format specification
// - Buffered reading for large files
// - Stream processing for continuous data
//
// The reader package supports multiple input formats and provides
// performance optimizations for large configuration files.
//
// Supported Sources:
// - Local file paths
// - Standard input (stdin)
// - Network streams (future)
//
// Usage:
//
//	content, err := reader.ReadFile("/path/to/config.yaml")
//	if err != nil {
//	    return err
//	}
package reader

import (
	"fmt"
	"konfigo/internal/errors"
	"os"
)

// maxFileSize is the maximum allowed size for configuration files (50 MiB).
const maxFileSize = 50 * 1024 * 1024

// ReadFile reads the contents of a file and returns the content as bytes.
// Files larger than maxFileSize (50 MiB) are rejected to prevent OOM.
func ReadFile(filePath string) ([]byte, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to stat file")
	}
	if info.IsDir() {
		return nil, errors.FileError(filePath, fmt.Errorf("path is a directory"), "cannot read directory as file")
	}
	if info.Size() > maxFileSize {
		return nil, errors.FileError(filePath, fmt.Errorf("file size %d exceeds limit %d", info.Size(), maxFileSize), "file too large")
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to read file")
	}
	return content, nil
}


