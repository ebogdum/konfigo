---
layout: home

hero:
  name: "Konfigo"
  text: "Smart Configuration Management"
  tagline: "Merge, transform, and validate configurations across any format. One tool, infinite possibilities."
  actions:
    - theme: brand
      text: Get Started in 5 Minutes
      link: /getting-started/quick-start
    - theme: alt
      text: View Examples
      link: /guide/recipes
    - theme: alt
      text: Explore Schemas
      link: /schema/

features:
  - icon: 🔄
    title: "Universal Format Support"
    details: "Seamlessly work with JSON, YAML, TOML, and .env files. Convert between formats effortlessly."
  - icon: 🧩
    title: "Intelligent Merging"
    details: "Combine configurations from multiple sources with smart conflict resolution and precedence rules."
  - icon: ⚡
    title: "Schema-Powered Processing"
    details: "Define variables, validate data, transform structures, and generate multiple outputs with flexible schemas."
  - icon: 🌍
    title: "Environment Integration"
    details: "Override any configuration value with environment variables. Perfect for CI/CD and containerized deployments."
  - icon: 🚀
    title: "Batch Generation"
    details: "Generate multiple tailored configuration files from a single source using powerful iteration features."
  - icon: 🛠️
    title: "Developer Friendly"
    details: "Rich CLI, comprehensive validation, detailed error messages, and extensive debugging capabilities."

---

## What is Konfigo?

Konfigo is a powerful command-line tool that solves the complexity of modern configuration management. Whether you're dealing with microservices, multi-environment deployments, or complex application settings, Konfigo provides a unified way to merge, validate, and transform your configurations.

## I want to...

<div class="quick-paths">

### 🚀 **Get started quickly**
→ [5-minute Quick Start](/getting-started/quick-start)  
Perfect for newcomers who want to see Konfigo in action immediately.

### 🔧 **Solve a specific problem**
→ [Common Tasks & Recipes](/guide/)  
Jump straight to solutions for merging, converting, or validating configurations.

### 📚 **Learn the concepts**
→ [Understanding Konfigo](/getting-started/concepts)  
Build a solid foundation before diving into advanced features.

### ⚡ **Master advanced features**
→ [Schema Guide](/schema/)  
Unlock the full power of schema-driven configuration processing.

</div>

## Real-World Examples

**Multi-Environment Deployment**
```bash
# Merge base config with environment-specific overrides
konfigo -s base.yaml,prod.yaml -of config.json
```

**Configuration Validation**
```bash
# Validate against schema and generate multiple outputs
konfigo -s config.yaml -S schema.yaml -V variables.yaml
```

**Format Conversion**
```bash
# Convert legacy .env files to modern YAML
konfigo -s legacy.env -oy -of modern.yaml
```

---

## Why Choose Konfigo?

- **🎯 Purpose-Built**: Designed specifically for configuration management challenges
- **🔒 Reliable**: Extensive testing with real-world configuration scenarios  
- **📖 Well-Documented**: Comprehensive guides and examples for every feature
- **🌟 Active Development**: Regular updates and community-driven improvements

Ready to simplify your configuration management? [Get started now](/getting-started/)!

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

---

## See Also
- [Core Concepts](./getting-started/concepts)
- [User Guide](./guide/)
- [Schema Guide](./schema/)
- [Reference](./reference/)
