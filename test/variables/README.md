# Variables Testing

This directory contains comprehensive tests for Konfigo's variable substitution feature.

## Overview

Variables in Konfigo support:
- **Three precedence levels**: `KONFIGO_VAR_*` environment variables (highest) > `-V` file variables > schema `vars` block (lowest)
- **Multiple variable sources in schema**: `value`, `fromEnv`, `fromPath` with optional `defaultValue`
- **Variable substitution**: `${VAR_NAME}` replacement in configuration values
- **Format-agnostic**: Works with all supported input/output formats

## Test Structure

```
variables/
├── input/           # Test input configurations
│   ├── base-config.json
│   ├── base-config.yaml
│   └── base-config.toml
├── config/          # Schema and variable files
│   ├── schema-basic.yaml      # Basic variable schema (YAML)
│   ├── schema-basic.json      # Basic variable schema (JSON)
│   ├── schema-error-test.yaml # Schema for testing error conditions
│   ├── variables-basic.yaml   # Variables file (YAML)
│   └── variables-basic.json   # Variables file (JSON)
├── output/          # Generated test outputs
├── expected/        # Expected reference outputs
├── test.sh          # Main test script
├── validate.sh      # Validation script
└── README.md        # This file
```

## Test Coverage

### 1. Variable Precedence Tests
- **Test 1**: Schema variables only (lowest precedence)
- **Test 2**: Variables file overrides schema variables
- **Test 3**: Environment variables override everything (highest precedence)

### 2. Variable Source Tests
- **Test 4**: `fromEnv` variable resolution from system environment
- Tests for `value`, `fromPath`, and `defaultValue` sources

### 3. Format Compatibility Tests
- **Tests 5-6**: Different input formats (YAML, TOML) with variables
- **Tests 7-9**: Different output formats (JSON, TOML, ENV) with variables
- **Test 10**: JSON schema with JSON variables file

### 4. Integration Tests
- **Test 11**: Variables without schema (basic substitution mode)
- **Test 13**: Integration with `KONFIGO_KEY_` environment variables
- **Test 14**: Complex nested variable substitution

### 5. Error Handling Tests
- **Test 12**: Missing required variables (expected to fail)

## Variable Configuration Examples

### Schema Variables (`vars` block)
```yaml
vars:
  # Literal value
  - name: "API_HOST"
    value: "api.example.com"
  
  # From environment with default
  - name: "API_PORT"
    fromEnv: "SERVICE_PORT"
    defaultValue: "8080"
    
  # From config path
  - name: "TARGET_NAMESPACE"
    fromPath: "deployment.namespace"
    
  # With default fallback
  - name: "TIMEOUT"
    defaultValue: "30s"
```

### Variables File (`-V` flag)
```yaml
# Variables that override schema defaults
API_HOST: "vars-file-api.example.com"
NESTED_VAR: "from-vars-file"
CUSTOM_VAR: "only-in-vars"
DATABASE_PASSWORD: "vars-file-password"
```

### Environment Variables (highest precedence)
```bash
export KONFIGO_VAR_API_HOST="env-override.example.com"
export KONFIGO_VAR_NESTED_VAR="from-environment"
export DB_PASS="secret123"
export SERVICE_PORT="9090"
```

## Variable Substitution in Configuration

Input configuration with variable placeholders:
```yaml
api:
  baseUrl: "${API_HOST}:${API_PORT}"
  timeout: "${TIMEOUT}"
deployment:
  namespace: "production"
  replicas: "${REPLICA_COUNT}"
```

After variable substitution:
```yaml
api:
  baseUrl: "env-override.example.com:9090"
  timeout: "30s"
deployment:
  namespace: "production"
  replicas: "2"  # from defaultValue in schema
```

## Running Tests

### Run All Tests
```bash
./test.sh
```

### Validate Against Expected Results
```bash
./validate.sh
```

### Run Individual Test
```bash
# Example: Test environment variable override
export KONFIGO_VAR_API_HOST="custom.example.com"
../../konfigo -s input/base-config.json -S config/schema-basic.yaml -V config/variables-basic.yaml -oy
```

## Test Results

All tests verify:
1. ✅ Correct variable precedence (env > file > schema)
2. ✅ Proper variable substitution in config values
3. ✅ Support for all input/output format combinations
4. ✅ Integration with `KONFIGO_KEY_` environment variables
5. ✅ Error handling for missing required variables
6. ✅ Variable resolution from multiple sources (`fromEnv`, `fromPath`, `value`, `defaultValue`)

## Notes

- Variable names are case-sensitive
- Unresolved variables are left as-is (e.g., `${UNDEFINED_VAR}`)
- Variables work in both basic mode (without schema) and schema mode
- `KONFIGO_KEY_` environment variables can inject configuration values directly
- Schema `config:` sections are not merged into configuration (schemas define processing operations, not config structure)
