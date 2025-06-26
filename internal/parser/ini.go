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
