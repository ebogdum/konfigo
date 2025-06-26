# Configuration Merging Test Suite

## Overview

This test suite comprehensively validates Konfigo's configuration merging functionality, including multiple source files, precedence rules, case sensitivity, immutable paths, environment variable overrides, and cross-format merging.

## Test Coverage

### ✅ Fully Tested Features

#### 1. Basic Merge Precedence (Tests 1-5)
- **JSON-to-JSON merging**: Later sources override earlier sources
- **Cross-format merging**: JSON + YAML, JSON + YAML + TOML, etc.
- **All format combinations**: JSON, YAML, TOML, ENV files in sequence
- **Order sensitivity**: Different results based on source order

#### 2. Case Sensitivity (Tests 6-7)
- **Case-insensitive merging** (default): `MyApp` and `myapp` treated as same key
- **Case-sensitive merging** (`-c` flag): `MyApp` and `myapp` treated as different keys

#### 3. Immutable Paths (Tests 8-10)
- **Schema-defined immutable paths**: Prevent later sources from overriding protected paths
- **Multiple schema formats**: YAML, JSON, TOML schema support
- **Deep path protection**: Nested paths like `application.name`, `database.port`

#### 4. Environment Variable Overrides (Tests 11-12)
- **Basic KONFIGO_KEY_ overrides**: Direct configuration value injection
- **Nested path overrides**: Complex paths like `database.pool.min`
- **Non-conflicting overrides**: Works when not conflicting with immutable paths

#### 5. Recursive File Discovery (Tests 13-15)
- **Recursive discovery** (`-r` flag): Finds config files in subdirectories
- **Combined with base files**: Merges discovered files with explicit sources
- **With schema processing**: Recursive discovery + immutable path protection

#### 6. Stdin Input Processing (Tests 16-19)
- **All input format support**: JSON (`-sj`), YAML (`-sy`), TOML (`-st`), ENV (`-se`)
- **Combined stdin + file merging**: Merges stdin input with file sources
- **Format specification**: Required format flags for stdin input

#### 7. Output Format Control (Tests 20-22)
- **All output formats**: JSON (`-oj`), YAML (`-oy`), TOML (`-ot`), ENV (`-oe`)
- **Cross-format processing**: Any input format combination to any output format

#### 8. Edge Cases & Error Handling (Tests 23-24)
- **Empty file handling**: Graceful processing of empty JSON objects
- **Complex deep merging**: Nested objects, arrays, multiple merge layers

#### 9. Cross-Format Integration (Tests 25-27)
- **All input formats to JSON**: Demonstrates format-agnostic merging
- **Schema processing with output control**: Immutable paths + format conversion
- **End-to-end processing**: Full pipeline with multiple features

## File Structure

```
test/merging/
├── input/                    # Test input files
│   ├── base-config.*        # Base configuration in all formats
│   ├── override-*.{json,yaml,toml,env}  # Override files for different environments
│   ├── case-*.json          # Case sensitivity test files
│   └── nested/              # Recursive discovery test files
│       ├── env/prod.json
│       └── services/
│           ├── web.yaml
│           └── background.toml
├── config/                   # Schema files
│   └── schema-immutable.*   # Immutable path schemas in all formats
├── output/                   # Test outputs (generated)
├── expected/                 # Expected test outputs
├── inspect/                  # Known issues and gaps
│   └── env-var-immutable-issues.md
├── test-core.sh             # Main test script (working tests only)
├── test.sh                  # Full test script (includes known failures)
├── validate.sh              # Output validation script
└── README.md                # This file
```

## Test Execution

### Run Core Tests (All Passing)
```bash
./test-core.sh
```

### Run Full Test Suite (Includes Known Issues)
```bash
./test.sh
```

### Validate Output Results
```bash
./validate.sh
```

## Test Results

- **Total Tests**: 27 core tests (all passing)
- **Coverage**: All major merging functionality
- **Validation**: 31 output files, 100% match rate with expected results

## Known Issues

### Environment Variables vs Immutable Paths

**Issue**: `KONFIGO_KEY_` environment variables do not override immutable paths as documented.

**Expected**: Per documentation, `KONFIGO_KEY_` variables should override immutable paths.

**Actual**: Immutable paths prevent `KONFIGO_KEY_` overrides.

**Status**: Documented in `inspect/env-var-immutable-issues.md` for investigation.

**Tests Affected**: 3 tests moved to inspect category

## Sample Test Cases

### Basic Merge Precedence
```bash
konfigo -s base-config.json,override-prod.json -oj
# Later source (override-prod.json) overrides base-config.json values
```

### Case Sensitivity
```bash
# Case-insensitive (default)
konfigo -s case-base.json,case-override.json -oj

# Case-sensitive 
konfigo -s case-base.json,case-override.json -c -oj
```

### Immutable Paths
```bash
konfigo -s base-config.json,override-prod.json -S schema-immutable.yaml -oj
# Protected paths in schema prevent overrides
```

### Environment Overrides
```bash
env KONFIGO_KEY_application.environment=production \
konfigo -s base-config.json -oj
# Environment variable directly sets config value
```

### Recursive Discovery
```bash
konfigo -s base-config.json,nested -r -oj
# Discovers and merges all config files in nested/ subdirectories
```

### Cross-Format Processing
```bash
konfigo -s config.json,override.yaml,final.toml -oe
# Merges JSON + YAML + TOML, outputs to ENV format
```

## Validation

All test outputs are validated against expected results:
- **Automated validation**: `validate.sh` compares all outputs
- **Deterministic results**: Tests produce consistent, predictable outputs
- **Format verification**: Outputs are valid in their respective formats

## Integration with Other Features

The merging functionality integrates with:
- **Variables**: Merged configuration used for variable substitution
- **Generators**: Merged data used as source for generators
- **Transformers**: Applied to merged configuration
- **Validators**: Validates final merged result

This test suite ensures the foundation of Konfigo's configuration processing pipeline works correctly across all supported formats and use cases.
