# Konfigo Comprehensive Testing Plan

## Overview
This document provides a structured approach to test all Konfigo features systematically. Each feature will have its own test directory with inputs, configs, outputs, and validation.

## Features Analysis

### 1. Core Functionality Features
- **Format Conversion**: JSON ↔ YAML ↔ TOML ↔ ENV ↔ INI
- **Configuration Merging**: Multiple source files with precedence
- **File Output**: Writing to files with format detection
- **Stdin Input**: Reading from stdin with format specification
- **Recursive Discovery**: Finding config files in directories

### 2. Schema-Driven Features
- **Variables**: Variable substitution with multiple sources
- **Generators**: Dynamic value generation (concat)
- **Transformers**: Data transformation (renameKey, changeCase, addKeyPrefix, setValue)
- **Validators**: Data validation (type, range, pattern, enum)
- **Input Schema**: Input structure validation
- **Output Schema**: Output filtering and structure enforcement
- **Immutable Paths**: Protected configuration paths

### 3. Advanced Features
- **Batch Processing**: konfigo_forEach for multiple outputs
- **Environment Variables**: KONFIGO_KEY_ and KONFIGO_VAR_ integration
- **Case Sensitivity**: Case-sensitive vs case-insensitive merging
- **Error Handling**: Graceful error reporting and validation failures

## Supported Formats

### Input Formats
- **JSON** (.json) - Full support for parsing and merging
- **YAML** (.yaml, .yml) - Full support for parsing and merging
- **TOML** (.toml) - Full support for parsing and merging
- **ENV** (.env) - Key-value pairs for environment variables
- **INI** (.ini) - Section-based configuration files

### Output Formats
- **JSON** - Structured data output
- **YAML** - Human-readable structured output
- **TOML** - Configuration-friendly structured output
- **ENV** - Environment variable format

### Schema/Config Formats
- **JSON** - For schema definitions and variable files
- **YAML** - For schema definitions and variable files
- **TOML** - For schema definitions and variable files

## Test Directory Structure

```
test/
├── testing.md                           # This file
├── inspect/                             # Failed tests requiring code fixes
├── format-conversion/                   # Basic format conversion tests
│   ├── input/
│   ├── config/
│   ├── output/
│   ├── expected/
│   └── test.sh
├── config-merging/                      # Multi-source merging tests
├── variables/                           # Variable substitution tests
├── generators/                          # Value generation tests
├── transformers/                        # Data transformation tests
├── validators/                          # Data validation tests
├── input-schema/                        # Input validation tests
├── output-schema/                       # Output filtering tests
├── immutable-paths/                     # Immutable path protection tests
├── batch-processing/                    # konfigo_forEach tests
├── environment-variables/               # ENV var integration tests
├── case-sensitivity/                    # Case handling tests
├── error-handling/                      # Error and failure tests
├── recursive-discovery/                 # Directory traversal tests
├── stdin-input/                         # Stdin processing tests
├── file-output/                         # File writing tests
└── complex-scenarios/                   # Multi-feature integration tests
```

## Test Implementation Steps

### Step 1: Create Base Test Structure
- Create main test directories
- Create subdirectories for each feature
- Create standardized folder structure in each test

### Step 2: Generate Input Files
For each test, create input files in all supported formats:
- `input.json` - JSON format input
- `input.yaml` - YAML format input  
- `input.toml` - TOML format input
- `input.env` - Environment format input
- `input.ini` - INI format input (where applicable)

### Step 3: Generate Config Files
Create configuration files for each test:
- `schema.json` - JSON schema configuration
- `schema.yaml` - YAML schema configuration
- `schema.toml` - TOML schema configuration
- `variables.json` - Variable definitions (where applicable)
- `variables.yaml` - Variable definitions (where applicable)

### Step 4: Generate Expected Outputs
Create expected output files for validation:
- `expected.json` - Expected JSON output
- `expected.yaml` - Expected YAML output
- `expected.toml` - Expected TOML output
- `expected.env` - Expected ENV output

