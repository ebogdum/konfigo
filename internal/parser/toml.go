package parser

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// TOMLParser handles TOML format parsing.
type TOMLParser struct{}

// Parse parses TOML content.
func (tp *TOMLParser) Parse(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}
	return data, nil
}

// Format returns the format name.
func (tp *TOMLParser) Format() string {
	return "toml"
}
