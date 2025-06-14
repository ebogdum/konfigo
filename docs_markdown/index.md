---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "Konfigo"
  text: "Documentation"
  tagline: "Simplify and Supercharge Your Configuration Management Workflow"
  actions:
    - theme: brand
      text: Quick Start
      link: /quick-start
    - theme: alt
      text: View User Guide
      link: /guide/
    - theme: alt
      text: Explore Schema
      link: /schema/

features:
  - title: "Multi-Format Support"
    details: "Seamlessly read, merge, and output JSON, YAML, TOML, and .env configuration files."
  - title: "Powerful Schema Processing"
    details: "Define variables, generate data, transform structures, and validate configurations with a flexible schema."
  - title: "Batch Output Generation"
    details: "Generate multiple tailored configuration files from a single schema using the 'konfigo_forEach' directive."
  - title: "Environment Integration"
    details: "Override configurations and supply variables directly through environment variables for dynamic setups."
  - title: "Flexible Merging"
    details: "Intelligently merge multiple configuration sources, respecting order and immutability rules."
  - title: "Comprehensive CLI"
    details: "Rich set of command-line options for fine-grained control over input, output, and processing."

---

<br>

Konfigo is a powerful and versatile command-line tool designed to simplify your configuration management workflow. It excels at reading various configuration file formats, merging them intelligently, and then processing the combined data against a user-defined schema. This schema can perform a wide array of operations, including:

*   **Variable Substitution**: Inject dynamic values from environment variables, dedicated variable files, or even other parts of your configuration.
*   **Data Generation**: Create new configuration values based on existing data (e.g., concatenating strings).
*   **Data Transformation**: Modify keys and values (e.g., renaming keys, changing string case, adding prefixes, setting static values).
*   **Data Validation**: Ensure your configuration adheres to specific rules and constraints (e.g., required fields, data types, numerical ranges, string patterns).
*   **Batch Processing**: Generate multiple output files from a single schema and a set of iterating variables, perfect for managing configurations across different environments or services.

Whether you're dealing with simple JSON files or complex, multi-layered YAML configurations with environment-specific overrides, Konfigo provides the tools to manage them efficiently and reliably.

## Key Features

*   **Multi-Format Support**: Reads JSON, YAML, TOML, and .env files.
*   **Flexible Merging**: Intelligently merges multiple configuration sources.
*   **Powerful Schema Processing**:
    *   Define variables with clear precedence (environment, vars file, schema defaults).
    *   Generate new data using `concat` and other potential generators.
    *   Transform data structures with operations like `renameKey`, `changeCase`, `addKeyPrefix`, and `setValue`.
    *   Validate configurations against a rich set of rules (`required`, `type`, `min`, `max`, `minLength`, `enum`, `regex`).
*   **Environment Variable Integration**:
    *   Override configuration values directly using `KONFIGO_KEY_path.to.key=value`.
    *   Supply variables for substitution using `KONFIGO_VAR_VARNAME=value`.
*   **Batch Output Generation**: Use the `konfigo_forEach` directive in your variables file to produce multiple tailored configuration outputs from a single run.
*   **Input/Output Control**:
    *   Read from files, directories (recursively), or stdin.
    *   Output to stdout or specified files.
    *   Control input and output formats (JSON, YAML, TOML, ENV).
*   **Schema Validation**: Validate input configurations against an `inputSchema` and filter outputs using an `outputSchema`.
*   **Immutability**: Protect specific configuration paths from being overridden during merges using the `immutable` schema directive.
*   **Customizable Behavior**: Options for case-sensitivity in key matching, verbose logging, and more.

## Getting Started

1.  **Installation**: (Add installation instructions here if available, e.g., `go install` or binary download links)
2.  **Basic Usage**: `konfigo -s source1.yml -s source2.json -of output.yml`
3.  **Using a Schema**: `konfigo -s config.json -S schema.yml -V staging-vars.yml -of staging_config.json`

Dive into the [Guide](./guide/) to learn more about the CLI and its features, or explore the [Schema](./schema/) documentation to unlock the full power of Konfigo's processing capabilities.
