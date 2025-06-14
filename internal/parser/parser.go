package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v3"
)

// Parse takes file content and parses it into a map.
// It uses the formatOverride if provided, otherwise it detects the format
// from the filePath extension.
func Parse(filePath string, content []byte, formatOverride string) (map[string]interface{}, error) {
	format := formatOverride
	if format == "" {
		format = strings.TrimPrefix(filepath.Ext(filePath), ".")
	}

	switch strings.ToLower(format) {
	case "json":
		return parseJSON(content)
	case "yaml", "yml":
		return parseYAML(content)
	case "toml":
		return parseTOML(content)
	case "ini":
		return parseINI(content)
	case "env":
		return parseENV(content)
	default:
		return nil, fmt.Errorf("unsupported file format: %s for file %s", format, filePath)
	}
}

// ... (rest of the file: parseJSON, parseYAML, etc. remain exactly the same)
func parseJSON(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return data, nil
}

func parseYAML(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return data, nil
}

func parseTOML(content []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}
	return data, nil
}

func parseINI(content []byte) (map[string]interface{}, error) {
	cfg, err := ini.Load(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse INI: %w", err)
	}

	data := make(map[string]interface{})
	defaultSection := cfg.Section(ini.DefaultSection)
	for _, key := range defaultSection.Keys() {
		data[key.Name()] = key.Value()
	}

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

func parseENV(content []byte) (map[string]interface{}, error) {
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

		if unquoted, err := strconv.Unquote(value); err == nil {
			value = unquoted
		}

		setNestedValue(data, key, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .env content: %w", err)
	}

	return data, nil
}

func setNestedValue(data map[string]interface{}, key string, value interface{}) {
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
