// Package reader provides file and stream reading capabilities for configuration data.
//
// This package handles reading configuration data from various sources including:
// - Local files with automatic format detection
// - Standard input (stdin) with explicit format specification
// - Directory discovery, optionally recursive
//
// Reads are size-bounded (see maxFileSize) to keep a malformed or oversized
// input from exhausting memory.
//
// Supported Sources:
// - Local file paths
// - Standard input (stdin)
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
	"io"
	"konfigo/internal/errors"
	"os"
)

// maxFileSize is the maximum allowed size for configuration files (50 MiB).
const maxFileSize = 50 * 1024 * 1024

// ReadFile reads the contents of a file and returns the content as bytes.
// Files larger than maxFileSize (50 MiB) are rejected to prevent OOM.
//
// The size limit is enforced on the bytes actually read rather than on the
// size reported by stat. Stat alone is not sufficient: it reports 0 for FIFOs,
// character devices and other non-regular files (so an unbounded stream would
// slip past a size-only check), and a regular file can also grow between the
// stat and the read.
func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to stat file")
	}
	defer f.Close()

	// Stat the open handle so the checks below describe the file actually being
	// read, not whatever the path resolved to a moment earlier.
	info, err := f.Stat()
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to stat file")
	}
	if info.IsDir() {
		return nil, errors.FileError(filePath, fmt.Errorf("path is a directory"), "cannot read directory as file")
	}
	if info.Mode().IsRegular() && info.Size() > maxFileSize {
		return nil, errors.FileError(filePath, fmt.Errorf("file size %d exceeds limit %d", info.Size(), maxFileSize), "file too large")
	}

	// Read at most one byte beyond the limit so an oversize input is detected
	// without buffering the whole thing.
	content, err := io.ReadAll(io.LimitReader(f, maxFileSize+1))
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to read file")
	}
	if len(content) > maxFileSize {
		return nil, errors.FileError(filePath, fmt.Errorf("content exceeds limit %d", maxFileSize), "file too large")
	}
	return content, nil
}
