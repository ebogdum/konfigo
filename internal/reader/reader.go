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
	"io"
	"konfigo/internal/errors"
	"os"
)

// ReadFile reads the contents of a file and returns the content as bytes.
// This centralizes file reading operations and provides consistent error handling.
func ReadFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to read file")
	}
	return content, nil
}

// ReadFiles reads multiple files and returns their contents.
// It returns a map of file paths to their content bytes.
// If any file fails to read, it returns an error for that specific file.
func ReadFiles(filePaths []string) (map[string][]byte, error) {
	contents := make(map[string][]byte)

	for _, filePath := range filePaths {
		content, err := ReadFile(filePath)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeFileRead, "failed to read file", err).WithContext("file", filePath)
		}
		contents[filePath] = content
	}

	return contents, nil
}

// FileExists checks if a file exists and is readable.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// ReadFileBuffered reads a file with buffered I/O for better performance with large files
func ReadFileBuffered(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to open file")
	}
	defer file.Close()

	// Get file size for efficient buffer allocation
	stat, err := file.Stat()
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to stat file")
	}

	// For small files, use regular reading
	if stat.Size() < 1024*1024 { // 1MB
		return ReadFile(filePath)
	}

	// For larger files, use buffered reading
	buffer := make([]byte, stat.Size())
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return nil, errors.FileError(filePath, err, "failed to read file")
	}

	return buffer, nil
}

// ReadFileStream reads a file as a stream for very large files
func ReadFileStream(filePath string, callback func([]byte) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeFileRead, "failed to open file", err).WithContext("file", filePath)
	}
	defer file.Close()

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			if callbackErr := callback(buffer[:n]); callbackErr != nil {
				return callbackErr
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.WrapError(errors.ErrorTypeFileRead, "failed to read file", err).WithContext("file", filePath)
		}
	}

	return nil
}
