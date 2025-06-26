package validator

// Registry holds all available validators.
type Registry struct {
	validators []Validator
}

// NewRegistry creates a new validator registry with all built-in validators.
func NewRegistry() *Registry {
	registry := &Registry{}
	
	// Register all built-in validators
	registry.Register(&TypeValidator{})
	registry.Register(&NumericValidator{})
	registry.Register(&StringValidator{})
	
	return registry
}

// Register adds a validator to the registry.
func (r *Registry) Register(validator Validator) {
	r.validators = append(r.validators, validator)
}

// GetValidators returns all registered validators.
func (r *Registry) GetValidators() []Validator {
	return r.validators
}
