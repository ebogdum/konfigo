package cli

import (
	"flag"
	"fmt"
)

// PrintHelp displays the comprehensive help message for Konfigo
func PrintHelp() {
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
	fmt.Fprintf(out, "    (Default behavior is quiet; no informational or debug logs are printed unless specified.)\n")
	fmt.Fprintf(out, "    -c\t\tUse case-sensitive key matching (default is case-insensitive).\n")
	fmt.Fprintf(out, "    -v\t\tEnable informational (INFO) logging.\n")
	fmt.Fprintf(out, "    -d\t\tEnable debug (DEBUG and INFO) logging. Overrides -v.\n")
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

// SetCustomUsage sets the flag.Usage to use our PrintHelp function
func SetCustomUsage() {
	flag.Usage = PrintHelp
}
