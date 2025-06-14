package loader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// supportedExtensions is a set of file extensions the tool can parse.
var supportedExtensions = map[string]struct{}{
	".json": {},
	".yaml": {},
	".yml":  {},
	".toml": {},
	".ini":  {},
	".env":  {},
}

// LoadFromPath discovers and returns a sorted list of configuration file paths.
// If recursive is false, it only reads files from the top level of the given directory.
// If recursive is true, it walks the entire directory tree.
// If the root path is a file, it returns a slice containing only that file.
func LoadFromPath(rootPath string, recursive bool) ([]string, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat path %s: %w", rootPath, err)
	}

	// If the path is a single file, the recursive flag has no effect.
	if !info.IsDir() {
		if _, ok := supportedExtensions[filepath.Ext(rootPath)]; ok {
			return []string{rootPath}, nil
		}
		return nil, fmt.Errorf("unsupported file type for single file input: %s", rootPath)
	}

	var files []string
	if recursive {
		// Use the original recursive walking method.
		err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if _, ok := supportedExtensions[filepath.Ext(path)]; ok {
				files = append(files, path)
			}
			return nil
		})
	} else {
		// Non-recursive: read only the top-level directory.
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", rootPath, err)
		}
		for _, entry := range entries {
			// Skip subdirectories.
			if entry.IsDir() {
				continue
			}
			if _, ok := supportedExtensions[filepath.Ext(entry.Name())]; ok {
				fullPath := filepath.Join(rootPath, entry.Name())
				files = append(files, fullPath)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error scanning directory %s: %w", rootPath, err)
	}

	// Both WalkDir and ReadDir provide lexicographical order, but we sort explicitly
	// to guarantee consistent behavior across all platforms and scenarios.
	sort.Strings(files)

	return files, nil
}
