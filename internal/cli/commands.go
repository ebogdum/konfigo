package cli

import (
	"konfigo/internal/errors"
	"konfigo/internal/logger"
)

// Command represents the main command execution logic
type Command struct {
	Config *Config
}

// NewCommand creates a new Command with the given configuration
func NewCommand(config *Config) *Command {
	return &Command{
		Config: config,
	}
}

// Execute runs the main command logic and returns the config for pipeline processing
func (cmd *Command) Execute() (*Config, error) {
	// Set custom usage for help
	SetCustomUsage()
	
	// Handle help or no arguments
	if cmd.Config.ShouldShowHelp() {
		PrintHelp()
		return nil, nil // nil error means help was shown successfully
	}
	
	// Validate configuration
	if err := cmd.Config.Validate(); err != nil {
		return nil, err
	}
	
	// Configure logger based on flags
	isDebug, isQuiet := cmd.Config.GetLoggerConfig()
	logger.Init(isDebug, isQuiet)
	
	// Validate that we have input sources
	sourcePaths := cmd.Config.GetSourcePaths()
	if sourcePaths == "" {
		return nil, errors.NewError(errors.ErrorTypeCLIValidation, "no input source specified. Use -s <paths> or pipe from stdin")
	}
	
	// Return the validated config for pipeline processing
	return cmd.Config, nil
}

// Run is the main entry point for the CLI application
func Run() (*Config, error) {
	config, err := ParseFlags()
	if err != nil {
		return nil, err
	}
	
	cmd := NewCommand(config)
	return cmd.Execute()
}
