// Package pipeline provides processing orchestration for Konfigo configuration management.
//
// This package coordinates the complete configuration processing workflow including:
// - Schema loading and validation
// - Source file discovery and parsing
// - Configuration merging with precedence rules
// - Feature processing (generators, transformers, validators, variables)
// - Output generation in multiple formats
// - Batch processing with konfigo_forEach
//
// The pipeline package follows a modular architecture:
// - coordinator.go: Processing mode coordination
// - pipeline.go: Main processing pipeline and workflow
// - single.go: Single configuration processing mode
// - batch.go: Batch processing with iteration support
// - optimized.go: Performance optimizations for large configurations
//
// Processing Flow:
//  1. Load schema and environment variables
//  2. Discover and parse source files
//  3. Merge configurations with precedence rules
//  4. Process schema features (variables, generators, transformers, validators)
//  5. Generate outputs in requested formats
//
// Usage:
//
//	pipeline := NewPipeline(cliConfig)
//	err := pipeline.Run()
package pipeline

import (
	"fmt"
	"konfigo/internal/cli"
	"konfigo/internal/config"
	"konfigo/internal/errors"
	"konfigo/internal/features/variables"
	"konfigo/internal/logger"
	"konfigo/internal/marshaller"
	"konfigo/internal/merger"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
	"konfigo/internal/schema"
	"konfigo/internal/writer"
	"sort"
	"strings"
)

// Pipeline represents the main processing pipeline
type Pipeline struct {
	Config *cli.Config
}

// NewPipeline creates a new pipeline with the given CLI configuration
func NewPipeline(config *cli.Config) *Pipeline {
	return &Pipeline{
		Config: config,
	}
}

// Run executes the complete processing pipeline
func (p *Pipeline) Run() error {
	// Load schema and environment variables early
	loadedSchema, immutablePaths, err := p.loadSchemaAndImmutablePaths()
	if err != nil {
		return err
	}

	env := config.NewEnvironment()
	envResult := env.Load()
	envConfig := envResult.Config.Data
	envVarsForSchema := envResult.Vars

	// Process sources and merge configurations
	baseFinalConfig, err := p.processSources(immutablePaths, envConfig)
	if err != nil {
		return err
	}

	// Load variables file and check for konfigo_forEach
	varsFromFileGlobal, forEachDirective, err := p.loadVariablesFile()
	if err != nil {
		return err
	}

	// Process schema (single or batch) or perform basic variable substitution
	if loadedSchema != nil {
		if forEachDirective != nil {
			// Batch processing mode
			return p.processBatch(baseFinalConfig, loadedSchema, forEachDirective, envVarsForSchema)
		} else {
			// Single processing mode with schema
			baseFinalConfig, err = p.processSingle(baseFinalConfig, loadedSchema, varsFromFileGlobal, envVarsForSchema)
			if err != nil {
				return err
			}
		}
	} else if forEachDirective != nil {
		return errors.NewError(errors.ErrorTypeSchemaLoad, "konfigo_forEach directive found, but no schema file (-S) was provided for processing")
	} else {
		// No schema provided - perform basic variable substitution
		baseFinalConfig, err = p.processBasicVariableSubstitution(baseFinalConfig, envVarsForSchema, varsFromFileGlobal)
		if err != nil {
			return err
		}
	}

	// Generate outputs (for single mode or if no schema)
	if forEachDirective == nil {
		return p.generateOutputs(baseFinalConfig)
	}

	return nil
}

// loadSchemaAndImmutablePaths loads the schema file and extracts immutable paths
func (p *Pipeline) loadSchemaAndImmutablePaths() (*schema.Schema, map[string]struct{}, error) {
	var loadedSchema *schema.Schema
	immutablePaths := make(map[string]struct{})

	if p.Config.SchemaFile != "" {
		logger.Log("Loading schema from %s", p.Config.SchemaFile)
		var err error
		loadedSchema, err = schema.Load(p.Config.SchemaFile)
		if err != nil {
			return nil, nil, err
		}

		for _, path := range loadedSchema.Immutable {
			immutablePaths[path] = struct{}{}
		}
	}

	return loadedSchema, immutablePaths, nil
}

// loadVariablesFile loads the variables file and extracts forEach directive
func (p *Pipeline) loadVariablesFile() (map[string]interface{}, *schema.KonfigoForEach, error) {
	var varsFromFileGlobal map[string]interface{}
	var forEachDirective *schema.KonfigoForEach

	if p.Config.VarsFile != "" {
		logger.Log("Loading variables from %s", p.Config.VarsFile)
		content, err := reader.ReadFile(p.Config.VarsFile)
		if err != nil {
			return nil, nil, errors.WrapError(errors.ErrorTypeFileRead, "failed to read vars file", err).WithContext("file", p.Config.VarsFile)
		}
		rawVarsFromFile, err := parser.Parse(p.Config.VarsFile, content, "")
		if err != nil {
			return nil, nil, errors.WrapError(errors.ErrorTypeParsing, "failed to parse vars file", err).WithContext("file", p.Config.VarsFile)
		}

		// Separate konfigo_forEach from other global vars using config package
		forEachDirective, varsFromFileGlobal, err = config.ExtractForEachFromVars(rawVarsFromFile)
		if err != nil {
			return nil, nil, errors.WrapError(errors.ErrorTypeSchemaProcess, "failed to extract konfigo_forEach directive", err)
		}
		if forEachDirective != nil {
			forEachDirective.GlobalVars = varsFromFileGlobal // Store global vars for resolver
		}
	}

	return varsFromFileGlobal, forEachDirective, nil
}

