# Core Concepts

Understanding the core concepts of Konfigo is essential for mastering configuration management with this tool. This section explains how Konfigo processes configurations, the principles behind its design, and the strategies it uses to ensure flexibility and reliability.

## At a Glance
- Learn about configuration sources and how they are discovered
- Understand the processing pipeline and variable resolution
- Explore merging strategies, format detection, and error handling

## Configuration Sources

Konfigo can read from multiple sources with automatic format detection:

- **Files**: JSON, YAML, TOML, ENV, INI formats
- **Directories**: Recursive discovery of configuration files
- **Stdin**: Piped input with format specification
- **Environment Variables**: Direct key-value overrides

## Processing Pipeline

Konfigo follows a structured processing pipeline:

1. **Discovery**: Find all configuration files from sources
2. **Parsing**: Parse each file according to its format
3. **Merging**: Merge configurations with precedence rules
4. **Environment**: Apply environment variable overrides
5. **Schema Processing**: Execute schema-driven operations
6. **Output**: Generate final outputs in requested formats

Each stage is optimized for performance and memory efficiency, with parallel processing capabilities for large configurations.

## Variable Resolution Precedence

Variables are resolved with the following priority (highest to lowest):

1. **Environment Variables** (`KONFIGO_VAR_*`)
2. **Variables File** (`-V` flag)
3. **Schema Variables** (`vars` section)
4. **Path References** (`fromPath` in variables)

This hierarchical approach ensures that runtime values can override defaults while maintaining predictable behavior.

## Configuration Merging Strategy

Konfigo uses intelligent deep merging with these principles:

- **Last wins**: Later sources override earlier ones
- **Deep merging**: Nested objects are merged recursively
- **Array handling**: Arrays can be replaced or merged based on configuration
- **Type safety**: Type conflicts are handled gracefully
- **Immutable paths**: Protected paths cannot be overridden (when specified in schema)

## Format Detection

Konfigo automatically detects configuration formats based on:

1. **File extension**: `.json`, `.yaml`, `.yml`, `.toml`, `.env`, `.ini`
2. **Content analysis**: Fallback parsing when extension is ambiguous
3. **Explicit flags**: Override detection with format-specific flags
4. **MIME headers**: For streamed content

## Error Handling

Konfigo provides comprehensive error reporting:

- **Parse errors**: Clear indication of syntax issues with line numbers
- **Schema validation**: Detailed validation failure messages
- **Missing files**: Helpful suggestions for file path issues
- **Type mismatches**: Clear explanation of data type conflicts

Understanding these core concepts will help you leverage Konfigo's full potential for managing complex configuration scenarios.

---

## See Also
- [Features](../features/)
- [User Guide](../guide/)
