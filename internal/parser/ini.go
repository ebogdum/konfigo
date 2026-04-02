package parser

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// INIParser handles INI format parsing.
type INIParser struct{}

// Parse parses INI content.
func (ip *INIParser) Parse(content []byte) (map[string]interface{}, error) {
	cfg, err := ini.Load(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse INI: %w", err)
	}

	data := make(map[string]interface{})

	// Handle default section
	defaultSection := cfg.Section(ini.DefaultSection)
	for _, key := range defaultSection.Keys() {
		data[key.Name()] = key.Value()
	}

	// Handle named sections
	for _, section := range cfg.Sections() {
		if section.Name() == ini.DefaultSection {
			continue
		}
		if _, exists := data[section.Name()]; exists {
			return nil, fmt.Errorf("INI section name %q collides with a key from the default section", section.Name())
		}
		sectionMap := make(map[string]interface{})
		for _, key := range section.Keys() {
			sectionMap[key.Name()] = key.Value()
		}
		data[section.Name()] = sectionMap
	}

	return data, nil
}

// Format returns the format name.
func (ip *INIParser) Format() string {
	return "ini"
}
