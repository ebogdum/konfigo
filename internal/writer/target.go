package writer

import (
	"fmt"
	"path/filepath"
	"strings"
)

// OutputTarget represents a single output destination.
type OutputTarget struct {
	Format   string
	Filename string // Empty for stdout
}

// DetermineOutputTargets determines the output targets based on flags and file specifications.
func DetermineOutputTargets(outputFile string, outJSON bool, outYAML bool, outTOML bool, outENV bool) []OutputTarget {
	var targets []OutputTarget
	
	// If output file is specified with extension, use that format
	if outputFile != "" && filepath.Ext(outputFile) != "" {
		format := strings.TrimPrefix(filepath.Ext(outputFile), ".")
		targets = append(targets, OutputTarget{Format: format, Filename: outputFile})
		return targets
	}
	
	// Add targets based on format flags
	if outJSON {
		targets = append(targets, OutputTarget{Format: "json"})
	}
	if outYAML {
		targets = append(targets, OutputTarget{Format: "yaml"})
	}
	if outTOML {
		targets = append(targets, OutputTarget{Format: "toml"})
	}
	if outENV {
		targets = append(targets, OutputTarget{Format: "env"})
	}
	
	// Default to JSON if no format specified
	if len(targets) == 0 {
		targets = append(targets, OutputTarget{Format: "json"})
	}
	
	// If output file is specified without extension, create filenames for each format
	if outputFile != "" {
		baseName := outputFile
		for i := range targets {
			targets[i].Filename = fmt.Sprintf("%s.%s", baseName, targets[i].Format)
		}
	}
	
	return targets
}

// WriteOutput writes content to the specified output target.
// If filename is empty, it writes to stdout.
func (ot OutputTarget) WriteOutput(content []byte) error {
	if ot.Filename == "" {
		return WriteToStdout(content)
	}
	return WriteFile(ot.Filename, content)
}

// IsStdout returns true if this target writes to stdout.
func (ot OutputTarget) IsStdout() bool {
	return ot.Filename == ""
}

// String returns a string representation of the output target.
func (ot OutputTarget) String() string {
	if ot.IsStdout() {
		return fmt.Sprintf("%s (stdout)", ot.Format)
	}
	return fmt.Sprintf("%s -> %s", ot.Format, ot.Filename)
}