### Step 5: Run Tests and Generate Outputs
Execute konfigo commands to generate actual outputs:
- `output.json` - Actual JSON output
- `output.yaml` - Actual YAML output
- `output.toml` - Actual TOML output
- `output.env` - Actual ENV output

### Step 6: Validation
Compare expected vs actual outputs:
- Run diff commands between expected/ and output/ directories
- Document any discrepancies
- Move failing tests to inspect/ directory

### Step 7: Create Test Scripts
Each test directory will have:
- `test.sh` - Automated test execution script
- `README.md` - Test description and expected behavior
- `validate.sh` - Output validation script

## Feature Testing Matrix

| Feature | JSON | YAML | TOML | ENV | INI | Schema | Variables |
|---------|------|------|------|-----|-----|--------|-----------|
| Format Conversion | ✓ | ✓ | ✓ | ✓ | ✓ | - | - |
| Config Merging | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| Variables | ✓ | ✓ | ✓ | ✓ | - | ✓ | ✓ |
| Generators | ✓ | ✓ | ✓ | ✓ | - | ✓ | ✓ |
| Transformers | ✓ | ✓ | ✓ | ✓ | - | ✓ | - |
| Validators | ✓ | ✓ | ✓ | ✓ | - | ✓ | - |
| Input Schema | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | - |
| Output Schema | ✓ | ✓ | ✓ | ✓ | - | ✓ | - |
| Batch Processing | ✓ | ✓ | ✓ | ✓ | - | ✓ | ✓ |

## Test Categories

### 1. Basic Format Tests
- **Purpose**: Validate format conversion capabilities
- **Input**: Simple configuration in each format
- **Expected**: Accurate conversion to all other formats
- **Validation**: Content preservation and format compliance

### 2. Schema Feature Tests  
- **Purpose**: Validate schema-driven processing
- **Input**: Configuration + schema + variables (where applicable)
- **Expected**: Processed configuration with applied transformations
- **Validation**: Correct application of schema directives

### 3. Integration Tests
- **Purpose**: Validate multiple features working together
- **Input**: Complex scenarios with multiple features
- **Expected**: Correctly processed multi-stage output
- **Validation**: End-to-end functionality verification

### 4. Error Handling Tests
- **Purpose**: Validate graceful error handling
- **Input**: Invalid configurations, malformed schemas
- **Expected**: Clear error messages and proper exit codes
- **Validation**: Error message clarity and appropriate failure modes

## Test Execution Strategy

### Phase 1: Individual Feature Testing
1. Test each feature in isolation
2. Validate basic functionality works correctly
3. Move failing tests to inspect/ for debugging

### Phase 2: Format Compatibility Testing
1. Test each feature across all supported formats
2. Ensure format-specific behaviors are correct
3. Validate cross-format compatibility

### Phase 3: Integration Testing
1. Test multiple features working together
2. Validate complex real-world scenarios
3. Ensure no feature interactions cause failures

### Phase 4: Error and Edge Case Testing
1. Test error conditions and edge cases
2. Validate graceful degradation
3. Ensure clear error reporting

## Success Criteria

### Individual Tests
- ✅ All format conversions preserve data accurately
- ✅ Schema features apply transformations correctly
- ✅ Error conditions are handled gracefully
- ✅ Output matches expected results exactly

### Overall Suite
- ✅ 100% of features tested across all applicable formats
- ✅ All tests pass or are documented in inspect/ with reasons
- ✅ Real-world scenarios work end-to-end
- ✅ Error handling is comprehensive and user-friendly

## Testing Progress

### ✅ Completed Features

#### 1. Format Conversion (test/format-conversion/)
- **Status**: ✅ COMPLETED - All tests passing
- **Coverage**: All format pairs (JSON ↔ YAML ↔ TOML ↔ ENV ↔ INI)
- **Files**: 25 input files, 400 conversion tests, 100% pass rate
- **Validation**: Automated with `validate.sh`

