package config

import (
	"errors"
	"fmt"
	"konfigo/internal/logger"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
	"konfigo/internal/schema"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// BatchProcessor handles konfigo_forEach batch processing.
type BatchProcessor struct{}

// NewBatchProcessor creates a new batch processor.
func NewBatchProcessor() *BatchProcessor {
	return &BatchProcessor{}
}

// ProcessResult contains the result of batch processing.
type ProcessResult struct {
	OutputFiles []string
	Errors      []error
}

// Process executes batch processing using the konfigo_forEach directive.
func (bp *BatchProcessor) Process(
	forEachConfig *schema.KonfigoForEach,
	baseConfig map[string]interface{},
	schemaFile *schema.Schema,
	varsFromFile map[string]interface{},
	envVarsForSchema map[string]string,
) (*ProcessResult, error) {
	
	if forEachConfig == nil {
		return nil, errors.New("konfigo_forEach configuration is nil")
	}

	// Validate the forEach configuration
	if err := bp.validateForEachConfig(forEachConfig); err != nil {
		return nil, err
	}

	logger.Log("Starting batch processing with konfigo_forEach...")

	result := &ProcessResult{
		OutputFiles: []string{},
		Errors:      []error{},
	}

	if len(forEachConfig.Items) > 0 {
		logger.Debug("Iterating using 'items' from konfigo_forEach.")
		bp.processItems(forEachConfig, baseConfig, schemaFile, varsFromFile, envVarsForSchema, result)
	} else if len(forEachConfig.ItemFiles) > 0 {
		logger.Debug("Iterating using 'itemFiles' from konfigo_forEach.")
		bp.processItemFiles(forEachConfig, baseConfig, schemaFile, varsFromFile, envVarsForSchema, result)
	}

	return result, nil
}

// validateForEachConfig validates the konfigo_forEach configuration.
func (bp *BatchProcessor) validateForEachConfig(forEachConfig *schema.KonfigoForEach) error {
	hasItems := len(forEachConfig.Items) > 0
	hasItemFiles := len(forEachConfig.ItemFiles) > 0

	if hasItems && hasItemFiles {
		return errors.New("konfigo_forEach cannot have both 'items' and 'itemFiles' defined simultaneously")
	}
	if !hasItems && !hasItemFiles {
		return errors.New("konfigo_forEach must define either 'items' or 'itemFiles'")
	}
	if forEachConfig.Output.FilenamePattern == "" {
		return errors.New("konfigo_forEach.output.filenamePattern is required")
	}
	
	return nil
}

// processItems processes the items array from konfigo_forEach.
func (bp *BatchProcessor) processItems(
	forEachConfig *schema.KonfigoForEach,
	baseConfig map[string]interface{},
	schemaFile *schema.Schema,
	varsFromFile map[string]interface{},
	envVarsForSchema map[string]string,
	result *ProcessResult,
) {
	for i, item := range forEachConfig.Items {
		if err := bp.processIteration(i, item, "", forEachConfig, baseConfig, schemaFile, varsFromFile, envVarsForSchema, result); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("iteration %d failed: %w", i, err))
		}
	}
}

// processItemFiles processes the itemFiles array from konfigo_forEach.
func (bp *BatchProcessor) processItemFiles(
	forEachConfig *schema.KonfigoForEach,
	baseConfig map[string]interface{},
	schemaFile *schema.Schema,
	varsFromFile map[string]interface{},
	envVarsForSchema map[string]string,
	result *ProcessResult,
) {
	for i, itemFile := range forEachConfig.ItemFiles {
		// Load item data from file
		itemContent, err := reader.ReadFile(itemFile)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to read item file %s: %w", itemFile, err))
			continue
		}

		itemData, err := parser.Parse(itemFile, itemContent, "")
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to parse item file %s: %w", itemFile, err))
			continue
		}

		// Get base filename without extension for pattern resolution
		itemFileBasename := strings.TrimSuffix(filepath.Base(itemFile), filepath.Ext(itemFile))

		if err := bp.processIteration(i, itemData, itemFileBasename, forEachConfig, baseConfig, schemaFile, varsFromFile, envVarsForSchema, result); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("iteration %d (file %s) failed: %w", i, itemFile, err))
		}
	}
}

// processIteration processes a single iteration of the batch.
func (bp *BatchProcessor) processIteration(
	index int,
	iterVars map[string]interface{},
	itemFileBasename string,
	forEachConfig *schema.KonfigoForEach,
	baseConfig map[string]interface{},
	schemaFile *schema.Schema,
	varsFromFile map[string]interface{},
	envVarsForSchema map[string]string,
	result *ProcessResult,
) error {
	// This is a placeholder for the actual iteration processing logic
	// In the real implementation, this would:
	// 1. Resolve the filename pattern with variables
	// 2. Create a copy of the base config
	// 3. Merge iteration variables
	// 4. Apply schema processing
	// 5. Write the output file
	
	logger.Debug("Processing iteration %d", index)
	
	// For now, just track that we processed this iteration
	// The actual implementation would need to be integrated with
	// the existing processing logic from main.go
	
	return nil
}

// ExtractForEachFromVars extracts konfigo_forEach directive from variables file.
func ExtractForEachFromVars(varsFromFile map[string]interface{}) (*schema.KonfigoForEach, map[string]interface{}, error) {
	if varsFromFile == nil {
		return nil, nil, nil
	}

	var forEachConfig *schema.KonfigoForEach
	globalVars := make(map[string]interface{})

	// Separate konfigo_forEach from other global vars
	for k, v := range varsFromFile {
		if k == "konfigo_forEach" {
			// Convert to KonfigoForEach struct
			yamlBytes, err := yaml.Marshal(v)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to marshal konfigo_forEach directive: %w", err)
			}
			
			forEachConfig = &schema.KonfigoForEach{}
			if err := yaml.Unmarshal(yamlBytes, forEachConfig); err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal konfigo_forEach directive: %w", err)
			}
			
			logger.Debug("Found konfigo_forEach directive.")
		} else {
			globalVars[k] = v
		}
	}

	if forEachConfig != nil {
		forEachConfig.GlobalVars = globalVars
	}

	return forEachConfig, globalVars, nil
}
