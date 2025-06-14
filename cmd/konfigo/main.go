package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"konfigo/internal/loader"
	"konfigo/internal/logger"
	"konfigo/internal/marshaller"
	"konfigo/internal/merger"
	"konfigo/internal/parser"
	"konfigo/internal/schema"
	"konfigo/internal/util"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3" // Used for parsing konfigo_forEach from vars file
)

type outputTarget struct {
	Format   string
	Filename string
}

type parseResult struct {
	FilePath string
	Data     map[string]interface{}
	Err      error
}

func printHelp() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Konfigo: A versatile tool for merging and converting configuration files.\n\n")
	fmt.Fprintf(out, "DESCRIPTION:\n")
	fmt.Fprintf(out, "  Konfigo reads configuration files, merges them, and processes them against a schema\n")
	fmt.Fprintf(out, "  to validate, transform, and generate final configuration values.\n\n")
	fmt.Fprintf(out, "USAGE:\n")
	fmt.Fprintf(out, "  konfigo [flags] -s <sources...>\n")
	fmt.Fprintf(out, "  cat config.yml | konfigo -sy -S schema.yml\n\n")
	fmt.Fprintf(out, "FLAGS:\n")
	fmt.Fprintf(out, "  Input & Sources:\n")
	fmt.Fprintf(out, "    -s <paths>\tComma-separated list of source files/directories. Use '-' for stdin.\n")
	fmt.Fprintf(out, "    -r\t\tRecursively search for configuration files in subdirectories.\n")
	fmt.Fprintf(out, "    -sj, -sy, -st, -se\n\t\tForce input to be parsed as a specific format (required for stdin).\n\n")
	fmt.Fprintf(out, "  Schema & Variables:\n")
	fmt.Fprintf(out, "    -S, --schema <path>\n\t\tPath to a schema file (YAML, JSON, TOML) for processing the config.\n")
	fmt.Fprintf(out, "    -V, --vars-file <path>\n\t\tPath to a file providing high-priority variables for substitution.\n\n")
	fmt.Fprintf(out, "    Variable Priority:\n")
	fmt.Fprintf(out, "    Variable values are resolved with the following priority (1 is highest):\n")
	fmt.Fprintf(out, "      1. Environment variables (KONFIGO_VAR_...).\n")
	fmt.Fprintf(out, "      2. Variables from the --vars-file (-V).\n")
	fmt.Fprintf(out, "      3. Variables defined in the schema's `vars:` section (-S).\n\n")
	fmt.Fprintf(out, "  Output & Formatting:\n")
	fmt.Fprintf(out, "    -of <path>\tWrite output to file. Extension determines format, or use with -oX flags.\n")
	fmt.Fprintf(out, "    -oj, -oy, -ot, -oe\n\t\tOutput in a specific format.\n\n")
	fmt.Fprintf(out, "  Behavior & Logging:\n")
	fmt.Fprintf(out, "    -c\t\tUse case-sensitive key matching (default is case-insensitive).\n")
	fmt.Fprintf(out, "    -v\t\tEnable verbose debug logging.\n")
	fmt.Fprintf(out, "    -q\t\tSuppress all logging except for final output; overrides -v.\n")
	fmt.Fprintf(out, "    -h\t\tShow this help message.\n\n")
	fmt.Fprintf(out, "ENVIRONMENT VARIABLES:\n")
	fmt.Fprintf(out, "  Konfigo reads two types of environment variables:\n")
	fmt.Fprintf(out, "  - KONFIGO_KEY_path.to.key=value\n")
	fmt.Fprintf(out, "    Sets a configuration value. Has the highest precedence over all file sources.\n")
	fmt.Fprintf(out, "    Example: KONFIGO_KEY_database.port=5432\n\n")
	fmt.Fprintf(out, "  - KONFIGO_VAR_VARNAME=value\n")
	fmt.Fprintf(out, "    Sets a substitution variable. Has the highest precedence for variables.\n")
	fmt.Fprintf(out, "    Example: KONFIGO_VAR_RELEASE_VERSION=1.2.3\n")
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// --- 1. Setup Flags ---
	schemaFile := flag.String("schema", "", "Path to a schema file for processing the config.")
	flag.StringVar(schemaFile, "S", "", "Path to a schema file (shorthand for --schema).")
	varsFileFlag := flag.String("vars-file", "", "Path to a file providing high-priority variables.")
	flag.StringVar(varsFileFlag, "V", "", "Path to a variables file (shorthand for --vars-file).")
	sourcePaths := flag.String("s", "", "Comma-separated list of source files/directories. Use '-' for stdin.")
	recursive := flag.Bool("r", false, "Recursively search for configuration files in subdirectories")
	caseSensitive := flag.Bool("c", false, "Use case-sensitive key matching (default is case-insensitive)")
	verbose := flag.Bool("v", false, "Enable verbose debug logging")
	quiet := flag.Bool("q", false, "Suppress all logging except for final output; overrides -v")
	inJSON := flag.Bool("sj", false, "Force input to be parsed as JSON (required for stdin)")
	inYAML := flag.Bool("sy", false, "Force input to be parsed as YAML (required for stdin)")
	inTOML := flag.Bool("st", false, "Force input to be parsed as TOML (required for stdin)")
	inENV := flag.Bool("se", false, "Force input to be parsed as ENV (required for stdin)")
	outputFile := flag.String("of", "", "Write output to file. Extension determines format, or use with -oX flags.")
	outJSON := flag.Bool("oj", false, "Output in JSON format")
	outYAML := flag.Bool("oy", false, "Output in YAML format")
	outTOML := flag.Bool("ot", false, "Output in TOML format")
	outENV := flag.Bool("oe", false, "Output in ENV format")
	help := flag.Bool("h", false, "Show this help message.")

	flag.Usage = printHelp
	flag.Parse()

	// --- 2. Initialize Logger and Handle Help/No-Argument Case ---
	if *quiet {
		*verbose = false
	}
	logger.Init(*verbose, *quiet)
	if *help || len(os.Args) == 1 {
		printHelp()
		return nil
	}

	// --- 3. Load Schema and Environment Variables Early ---
	var loadedSchema *schema.Schema
	if *schemaFile != "" {
		var err error
		logger.Log("Loading schema from %s", *schemaFile)
		loadedSchema, err = schema.Load(*schemaFile)
		if err != nil {
			return err
		}
	}

	immutablePaths := make(map[string]struct{})
	if loadedSchema != nil {
		for _, path := range loadedSchema.Immutable {
			immutablePaths[path] = struct{}{}
		}
	}

	envConfig, envVarsForSchema := loadFromEnv() // envVarsForSchema are KONFIGO_VAR_...

	// --- 4. Load, Parse, and Merge Configurations ---
	if *sourcePaths == "" && flag.NArg() > 0 {
		*sourcePaths = strings.Join(flag.Args(), ",")
	}
	if *sourcePaths == "" {
		return errors.New("no input source specified. Use -s <paths> or pipe from stdin")
	}

	inputFormatOverride := ""
	if *inJSON {
		inputFormatOverride = "json"
	} else if *inYAML {
		inputFormatOverride = "yaml"
	} else if *inTOML {
		inputFormatOverride = "toml"
	} else if *inENV {
		inputFormatOverride = "env"
	}

	baseFinalConfig, err := processSources(*sourcePaths, *recursive, *caseSensitive, inputFormatOverride, immutablePaths, envConfig)
	if err != nil {
		return err
	}

	// --- 5. Load Variables File and Check for konfigo_forEach ---
	var varsFromFileGlobal map[string]interface{}
	var forEachDirective *schema.KonfigoForEach

	if *varsFileFlag != "" {
		logger.Log("Loading variables from %s", *varsFileFlag)
		content, err := os.ReadFile(*varsFileFlag)
		if err != nil {
			return fmt.Errorf("failed to read vars file %s: %w", *varsFileFlag, err)
		}
		rawVarsFromFile, err := parser.Parse(*varsFileFlag, content, "")
		if err != nil {
			return fmt.Errorf("failed to parse vars file %s: %w", *varsFileFlag, err)
		}

		// Separate konfigo_forEach from other global vars
		varsFromFileGlobal = make(map[string]interface{})
		for k, v := range rawVarsFromFile {
			if k == "konfigo_forEach" {
				// Marshal to YAML and then Unmarshal to KonfigoForEach struct
				// This is a common way to convert map[string]interface{} to a struct
				yamlBytes, err := yaml.Marshal(v)
				if err != nil {
					return fmt.Errorf("failed to marshal konfigo_forEach directive: %w", err)
				}
				forEachDirective = &schema.KonfigoForEach{}
				if err := yaml.Unmarshal(yamlBytes, forEachDirective); err != nil {
					return fmt.Errorf("failed to unmarshal konfigo_forEach directive: %w", err)
				}
				logger.Debug("Found konfigo_forEach directive.")
			} else {
				varsFromFileGlobal[k] = v
			}
		}
		if forEachDirective != nil {
			forEachDirective.GlobalVars = varsFromFileGlobal // Store global vars for resolver
		}
	}

	// --- 6. Process Schema (Single or Batch) ---
	if loadedSchema == nil && forEachDirective != nil {
		return errors.New("konfigo_forEach directive found in variables file, but no schema file (-S) was provided for processing")
	}

	if loadedSchema != nil {
		if forEachDirective != nil {
			// Batch Processing Mode
			logger.Log("Starting batch processing with konfigo_forEach...")

			if len(forEachDirective.Items) > 0 && len(forEachDirective.ItemFiles) > 0 {
				return errors.New("konfigo_forEach cannot have both 'items' and 'itemFiles' defined simultaneously")
			}
			if len(forEachDirective.Items) == 0 && len(forEachDirective.ItemFiles) == 0 {
				return errors.New("konfigo_forEach must define either 'items' or 'itemFiles'")
			}
			if forEachDirective.Output.FilenamePattern == "" {
				return errors.New("konfigo_forEach.output.filenamePattern is required")
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
					if !filepath.IsAbs(itemFilePath) && *varsFileFlag != "" {
						fullItemFilePath = filepath.Join(filepath.Dir(*varsFileFlag), itemFilePath)
					}
					logger.Debug("Loading iteration variables from itemFile: %s", fullItemFilePath)
					content, err := os.ReadFile(fullItemFilePath)
					if err != nil {
						return fmt.Errorf("failed to read itemFile %s: %w", fullItemFilePath, err)
					}
					itemVars, err := parser.Parse(fullItemFilePath, content, "")
					if err != nil {
						return fmt.Errorf("failed to parse itemFile %s: %w", fullItemFilePath, err)
					}
					iterationSources = append(iterationSources, itemVars)
					itemFileBasenames = append(itemFileBasenames, strings.TrimSuffix(filepath.Base(fullItemFilePath), filepath.Ext(fullItemFilePath)))
				}
			}

			for i, iterVars := range iterationSources {
				logger.Log("Processing iteration %d...", i)

				// Create a deep copy of the base configuration for this iteration
				currentConfig, err := util.DeepCopyMap(baseFinalConfig)
				if err != nil {
					return fmt.Errorf("failed to deep copy base config for iteration %d: %w", i, err)
				}

				// Prepare variables for this iteration based on precedence:
				// Iteration Vars > Global Vars from Vars File > Schema Vars (handled by resolver) > Env Vars (handled by resolver)
				varsForThisIteration := make(map[string]interface{})
				if forEachDirective.GlobalVars != nil { // Start with globals from the main vars file
					for k, v := range forEachDirective.GlobalVars {
						varsForThisIteration[k] = v
					}
				}
				for k, v := range iterVars { // Override with iteration-specific vars
					varsForThisIteration[k] = v
				}

				// Add ITEM_INDEX and ITEM_FILE_BASENAME to iteration vars
				varsForThisIteration["ITEM_INDEX"] = strconv.Itoa(i)
				if i < len(itemFileBasenames) { // Check bounds for safety
					varsForThisIteration["ITEM_FILE_BASENAME"] = itemFileBasenames[i]
				}

				processedConfig, err := schema.Process(currentConfig, loadedSchema, varsForThisIteration, envVarsForSchema)
				if err != nil {
					return fmt.Errorf("schema processing failed for iteration %d: %w", i, err)
				}

				// Determine output filename for this iteration
				outputFilename, err := resolveFilenamePattern(forEachDirective.Output.FilenamePattern, varsForThisIteration, envVarsForSchema, loadedSchema.Vars, i, itemFileBasenames[i])
				if err != nil {
					return fmt.Errorf("failed to resolve output filename for iteration %d: %w", i, err)
				}

				// Ensure output directory exists
				outputDir := filepath.Dir(outputFilename)
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("failed to create output directory %s for iteration %d: %w", outputDir, i, err)
				}

				outputFormat := strings.ToLower(strings.TrimPrefix(filepath.Ext(outputFilename), "."))
				if forEachDirective.Output.Format != "" {
					outputFormat = strings.ToLower(forEachDirective.Output.Format)
				}
				if outputFormat == "" { // Fallback if no extension and no explicit format
					logger.Warn("Output format for iteration %d (%s) is ambiguous, defaulting to YAML. Specify format in konfigo_forEach.output.format or use a file extension.", i, outputFilename)
					outputFormat = "yaml"
				}

				outputBytes, err := marshaller.Marshal(processedConfig, outputFormat)
				if err != nil {
					return fmt.Errorf("error marshalling to %s for iteration %d (%s): %w", outputFormat, i, outputFilename, err)
				}
				logger.Log("Writing output for iteration %d to %s (format: %s)", i, outputFilename, outputFormat)
				if err := os.WriteFile(outputFilename, outputBytes, 0644); err != nil {
					return fmt.Errorf("error writing to file %s for iteration %d: %w", outputFilename, i, err)
				}
			}
			logger.Log("Batch processing completed.")
			return nil // Batch processing handles its own output, so we can return early.

		} else {
			// Single Processing Mode (original logic)
			varsToProcess := varsFromFileGlobal // Use only global vars if no forEach
			if varsToProcess == nil {
				varsToProcess = make(map[string]interface{}) // Ensure not nil for schema.Process
			}
			processedConfig, err := schema.Process(baseFinalConfig, loadedSchema, varsToProcess, envVarsForSchema)
			if err != nil {
				return err
			}
			baseFinalConfig = processedConfig // Update baseFinalConfig with the processed version
		}
	} else {
		// No schema provided, use baseFinalConfig as is
		logger.Debug("No schema provided, skipping schema processing.")
	}

	// --- 7. Determine and Generate Outputs (for single mode or if no schema) ---
	// This part is skipped if batch processing already handled outputs.
	// If forEachDirective was nil, this is the standard output path.
	// If forEachDirective was not nil, we returned early from the batch processing block.
	// So, this code only runs for the single output scenario.

	targets := determineOutputTargets(*outputFile, *outJSON, *outYAML, *outTOML, *outENV)
	if len(targets) == 0 && *outputFile == "" { // Default to YAML stdout if no other output specified
		targets = append(targets, outputTarget{Format: "yaml", Filename: ""})
	}

	for i, target := range targets {
		outputBytes, err := marshaller.Marshal(baseFinalConfig, target.Format)
		if err != nil {
			return fmt.Errorf("error marshalling to %s: %w", target.Format, err)
		}
		if target.Filename == "" { // Output to stdout
			if i > 0 && len(targets) > 1 { // Add separator for multiple stdout formats
				fmt.Println("\\n---")
			}
			fmt.Println(string(outputBytes))
		} else {
			logger.Log("Writing output to %s (format: %s)", target.Filename, target.Format)
			// Ensure output directory exists
			outputDir := filepath.Dir(target.Filename)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
			}
			if err := os.WriteFile(target.Filename, outputBytes, 0644); err != nil {
				return fmt.Errorf("error writing to file %s: %w", target.Filename, err)
			}
		}
	}
	return nil
}

