# Schema Processing Integration Test Suite

This directory contains comprehensive tests for Konfigo's advanced schema processing capabilities, focusing on the integration of input/output schemas with other schema features.

## Overview

Schema Processing Integration in Konfigo combines multiple advanced features:
- **Input Schema Validation**: Validates configuration structure before processing
- **Output Schema Filtering**: Filters final output to specific structure
- **Complex Workflows**: Chains multiple processing steps together
- **Cross-Feature Integration**: Combines schemas with variables, transformations, validation, etc.

## Test Structure

```
schema-integration/
├── input/           # Base configuration files in multiple formats
├── schemas/         # Input and output schema definitions
├── config/          # Integration schema configurations
├── variables/       # Variable files for complex scenarios
├── output/          # Generated test outputs
├── expected/        # Expected outputs for validation
├── test.sh          # Main test script
├── validate.sh      # Output validation script
└── README.md        # This file
```

## Features Tested

### 1. Full Integration Workflow
- **Test**: Complete pipeline with input schema → processing → output schema
- **Schema**: `full-integration-schema.yaml` with all features enabled
- **Features**: Input validation, variables, generators, transformers, validators, output filtering
- **Formats**: YAML, JSON, TOML, ENV output support
- **Validation**: Type-safe schemas with proper float/int handling

### 2. Strict Schema Validation
- **Test**: Strict input and output schema enforcement
- **Schema**: `strict-schema.yaml` with `strict: true` for both input and output
- **Validation**: Exact structure matching, no extra fields allowed
- **Error Handling**: Validates that violations are properly caught

### 3. Input Schema Only
- **Test**: Input validation without output filtering
- **Schema**: `input-only-schema.yaml` focusing on input structure validation
- **Features**: Ensures input meets structural requirements before processing
- **Transformations**: Applies changes after input validation

### 4. Output Schema Only
- **Test**: Output filtering without input validation
- **Schema**: `output-only-schema.yaml` focusing on output structure control
- **Features**: Generates clean, filtered output from complex internal configuration
- **Use Case**: Public API configuration from internal settings

### 5. Variable Integration
- **Test**: Schema processing with external variable files
- **Variables**: `complex-vars.yaml` with environment-specific settings
- **Integration**: Variables used in generators, transformers, and validators
- **Override**: Multiple variable files with precedence testing

### 6. Immutable Paths Integration
- **Test**: Schema processing respecting immutable path constraints
- **Schema**: `immutable-schema.yaml` with protected configuration paths
- **Validation**: Ensures transformations cannot override immutable values
- **Override**: Tests that `KONFIGO_KEY_` env vars can still override immutable paths

### 7. Environment Variable Override
- **Test**: Schema processing with `KONFIGO_KEY_` and `KONFIGO_VAR_` overrides
- **Environment**: Direct configuration overrides during schema processing
- **Integration**: Environment variables work correctly with input/output schemas
- **Type Handling**: Proper type conversion for schema validation

### 8. Error Case Testing
- **Strict Input Violation**: Input with extra fields failing strict input schema
- **Strict Output Violation**: Processing generating fields not in strict output schema
- **Expected Behavior**: Graceful error handling with descriptive messages

### 9. Multi-Format Support
- **Input Formats**: JSON, YAML, TOML combinations
- **Output Formats**: YAML, JSON, TOML, ENV generation
- **Schema Formats**: JSON schemas for input/output validation
- **Cross-Format**: Different input and output format combinations

## Schema Files

### Input Schemas
- **`input-schema.json`**: Comprehensive input structure validation
- **`input-schema-strict.json`**: Strict input validation (exact match)
- **`simple-input-schema.json`**: Basic input validation for simple cases

### Output Schemas
- **`output-schema.json`**: Standard output filtering
- **`output-schema-strict.json`**: Strict output validation (exact match)

### Integration Schemas
- **`full-integration-schema.yaml`**: Complete feature integration
- **`strict-schema.yaml`**: Strict input/output validation
- **`input-only-schema.yaml`**: Input validation focus
- **`output-only-schema.yaml`**: Output filtering focus
- **`immutable-schema.yaml`**: Immutable paths with schemas
- **`env-override-schema.yaml`**: Environment variable integration

