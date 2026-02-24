package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSupported_LowercaseExtensions(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"config.json", true},
		{"config.yaml", true},
		{"config.yml", true},
		{"config.toml", true},
		{"config.env", true},
		{"config.ini", true},
		{"config.txt", false},
		{"config.xml", false},
	}
	for _, tc := range tests {
		got := IsSupported(tc.path)
		if got != tc.expected {
			t.Errorf("IsSupported(%q) = %v, want %v", tc.path, got, tc.expected)
		}
	}
}

func TestIsSupported_UppercaseExtensions(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"config.JSON", true},
		{"config.YAML", true},
		{"config.YML", true},
		{"config.TOML", true},
		{"config.ENV", true},
		{"config.Json", true},
		{"config.Yaml", true},
	}
	for _, tc := range tests {
		got := IsSupported(tc.path)
		if got != tc.expected {
			t.Errorf("IsSupported(%q) = %v, want %v", tc.path, got, tc.expected)
		}
	}
}

func TestDiscoverFiles_SingleFileUppercase(t *testing.T) {
	// Create a temp file with uppercase extension
	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.JSON")
	if err := os.WriteFile(filePath, []byte(`{"key":"val"}`), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := DiscoverFiles(filePath, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0] != filePath {
		t.Errorf("expected %s, got %s", filePath, files[0])
	}
}

func TestDiscoverFiles_DirectoryMixedCase(t *testing.T) {
	dir := t.TempDir()

	// Create files with various case extensions
	testFiles := []string{"a.json", "b.YAML", "c.Toml", "d.txt"}
	for _, name := range testFiles {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := DiscoverFiles(dir, false)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Should find json, YAML, and Toml but not txt
	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(files), files)
	}
}

func TestDiscoverFiles_UnsupportedSingleFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "config.txt")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := DiscoverFiles(filePath, false)
	if err == nil {
		t.Fatal("expected error for unsupported file type, got nil")
	}
}

func TestGetSupportedExtensions(t *testing.T) {
	exts := GetSupportedExtensions()
	if len(exts) == 0 {
		t.Fatal("expected non-empty list of supported extensions")
	}
	// Check that .json is in the list
	found := false
	for _, ext := range exts {
		if ext == ".json" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected .json in supported extensions")
	}
}
