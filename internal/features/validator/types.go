package validator

// Group represents a validation group with a path and rules.
type Group struct {
	Path  string `yaml:"path" json:"path"`
	Rules Rule   `yaml:"rules" json:"rules"`
}

// Rule represents validation rules that can be applied to a value.
type Rule struct {
	Required  bool     `yaml:"required" json:"required"`
	Type      string   `yaml:"type" json:"type"`
	Min       *float64 `yaml:"min" json:"min"`
	Max       *float64 `yaml:"max" json:"max"`
	MinLength *int     `yaml:"minLength" json:"minLength"`
	Enum      []string `yaml:"enum" json:"enum"`
	Regex     string   `yaml:"regex" json:"regex"`
}

// Validator interface defines a validation function.
type Validator interface {
	// Validate performs validation on the given value at the specified path.
	Validate(value interface{}, path string, rule Rule) error
}
