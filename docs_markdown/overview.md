# Konfigo Documentation

> **Complete documentation for Konfigo - a powerful configuration management tool with advanced merging, transformation, validation, and batch processing capabilities.**

## Quick Start

```bash
# Basic configuration merging
konfigo -s base.yaml,prod.yaml -oy

# With schema processing and variables
konfigo -s config.yaml -S schema.yaml -V variables.yaml -of output.json

# Batch processing for multiple environments
konfigo -s template.yaml -S schema.yaml -V batch-vars.yaml
```

## Documentation Structure

### 📖 [User Guide](./guide/)
Comprehensive guide covering all Konfigo features, from basic usage to advanced configurations. **Start here for a complete overview.**

### 🛠️ Core Features

#### **[Configuration Merging](./features/merging.md)**
- Deep merging with precedence rules
- Array and object handling strategies
- Immutable path protection
- Real-world merging examples

#### **[Format Conversion](./features/format-conversion.md)**
- Support for JSON, YAML, TOML, ENV, and INI formats
- Automatic format detection
- Multiple output format generation
- Format-specific best practices

#### **[Environment Integration](./features/env-integration.md)**
- Direct configuration key overrides (`KONFIGO_KEY_*`)
- Schema variable overrides (`KONFIGO_VAR_*`)
- CI/CD pipeline integration patterns
- Environment-specific configuration management

#### **[Recursive Discovery](./features/recursive-discovery.md)**
- Automatic configuration file discovery
- Directory structure best practices
- Performance optimization for large codebases
- File filtering and organization strategies

#### **[Batch Processing](./features/batch-processing.md)**
- Multi-environment configuration generation
- `konfigo_forEach` directive usage
- Template-based deployment configurations
- Advanced batch processing patterns

### 🔧 Schema System

#### **[Schema Overview](./schema/index.md)**
Complete introduction to Konfigo's schema system with processing pipeline overview.

#### **[Variables & Substitution](./schema/variables.md)**
- Variable definition and resolution
- Multiple source types (`value`, `fromEnv`, `fromPath`)
- Variable precedence and inheritance
- Dynamic variable patterns

#### **[Data Generation](./schema/generation.md)**
- `concat` generator for dynamic value creation
- Configuration-based data composition
- Variable integration in generators
- Complex generation patterns

#### **[Data Transformation](./schema/transformation.md)**
- Key renaming and restructuring (`renameKey`)
- String case transformation (`changeCase`)
- Key prefix addition (`addKeyPrefix`)
- Value setting and updates (`setValue`)

#### **[Data Validation](./schema/validation.md)**
- Type validation and constraints
- String patterns and enumerations
- Numeric ranges and bounds
- Complex validation scenarios

#### **[Advanced Features](./schema/advanced.md)**
- Immutable field protection
- Input schema validation
- Output schema filtering
- Strict vs. flexible validation modes

### 📚 Feature Index

#### **[Features Overview](./features/)**
Index of all Konfigo features with test coverage information and implementation patterns.

## Key Concepts

### Configuration Processing Pipeline

```
Input Sources → Discovery → Parsing → Merging → Schema Processing → Output
```

1. **Discovery**: Find configuration files automatically or from specified sources
2. **Parsing**: Parse multiple formats (JSON, YAML, TOML, ENV, INI)
3. **Merging**: Intelligent deep merging with precedence rules
4. **Schema Processing**: Apply variables, generators, transformations, and validation
5. **Output**: Generate clean, validated configuration in desired formats

### Variable Resolution Hierarchy

```
Environment Variables (KONFIGO_VAR_*) [Highest]
        ↓
Variables File (-V flag)
        ↓
Schema Variables (vars block)
        ↓
Configuration Paths (fromPath) [Lowest]
```

### Schema Processing Order

```
Input Validation → Variables → Generators → Transformations → Variable Substitution → Validation → Output Filtering
```

## Common Use Cases

### 1. Multi-Environment Configuration
```bash
# Development
konfigo -s base.yaml,environments/dev.yaml

# Production with secrets
export KONFIGO_VAR_DB_PASSWORD=secret
konfigo -s base.yaml,environments/prod.yaml -S schema.yaml
```

### 2. Kubernetes Deployment Generation
```bash
# Generate multiple deployment manifests
konfigo -s k8s-template.yaml -S k8s-schema.yaml -V services.yaml
```

### 3. Configuration Migration
```bash
# Transform legacy configurations to new format
konfigo -s legacy-config.json -S migration-schema.yaml -of modern-config.yaml
```

