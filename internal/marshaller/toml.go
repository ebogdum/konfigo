package marshaller

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// TOMLMarshaller handles TOML format marshalling.
type TOMLMarshaller struct{}

// Marshal marshals data to TOML format.
func (tm *TOMLMarshaller) Marshal(data map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Format returns the format name.
func (tm *TOMLMarshaller) Format() string {
	return "toml"
}
