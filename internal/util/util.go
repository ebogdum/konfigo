package util

import (
	"encoding/json" // For deep copy fallback
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
	if len(keys) == 0 {
		return
	}
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
// Uses optimized native copying for better performance, falls back to JSON for complex types.
func DeepCopyMap(originalMap map[string]interface{}) (map[string]interface{}, error) {
	if originalMap == nil {
		return nil, nil
	}
	
	// Try native deep copy first (much faster)
	copiedMap, err := deepCopyMapNative(originalMap)
	if err == nil {
		return copiedMap, nil
	}
	
	// Fallback to JSON-based copy for complex types
	bytes, err := json.Marshal(originalMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map for deep copy: %w", err)
	}
	var jsonCopiedMap map[string]interface{}
	err = json.Unmarshal(bytes, &jsonCopiedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal map for deep copy: %w", err)
	}
	return jsonCopiedMap, nil
}

// deepCopyMapNative performs native deep copy for common configuration data types.
// Returns error if encounters unsupported types that require JSON fallback.
func deepCopyMapNative(originalMap map[string]interface{}) (map[string]interface{}, error) {
	copiedMap := make(map[string]interface{}, len(originalMap))
	
	for key, value := range originalMap {
		copiedValue, err := deepCopyValue(value)
		if err != nil {
			return nil, err // Trigger JSON fallback
		}
		copiedMap[key] = copiedValue
	}
	
	return copiedMap, nil
}

// deepCopyValue recursively copies a value, returning error for unsupported types
func deepCopyValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case bool:
		return v, nil
	case int:
		return v, nil
	case int32:
		return v, nil
	case int64:
		return v, nil
	case float32:
		return v, nil
	case float64:
		return v, nil
	case string:
		return v, nil
	case map[string]interface{}:
		copiedMap := make(map[string]interface{}, len(v))
		for k, val := range v {
			copiedVal, err := deepCopyValue(val)
			if err != nil {
				return nil, err
			}
			copiedMap[k] = copiedVal
		}
		return copiedMap, nil
	case []interface{}:
		copiedSlice := make([]interface{}, len(v))
		for i, val := range v {
			copiedVal, err := deepCopyValue(val)
			if err != nil {
				return nil, err
			}
			copiedSlice[i] = copiedVal
		}
		return copiedSlice, nil
	default:
		// Unsupported type, trigger JSON fallback
		return nil, fmt.Errorf("unsupported type for native copy: %T", v)
	}
}

// StringsPool provides a simple string pool to reduce memory allocations
type StringsPool struct {
	pool map[string]string
}

// NewStringsPool creates a new string pool
func NewStringsPool() *StringsPool {
	return &StringsPool{
		pool: make(map[string]string),
	}
}

// Get returns a pooled string, adding it if not present
func (sp *StringsPool) Get(s string) string {
	if pooled, exists := sp.pool[s]; exists {
		return pooled
	}
	sp.pool[s] = s
	return s
}

// Clear clears the pool to free memory
func (sp *StringsPool) Clear() {
	sp.pool = make(map[string]string)
}
