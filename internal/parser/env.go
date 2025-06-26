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

		// Unquote if necessary
		if unquoted, err := strconv.Unquote(value); err == nil {
			value = unquoted
		}

		ep.setNestedValue(data, key, value)
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

// setNestedValue sets a nested value in the data map using dot notation.
func (ep *ENVParser) setNestedValue(data map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
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
