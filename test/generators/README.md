# Generators & Data Generation Tests

This directory contains comprehensive tests for Konfigo's data generation functionality, specifically the `concat` generator.

## Overview

The `concat` generator creates new configuration values by concatenating source values from the configuration with static text and variables. This test suite validates all aspects of generator functionality across all supported formats.

## Test Structure

```
generators/
├── input/              # Input configuration files
├── config/             # Schema and variable files
├── output/             # Generated test outputs (cleaned on each run)
├── expected/           # Expected outputs for validation
├── test.sh            # Main test script
├── validate.sh        # Output validation script
└── README.md          # This file
```

## Test Categories

### 1. Basic Generator Tests
- **Purpose**: Test simple concat generator functionality
- **Files**: `schema-basic.json`, `schema-basic.yaml`
- **Coverage**: Single generator with placeholders from config paths

### 2. Multiple Generators Tests  
- **Purpose**: Test multiple generators in sequence
- **Files**: `schema-multiple.yaml`
- **Coverage**: Multiple generators creating different values, variable substitution

### 3. Mixed Content Tests
- **Purpose**: Test combination of placeholders, variables, and static text
- **Files**: `schema-mixed.yaml`
- **Coverage**: Deep nested paths, number/boolean formatting, variable mixing

### 4. Variables Only Tests
- **Purpose**: Test generators that only use variables (no config placeholders)
- **Files**: `schema-vars-only.yaml`
- **Coverage**: Pure variable substitution in generators

### 5. External Variables Tests
- **Purpose**: Test variable precedence and external variable files
- **Files**: `schema-cascading.yaml`, `variables.yaml`, `variables.json`
- **Coverage**: Variable file loading, precedence validation

### 6. Cascading Generators Tests
- **Purpose**: Test complex scenarios with variables and multiple generators
- **Files**: `schema-cascading.yaml`
- **Coverage**: Integration between generators and variable systems

### 7. Error Handling Tests
- **Purpose**: Test proper error handling for invalid configurations
- **Files**: `schema-error-missing-path.yaml`
- **Coverage**: Missing source paths, proper error reporting

### 8. Environment Variable Tests
- **Purpose**: Test `KONFIGO_VAR_*` environment variable integration
- **Coverage**: Environment variable precedence over file variables

### 9. Edge Case Tests
- **Purpose**: Test special scenarios like stdin input
- **Coverage**: Stdin processing with generators

## Generator Features Tested

### Core Functionality
- ✅ `{placeholder}` substitution from config paths
- ✅ `${VARIABLE}` substitution from variables
- ✅ Static text concatenation
- ✅ Multiple generators in sequence
- ✅ Deep nested path references
- ✅ Variable precedence (env > file > schema)

### Data Type Handling
- ✅ String values
- ✅ Number formatting (integers and floats)
- ✅ Boolean formatting
- ✅ Nested object access

### Error Conditions
- ✅ Missing source paths (properly fails)
- ⚠️ Empty target paths (validation gap identified)
- ⚠️ Empty format strings (validation gap identified)
- ⚠️ No sources defined (validation gap identified)

### Format Compatibility
- ✅ JSON input/output
- ✅ YAML input/output  
- ✅ TOML input/output
- ✅ Cross-format processing (any input → any output)

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
# Test specific scenario with debug output
/Users/bogdan/go/src/konfigo/konfigo -s input/base-config.yaml -S config/schema-basic.yaml -oy -d
```

## Test Results

As of last run:
- **Total Tests**: 17 test scenarios
- **Total Files Generated**: 44 output files (JSON, YAML, TOML for each test)
- **Pass Rate**: 100% (17/17 passed)
- **Coverage**: All major generator features and format combinations

## Example Generator Configuration

### Basic Concat Generator
```yaml
generators:
  - type: "concat"
    targetPath: "service.identifier"
    format: "Service: {name} (ID: {id}) running in {region}"
    sources:
      name: "service.name"
      id: "service.instanceId"
      region: "region"
```

### With Variables
```yaml
vars:
  - name: "APP_VERSION"
    value: "1.2.3"
    
generators:
  - type: "concat"
    targetPath: "service.fullName"
    format: "{name}-${APP_VERSION}"
    sources:
      name: "service.name"
```

### Multiple Generators
```yaml
generators:
  - type: "concat"
    targetPath: "database.connectionString"
    format: "postgresql://{host}:{port}/{db}"
    sources:
      host: "database.host"
      port: "database.port"
      db: "database.name"
  - type: "concat"
    targetPath: "service.url"
    format: "https://{service}.${DOMAIN}:{port}"
    sources:
      service: "service.name"
      port: "service.port"
```

## Known Issues

Several validation gaps were identified and documented in `/test/inspect/`:
1. Empty target path validation not enforced
2. Empty format string validation not enforced  
3. No sources validation not enforced

These are functional but represent areas for future improvement.

## Integration

This test suite integrates with:
- **Variable System**: Tests variable precedence and substitution
- **Format Conversion**: Tests cross-format generator processing
- **Schema Processing**: Tests generator integration with overall schema pipeline

## Next Steps

This completes the generator testing. Next features to test:
1. **Transformers** (`renameKey`, `changeCase`, `addKeyPrefix`, `setValue`)
2. **Validators** (type, range, pattern, enum validation)
3. **Configuration Merging** (multi-source file merging)
4. **Batch Processing** (`konfigo_forEach` functionality)
