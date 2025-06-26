package pipeline

import (
	"konfigo/internal/cli"
)

// Coordinator coordinates different processing modes and pipelines
type Coordinator struct {
	Config *cli.Config
}

// NewCoordinator creates a new coordinator with the given CLI configuration
func NewCoordinator(config *cli.Config) *Coordinator {
	return &Coordinator{
		Config: config,
	}
}

// Execute runs the appropriate processing pipeline based on the configuration
func (c *Coordinator) Execute() error {
	// For now, we only have one pipeline type, but this allows for
	// future extension to support different processing modes
	pipeline := NewPipeline(c.Config)
	return pipeline.Run()
}

// Run is a convenience function that creates a coordinator and executes it
func Run(config *cli.Config) error {
	coordinator := NewCoordinator(config)
	return coordinator.Execute()
}
