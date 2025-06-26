package pipeline

import (
	"fmt"
	"konfigo/internal/errors"
	"konfigo/internal/features/variables"
	"konfigo/internal/logger"
	"konfigo/internal/marshaller"
	"konfigo/internal/parser"
	"konfigo/internal/reader"
	"konfigo/internal/schema"
	"konfigo/internal/util"
	"konfigo/internal/writer"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// processBatch handles batch processing mode with konfigo_forEach
func (p *Pipeline) processBatch(baseFinalConfig map[string]interface{}, loadedSchema *schema.Schema, forEachDirective *schema.KonfigoForEach, envVarsForSchema map[string]string) error {
	logger.Log("Starting batch processing with konfigo_forEach...")

	if len(forEachDirective.Items) > 0 && len(forEachDirective.ItemFiles) > 0 {
		return errors.NewError(errors.ErrorTypeConfigMerge, "konfigo_forEach cannot have both 'items' and 'itemFiles' defined simultaneously")
	}
	if len(forEachDirective.Items) == 0 && len(forEachDirective.ItemFiles) == 0 {
		return errors.NewError(errors.ErrorTypeConfigMerge, "konfigo_forEach must define either 'items' or 'itemFiles'")
	}
	if forEachDirective.Output.FilenamePattern == "" {
		return errors.NewError(errors.ErrorTypeCLIValidation, "konfigo_forEach.output.filenamePattern is required")
	}

	iterationSources := []map[string]interface{}{}
	itemFileBasenames := []string{} // For ${ITEM_FILE_BASENAME}

	if len(forEachDirective.Items) > 0 {
		logger.Debug("Iterating using 'items' from konfigo_forEach.")
		iterationSources = forEachDirective.Items
		for range forEachDirective.Items { // Populate basenames with empty strings for 'items'
			itemFileBasenames = append(itemFileBasenames, "")
		}
	} else { // len(forEachDirective.ItemFiles) > 0
		logger.Debug("Iterating using 'itemFiles' from konfigo_forEach.")
		for _, itemFilePath := range forEachDirective.ItemFiles {
			fullItemFilePath := itemFilePath
			if !filepath.IsAbs(itemFilePath) && p.Config.VarsFile != "" {
				fullItemFilePath = filepath.Join(filepath.Dir(p.Config.VarsFile), itemFilePath)
			}
			logger.Debug("Loading iteration variables from itemFile: %s", fullItemFilePath)
			content, err := reader.ReadFile(fullItemFilePath)
			if err != nil {
				return errors.WrapError(errors.ErrorTypeFileRead, "failed to read itemFile", err).WithContext("file", fullItemFilePath)
			}
			itemVars, err := parser.Parse(fullItemFilePath, content, "")
			if err != nil {
				return errors.WrapError(errors.ErrorTypeParsing, "failed to parse itemFile", err).WithContext("file", fullItemFilePath)
			}
			iterationSources = append(iterationSources, itemVars)
			itemFileBasenames = append(itemFileBasenames, strings.TrimSuffix(filepath.Base(fullItemFilePath), filepath.Ext(fullItemFilePath)))
		}
	}

	for i, iterVars := range iterationSources {
		logger.Log("Processing iteration %d...", i)

		currentConfig, err := util.DeepCopyMap(baseFinalConfig)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeDeepCopy, "failed to deep copy base config for iteration", err).WithContext("iteration", i)
		}

		varsForThisIteration := make(map[string]interface{})
		if forEachDirective.GlobalVars != nil {
			for k, v := range forEachDirective.GlobalVars {
				varsForThisIteration[k] = v
			}
		}
		for k, v := range iterVars {
			varsForThisIteration[k] = v
		}

		varsForThisIteration["ITEM_INDEX"] = strconv.Itoa(i)
		if i < len(itemFileBasenames) {
			varsForThisIteration["ITEM_FILE_BASENAME"] = itemFileBasenames[i]
		}

		processedConfig, err := schema.Process(currentConfig, loadedSchema, varsForThisIteration, envVarsForSchema)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeSchemaProcess, "schema processing failed for iteration", err).WithContext("iteration", i)
		}

		outputFilename, err := resolveFilenamePattern(forEachDirective.Output.FilenamePattern, varsForThisIteration, envVarsForSchema, loadedSchema.Vars, i, itemFileBasenames[i])
		if err != nil {
			return errors.WrapError(errors.ErrorTypeInternal, "failed to resolve output filename for iteration", err).WithContext("iteration", i)
		}

		outputFormat := strings.ToLower(strings.TrimPrefix(filepath.Ext(outputFilename), "."))
		if forEachDirective.Output.Format != "" {
			outputFormat = strings.ToLower(forEachDirective.Output.Format)
		}
		if outputFormat == "" {
			logger.Warn("Output format for iteration %d (%s) is ambiguous, defaulting to YAML. Specify format in konfigo_forEach.output.format or use a file extension.", i, outputFilename)
			outputFormat = "yaml"
		}

		outputBytes, err := marshaller.Marshal(processedConfig, outputFormat)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeInternal, "error marshalling", err).WithContext("format", outputFormat).WithContext("iteration", i).WithContext("file", outputFilename)
		}
		logger.Log("Writing output for iteration %d to %s (format: %s)", i, outputFilename, outputFormat)
		if err := writer.WriteFile(outputFilename, outputBytes); err != nil {
			return errors.WrapError(errors.ErrorTypeFileWrite, "error writing to file", err).WithContext("file", outputFilename).WithContext("iteration", i)
		}
	}
	logger.Log("Batch processing completed.")
	return nil
}

// resolveFilenamePattern substitutes placeholders in the filename pattern.
// Placeholders: ${VAR_NAME}, ${ITEM_INDEX}, ${ITEM_FILE_BASENAME}
func resolveFilenamePattern(pattern string, iterVars map[string]interface{}, envVarsForSchema map[string]string, schemaVars []variables.Definition, itemIndex int, itemFileBasename string) (string, error) {
	resolvedPattern := pattern

	// Substitute ${ITEM_INDEX}
	resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_INDEX}", strconv.Itoa(itemIndex))

	// Substitute ${ITEM_FILE_BASENAME}
	if itemFileBasename != "" {
		resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_FILE_BASENAME}", itemFileBasename)
	} else {
		// If itemFileBasename is empty (e.g. when using 'items'), remove the placeholder
		resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_FILE_BASENAME}", "")
	}

	// Regex to find ${VAR_NAME}
	varRegex := regexp.MustCompile(`\$\{[A-Z0-9_]+\}`)

	resolvedPattern = varRegex.ReplaceAllStringFunc(resolvedPattern, func(match string) string {
		varName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")

		// 1. Check iterVars (highest priority for filename context)
		if val, ok := iterVars[varName]; ok {
			return fmt.Sprintf("%v", val)
		}
		// 2. Check envVarsForSchema (KONFIGO_VAR_...)
		if val, ok := envVarsForSchema[varName]; ok {
			return val
		}
		// 3. Check schemaVars (default values or simple values if not from env/path)
		for _, sv := range schemaVars {
			if sv.Name == varName {
				if sv.Value != "" { // Only use direct value or default for simplicity in filename
					return sv.Value
				}
				if sv.DefaultValue != "" {
					return sv.DefaultValue
				}
			}
		}
		// Log error and replace with empty string if not found
		logger.Log("ERROR: Variable %s in filenamePattern not found, replacing with empty string.", match)
		return "" // Replace unresolved variable with an empty string
	})

	// Clean up path, e.g. remove double slashes if a variable was empty
	resolvedPattern = filepath.Clean(resolvedPattern)

	return resolvedPattern, nil
}
