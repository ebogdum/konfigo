package util

import (
	"encoding/json" // For deep copy
	"fmt"
	"strings"
)

// GetNestedValue retrieves a value from a nested map using a dot-separated path.
func GetNestedValue(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	current := interface{}(data)

	for _, key := range keys {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				current = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return current, true
}

// SetNestedValue sets a value in a nested map using a dot-separated path.
// It creates nested maps as needed.
func SetNestedValue(data map[string]interface{}, path string, value interface{}) {
	keys := strings.Split(path, ".")
	currentMap := data

	for i, k := range keys {
		if i == len(keys)-1 {
			currentMap[k] = value
			return
		}

		if nextMap, ok := currentMap[k].(map[string]interface{}); ok {
			currentMap = nextMap
		} else {
			newMap := make(map[string]interface{})
			currentMap[k] = newMap
			currentMap = newMap
		}
	}
}

// DeleteNestedValue deletes a key from a nested map.
func DeleteNestedValue(data map[string]interface{}, path string) {
	keys := strings.Split(path, ".")
	if len(keys) == 1 {
		delete(data, keys[0])
		return
	}

	parentPath := strings.Join(keys[:len(keys)-1], ".")
	lastKey := keys[len(keys)-1]

	if parent, ok := GetNestedValue(data, parentPath); ok {
		if parentMap, ok := parent.(map[string]interface{}); ok {
			delete(parentMap, lastKey)
		}
	}
}

// WalkAndReplace recursively walks through the map/slice structure and applies the replacer function to all string values.
func WalkAndReplace(data interface{}, replacerFunc func(string) string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		newData := make(map[string]interface{}, len(v))
		for key, val := range v {
			newData[key] = WalkAndReplace(val, replacerFunc)
		}
		return newData
	case []interface{}:
		newData := make([]interface{}, len(v))
		for i, val := range v {
			newData[i] = WalkAndReplace(val, replacerFunc)
		}
		return newData
	case string:
		return replacerFunc(v)
	default:
		return v
	}
}

// DeepCopyMap creates a deep copy of a map[string]interface{}.
func DeepCopyMap(originalMap map[string]interface{}) (map[string]interface{}, error) {
	if originalMap == nil {
		return nil, nil
	}
	// A common way to deep copy arbitrary structures in Go is to marshal and unmarshal them,
	// typically using JSON or another format like Gob if more Go-specific types are involved.
	// JSON is suitable here as config data is generally JSON-compatible.
	bytes, err := json.Marshal(originalMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map for deep copy: %w", err)
	}
	var copiedMap map[string]interface{}
	err = json.Unmarshal(bytes, &copiedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map for deep copy: %w", err)
	}
	return copiedMap, nil
}