#### 2. Variables & Substitution (test/variables/)
- **Status**: ✅ COMPLETED - All tests passing  
- **Coverage**: Variable precedence, substitution, all formats, error handling
- **Tests**: 14 comprehensive test scenarios
- **Features Tested**:
  - Variable precedence (`KONFIGO_VAR_*` > `-V` file > schema `vars`)
  - Variable sources (`value`, `fromEnv`, `fromPath`, `defaultValue`)
  - `${VAR_NAME}` substitution in configuration values
  - Integration with `KONFIGO_KEY_*` environment variables
  - Error handling for missing variables
  - All input/output format combinations
- **Validation**: Automated with `validate.sh`

### 🔄 In Progress Features

*None currently*

### 📋 Pending Features

#### 3. Generators & Data Generation (test/generators/)
- **Status**: ✅ COMPLETED - All tests passing
- **Coverage**: Concat generator functionality, all formats, complex scenarios
- **Tests**: 17 comprehensive test scenarios covering:
  - Basic concat generation with config path placeholders
  - Multiple generators in sequence  
  - Mixed content (placeholders + variables + static text)
  - Variables-only generators (no config placeholders)
  - External variable files and precedence
  - Cascading generators with complex variable interactions
  - Error handling for missing source paths
  - Environment variable overrides (`KONFIGO_VAR_*`)
  - Edge cases (stdin input)
  - Cross-format processing (any input → any output)
- **Validation**: Automated with `validate.sh` - 44 output files, 100% match rate
- **Issues Found**: 3 validation gaps identified and documented in `inspect/`:
  - Empty target path validation not enforced
  - Empty format string validation not enforced
  - No sources validation not enforced (for static-only generators)

### 🔄 In Progress Features

*None currently*

### 📋 Pending Features

#### 4. Transformers & Data Transformation (test/transformers/)
- **Status**: ✅ COMPLETED - All tests passing
- **Coverage**: All transformer types, all formats, complex scenarios
- **Tests**: 33 comprehensive test scenarios covering:
  - **renameKey**: Moving values between paths, deleting originals
  - **changeCase**: All case types (upper, lower, snake, camel, kebab, pascal)
  - **addKeyPrefix**: Map key prefixing with variable substitution
  - **setValue**: Simple and complex value setting with variable substitution
  - Combined transformation pipelines with interdependencies
  - Error handling for missing paths, type mismatches, invalid parameters
  - Environment variable substitution (`KONFIGO_VAR_*`)
  - Edge cases (stdin input, cross-format processing)
- **Validation**: Automated with `validate.sh` - 95 output files, 100% match rate
- **Features Validated**: All 4 transformer types working correctly across all formats

### 🔄 In Progress Features

*None currently*

### 📋 Pending Features

#### 5. Validators & Data Validation (test/validators/)
- **Status**: ✅ COMPLETED - 19/21 tests passing, 2 known issues documented
- **Coverage**: All validation rules, comprehensive error testing, cross-format processing
- **Tests**: 21 comprehensive test scenarios covering:
  - **Core Validation Rules**: required, type, min/max, minLength, enum, regex
  - **String Validation**: All string constraints across all formats
  - **Numeric Validation**: min/max constraints (JSON inputs only due to type limitations)
  - **Boolean Validation**: Type checking across all formats
  - **Error Handling**: Missing required fields, type mismatches, constraint violations
  - **Cross-Format Processing**: JSON input → all output formats
  - **Edge Cases**: Boundary values, regex patterns, enum validation
  - **Format-Specific Issues**: YAML/TOML type system limitations
- **Validation**: Automated with `validate.sh` - 24 output files, 100% match rate
- **Issues Found**: 2 type system limitations identified and documented in `inspect/`:
  - YAML/TOML integer validation (int vs float64 type conflicts)
  - Numeric validation only accepts float64 types

