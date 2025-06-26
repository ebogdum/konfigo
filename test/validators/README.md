# Validators & Data Validation Test Suite

This test suite comprehensively tests Konfigo's data validation capabilities across all supported input/output formats.

## Features Tested

### ✅ Core Validation Rules
- **Required Fields**: `required: true/false`
- **Type Validation**: `type: "string|bool|integer|float64"`
- **String Constraints**: `minLength`, `enum`, `regex`  
- **Numeric Constraints**: `min`, `max` (JSON inputs only)
- **Boolean Validation**: Proper type checking

### ✅ Format Coverage
- **Input Formats**: JSON, YAML, TOML
- **Output Formats**: JSON, YAML, TOML  
- **Schema Formats**: YAML, JSON

### ✅ Test Scenarios
- **Valid Configurations**: Proper validation success
- **Invalid Configurations**: Expected validation failures
- **Cross-Format Processing**: Input format A → Output format B
- **Error Handling**: Missing fields, type mismatches, constraint violations
- **Edge Cases**: Boundary values, regex patterns, enum validation

## Test Structure

```
input/               # Test input files
├── base-config.json     # Valid configuration (JSON)
├── base-config.yaml     # Valid configuration (YAML)  
├── base-config.toml     # Valid configuration (TOML)
├── invalid-config.json  # Invalid data for error testing
└── missing-fields.yaml  # Missing required fields

config/              # Validation schemas
├── schema-basic.yaml      # Comprehensive validation rules (JSON-compatible)
├── schema-basic.json      # Same rules in JSON format
├── schema-safe.yaml       # Safe rules (all formats)
├── schema-complex.yaml    # Complex patterns and nested validation
├── schema-optional.yaml   # Optional field testing
├── schema-edge-cases.yaml # Boundary and edge case testing
└── schema-error-tests.yaml # Error scenario testing

output/              # Generated test outputs (24 files)
expected/            # Reference outputs for validation (24 files)
```

## Test Results

### ✅ Passing Tests (19/21)
- **Basic Validation**: JSON input with comprehensive rules ✓
- **Safe Validation**: All formats with compatible rules ✓
- **Cross-Format**: JSON → All output formats ✓
- **Type-Specific**: Individual rule validation ✓
- **Enum Validation**: Valid and invalid enum values ✓
- **Regex Validation**: Pattern matching ✓
- **Boundary Values**: Min/max edge cases ✓
- **Error Scenarios**: Expected validation failures ✓

### ⚠️ Known Issues (2/21 - Documented)
- **YAML/TOML Integer Validation**: Type system limitations
- **Min/Max on Non-Float64**: Numeric validator constraints

## Validation Rules Tested

### String Validation
```yaml
rules:
  required: true
  type: "string"
  minLength: 8
  enum: ["dev", "staging", "prod"]
  regex: "^[a-zA-Z0-9_-]+$"
```

### Numeric Validation (JSON)
```yaml  
rules:
  type: "integer"     # Special case: checks whole numbers
  type: "float64"     # Direct type match
  min: 0
  max: 100
```

### Boolean Validation
```yaml
rules:
  type: "bool"        # Go reflection type name
```

### Required Field Validation
```yaml
rules:
  required: true      # Fails if path not found
  required: false     # Skips validation if path not found
```

## Format-Specific Considerations

### JSON Input (Recommended)
- All numeric values become `float64`
- Compatible with `integer` type (special validation)
- Full min/max validation support
- Complete test coverage

### YAML/TOML Input (Limited)
- Integers become `int`, floats become `float64`
- `integer` type validation expects `float64` (fails)
- `min`/`max` validation only accepts `float64` (fails on `int`)
- String and boolean validation works well

## Usage Examples

### Run Full Test Suite
```bash
./test.sh
```

### Validate Results
```bash
./validate.sh
```

### Test Specific Validation
```bash
../../konfigo -s input/base-config.json -S config/schema-basic.yaml -oj
```

### Test Error Cases
```bash
../../konfigo -s input/invalid-config.json -S config/schema-basic.yaml -oj
# Expected: validation failure
```

## Key Findings

### ✅ What Works Well
1. **String Validation**: minLength, enum, regex work across all formats
2. **Boolean Validation**: Consistent type checking
3. **Required Field Validation**: Reliable across formats
4. **Error Reporting**: Clear, actionable error messages
5. **JSON Processing**: Full validation feature support

### ⚠️ Limitations Discovered
1. **Type System Inconsistency**: Different Go types per input format
2. **Numeric Validation Gaps**: min/max only works with float64
3. **Documentation Mismatch**: Docs show logical names, code uses reflection names

### 📋 Recommendations
1. **Use JSON for Full Validation**: Most reliable format for comprehensive rules
2. **Limit YAML/TOML Validation**: Stick to string/boolean validation  
3. **Normalize Type Handling**: Consider type conversion during parsing
4. **Update Documentation**: Clarify actual type names required

## Files Generated

- **24 Output Files**: Covering all test scenarios and formats
- **100% Validation Match Rate**: All expected outputs verified
- **Comprehensive Coverage**: 19/21 test scenarios passing
- **Clear Issue Documentation**: Known limitations documented in `../inspect/`

This test suite provides robust validation testing within the current system limitations and clearly documents areas for potential improvement.
