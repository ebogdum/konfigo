// Package writer handles all output operations for configuration files.
// This package centralizes file writing, directory creation, and output target management.
package writer

import (
	"konfigo/internal/errors"
	"os"
)

// WriteFile writes content to a file with proper error handling.
// It ensures the directory exists before writing the file.
func WriteFile(filePath string, content []byte) error {
	// Ensure the directory exists
	if err := EnsureDirectory(filePath); err != nil {
		return errors.WrapError(errors.ErrorTypeFileWrite, "failed to ensure directory for file", err).WithContext("file", filePath)
	}
	
	// Write the file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return errors.WrapError(errors.ErrorTypeFileWrite, "failed to write file", err).WithContext("file", filePath)
	}
	
	return nil
}

// WriteToStdout writes content to standard output.
func WriteToStdout(content []byte) error {
	_, err := os.Stdout.Write(content)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeFileWrite, "failed to write to stdout", err)
	}
	return nil
}

// WriteMultipleFiles writes content to multiple files.
// It returns an error if any file fails to write.
func WriteMultipleFiles(files map[string][]byte) error {
	for filePath, content := range files {
		if err := WriteFile(filePath, content); err != nil {
			return errors.WrapError(errors.ErrorTypeFileWrite, "failed to write file", err).WithContext("file", filePath)
		}
	}
	return nil
}
