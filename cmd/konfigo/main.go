package main

import (
	"konfigo/internal/cli"
	"konfigo/internal/pipeline"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Parse CLI flags and get configuration
	config, err := cli.Run()
	if err != nil {
		return err
	}
	
	// If config is nil, help was shown and we should exit successfully
	if config == nil {
		return nil
	}
	
	// Run the processing pipeline
	return pipeline.Run(config)
}