package writer

import (
	"konfigo/internal/errors"
	"os"
	"path/filepath"
)

// EnsureDirectory ensures that the directory for the given file path exists.
// It creates the directory structure recursively if it doesn't exist.
func EnsureDirectory(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "/" {
		return nil // No directory to create
	}

	// Create directory recursively (MkdirAll is idempotent for existing directories)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WrapError(errors.ErrorTypeFileWrite, "failed to create directory", err).WithContext("directory", dir)
	}

	return nil
}

// CreateDirectory creates a directory at the specified path.
// It creates parent directories as needed.
func CreateDirectory(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return errors.WrapError(errors.ErrorTypeFileWrite, "failed to create directory", err).WithContext("directory", dirPath)
	}
	return nil
}

// DirectoryExists checks if a directory exists.
func DirectoryExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
