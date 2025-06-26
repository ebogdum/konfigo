# Transformers & Data Transformation Tests

This directory contains comprehensive tests for Konfigo's data transformation functionality, covering all transformer types: `renameKey`, `changeCase`, `addKeyPrefix`, and `setValue`.

## Overview

Transformers modify the structure and content of merged configuration data after initial merging and variable resolution. This test suite validates all aspects of transformer functionality across all supported formats.

## Test Structure

```
transformers/
├── input/              # Input configuration files
├── config/             # Schema files with transform definitions
├── output/             # Generated test outputs (cleaned on each run)
├── expected/           # Expected outputs for validation
├── test.sh            # Main test script
├── validate.sh        # Output validation script
└── README.md          # This file
```

## Test Categories

### 1. Rename Key Tests (`renameKey`)
- **Purpose**: Test key renaming functionality
- **Files**: `schema-rename-basic.yaml`, `schema-rename-basic.json`
- **Coverage**: Moving values from one path to another, deleting original paths

### 2. Change Case Tests (`changeCase`)
- **Purpose**: Test string case transformation
- **Files**: `schema-changecase-basic.yaml`, `schema-changecase-all-types.yaml`
- **Coverage**: All case types (upper, lower, snake, camel, kebab)

### 3. Add Key Prefix Tests (`addKeyPrefix`)
- **Purpose**: Test map key prefixing functionality
- **Files**: `schema-addprefix-basic.yaml`, `schema-addprefix-vars.yaml`
- **Coverage**: Static prefixes, variable substitution in prefixes

### 4. Set Value Tests (`setValue`)
- **Purpose**: Test setting values at paths
- **Files**: `schema-setvalue-basic.yaml`, `schema-setvalue-vars.yaml`, `schema-setvalue-complex.yaml`
- **Coverage**: Simple values, variable substitution, complex nested structures

### 5. Combined Transformations Tests
- **Purpose**: Test multiple transformers working together
- **Files**: `schema-combined.yaml`
- **Coverage**: Sequential transformation pipeline, interdependencies

### 6. Error Handling Tests
- **Purpose**: Test proper error handling for invalid configurations
- **Files**: `schema-error-*.yaml`
- **Coverage**: Missing paths, type mismatches, invalid case types

### 7. Environment Variable Tests
- **Purpose**: Test `KONFIGO_VAR_*` environment variable substitution in transformers
- **Coverage**: Variable substitution in transformer definitions

### 8. Edge Case Tests
- **Purpose**: Test special scenarios like stdin input
- **Coverage**: Stdin processing with transformers

## Transformer Features Tested

### renameKey Transformer
- ✅ Moving values between paths
- ✅ Creating intermediate paths as needed
- ✅ Deleting original paths after move
- ✅ Error handling for missing source paths
- ✅ Cross-format compatibility

### changeCase Transformer
- ✅ `upper` case conversion (UPPERCASE)
- ✅ `lower` case conversion (lowercase)
- ✅ `snake` case conversion (snake_case)
- ✅ `camel` case conversion (camelCase)
- ✅ `kebab` case conversion (kebab-case)
- ✅ `pascal` case conversion (PascalCase)
- ✅ Error handling for non-string values
- ✅ Error handling for invalid case types

### addKeyPrefix Transformer
- ✅ Adding prefixes to map keys
- ✅ Variable substitution in prefixes
- ✅ Error handling for non-map values
- ✅ Preserving nested map structures

### setValue Transformer
- ✅ Setting simple values (strings, numbers, booleans)
- ✅ Setting complex nested structures (maps, arrays)
- ✅ Variable substitution in string values
- ✅ Creating paths that don't exist
- ✅ Overwriting existing values

### Advanced Features
- ✅ Sequential transformation pipelines
- ✅ Variable substitution in transformer definitions
- ✅ Environment variable integration (`KONFIGO_VAR_*`)
- ✅ Cross-format processing (any input → any output)
- ✅ Error propagation and clear error messages

## Running Tests

### Full Test Suite
```bash
./test.sh
```

### Validation Only (after running tests)
```bash
./validate.sh
```

### Individual Test Debugging
```bash
# Test specific transformer with debug output
/Users/bogdan/go/src/konfigo/konfigo -s input/base-config.yaml -S config/schema-rename-basic.yaml -oy -d
```

## Test Results

As of last run:
- **Total Tests**: 33 test scenarios
- **Total Files Generated**: 95 output files (JSON, YAML, TOML for each test)
- **Pass Rate**: 100% (33/33 passed)
- **Coverage**: All transformer types and format combinations

## Example Transformer Configurations

### Basic Rename Key
```yaml
transform:
  - type: "renameKey"
    from: "user.name"
    to: "user.fullName"
```

### Multiple Case Changes
```yaml
transform:
  - type: "changeCase"
    path: "strings.CamelCase"
    case: "upper"
  - type: "changeCase"
    path: "strings.snake_case"
    case: "camel"
```

### Key Prefix with Variables
```yaml
vars:
  - name: "ENV_PREFIX"
    value: "prod"
transform:
  - type: "addKeyPrefix"
    path: "database"
    prefix: "${ENV_PREFIX}_"
```

### Complex Value Setting
```yaml
transform:
  - type: "setValue"
    path: "feature.newFlags"
    value:
      advancedUI: true
      betaAccess: false
      experiments:
        - "feature-a"
        - "feature-b"
```

### Combined Transformation Pipeline
```yaml
vars:
  - name: "ENV_PREFIX"
    value: "prod"
transform:
  - type: "renameKey"
    from: "legacy.api_endpoint"
    to: "service.url"
  - type: "changeCase"
    path: "service.url"
    case: "lower"
  - type: "setValue"
    path: "service.environment"
    value: "${ENV_PREFIX}"
  - type: "addKeyPrefix"
    path: "service"
    prefix: "${ENV_PREFIX}_"
```

## Input Data Structure

The test input files contain a comprehensive structure designed to test all transformer scenarios:

```yaml
user:
  name: "Alice"
  id: 123
  email: "alice@example.com"
apiSettings:
  RequestTimeout: "ThirtySeconds"
  MaxRetries: "Five"
  BaseURL: "HTTP://API.EXAMPLE.COM"
settings:
  timeout: 30
  retries: 3
  debug: true
legacy:
  api_endpoint: "HTTP://OLD-DOMAIN.COM/api"
  auth_token: "legacy-token-123"
strings:
  CamelCase: "testValue"
  snake_case: "another_value"
  UPPER_CASE: "THIRD_VALUE"
  kebab-case: "fourth-value"
  Mixed_Format: "mixed-Format_STRING"
# ... additional test data
```

## Integration

This test suite integrates with:
- **Variable System**: Tests variable substitution in transformer definitions
- **Generator System**: Tests transformation after generation
- **Format Conversion**: Tests cross-format transformer processing
- **Schema Processing**: Tests transformer integration with overall schema pipeline

## Next Steps

This completes the transformer testing. Next features to test:
1. **Validators** (type, range, pattern, enum validation)
2. **Configuration Merging** (multi-source file merging)
3. **Input/Output Schema** (structure validation and filtering)
4. **Batch Processing** (`konfigo_forEach` functionality)
