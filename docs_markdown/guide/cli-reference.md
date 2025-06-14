# CLI Reference

Konfigo is a versatile command-line tool for merging, processing, and transforming configuration files. This page details all available command-line flags and options.

## Synopsis

```bash
konfigo [flags] -s <sources...>
konfigo [flags] -s - -s<format_flag> # Reading from stdin
cat config.yml | konfigo -sy -S schema.yml # Example with stdin and schema
```

## Flags

Flags are used to control Konfigo's behavior, from specifying input sources and output formats to enabling schema processing and managing logging.

### Input & Sources

These flags control how Konfigo discovers and parses your input configuration files.

*   `-s <paths>`:
    *   **Description**: A comma-separated list of source files or directories. Konfigo will read and merge these sources in the order they are provided.
    *   Use `-` to specify reading from standard input (stdin). When using stdin, you **must** also specify the input format using one of the `-s<format>` flags (e.g., `-sy` for YAML).
    *   **Example**: `konfigo -s base.json,env/dev.yml,secrets.env`
    *   **Example (stdin)**: `cat my_config.json | konfigo -s - -sj`

*   `-r`:
    *   **Description**: Recursively search for configuration files in subdirectories of any directories specified in `-s`.
    *   Konfigo identifies files by common configuration extensions (e.g., `.json`, `.yaml`, `.yml`, `.toml`, `.env`).
    *   **Example**: `konfigo -s ./configs -r`

*   `-sj`:
    *   **Description**: Force input to be parsed as JSON.
    *   This is **required** if reading JSON content from stdin (`-s -`).
    *   **Example**: `echo '{"key": "value"}' | konfigo -s - -sj`

*   `-sy`:
    *   **Description**: Force input to be parsed as YAML.
    *   This is **required** if reading YAML content from stdin (`-s -`).
    *   **Example**: `echo 'key: value' | konfigo -s - -sy`

*   `-st`:
    *   **Description**: Force input to be parsed as TOML.
    *   This is **required** if reading TOML content from stdin (`-s -`).
    *   **Example**: `echo 'key = "value"' | konfigo -s - -st`

*   `-se`:
    *   **Description**: Force input to be parsed as an ENV file.
    *   This is **required** if reading ENV content from stdin (`-s -`).
    *   **Example**: `echo 'KEY=value' | konfigo -s - -se`

### Schema & Variables

These flags enable Konfigo's powerful schema-driven processing and variable substitution features.

*   `-S, --schema <path>`:
    *   **Description**: Path to a schema file (must be YAML, JSON, or TOML). This schema defines how the merged configuration should be processed, including variable resolution, data generation, transformations, and validation.
    *   Refer to the [Schema Documentation](../schema/index.md) for details on schema structure and capabilities.
    *   **Example**: `konfigo -s config.yml -S schema.yml`

*   `-V, --vars-file <path>`:
    *   **Description**: Path to a file (YAML, JSON, or TOML) providing high-priority variables for substitution within your schema and configuration.
    *   Variables from this file override those defined in the schema's `vars` block but are themselves overridden by `KONFIGO_VAR_...` environment variables.
    *   This file can also contain the `konfigo_forEach` directive for batch processing.
    *   **Example**: `konfigo -s config.yml -S schema.yml -V prod-vars.yml`
    *   See [Variable Precedence](#variable-precedence) and [Batch Processing with `konfigo_forEach`](../schema/variables.md#batch-processing-with-konfigo_foreach) for more details.

#### Variable Precedence

Konfigo resolves variables used in `${VAR_NAME}` substitutions with the following priority (1 is highest):

1.  **Environment Variables**: Set as `KONFIGO_VAR_VARNAME=value`. (See [Environment Variables](./environment-variables.md))
2.  **Variables File**: Variables defined in the file specified by `-V` or `--vars-file`.
    *   In batch mode (`konfigo_forEach`), iteration-specific variables take precedence over global variables within this file.
3.  **Schema `vars` Block**: Variables defined within the `vars:` section of the schema file specified by `-S`.

### Output & Formatting

These flags control the format and destination of Konfigo's output.

*   `-of <path>`:
    *   **Description**: Write the final processed configuration to the specified file path.
    *   If the filename has an extension (e.g., `.json`, `.yaml`, `.toml`, `.env`), Konfigo will use that extension to determine the output format.
    *   If used in conjunction with specific format flags (`-oj`, `-oy`, etc.), this path acts as a base name, and the format flag's extension will be appended. For example, `konfigo -s c.json -of out/config -oy` would write to `out/config.yaml`.
    *   If this flag is not provided, output is sent to standard output (stdout), defaulting to YAML format unless overridden by an `-o<format>` flag.
    *   **Example (extension determines format)**: `konfigo -s c.json -of config.yaml`
    *   **Example (used as base name)**: `konfigo -s c.json -of config -oj -oy` (writes `config.json` and `config.yaml`)

*   `-oj`:
    *   **Description**: Output the final configuration in JSON format.
    *   **Example**: `konfigo -s c.yml -oj` (outputs JSON to stdout)

*   `-oy`:
    *   **Description**: Output the final configuration in YAML format. This is the default output format if no other output flags are specified.
    *   **Example**: `konfigo -s c.json -oy` (outputs YAML to stdout)

*   `-ot`:
    *   **Description**: Output the final configuration in TOML format.
    *   **Example**: `konfigo -s c.json -ot` (outputs TOML to stdout)

*   `-oe`:
    *   **Description**: Output the final configuration in ENV file format.
    *   **Example**: `konfigo -s c.json -oe` (outputs ENV to stdout)

### Behavior & Logging

These flags adjust Konfigo's operational behavior and the verbosity of its logging.

*   `-c`:
    *   **Description**: Use case-sensitive key matching during merging.
    *   By default, Konfigo performs case-insensitive key matching (e.g., `key` and `Key` would be treated as the same key, with the latter overriding the former if it appears later in the merge sequence).
    *   **Example**: `konfigo -s c.json -c`

*   `-v`:
    *   **Description**: Enable verbose debug logging. This provides detailed information about the steps Konfigo is taking, which can be helpful for troubleshooting.
    *   This flag is overridden by `-q`.
    *   **Example**: `konfigo -s c.json -v`

*   `-q`:
    *   **Description**: Suppress all logging output except for the final configuration data (if outputting to stdout) or critical errors.
    *   This flag overrides `-v`.
    *   **Example**: `konfigo -s c.json -q`

*   `-h`:
    *   **Description**: Show the help message, which includes a summary of all available flags and basic usage instructions.
    *   **Example**: `konfigo -h`

## Exit Codes

*   **0**: Successful execution.
*   **Non-zero**: An error occurred. Error messages are typically printed to stderr.

For more details on how Konfigo uses environment variables, see the [Environment Variables](./environment-variables.md) page. For a deep dive into schema capabilities, refer to the [Schema documentation](../schema/index.md).
