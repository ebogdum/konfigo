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
	"os"
	"strings"

	"konfigo/internal/errors"
)

// flagSet is the package-level flag set used for parsing.
// Using a dedicated FlagSet instead of the global flag.CommandLine
// allows tests to call ParseFlags multiple times without panicking.
var flagSet *flag.FlagSet

// Config holds all CLI flag values
type Config struct {
	// Schema and Variables
	SchemaFile string
	VarsFile   string

	// Sources and Input
	SourcePaths   string
	Recursive     bool
	CaseSensitive bool
	InputJSON     bool
	InputYAML     bool
	InputTOML     bool
	InputENV      bool

	// Output
	OutputFile string
	OutputJSON bool
	OutputYAML bool
	OutputTOML bool
	OutputENV  bool

	// Behavior and Logging
	MergeArrays bool
	Verbose     bool
	Debug       bool
	Help        bool
}

// ParseFlags parses command line flags and returns a Config struct.
// It creates a fresh FlagSet each time so the function is safe to call repeatedly in tests.
func ParseFlags() (*Config, error) {
	config := &Config{}

	flagSet = flag.NewFlagSet("konfigo", flag.ContinueOnError)

	// Schema & Variables
	flagSet.StringVar(&config.SchemaFile, "schema", "", "Path to a schema file for processing the config.")
	flagSet.StringVar(&config.SchemaFile, "S", "", "Path to a schema file (shorthand for --schema).")
	flagSet.StringVar(&config.VarsFile, "vars-file", "", "Path to a file providing high-priority variables.")
	flagSet.StringVar(&config.VarsFile, "V", "", "Path to a variables file (shorthand for --vars-file).")

	// Sources and Input
	flagSet.StringVar(&config.SourcePaths, "s", "", "Comma-separated list of source files/directories. Use '-' for stdin.")
	flagSet.BoolVar(&config.Recursive, "r", false, "Recursively search for configuration files in subdirectories")
	flagSet.BoolVar(&config.CaseSensitive, "c", false, "Use case-sensitive key matching (default is case-insensitive)")
	flagSet.BoolVar(&config.InputJSON, "sj", false, "Force input to be parsed as JSON (required for stdin)")
	flagSet.BoolVar(&config.InputYAML, "sy", false, "Force input to be parsed as YAML (required for stdin)")
	flagSet.BoolVar(&config.InputTOML, "st", false, "Force input to be parsed as TOML (required for stdin)")
	flagSet.BoolVar(&config.InputENV, "se", false, "Force input to be parsed as ENV (required for stdin)")

	// Output
	flagSet.StringVar(&config.OutputFile, "of", "", "Write output to file. Extension determines format, or use with -oX flags.")
	flagSet.BoolVar(&config.OutputJSON, "oj", false, "Output in JSON format")
	flagSet.BoolVar(&config.OutputYAML, "oy", false, "Output in YAML format")
	flagSet.BoolVar(&config.OutputTOML, "ot", false, "Output in TOML format")
	flagSet.BoolVar(&config.OutputENV, "oe", false, "Output in ENV format")

	// Behavior and Logging
	flagSet.BoolVar(&config.MergeArrays, "m", false, "Merge arrays by union with deduplication instead of replacing.")
	flagSet.BoolVar(&config.Verbose, "v", false, "Enable informational (INFO) logging. Overrides default quiet behavior.")
	flagSet.BoolVar(&config.Debug, "d", false, "Enable debug (DEBUG and INFO) logging. Overrides -v and default quiet behavior.")
	flagSet.BoolVar(&config.Help, "h", false, "Show this help message.")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, errors.WrapError(errors.ErrorTypeCLIFlag, "failed to parse flags", err)
	}

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
	if c.SourcePaths == "" && flagSet != nil && flagSet.NArg() > 0 {
		return strings.Join(flagSet.Args(), ",")
	}
	return c.SourcePaths
}

// ShouldShowHelp returns true if help should be displayed
func (c *Config) ShouldShowHelp() bool {
	if flagSet == nil {
		return c.Help
	}
	return c.Help || (flagSet.NFlag() == 0 && flagSet.NArg() == 0)
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