// resolveFilenamePattern substitutes placeholders in the filename pattern.
// Placeholders: ${VAR_NAME}, ${ITEM_INDEX}, ${ITEM_FILE_BASENAME}
func resolveFilenamePattern(pattern string, iterVars map[string]interface{}, envVarsForSchema map[string]string, schemaVars []schema.VarDef, itemIndex int, itemFileBasename string) (string, error) {
	// Create a temporary resolver for this specific iteration's context
	// The 'config' map for resolver can be nil as we are only resolving based on explicit vars.

	// Priority for resolving VAR_NAME:
	// 1. Iteration-specific vars (iterVars)
	// 2. Environment KONFIGO_VAR_ (envVarsForSchema)
	// 3. Schema vars (default values or simple values if not from env/path)
	// For simplicity in filename, we'll manually check schema vars if not in iterVars or envVars.

	resolvedPattern := pattern

	// Substitute ${ITEM_INDEX}
	resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_INDEX}", strconv.Itoa(itemIndex))

	// Substitute ${ITEM_FILE_BASENAME}
	if itemFileBasename != "" {
		resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_FILE_BASENAME}", itemFileBasename)
	} else {
		// If itemFileBasename is empty (e.g. when using 'items'), remove the placeholder or leave it if not critical
		// For now, let's remove it to avoid literal "${ITEM_FILE_BASENAME}" in filenames.
		resolvedPattern = strings.ReplaceAll(resolvedPattern, "${ITEM_FILE_BASENAME}", "")
	}

	// Regex to find ${VAR_NAME}
	varRegex := regexp.MustCompile(`\\$\\{[A-Z0-9_]+\\}`)

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
		logger.Warn("Variable %s in filenamePattern not found, leaving as is.", match)
		return match // Leave unresolved
	})

	// Clean up path, e.g. remove double slashes if a variable was empty
	resolvedPattern = filepath.Clean(resolvedPattern)

	return resolvedPattern, nil
}

