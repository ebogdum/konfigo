package config

// Config represents the main configuration structure.
type Config struct {
	// Data holds the parsed configuration data
	Data map[string]interface{}
	// Sources tracks where each configuration value came from
	Sources map[string]string
}

// New creates a new Config instance.
func New() *Config {
	return &Config{
		Data:    make(map[string]interface{}),
		Sources: make(map[string]string),
	}
}

// Merge merges another configuration into this one.
func (c *Config) Merge(other *Config) {
	for key, value := range other.Data {
		c.Data[key] = value
	}
	for key, source := range other.Sources {
		c.Sources[key] = source
	}
}

// Set sets a configuration value with its source.
func (c *Config) Set(key string, value interface{}, source string) {
	c.Data[key] = value
	c.Sources[key] = source
}

// Get retrieves a configuration value.
func (c *Config) Get(key string) (interface{}, bool) {
	value, exists := c.Data[key]
	return value, exists
}

// GetSource retrieves the source of a configuration value.
func (c *Config) GetSource(key string) (string, bool) {
	source, exists := c.Sources[key]
	return source, exists
}
