// Package transformer provides functionality for transforming configuration data.
// It supports various types of transformations including key renaming, case changes, and value setting.
package transformer

// Definition represents a transformation configuration.
type Definition struct {
	Type   string      `yaml:"type" json:"type"`
	Path   string      `yaml:"path" json:"path"`
	From   string      `yaml:"from" json:"from"`
	To     string      `yaml:"to" json:"to"`
	Case   string      `yaml:"case" json:"case"`
	Prefix string      `yaml:"prefix" json:"prefix"`
	Value  interface{} `yaml:"value" json:"value"`
}

// Transformer represents a function that can transform configuration values.
type Transformer interface {
	// Transform applies the transformation logic to the configuration.
	Transform(config map[string]interface{}, def Definition) error
	
	// Type returns the transformer type name.
	Type() string
}

// VariableResolver provides an interface for variable substitution.
type VariableResolver interface {
	// SubstituteString performs variable substitution on a string.
	SubstituteString(input string) string
}

// Registry holds registered transformers by type.
type Registry map[string]Transformer

// NewRegistry creates a new transformer registry with default transformers.
func NewRegistry() Registry {
	registry := make(Registry)
	
	// Register built-in transformers
	registry[RenameKeyType] = &RenameKeyTransformer{}
	registry[ChangeCaseType] = &ChangeCaseTransformer{}
	registry[AddKeyPrefixType] = &AddKeyPrefixTransformer{}
	registry[SetValueType] = &SetValueTransformer{}
	
	return registry
}

// Register adds a transformer to the registry.
func (r Registry) Register(transformer Transformer) {
	r[transformer.Type()] = transformer
}

// Get retrieves a transformer by type.
func (r Registry) Get(transformerType string) (Transformer, bool) {
	transformer, exists := r[transformerType]
	return transformer, exists
}

// GetTypes returns all registered transformer types.
func (r Registry) GetTypes() []string {
	types := make([]string, 0, len(r))
	for t := range r {
		types = append(types, t)
	}
	return types
}