### 4. CI/CD Integration
```bash
# Automated configuration processing in pipelines
export KONFIGO_VAR_VERSION=${CI_COMMIT_TAG}
export KONFIGO_VAR_ENVIRONMENT=${CI_ENVIRONMENT_NAME}
konfigo -s base-config.yaml -S deployment-schema.yaml -of deployment.yaml
```

## Examples from Tests

All documentation includes real examples from Konfigo's comprehensive test suite:

- **[Test-Based Generation Examples](./schema/generation.md#examples-from-tests)**
- **[Test-Based Transformation Examples](./schema/transformation.md#examples-from-tests)**
- **[Test-Based Validation Examples](./schema/validation.md#examples-from-tests)**
- **[Test-Based Merging Examples](./features/merging.md#real-world-examples)**
- **[Test-Based Batch Processing Examples](./features/batch-processing.md#examples-from-tests)**

## Quick Reference

### Essential Commands
```bash
# Basic merging
konfigo -s config1.yaml,config2.yaml

# With schema processing
konfigo -s config.yaml -S schema.yaml

# Multiple outputs
konfigo -s config.yaml -oj -oy -of config.toml

# Environment override
KONFIGO_KEY_app.debug=true konfigo -s config.yaml

# Variable override
KONFIGO_VAR_ENVIRONMENT=prod konfigo -s config.yaml -S schema.yaml

# Batch processing
konfigo -s template.yaml -S schema.yaml -V batch-vars.yaml

# Debug mode
konfigo -d -s config.yaml -S schema.yaml
```

### Common Schema Patterns
```yaml
# Variables with fallbacks
vars:
  - name: "ENVIRONMENT"
    fromEnv: "NODE_ENV"
    defaultValue: "development"

# Generate connection strings
generators:
  - type: "concat"
    targetPath: "database.url"
    format: "postgresql://{user}:${PASSWORD}@{host}:5432/{db}"
    sources:
      user: "database.username"
      host: "database.host"
      db: "database.name"

# Validate configuration
validate:
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
```

## Contributing to Documentation

When adding new features or tests:

1. **Update Feature Documentation**: Add examples to relevant feature docs
2. **Include Test Examples**: Reference actual test cases in examples
3. **Update Schema Docs**: Document new schema capabilities
4. **Add Use Cases**: Include real-world usage patterns
5. **Test Documentation**: Verify all examples work as documented

## Getting Help

- **[User Guide](./guide/)**: Comprehensive coverage of all features
- **[Feature-Specific Docs](./features/)**: Detailed feature documentation
- **[Schema Reference](./schema/)**: Complete schema system documentation
- **Built-in Help**: `konfigo -h` for command-line help
- **Debug Mode**: `konfigo -d` for detailed processing information

---

This documentation is organized to support both learning Konfigo from scratch and serving as a comprehensive reference for advanced users. Each section includes practical examples derived from Konfigo's extensive test suite to ensure accuracy and real-world applicability.
- **[Merging](./features/merging.md)** - Configuration merging and precedence rules
- **[Environment Integration](./features/env-integration.md)** - Environment variable handling
- **[Recursive Discovery](./features/recursive-discovery.md)** - File discovery patterns

### Schema Features
- **[Schema Overview](./schema/index.md)** - Schema structure and concepts
- **[Variables](./schema/variables.md)** - Variable definition and substitution
- **[Generators](./schema/generation.md)** - Data generation capabilities
- **[Transformers](./schema/transformation.md)** - Data transformation operations
- **[Validators](./schema/validation.md)** - Configuration validation rules

### Advanced Features
- **[Batch Processing](./features/batch-processing.md)** - Multi-output generation with `konfigo_forEach`
- **[Advanced Schema Features](./advanced/schema-features)** - Advanced schema patterns
- **[Use Cases](./guide/use-cases.md)** - Real-world examples and patterns

## Quick Examples

### Basic Merging
```bash
# Merge configuration files
konfigo -s base.json,environment.yml,local.toml

# Read from stdin
cat config.yml | konfigo -sy
```

### Schema Processing
```bash
# Process with schema
konfigo -s config.yml -S schema.yml

# Include variables file
konfigo -s config.yml -S schema.yml -V variables.yml
```

### Batch Processing
```bash
# Generate multiple outputs
konfigo -s base.yml -S deployment-schema.yml -V batch-vars.yml
```

### Format Conversion
```bash
# Convert YAML to JSON
konfigo -s config.yml -oj

# Output to file
konfigo -s config.yml -of output.json
```

## Test Examples

Throughout this documentation, examples are drawn from Konfigo's comprehensive test suite located in the `test/` directory. Each feature area has corresponding test cases that demonstrate real-world usage patterns.

## Getting Help

Use `konfigo -h` for detailed command-line help, or refer to the specific feature documentation for in-depth explanations and examples.
