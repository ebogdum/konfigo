// Package merger provides configuration data merging capabilities with precedence rules.
//
// This package handles merging multiple configuration maps with support for:
// - Case-sensitive and case-insensitive key merging
// - Immutable path protection (prevents overwriting critical values)
// - Deep recursive merging of nested maps and arrays
// - Merge conflict detection and resolution
//
// Merge Rules:
// - Later sources override earlier sources (left-to-right precedence)
// - Immutable paths are protected from being overwritten
// - Arrays are replaced entirely, not merged element by element
// - Maps are merged recursively
//
// Usage:
//
//	base := map[string]interface{}{"db": map[string]interface{}{"host": "localhost"}}
//	override := map[string]interface{}{"db": map[string]interface{}{"port": 5432}}
//	merger.Merge(base, override, true, nil)
//	// Result: {"db": {"host": "localhost", "port": 5432}}
package merger

import (
	"konfigo/internal/logger"
	"strings"
)

// Merge recursively merges a source map into a destination map, respecting immutable paths.
func Merge(dst, src map[string]interface{}, caseSensitive bool, immutablePaths map[string]struct{}) {
	if caseSensitive {
		mergeCaseSensitive(dst, src, "", immutablePaths)
	} else {
		mergeCaseInsensitive(dst, src, "", immutablePaths)
	}
}

// mergeCaseSensitive performs a merge where keys are matched exactly.
func mergeCaseSensitive(dst, src map[string]interface{}, path string, immutablePaths map[string]struct{}) {
	for key, srcVal := range src {
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}

		// Check for immutability
		if _, isImmutable := immutablePaths[currentPath]; isImmutable {
			if _, exists := dst[key]; exists {
				logger.Debug("  - Skipping overwrite of immutable key: %s", currentPath)
				continue
			}
		}

		if dstVal, ok := dst[key]; ok {
			if dstMap, dstOk := dstVal.(map[string]interface{}); dstOk {
				if srcMap, srcOk := srcVal.(map[string]interface{}); srcOk {
					mergeCaseSensitive(dstMap, srcMap, currentPath, immutablePaths)
					continue
				}
			}
		}
		dst[key] = srcVal
	}
}

// mergeCaseInsensitive performs a merge that ignores key casing for matching.
func mergeCaseInsensitive(dst, src map[string]interface{}, path string, immutablePaths map[string]struct{}) {
	for srcKey, srcVal := range src {
		existingDstKey, found := findCaseInsensitiveKey(dst, srcKey)
		currentPath := srcKey
		if path != "" {
			currentPath = path + "." + srcKey
		}

		// Check for immutability
		if found {
			// Use the original path from dst for the immutability check
			immutableCheckPath := existingDstKey
			if path != "" {
				immutableCheckPath = path + "." + existingDstKey
			}
			if _, isImmutable := immutablePaths[immutableCheckPath]; isImmutable {
				logger.Debug("  - Skipping overwrite of immutable key: %s", immutableCheckPath)
				continue
			}
		}

		if !found {
			dst[srcKey] = srcVal
			continue
		}

		dstVal := dst[existingDstKey]
		delete(dst, existingDstKey)

		dstMap, dstOk := dstVal.(map[string]interface{})
		srcMap, srcOk := srcVal.(map[string]interface{})

		if dstOk && srcOk {
			mergeCaseInsensitive(dstMap, srcMap, currentPath, immutablePaths)
			dst[srcKey] = dstMap
		} else {
			dst[srcKey] = srcVal
		}
	}
}

// findCaseInsensitiveKey iterates over a map's keys and returns the original key
// if a case-insensitive match is found.
func findCaseInsensitiveKey(m map[string]interface{}, key string) (string, bool) {
	for k := range m {
		if strings.EqualFold(k, key) {
			return k, true
		}
	}
	return "", false
}
