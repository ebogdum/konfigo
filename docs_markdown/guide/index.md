# Konfigo User Guide

Welcome to the Konfigo User Guide! This guide provides detailed information on how to use Konfigo, from basic command-line operations to advanced schema-driven configuration processing.

## Table of Contents

*   **[CLI Reference](./cli-reference.md)**
    *   A comprehensive reference for all command-line flags and options available in Konfigo. Learn how to specify input sources, control output formats, enable schema processing, and manage logging.

*   **[Environment Variables](./environment-variables.md)**
    *   Understand how Konfigo utilizes environment variables for both direct configuration overrides (`KONFIGO_KEY_...`) and for supplying dynamic values for variable substitution (`KONFIGO_VAR_...`).

*   **[Use Cases](./use-cases.md)** (To be expanded)
    *   Explore practical examples and scenarios where Konfigo can simplify and enhance your configuration management workflows.
        *   Managing configurations for different environments (dev, staging, prod).
        *   Generating multiple similar configuration files (e.g., for microservices).
        *   Validating configuration against a strict contract.
        *   Transforming legacy configuration structures to a new format.

## Core Workflow

The typical workflow with Konfigo involves these steps:

1.  **Prepare Configuration Sources**:
    *   Your configurations can be spread across multiple files (JSON, YAML, TOML, .env) and directories.
    *   You might also have values set via `KONFIGO_KEY_...` environment variables.

2.  **Define a Schema (Optional but Recommended for Advanced Use)**:
    *   Create a schema file (YAML, JSON, or TOML) that specifies how Konfigo should process your data. This includes:
        *   [Variable definitions and substitution rules](../schema/variables.md)
        *   [Data generation rules](../schema/generation.md)
        *   [Data transformation rules](../schema/transformation.md)
        *   [Data validation rules](../schema/validation.md)
        *   Optionally, `inputSchema` for pre-validation and `outputSchema` for post-filtering.
        *   Optionally, `immutable` paths to protect certain keys.

3.  **Run Konfigo**:
    *   Use the `konfigo` command with appropriate flags:
        *   `-s` to specify your source files/directories.
        *   `-S` to apply your schema.
        *   `-V` to provide a variables file (which can also trigger `konfigo_forEach` batch mode).
        *   `-of` or `-o<format>` to control the output.
        *   Other flags to control behavior like recursion (`-r`), case-sensitivity (`-c`), etc.

4.  **Processing Steps (Simplified View)**:
    a.  **Load Sources**: Konfigo reads and parses all specified source files.
    b.  **Merge Configurations**: The parsed data is merged into a single configuration map. `KONFIGO_KEY_...` environment variables are applied as overrides at this stage. Immutable paths are respected.
    c.  **(If Schema Provided)**:
        i.  **Input Schema Validation**: If `inputSchema` is defined, the merged config is validated against it.
        ii. **Variable Resolution**: Variables are resolved based on precedence (`KONFIGO_VAR_...` > `-V` file > schema `vars`).
        iii. **Generators**: Data generation rules are applied.
        iv. **Transformers**: Data transformation rules are applied.
        v.  **Global Variable Substitution**: `${VAR_NAME}` placeholders are substituted throughout the configuration.
        vi. **Validation**: The processed configuration is validated against the `validate` rules in the schema.
        vii. **Output Schema Filtering**: If `outputSchema` is defined, the configuration is filtered.
    d.  **(If Batch Mode with `konfigo_forEach` in `-V` file and Schema Provided)**:
        *   Steps c.ii through c.vii are performed for *each iteration* defined in `konfigo_forEach`, using a deep copy of the merged configuration from step 4.b and iteration-specific variables. Each iteration produces its own output file.
    e.  **Output**: The final configuration (or multiple configurations in batch mode) is written to the specified output file(s) or to stdout in the chosen format.

This guide, along with the [Schema Documentation](../schema/index.md), aims to provide you with all the information needed to master Konfigo.