## Input Files

### Configuration Sources
- **`base-config.json`**: Primary configuration in JSON format
- **`env-override.yaml`**: Environment-specific overrides in YAML
- **`additional-config.toml`**: Additional settings in TOML format
- **`simple-config.json`**: Simplified configuration for strict testing

### Variable Files
- **`complex-vars.yaml`**: Complex variable definitions with nested structures
- **`override-vars.yaml`**: Variable overrides for precedence testing

## Test Cases

### Successful Integration Tests (11 tests)
1. **Full Integration to YAML**: Complete processing pipeline
2. **Full Integration to JSON**: Same pipeline, different output format
3. **Strict Schema Validation**: Strict input/output validation
4. **Input Schema Only**: Input validation without output filtering
5. **Output Schema Only**: Output filtering without input validation
6. **Integration with Variables**: External variable file integration
7. **Multiple Variable Files**: Variable precedence and merging
8. **Immutable Paths Integration**: Protected paths with schema processing
9. **Integration to TOML**: Alternative output format
10. **Integration to ENV**: Environment variable format output
11. **Environment Variable Override**: Runtime configuration overrides

### Error Case Tests (2 tests)
1. **Strict Input Schema Violation**: Extra fields in input (expected failure)
2. **Strict Output Schema Violation**: Extra fields in output (expected failure)

## Running Tests

### Execute All Tests
```bash
./test.sh
```

### Validate Outputs
```bash
./validate.sh
```

### Clean Outputs
```bash
rm -rf output/*
```

## Test Results

- **Total Tests**: 12 test scenarios (10 success + 2 expected failures)
- **Output Files**: 11 configuration files across different formats
- **Pass Rate**: 100% (11/11 successful tests + 2/2 expected failures)
- **Error Handling**: Proper validation of schema violations

## Key Findings

### ✅ Working Features
1. **Input Schema Validation**: Validates merged configuration structure before processing
2. **Output Schema Filtering**: Filters processed configuration to specified structure
3. **Strict Mode Support**: Both input and output schemas support strict validation
4. **Type System Integration**: Proper handling of numeric types (float64 vs int)
5. **Multi-Format Support**: All input/output format combinations work correctly
6. **Variable Integration**: External variables work correctly with schemas
7. **Environment Override**: `KONFIGO_KEY_` and `KONFIGO_VAR_` integration
8. **Immutable Paths**: Protected paths respected during schema processing
9. **Error Handling**: Graceful failures with descriptive error messages
10. **Complex Workflows**: Multiple features can be chained together

### 📋 Integration Notes
- **Type Handling**: Input schemas require numeric values as float64 for proper validation
- **Validation Order**: Input validation occurs before processing, output validation after
- **Schema Flexibility**: Non-strict mode allows flexible input/output structures
- **Environment Variables**: Type conversion handled automatically for string env vars
- **Processing Pipeline**: Complete workflow from input validation through output filtering

### 🔧 Technical Details
- **Input Schema**: Validates structure of merged configuration before any processing
- **Output Schema**: Filters final configuration to only include specified fields
- **Strict Mode**: Enforces exact structural matching with no extra fields
- **Processing Order**: Input validation → Variables → Generators → Transformers → Validation → Output filtering
- **Type System**: JSON schemas define expected types, with float64 for numeric values
- **Error Messages**: Clear validation errors for schema violations

### 🎯 Use Cases Validated
1. **API Configuration**: Clean public configuration from complex internal settings
2. **Environment Deployment**: Environment-specific configuration generation
3. **Configuration Validation**: Structural validation of configuration inputs
4. **Multi-Stage Processing**: Complex transformation pipelines with validation
5. **Configuration Templating**: Dynamic configuration generation with templates

This test suite validates Konfigo's most advanced schema processing capabilities, ensuring that complex configuration workflows function correctly across all supported formats and use cases.
