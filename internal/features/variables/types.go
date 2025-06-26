package variables

import (
	"regexp"
)

// VarRegex is the regular expression for matching variable placeholders like ${VAR_NAME}.
var VarRegex = regexp.MustCompile(`\$\{[A-Z0-9_]+\}`)

// Definition defines a variable that can be used for substitution.
type Definition struct {
	Name         string `yaml:"name" json:"name"`
	Value        string `yaml:"value,omitempty" json:"value,omitempty"`
	FromEnv      string `yaml:"fromEnv,omitempty" json:"fromEnv,omitempty"`
	FromPath     string `yaml:"fromPath,omitempty" json:"fromPath,omitempty"`
	DefaultValue string `yaml:"defaultValue,omitempty" json:"defaultValue,omitempty"`
}

// Resolver interface defines the contract for variable resolution.
type Resolver interface {
	// SubstituteString performs variable substitution on a single string.
	SubstituteString(input string) string
}
