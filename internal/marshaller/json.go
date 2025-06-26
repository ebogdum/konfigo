package marshaller

import (
	"encoding/json"
)

// JSONMarshaller handles JSON format marshalling.
type JSONMarshaller struct{}

// Marshal marshals data to JSON format.
func (jm *JSONMarshaller) Marshal(data map[string]interface{}) ([]byte, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}
	// Add trailing newline for consistency with other formats
	return append(bytes, '\n'), nil
}

// Format returns the format name.
func (jm *JSONMarshaller) Format() string {
	return "json"
}
