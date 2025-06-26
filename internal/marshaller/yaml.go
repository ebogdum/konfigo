package marshaller

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

// YAMLMarshaller handles YAML format marshalling.
type YAMLMarshaller struct{}

// Marshal marshals data to YAML format.
func (ym *YAMLMarshaller) Marshal(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	encoder.Close()
	return buf.Bytes(), nil
}

// Format returns the format name.
func (ym *YAMLMarshaller) Format() string {
	return "yaml"
}
