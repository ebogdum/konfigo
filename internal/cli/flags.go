// Package cli provides command-line interface functionality for Konfigo.
//
// This package handles all aspects of the command-line interface including:
// - Flag parsing and validation
// - Help text generation and display
// - Command execution coordination
// - Configuration object creation
//
// The CLI package follows a clean separation of concerns:
// - flags.go: Flag definitions, parsing, and validation
// - help.go: Help text generation and display
// - commands.go: Command execution logic and coordination
//
// Usage:
//
//	config, err := cli.Run()
//	if err != nil {
//	    return err
//	}
//	if config == nil {
//	    // Help was shown, exit successfully
//	    return nil
//	}
//	// Use config for processing...
package cli

import (
	"flag"
	"strings"

	"konfigo/internal/errors"
)

// Config holds all CLI flag values
type Config struct {
	// Schema and Variables
	SchemaFile   string
	VarsFile     string
	
	// Sources and Input
	SourcePaths       string
	Recursive         bool
	CaseSensitive     bool
	InputJSON         bool
	InputYAML         bool
	InputTOML         bool
	InputENV          bool
	
	// Output
	OutputFile        string
	OutputJSON        bool
	OutputYAML        bool
	OutputTOML        bool
	OutputENV         bool
	
	// Behavior and Logging
	Verbose           bool
	Debug             bool
	Help              bool
}

// ParseFlags parses command line flags and returns a Config struct
func ParseFlags() (*Config, error) {
	config := &Config{}
	
	// Schema & Variables
	flag.StringVar(&config.SchemaFile, "schema", "", "Path to a schema file for processing the config.")
	flag.StringVar(&config.SchemaFile, "S", "", "Path to a schema file (shorthand for --schema).")
	flag.StringVar(&config.VarsFile, "vars-file", "", "Path to a file providing high-priority variables.")
	flag.StringVar(&config.VarsFile, "V", "", "Path to a variables file (shorthand for --vars-file).")
	
	// Sources and Input
	flag.StringVar(&config.SourcePaths, "s", "", "Comma-separated list of source files/directories. Use '-' for stdin.")
	flag.BoolVar(&config.Recursive, "r", false, "Recursively search for configuration files in subdirectories")
	flag.BoolVar(&config.CaseSensitive, "c", false, "Use case-sensitive key matching (default is case-insensitive)")
	flag.BoolVar(&config.InputJSON, "sj", false, "Force input to be parsed as JSON (required for stdin)")
	flag.BoolVar(&config.InputYAML, "sy", false, "Force input to be parsed as YAML (required for stdin)")
	flag.BoolVar(&config.InputTOML, "st", false, "Force input to be parsed as TOML (required for stdin)")
	flag.BoolVar(&config.InputENV, "se", false, "Force input to be parsed as ENV (required for stdin)")
	
	// Output
	flag.StringVar(&config.OutputFile, "of", "", "Write output to file. Extension determines format, or use with -oX flags.")
	flag.BoolVar(&config.OutputJSON, "oj", false, "Output in JSON format")
	flag.BoolVar(&config.OutputYAML, "oy", false, "Output in YAML format")
	flag.BoolVar(&config.OutputTOML, "ot", false, "Output in TOML format")
	flag.BoolVar(&config.OutputENV, "oe", false, "Output in ENV format")
	
	// Behavior and Logging
	flag.BoolVar(&config.Verbose, "v", false, "Enable informational (INFO) logging. Overrides default quiet behavior.")
	flag.BoolVar(&config.Debug, "d", false, "Enable debug (DEBUG and INFO) logging. Overrides -v and default quiet behavior.")
	flag.BoolVar(&config.Help, "h", false, "Show this help message.")
	
	flag.Parse()
	
	return config, nil
}

// GetInputFormat returns the specified input format override, or empty string if none
func (c *Config) GetInputFormat() string {
	if c.InputJSON {
		return "json"
	} else if c.InputYAML {
		return "yaml"
	} else if c.InputTOML {
		return "toml"
	} else if c.InputENV {
		return "env"
	}
	return ""
}

// GetSourcePaths returns the list of source paths, handling both -s flag and positional args
func (c *Config) GetSourcePaths() string {
	if c.SourcePaths == "" && flag.NArg() > 0 {
		return strings.Join(flag.Args(), ",")
	}
	return c.SourcePaths
}

// ShouldShowHelp returns true if help should be displayed
func (c *Config) ShouldShowHelp() bool {
	return c.Help || flag.NFlag() == 0
}

// GetLoggerConfig returns the logger configuration based on debug/verbose flags
func (c *Config) GetLoggerConfig() (isDebug bool, isQuiet bool) {
	if c.Debug {
		return true, false
	} else if c.Verbose {
		return false, false
	}
	return false, true // Default: no debug, quiet
}

// Validate performs basic validation on the flag configuration
func (c *Config) Validate() error {
	// Check for conflicting input format flags
	inputFormats := []bool{c.InputJSON, c.InputYAML, c.InputTOML, c.InputENV}
	inputCount := 0
	for _, set := range inputFormats {
		if set {
			inputCount++
		}
	}
	if inputCount > 1 {
		return errors.NewError(errors.ErrorTypeCLIFlag, "only one input format flag (-sj, -sy, -st, -se) can be specified")
	}
	
	return nil
}

// ValidateSourcePaths validates the source paths for common issues
func (c *Config) ValidateSourcePaths() error {
	sourcePaths := c.GetSourcePaths()
	if sourcePaths == "" {
		return errors.NewError(errors.ErrorTypeCLIValidation, "no input source specified. Use -s <paths> or pipe from stdin")
	}
	
	sources := strings.Split(sourcePaths, ",")
	for _, source := range sources {
		source = strings.TrimSpace(source)
		if source == "" {
			continue
		}
		
		// If it's stdin, validate that a format is specified
		if source == "-" {
			inputFormat := c.GetInputFormat()
			if inputFormat == "" {
				return errors.NewError(errors.ErrorTypeCLIValidation, "when using stdin (-), an input format must be specified (-sj, -sy, -st, or -se)")
			}
		}
	}
	
	return nil
}

// ValidateOutputConfiguration validates output configuration for consistency
func (c *Config) ValidateOutputConfiguration() error {
	// Check for conflicting output format flags
	outputFormats := []bool{c.OutputJSON, c.OutputYAML, c.OutputTOML, c.OutputENV}
	outputCount := 0
	for _, set := range outputFormats {
		if set {
			outputCount++
		}
	}
	
	// Multiple output formats are allowed, but warn if both file and stdout outputs are mixed
	if c.OutputFile != "" && outputCount > 0 {
		// This is actually allowed - file output and stdout output can coexist
		// Just ensure it's not confusing
	}
	
	return nil
}
