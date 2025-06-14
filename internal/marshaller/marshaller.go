package marshaller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Marshal takes the final merged data and a format string, returning the
// data as a byte slice in the specified format.
func Marshal(data map[string]interface{}, format string) ([]byte, error) {
	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(data, "", "  ")
	case "yaml", "yml":
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(2)
		if err := encoder.Encode(data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "toml":
		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case "env":
		return marshalENV(data)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

// marshalENV flattens a nested map into a .env file format.
// Nested keys are joined with an underscore and converted to uppercase.
func marshalENV(data map[string]interface{}) ([]byte, error) {
	var lines []string
	flattened := make(map[string]string)
	flattenMap("", data, flattened)

	// Sort keys for deterministic output
	keys := make([]string, 0, len(flattened))
	for k := range flattened {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s", k, flattened[k]))
	}

	return []byte(strings.Join(lines, "\n")), nil
}

// flattenMap is a recursive helper to flatten the map for .env output.
func flattenMap(prefix string, data map[string]interface{}, flattened map[string]string) {
	for k, v := range data {
		var sb strings.Builder
		if prefix != "" {
			sb.WriteString(prefix)
			sb.WriteString("_")
		}
		sb.WriteString(k)
		newKey := strings.ToUpper(sb.String())

		switch val := v.(type) {
		case map[string]interface{}:
			// If it's a map, recurse
			flattenMap(newKey, val, flattened)
		case string:
			// If it's a string that needs quoting, add quotes
			if strings.Contains(val, " ") || strings.Contains(val, "#") {
				flattened[newKey] = fmt.Sprintf("%q", val)
			} else {
				flattened[newKey] = val
			}
		default:
			// For other types (int, bool, etc.), convert to string
			flattened[newKey] = fmt.Sprintf("%v", v)
		}
	}
}
