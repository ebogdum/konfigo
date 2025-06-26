package parser

// Parser interface defines the contract for format parsers.
type Parser interface {
	// Parse parses content and returns the resulting map.
	Parse(content []byte) (map[string]interface{}, error)
	
	// Format returns the format name this parser handles.
	Format() string
}

// Registry holds all available parsers.
type Registry struct {
	parsers map[string]Parser
}

// NewRegistry creates a new parser registry with all built-in parsers.
func NewRegistry() *Registry {
	registry := &Registry{
		parsers: make(map[string]Parser),
	}
	
	// Register all built-in parsers
	registry.Register(&JSONParser{})
	registry.Register(&YAMLParser{})
	registry.Register(&TOMLParser{})
	registry.Register(&INIParser{})
	registry.Register(&ENVParser{})
	
	return registry
}

// Register adds a parser to the registry.
func (r *Registry) Register(parser Parser) {
	r.parsers[parser.Format()] = parser
}

// Get retrieves a parser by format name.
func (r *Registry) Get(format string) (Parser, bool) {
	parser, exists := r.parsers[NormalizeFormat(format)]
	return parser, exists
}

// GetFormats returns all supported format names.
func (r *Registry) GetFormats() []string {
	formats := make([]string, 0, len(r.parsers))
	for format := range r.parsers {
		formats = append(formats, format)
	}
	return formats
}
