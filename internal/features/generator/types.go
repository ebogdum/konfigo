// Package generator provides functionality for generating new configuration data.
// It supports various types of generators that can create values based on existing configuration.
package generator

// Definition represents a generator configuration.
type Definition struct {
	Type       string            `yaml:"type" json:"type"`
	TargetPath string            `yaml:"targetPath" json:"targetPath"`
	Format     string            `yaml:"format" json:"format"`
	Sources    map[string]string `yaml:"sources" json:"sources"`
}

// Generator represents a function that can generate configuration values.
type Generator interface {
	// Generate applies the generator logic to the configuration.
	Generate(config map[string]interface{}, def Definition, resolver VariableResolver) error

	// Type returns the generator type name.
	Type() string
}

// VariableResolver provides an interface for variable substitution.
// This allows generators to substitute variables in their output.
type VariableResolver interface {
	// SubstituteString performs variable substitution on a string.
	SubstituteString(input string) string
}

// Registry holds registered generators by type.
type Registry map[string]Generator

// NewRegistry creates a new generator registry with default generators.
func NewRegistry() Registry {
	registry := make(Registry)

	// Register built-in generators
	registry[ConcatGeneratorType] = &ConcatGenerator{}
	registry[TimestampGeneratorType] = &TimestampGenerator{}
	registry[RandomGeneratorType] = &RandomGenerator{}
	registry[IdGeneratorType] = &IdGenerator{}

	return registry
}

// Register adds a generator to the registry.
func (r Registry) Register(generator Generator) {
	r[generator.Type()] = generator
}

// Get retrieves a generator by type.
func (r Registry) Get(generatorType string) (Generator, bool) {
	gen, exists := r[generatorType]
	return gen, exists
}

// GetTypes returns all registered generator types.
func (r Registry) GetTypes() []string {
	types := make([]string, 0, len(r))
	for t := range r {
		types = append(types, t)
	}
	return types
}
