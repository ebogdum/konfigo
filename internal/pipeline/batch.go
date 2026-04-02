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

// filenameVarRegex matches ${VAR_NAME} placeholders in filename patterns.
// Compiled once at package level to avoid per-iteration regex compilation.
var filenameVarRegex = regexp.MustCompile(`\$\{[A-Za-z0-9_]+\}`)

// processBatch handles batch processing mode with forEach
func (p *Pipeline) processBatch(baseFinalConfig map[string]interface{}, loadedSchema *schema.Schema, forEachDirective *schema.KonfigoForEach, envVarsForSchema map[string]string) error {
	logger.Log("Starting batch processing with forEach...")

	if len(forEachDirective.Items) > 0 && len(forEachDirective.ItemFiles) > 0 {
		return errors.NewError(errors.ErrorTypeConfigMerge, "forEach cannot have both 'items' and 'itemFiles' defined simultaneously")
	}
	if len(forEachDirective.Items) == 0 && len(forEachDirective.ItemFiles) == 0 {
		return errors.NewError(errors.ErrorTypeConfigMerge, "forEach must define either 'items' or 'itemFiles'")
	}
	if forEachDirective.Output.FilenamePattern == "" {
		return errors.NewError(errors.ErrorTypeCLIValidation, "forEach.output.filenamePattern is required")
	}

	iterationSources := []map[string]interface{}{}
	itemFileBasenames := []string{} // For ${ITEM_FILE_BASENAME}

	if len(forEachDirective.Items) > 0 {
		logger.Debug("Iterating using 'items' from forEach.")
		if strings.Contains(forEachDirective.Output.FilenamePattern, "${ITEM_FILE_BASENAME}") {
			logger.Warn("forEach.output.filenamePattern uses ${ITEM_FILE_BASENAME} but 'items' mode is active — this placeholder will resolve to an empty string. Use 'itemFiles' instead or remove the placeholder.")
		}
		iterationSources = forEachDirective.Items
		for range forEachDirective.Items {
			itemFileBasenames = append(itemFileBasenames, "")
		}
	} else { // len(forEachDirective.ItemFiles) > 0
		logger.Debug("Iterating using 'itemFiles' from forEach.")
		for _, itemFilePath := range forEachDirective.ItemFiles {
			fullItemFilePath, err := p.resolveItemFilePath(itemFilePath)
			if err != nil {
				return err
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
			logger.Warn("Output format for iteration %d (%s) is ambiguous, defaulting to YAML. Specify format in forEach.output.format or use a file extension.", i, outputFilename)
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

// resolveItemFilePath resolves and validates an itemFile path, preventing path traversal.
func (p *Pipeline) resolveItemFilePath(itemFilePath string) (string, error) {
	fullItemFilePath := itemFilePath
	if !filepath.IsAbs(itemFilePath) {
		if p.Config.VarsFile == "" {
			return "", errors.NewError(errors.ErrorTypeCLIValidation, "itemFiles with relative paths require a vars file (-V) to resolve against")
		}
		fullItemFilePath = filepath.Join(filepath.Dir(p.Config.VarsFile), itemFilePath)
	}

	// Resolve symlinks to prevent containment bypass, then validate
	absItemPath, err := filepath.Abs(fullItemFilePath)
	if err != nil {
		return "", errors.WrapError(errors.ErrorTypeFileRead, "failed to resolve itemFile path", err).WithContext("file", itemFilePath)
	}
	// Resolve symlinks on the parent directory (the file itself may not exist yet for validation)
	realItemPath, err := filepath.EvalSymlinks(filepath.Dir(absItemPath))
	if err != nil {
		return "", errors.WrapError(errors.ErrorTypeFileRead, "failed to resolve itemFile symlinks", err).WithContext("file", itemFilePath)
	}
	realItemPath = filepath.Join(realItemPath, filepath.Base(absItemPath))

	// Only enforce containment when we have a vars file to contain against
	if p.Config.VarsFile != "" {
		allowedBase, err := filepath.Abs(filepath.Dir(p.Config.VarsFile))
		if err != nil {
			return "", errors.WrapError(errors.ErrorTypeFileRead, "failed to resolve vars file directory", err)
		}
		realAllowedBase, err := filepath.EvalSymlinks(allowedBase)
		if err != nil {
			return "", errors.WrapError(errors.ErrorTypeFileRead, "failed to resolve vars directory symlinks", err)
		}
		if !strings.HasPrefix(realItemPath, realAllowedBase+string(filepath.Separator)) && realItemPath != realAllowedBase {
			return "", errors.NewErrorf(errors.ErrorTypeCLIValidation, "itemFile path %q escapes the vars file directory %q", itemFilePath, realAllowedBase)
		}
	}

	return realItemPath, nil
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

	var unresolvedVars []string
	resolvedPattern = filenameVarRegex.ReplaceAllStringFunc(resolvedPattern, func(match string) string {
		varName := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")

		// 1. Check iterVars (highest priority for filename context)
		if val, ok := iterVars[varName]; ok {
			strVal := fmt.Sprintf("%v", val)
			// Reject variable values containing path traversal components
			if strings.Contains(strVal, "..") || filepath.IsAbs(strVal) {
				unresolvedVars = append(unresolvedVars, fmt.Sprintf("%s (unsafe path component)", match))
				return ""
			}
			return strVal
		}
		// 2. Check envVarsForSchema (KONFIGO_VAR_...)
		if val, ok := envVarsForSchema[varName]; ok {
			if strings.Contains(val, "..") || filepath.IsAbs(val) {
				unresolvedVars = append(unresolvedVars, fmt.Sprintf("%s (unsafe path component)", match))
				return ""
			}
			return val
		}
		// 3. Check schemaVars (default values or simple values if not from env/path)
		for _, sv := range schemaVars {
			if sv.Name == varName {
				val := sv.Value
				if val == "" {
					val = sv.DefaultValue
				}
				if val != "" {
					if strings.Contains(val, "..") || filepath.IsAbs(val) {
						unresolvedVars = append(unresolvedVars, fmt.Sprintf("%s (unsafe path component from schema var)", match))
						return ""
					}
					return val
				}
			}
		}
		// Collect all unresolved variable names
		unresolvedVars = append(unresolvedVars, match)
		return ""
	})

	if len(unresolvedVars) > 0 {
		return "", fmt.Errorf("unresolved variables in filenamePattern: %s", strings.Join(unresolvedVars, ", "))
	}

	// Clean up path
	resolvedPattern = filepath.Clean(resolvedPattern)

	// Final safety check: reject patterns that resolve to absolute paths or escape upward
	if filepath.IsAbs(resolvedPattern) {
		return "", fmt.Errorf("resolved filenamePattern %q must be a relative path", resolvedPattern)
	}
	if strings.HasPrefix(resolvedPattern, "..") {
		return "", fmt.Errorf("resolved filenamePattern %q escapes the output directory", resolvedPattern)
	}

	return resolvedPattern, nil
}
