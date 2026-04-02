package pipeline

import (
	"konfigo/internal/cli"
)

// Run creates and executes the processing pipeline with the given CLI configuration.
func Run(config *cli.Config) error {
	pipeline := NewPipeline(config)
	return pipeline.Run()
}
