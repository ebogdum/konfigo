package parser

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// ENVParser handles ENV format parsing.
type ENVParser struct{}

// Parse parses ENV content.
func (ep *ENVParser) Parse(content []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	scanner.Buffer(make([]byte, 1<<20), 1<<20) // 1 MiB buffer for large values

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Unquote if the value is quoted
		if unquoted, err := strconv.Unquote(value); err == nil {
			value = unquoted
		} else {
			// Strip inline comments (only for unquoted values)
			if idx := strings.Index(value, " #"); idx >= 0 {
				value = strings.TrimSpace(value[:idx])
			}
		}

		if err := ep.setNestedValue(data, key, value); err != nil {
			return nil, fmt.Errorf("error processing key %q: %w", key, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .env content: %w", err)
	}

	return data, nil
}

// Format returns the format name.
func (ep *ENVParser) Format() string {
	return "env"
}

// setNestedValue sets a nested value in the data map using dot or underscore notation.
// Note: ENV files use dots for nesting (e.g., SERVICE.HOST=x -> {SERVICE:{HOST:x}})
// while the ENV marshaller uses underscores. This means ENV round-trips are not lossless
// for keys containing dots. This is by design: dots in ENV keys are uncommon, and the
// nesting behavior matches tools like Spring Boot and Quarkus.
func (ep *ENVParser) setNestedValue(data map[string]interface{}, key string, value interface{}) error {
	keys := strings.Split(key, ".")
	currentMap := data

	for i, k := range keys {
		if i == len(keys)-1 {
			currentMap[k] = value
			return nil
		}

		existing, exists := currentMap[k]
		if !exists {
			newMap := make(map[string]interface{})
			currentMap[k] = newMap
			currentMap = newMap
			continue
		}

		if nextMap, ok := existing.(map[string]interface{}); ok {
			currentMap = nextMap
		} else {
			return fmt.Errorf("key conflict: %q is already a scalar value, cannot create nested key under it", k)
		}
	}
	return nil
}
