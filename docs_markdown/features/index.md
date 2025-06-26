# Features

Welcome to the Features section! Here you'll find detailed documentation on Konfigo's core capabilities, each with practical examples and references to the test suite. Use this section to understand what Konfigo can do and how to leverage its advanced features in your own workflows.

## At a Glance
- Learn about supported file formats and conversion
- Understand merging strategies and environment integration
- Discover batch processing and recursive file discovery
- See how each feature is tested and validated

## Feature Documentation

### Input and Processing
- **[Format Conversion](./format-conversion.md)** - File format support and conversion between JSON, YAML, TOML, ENV, and INI
- **[Merging](./merging.md)** - Configuration merging strategies and precedence rules
- **[Environment Integration](./env-integration.md)** - Environment variable handling and `KONFIGO_KEY_` prefixes
- **[Recursive Discovery](./recursive-discovery.md)** - Automatic file discovery in directory trees

### Advanced Processing
- **[Batch Processing](./batch-processing.md)** - Multi-output generation using `konfigo_forEach`

## Test Coverage

Each feature is extensively tested in the `test/` directory:

```
test/
├── format-conversion/     # Format conversion tests
├── merging/              # Configuration merging tests  
├── env-integration/      # Environment variable tests
├── recursive-discovery/  # File discovery tests
├── batch/               # Batch processing tests
```

Most feature tests follow this structure:
- `input/` - Test input files in various formats
- `config/` - Schema and configuration files
- `expected/` - Expected output files
- `output/` - Generated output files
- `test.sh` - Test execution script
- `validate.sh` - Output validation script

See individual feature documentation for specific examples and usage patterns.

---

## See Also
- [Core Concepts](../concepts/)
- [User Guide](../guide/)
