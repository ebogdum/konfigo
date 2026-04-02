package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// maxYAMLInputSize is the maximum allowed size for YAML input to prevent
// denial-of-service via alias/anchor expansion (billion-laughs attack).
const maxYAMLInputSize = 10 * 1024 * 1024 // 10 MiB

// YAMLParser handles YAML format parsing.
type YAMLParser struct{}

// Parse parses YAML content.
func (yp *YAMLParser) Parse(content []byte) (map[string]interface{}, error) {
	if len(content) > maxYAMLInputSize {
		return nil, fmt.Errorf("YAML input exceeds maximum allowed size of %d bytes", maxYAMLInputSize)
	}
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