#### 6. Configuration Merging (test/merging/)
- **Status**: ✅ COMPLETED - 27/27 core tests passing, 3 known issues documented
- **Coverage**: Complete merging functionality validation across all supported scenarios
- **Tests**: 27 comprehensive test scenarios covering:
  - **Basic Merge Precedence**: Multi-source merging with proper order-based precedence
  - **Cross-Format Merging**: All format combinations (JSON ↔ YAML ↔ TOML ↔ ENV ↔ INI)
  - **Case Sensitivity**: Case-sensitive (`-c`) vs case-insensitive (default) merging
  - **Immutable Paths**: Schema-defined paths that resist file-based overrides
  - **Environment Variable Overrides**: `KONFIGO_KEY_*` direct configuration injection
  - **Recursive Discovery**: `-r` flag for finding and merging files in subdirectories
  - **Stdin Integration**: All input formats via stdin with proper format flags
  - **Output Format Control**: Any input combination to any output format
  - **Edge Cases**: Empty files, complex nested structures, error handling
- **Validation**: Automated with `validate.sh` - 31 output files, 100% match rate
- **Issues Found**: 1 environment variable vs immutable path behavior documented in `inspect/`:
  - `KONFIGO_KEY_*` variables not overriding immutable paths as documented
- **Features Validated**: 
  - Multi-source precedence rules working correctly
  - All case sensitivity modes working correctly  
  - Immutable path protection working correctly for file-based sources
  - Recursive file discovery working correctly
  - Cross-format merging working correctly across all format combinations

### 🔄 In Progress Features

*None currently*

### 📋 Pending Features

#### 7. Batch Processing (`konfigo_forEach`)
- **Status**: ✅ COMPLETED - 4/4 batch categories passing, 11 output files generated
- **Coverage**: Complete batch processing functionality validation across all supported scenarios
- **Tests**: 10 comprehensive test scenarios covering:
  - **Items-Based Batching**: Inline array processing for service configurations
  - **ItemFiles-Based Batching**: External file references for environment configurations
  - **Multi-Format Output**: YAML, JSON, TOML, ENV output support
  - **Schema Integration**: Batch processing with schema validation
  - **Nested Output Paths**: Complex directory structures and path templating
  - **Multiple Variable Sources**: Combining multiple batch variable files
  - **Complex Multi-Level**: Nested batch processing scenarios
  - **Error Handling**: Missing itemFiles graceful failure
  - **Cross-Format Processing**: All input/output format combinations
  - **Variable Substitution**: Template variables in output paths
- **Output Structure**: Directory-based organization with multiple files per batch
- **Validation**: Automated with `validate.sh` - 11 output files, 100% match rate
- **Features Validated**: All batch processing patterns working correctly

#### 8. Schema Processing Integration
- **Status**: ✅ COMPLETED - 11/11 integration tests passing, 2/2 error cases handled correctly
- **Coverage**: Complete schema processing integration validation across all advanced workflows
- **Tests**: 12 comprehensive test scenarios covering:
  - **Full Integration Workflow**: Input schema → processing → output schema complete pipeline
  - **Strict Schema Validation**: Strict input/output validation with exact structure matching
  - **Input Schema Only**: Input validation without output filtering
  - **Output Schema Only**: Output filtering without input validation  
  - **Variable Integration**: External variable files with schema processing
  - **Multiple Variable Files**: Variable precedence and merging with schemas
  - **Immutable Paths Integration**: Protected paths respected during schema processing
  - **Environment Variable Override**: `KONFIGO_KEY_` and `KONFIGO_VAR_` integration with schemas
  - **Multi-Format Support**: All input/output format combinations (JSON, YAML, TOML, ENV)
  - **Error Handling**: Graceful failures for strict schema violations
  - **Complex Workflows**: Chaining multiple processing steps with validation
  - **Type System Integration**: Proper handling of numeric types in schemas
- **Output Structure**: 11 validated configuration files across multiple formats
- **Validation**: Automated with `validate.sh` - 11 output files, 100% match rate
- **Features Validated**: All advanced schema processing patterns working correctly
- **Error Cases**: 2 expected failures properly caught and handled