// ... existing code ...
// loadFromEnv, processSources, determineOutputTargets, etc. remain the same or with minor adjustments if needed.
// Make sure processSources returns baseFinalConfig which is the merged config *before* schema processing.
// The schema.Process function will then be called either once (single mode) or multiple times (batch mode)
// on a *copy* of this baseFinalConfig.

func loadFromEnv() (map[string]interface{}, map[string]string) {
	config := make(map[string]interface{})
	vars := make(map[string]string)

	keyPrefix := "KONFIGO_KEY_"
	varPrefix := "KONFIGO_VAR_"

	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		key, value := parts[0], parts[1]

		if strings.HasPrefix(key, keyPrefix) {
			configKey := strings.TrimPrefix(key, keyPrefix)
			logger.Debug("  - Loading from env config: %s -> %s", key, configKey)
			util.SetNestedValue(config, configKey, value)
		} else if strings.HasPrefix(key, varPrefix) {
			varName := strings.TrimPrefix(key, varPrefix)
			logger.Debug("  - Loading from env var: %s -> %s", key, varName)
			vars[varName] = value
		}
	}
	return config, vars
}

func processSources(sourcePaths string, recursive bool, caseSensitive bool, inputFormatOverride string, immutablePaths map[string]struct{}, envConfig map[string]interface{}) (map[string]interface{}, error) {
	var allFiles []string
	var stdinData []byte
	sources := strings.Split(sourcePaths, ",")
	logger.Log("Discovering configuration files...")
	for _, source := range sources {
		source = strings.TrimSpace(source)
		if source == "" {
			continue
		}
		if source == "-" {
			logger.Debug("Reading from standard input (stdin)")
			info, _ := os.Stdin.Stat()
			if (info.Mode() & os.ModeCharDevice) != 0 {
				return nil, errors.New("stdin is a terminal, not a pipe")
			}
			var err error
			stdinData, err = io.ReadAll(os.Stdin)
			if err != nil {
				return nil, fmt.Errorf("failed to read from stdin: %w", err)
			}
			if inputFormatOverride == "" {
				return nil, errors.New("reading from stdin requires an input format flag")
			}
			continue
		}
		files, err := loader.LoadFromPath(source, recursive)
		if err != nil {
			return nil, fmt.Errorf("error loading from source '%s': %w", source, err)
		}
		logger.Debug("Found %d file(s) in source: %s", len(files), source)
		allFiles = append(allFiles, files...)
	}

	results := parseFilesParallel(allFiles, inputFormatOverride)
	sort.Slice(results, func(i, j int) bool { return results[i].FilePath < results[j].FilePath })

	finalConfig := make(map[string]interface{})
	logger.Log("Merging %d configuration file(s)...", len(results))
	for _, res := range results {
		if res.Err != nil {
			logger.Log("  - Warning: Skipping file %s due to parse error: %v", res.FilePath, res.Err)
			continue
		}
		merger.Merge(finalConfig, res.Data, caseSensitive, immutablePaths)
	}

	if len(stdinData) > 0 {
		logger.Log("Merging configuration from stdin...")
		data, err := parser.Parse("stdin", stdinData, inputFormatOverride)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stdin: %w", err)
		}
		merger.Merge(finalConfig, data, caseSensitive, immutablePaths)
	}

	if len(envConfig) > 0 {
		logger.Log("Merging %d configuration key(s) from environment variables...", len(envConfig))
		merger.Merge(finalConfig, envConfig, caseSensitive, immutablePaths)
	}

	return finalConfig, nil
}

