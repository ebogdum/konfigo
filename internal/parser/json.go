package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JSONParser handles JSON format parsing.
type JSONParser struct{}

// Parse parses JSON content using json.Number to preserve integer fidelity.
func (jp *JSONParser) Parse(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.UseNumber()
	if err := dec.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	normalizeJSONNumbers(data)
	return data, nil
}

// normalizeJSONNumbers walks a parsed JSON map and converts json.Number values
// to int64 (if the number has no fractional part) or float64.
func normalizeJSONNumbers(m map[string]interface{}) {
	for k, v := range m {
		switch val := v.(type) {
		case json.Number:
			if i, err := val.Int64(); err == nil {
				m[k] = i
			} else if f, err := val.Float64(); err == nil {
				m[k] = f
			}
		case map[string]interface{}:
			normalizeJSONNumbers(val)
		case []interface{}:
			normalizeJSONSlice(val)
		}
	}
}

// normalizeJSONSlice walks a slice and converts json.Number values.
func normalizeJSONSlice(s []interface{}) {
	for i, v := range s {
		switch val := v.(type) {
		case json.Number:
			if n, err := val.Int64(); err == nil {
				s[i] = n
			} else if f, err := val.Float64(); err == nil {
				s[i] = f
			}
		case map[string]interface{}:
			normalizeJSONNumbers(val)
		case []interface{}:
			normalizeJSONSlice(val)
		}
	}
}

// Format returns the format name.
func (jp *JSONParser) Format() string {
	return "json"
}
