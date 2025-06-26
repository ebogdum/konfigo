package marshaller

import (
	"strings"
)

// Marshaller interface defines the contract for format marshallers.
type Marshaller interface {
	// Marshal marshals data to the format-specific byte representation.
	Marshal(data map[string]interface{}) ([]byte, error)
	
	// Format returns the format name this marshaller handles.
	Format() string
}

// Registry holds all available marshallers.
type Registry struct {
	marshallers map[string]Marshaller
}

// NewRegistry creates a new marshaller registry with all built-in marshallers.
func NewRegistry() *Registry {
	registry := &Registry{
		marshallers: make(map[string]Marshaller),
	}
	
	// Register all built-in marshallers
	registry.Register(&JSONMarshaller{})
	registry.Register(&YAMLMarshaller{})
	registry.Register(&TOMLMarshaller{})
	registry.Register(&ENVMarshaller{})
	
	return registry
}

// Register adds a marshaller to the registry.
func (r *Registry) Register(marshaller Marshaller) {
	r.marshallers[marshaller.Format()] = marshaller
}

// Get retrieves a marshaller by format name.
func (r *Registry) Get(format string) (Marshaller, bool) {
	normalizedFormat := strings.ToLower(format)
	if normalizedFormat == "yml" {
		normalizedFormat = "yaml"
	}
	marshaller, exists := r.marshallers[normalizedFormat]
	return marshaller, exists
}

// GetFormats returns all supported format names.
func (r *Registry) GetFormats() []string {
	formats := make([]string, 0, len(r.marshallers))
	for format := range r.marshallers {
		formats = append(formats, format)
	}
	return formats
}

// IsFormatSupported checks if the given format is supported.
func IsFormatSupported(format string) bool {
	defaultRegistry := NewRegistry()
	_, exists := defaultRegistry.Get(format)
	return exists
}