#### 9. Environment Variable Integration (`KONFIGO_KEY_` and `KONFIGO_VAR_`)
- **Priority**: Medium (critical for deployment scenarios)
- **Features**: Direct config overrides, nested path handling, variable integration
- **Status**: ✅ **COMPLETED**
- **Test Coverage**: 
  - **Basic Overrides**: Simple key-value environment variable overrides
  - **Nested Paths**: Deep nested path modifications (e.g., `database.connection.timeout`)
  - **Multiple Input Files**: Environment overrides with multiple configuration sources
  - **Variable Integration**: `KONFIGO_VAR_` for external variable injection
  - **Combined Usage**: Both `KONFIGO_KEY_` and `KONFIGO_VAR_` simultaneously
  - **Schema Validation**: Environment overrides with schema constraints
  - **Immutable Paths**: Environment variables vs immutable path protection
  - **New Key Creation**: Creating new configuration branches via environment
  - **Array Index Access**: Overriding array elements (converts to object keys)
  - **Type Conversion**: String conversion behavior for all environment values
  - **Output Formats**: JSON, YAML, TOML, ENV format outputs with environment overrides
  - **Precedence Testing**: Environment variables vs file values precedence
  - **Complex Key Names**: Special characters, hyphens, underscores, numeric prefixes
  - **Edge Cases**: Complex key patterns and error handling
- **Output Structure**: 15 validated configuration files across multiple scenarios
- **Validation**: Automated with `validate.sh` - 15 output files, 100% match rate
- **Features Validated**: All environment variable integration patterns working correctly
- **Known Issues**: Type conversion limitations documented in `inspect/env-type-conversion-issues.md`

#### 10. Recursive Discovery
- **Priority**: High (core functionality for directory-based configurations)
- **Features**: Automatic discovery of configuration files in directory trees
- **Status**: ✅ **COMPLETED**
- **Test Coverage**: 
  - **Basic Discovery**: Automatic finding and merging of all config files in directory tree
  - **File Type Recognition**: Supports JSON, YAML, TOML, ENV formats, ignores non-config files
  - **Directory Traversal**: Deep recursive traversal with 10 files across 6 directories tested
  - **Schema Integration**: Recursive discovery with schema validation and transformations
  - **Multiple Output Formats**: JSON, YAML, TOML, ENV output format support
  - **Selective Discovery**: Subdirectory-specific and multi-directory discovery
  - **Non-Recursive Comparison**: Single file and explicit file list processing
  - **Environment Integration**: KONFIGO_KEY_ overrides with recursive discovery
  - **Case Sensitivity**: Case-sensitive vs case-insensitive key matching
  - **Debug Information**: Detailed logging of discovery and processing pipeline
  - **Merge Behavior**: Key conflict resolution and nested object merging
  - **Error Handling**: Graceful handling of non-config files and directory structure
- **Output Structure**: 13 validated files including debug logs
- **Validation**: Automated with `validate.sh` - 13 output files, 100% match rate
- **Features Validated**: Complete recursive discovery pipeline working correctly
- **Performance**: Successfully processed 10 config files across complex directory structure

#### 11. Advanced Features (Case Sensitivity & Error Reporting)
- **Priority**: Low
- **Features**: Case sensitivity options, error reporting
- **Status**: Not started

## Implementation Checklist

- [ ] Create directory structure
- [ ] Generate input files for all formats
- [ ] Create schema configurations
- [ ] Generate expected outputs
- [ ] Implement test scripts
- [ ] Run tests and collect outputs
- [ ] Validate results and document failures
- [ ] Move failing tests to inspect/
- [ ] Create summary documentation
- [ ] Provide recommendations for code fixes

## Additional Steps Needed

1. **Performance Testing**: Add performance benchmarks for large files
2. **Security Testing**: Test with malicious input files
3. **Memory Testing**: Validate memory usage with large configurations
4. **Concurrency Testing**: Test parallel processing capabilities
5. **Regression Testing**: Ensure new features don't break existing functionality
6. **Documentation Testing**: Validate all examples in documentation work
7. **CLI Testing**: Test all command-line flag combinations
8. **Environment Testing**: Test in different operating systems
9. **Version Testing**: Test with different Go versions
10. **Package Testing**: Test installation methods and distribution

This plan provides a comprehensive approach to testing all Konfigo functionality systematically and thoroughly.
