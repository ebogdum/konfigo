package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAMLParser handles YAML format parsing.
type YAMLParser struct{}

// Parse parses YAML content.
func (yp *YAMLParser) Parse(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return data, nil
}

// Format returns the format name.
func (yp *YAMLParser) Format() string {
	return "yaml"
}