func parseFilesParallel(files []string, formatOverride string) []parseResult {
	if len(files) == 0 {
		return nil
	}
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, len(files))
	results := make(chan parseResult, len(files))
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				content, err := os.ReadFile(path)
				if err != nil {
					results <- parseResult{FilePath: path, Err: err}
					continue
				}
				data, err := parser.Parse(path, content, formatOverride)
				results <- parseResult{FilePath: path, Data: data, Err: err}
			}
		}()
	}
	for _, file := range files {
		jobs <- file
	}
	close(jobs)
	wg.Wait()
	close(results)
	var parsedResults []parseResult
	for res := range results {
		parsedResults = append(parsedResults, res)
	}
	return parsedResults
}

func determineOutputTargets(outputFile string, outJSON bool, outYAML bool, outTOML bool, outENV bool) []outputTarget {
	var targets []outputTarget
	if outputFile != "" && filepath.Ext(outputFile) != "" {
		format := strings.TrimPrefix(filepath.Ext(outputFile), ".")
		targets = append(targets, outputTarget{Format: format, Filename: outputFile})
		return targets
	}
	if outJSON {
		targets = append(targets, outputTarget{Format: "json"})
	}
	if outYAML {
		targets = append(targets, outputTarget{Format: "yaml"})
	}
	if outTOML {
		targets = append(targets, outputTarget{Format: "toml"})
	}
	if outENV {
		targets = append(targets, outputTarget{Format: "env"})
	}
	if len(targets) == 0 {
		targets = append(targets, outputTarget{Format: "json"})
	}
	if outputFile != "" {
		baseName := outputFile
		for i := range targets {
			targets[i].Filename = fmt.Sprintf("%s.%s", baseName, targets[i].Format)
		}
	}
	return targets
}

var (
	varRegex = regexp.MustCompile(`\\$\\{[A-Z0-9_]+\\}`) // Added for resolveFilenamePattern
)

// Make sure this regex is defined if not already (it was in schema/vars.go)
// var varRegex = regexp.MustCompile(`\\$\\{[A-Z0-9_]+\\}`)
