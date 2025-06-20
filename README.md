# Konfigo: Versatile Configuration Management

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Konfigo is a powerful command-line tool designed to streamline your configuration workflow. It excels at reading various configuration file formats (JSON, YAML, TOML, .env), merging them intelligently, and processing the combined data against a user-defined schema for validation, transformation, variable substitution, and even batch output generation.

Whether you're managing simple settings or complex, multi-layered configurations with environment-specific overrides, Konfigo provides the tools to do so efficiently and reliably.

## Key Features

*   **Multi-Format Support**: Reads and writes JSON, YAML, TOML, and .env files.
*   **Flexible Merging**: Intelligently merges multiple configuration sources, respecting order and immutability rules.
*   **Powerful Schema Processing**:
    *   **Variable Substitution**: Inject dynamic values from environment variables (`KONFIGO_VAR_...`), dedicated variable files (`-V`), or schema defaults.
    *   **Data Generation**: Create new configuration values (e.g., `concat` strings).
    *   **Data Transformation**: Modify keys and values (e.g., `renameKey`, `changeCase`, `addKeyPrefix`, `setValue`).
    *   **Data Validation**: Enforce rules (`required`, `type`, `min`, `max`, `minLength`, `enum`, `regex`).
    *   **Input/Output Schemas**: Validate incoming data and filter outgoing data against defined structures.
*   **Batch Processing**: Use the `konfigo_forEach` directive in a variables file to generate multiple tailored configuration outputs from a single schema and run.
*   **Environment Variable Integration**:
    *   Override any configuration value directly using `KONFIGO_KEY_path.to.key=value`.
*   **Comprehensive CLI**: Rich set of command-line options for fine-grained control over input, output, and processing behavior.

## Getting Started

### 1. Installation

The primary way to install Konfigo is downloading a pre-built binary from the release page:

[Releases](https://github.com/ebogdum/konfigo/releases)

For other installation methods, please refer to the [Installation Guide](docs_markdown/installation.md) in our documentation.

### 2. Basic Usage

Merge two configuration files (`config.json` and `overrides.yml`) and output the result to `final.yml`:

```bash
konfigo -s config.json,overrides.yml -of final.yml
```

### 3. Using a Schema

Merge `config.json`, process it with `schema.yml`, use variables from `staging-vars.yml`, and output to `staging_config.json`:

```bash
konfigo -s config.json -S schema.yml -V staging-vars.yml -of staging_config.json
```

## Command-Line Options

Below is a summary of the available command-line options. For more details, run `konfigo -h`.

**Input & Sources**

| Flag(s)                 | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `-s <paths>`            | Comma-separated list of source files/directories. Use '-' for stdin.        |
| `-r`                    | Recursively search for configuration files in subdirectories.               |
| `-sj`                   | Force input to be parsed as JSON (required for stdin).                      |
| `-sy`                   | Force input to be parsed as YAML (required for stdin).                      |
| `-st`                   | Force input to be parsed as TOML (required for stdin).                      |
| `-se`                   | Force input to be parsed as ENV (required for stdin).                       |
| `-si`                   | Force input to be parsed as INI (required for stdin).                       |

**Schema & Variables**

| Flag(s)                 | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `-S, --schema <path>`   | Path to a schema file for processing the config.                            |
| `-V, --vars-file <path>`| Path to a file providing high-priority variables.                           |

**Output & Formatting**

| Flag(s)                 | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `-of <path>`            | Write output to file. Extension determines format, or use with -oX flags.   |
| `-oj`                   | Output in JSON format.                                                      |
| `-oy`                   | Output in YAML format.                                                      |
| `-ot`                   | Output in TOML format.                                                      |
| `-oe`                   | Output in ENV format.                                                       |
| `-oi`                   | Output in INI format.                                                       |

**Behavior & Logging**

| Flag(s)                 | Description                                                                 |
|-------------------------|-----------------------------------------------------------------------------|
| `-c`                    | Use case-sensitive key matching (default is case-insensitive).              |
| `-v`                    | Enable informational (INFO) logging. Overrides default quiet behavior.      |
| `-d`                    | Enable debug (DEBUG and INFO) logging. Overrides -v and default quiet behavior. |
| `-h`                    | Show this help message.                                                     |

## Documentation

For detailed information on all features, CLI options, and schema capabilities, please visit our full documentation site:

**[Konfigo Documentation Site](https://ebogdum.github.io/konfigo/)**

Alternatively, you can browse the Markdown files directly in the [`/docs_markdown`](docs_markdown) directory.

Key sections:
*   [User Guide](docs_markdown/guide/index.md)
*   [Schema Guide](docs_markdown/schema/index.md)

## Contributing

Contributions are welcome!

## License

Konfigo is licensed under the [MIT License](./LICENSE).

