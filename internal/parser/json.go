package parser

import (
	"encoding/json"
	"fmt"
)

// JSONParser handles JSON format parsing.
type JSONParser struct{}

// Parse parses JSON content.
func (jp *JSONParser) Parse(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return data, nil
}

// Format returns the format name.
func (jp *JSONParser) Format() string {
	return "json"
}