// processBasicVariableSubstitution performs variable substitution without schema
func (p *Pipeline) processBasicVariableSubstitution(baseConfig map[string]interface{}, envVarsForSchema map[string]string, varsFromFileGlobal map[string]interface{}) (map[string]interface{}, error) {
	logger.Debug("No schema provided. Performing basic variable substitution from environment and -V file if present.")

	// Create a resolver with available variable sources
	resolver, err := variables.NewResolver(envVarsForSchema, varsFromFileGlobal, []variables.Definition{}, baseConfig)
	if err != nil {
		return nil, errors.WrapError(errors.ErrorTypeVarResolution, "failed to create variable resolver without schema", err)
	}
	
	return variables.Substitute(baseConfig, resolver), nil
}

// generateOutputs handles output generation for single processing mode
func (p *Pipeline) generateOutputs(finalConfig map[string]interface{}) error {
	targets := writer.DetermineOutputTargets(p.Config.OutputFile, p.Config.OutputJSON, p.Config.OutputYAML, p.Config.OutputTOML, p.Config.OutputENV)
	if len(targets) == 0 && p.Config.OutputFile == "" {
		// Default to YAML stdout if no other output specified
		targets = append(targets, writer.OutputTarget{Format: "yaml", Filename: ""})
	}

	for i, target := range targets {
		outputBytes, err := marshaller.Marshal(finalConfig, target.Format)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeInternal, "error marshalling", err).WithContext("format", target.Format)
		}
		if target.Filename == "" { // Output to stdout
			if i > 0 && len(targets) > 1 { // Add separator for multiple stdout formats
				fmt.Println("---")
			}
			fmt.Println(string(outputBytes))
		} else {
			logger.Log("Writing output to %s (format: %s)", target.Filename, target.Format)
			if err := writer.WriteFile(target.Filename, outputBytes); err != nil {
				return errors.WrapError(errors.ErrorTypeFileWrite, "error writing to file", err).WithContext("file", target.Filename)
			}
		}
	}
	
	return nil
}

// parseResult holds the result of parsing a single file
type parseResult struct {
	FilePath string
	Data     map[string]interface{}
	Err      error
}

// processSources handles discovery, parsing, and merging of source files
func (p *Pipeline) processSources(immutablePaths map[string]struct{}, envConfig map[string]interface{}) (map[string]interface{}, error) {
	var allFiles []string
	var stdinData []byte
	sourcePaths := p.Config.GetSourcePaths()
	inputFormatOverride := p.Config.GetInputFormat()
	
	sources := strings.Split(sourcePaths, ",")
	logger.Log("Discovering configuration files...")
	
	for _, source := range sources {
		source = strings.TrimSpace(source)
		if source == "" {
			continue
		}
		if source == "-" {
			logger.Debug("Reading from standard input (stdin)")
			if err := reader.ValidateStdinFormat(inputFormatOverride); err != nil {
				return nil, err
			}
			var err error
			stdinData, err = reader.ReadStdin()
			if err != nil {
				return nil, err
			}
			continue
		}
		files, err := reader.DiscoverFiles(source, p.Config.Recursive)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeFileRead, "error loading from source", err).WithContext("source", source)
		}
		logger.Debug("Found %d file(s) in source: %s", len(files), source)
		allFiles = append(allFiles, files...)
	}

	results := p.parseFilesParallel(allFiles, inputFormatOverride)
	sort.Slice(results, func(i, j int) bool { return results[i].FilePath < results[j].FilePath })

	finalConfig := make(map[string]interface{})
	logger.Log("Merging %d configuration file(s)...", len(results))
	for _, res := range results {
		if res.Err != nil {
			logger.Log("  - Warning: Skipping file %s due to parse error: %v", res.FilePath, res.Err)
			continue
		}
		merger.Merge(finalConfig, res.Data, p.Config.CaseSensitive, immutablePaths)
	}

	if len(stdinData) > 0 {
		logger.Log("Merging configuration from stdin...")
		data, err := parser.Parse("stdin", stdinData, inputFormatOverride)
		if err != nil {
			return nil, errors.WrapError(errors.ErrorTypeStdinRead, "failed to parse stdin", err)
		}
		merger.Merge(finalConfig, data, p.Config.CaseSensitive, immutablePaths)
	}

	if len(envConfig) > 0 {
		logger.Log("Merging %d configuration key(s) from environment variables...", len(envConfig))
		merger.Merge(finalConfig, envConfig, p.Config.CaseSensitive, immutablePaths)
	}

	return finalConfig, nil
}

// parseFilesParallel parses multiple files in parallel using optimized processing
func (p *Pipeline) parseFilesParallel(files []string, formatOverride string) []parseResult {
	if len(files) == 0 {
		return nil
	}
	
	// Use optimized file processor for better performance
	processor := NewOptimizedFileProcessor()
	return processor.ProcessFiles(files, formatOverride)
}
