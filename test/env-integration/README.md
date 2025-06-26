# Environment Variable Integration Test Suite

This test suite validates Konfigo's environment variable integration features, specifically `KONFIGO_KEY_` for direct configuration overrides and `KONFIGO_VAR_` for variable substitution.

## Features Tested

### 1. KONFIGO_KEY_ Direct Configuration Overrides

#### Basic Overrides
- **Test**: `basic-overrides.yaml`
- **Features**: Simple key-value overrides using `KONFIGO_KEY_`
- **Environment Variables**:
  - `KONFIGO_KEY_app.name=env-override-app`
  - `KONFIGO_KEY_app.port=9090`
  - `KONFIGO_KEY_database.ssl=true`
- **Result**: All specified paths are overridden with string values

#### Nested Path Overrides
- **Test**: `nested-overrides.yaml`
- **Features**: Deep nested path overrides and array index access
- **Environment Variables**:
  - `KONFIGO_KEY_database.connection.timeout=60`
  - `KONFIGO_KEY_nested.deep.very.deep.value=env-modified`
  - `KONFIGO_KEY_logging.outputs.0=syslog`
- **Result**: Deep paths work correctly, array indices create object keys

#### Multiple Input Files
- **Test**: `multi-file-overrides.yaml`
- **Features**: Environment overrides with multiple input sources
- **Input**: base-config.json + override-config.yaml + additional-config.toml
- **Result**: Environment variables override all input file values

#### New Key Creation
- **Test**: `new-keys.yaml`
- **Features**: Creating entirely new configuration paths
- **Environment Variables**:
  - `KONFIGO_KEY_runtime.environment=staging`
  - `KONFIGO_KEY_new_section.enabled=true`
- **Result**: New nested structures are created successfully

#### Array Index Behavior
- **Test**: `array-overrides.yaml`
- **Features**: Overriding array elements by index
- **Environment Variables**:
  - `KONFIGO_KEY_logging.outputs.0=stderr`
  - `KONFIGO_KEY_logging.outputs.1=journald`
- **Result**: Arrays become objects with string keys ("0", "1", etc.)

#### Precedence Testing
- **Test**: `precedence-test.yaml`
- **Features**: Environment variables vs file values precedence
- **Result**: Environment variables always win over file values

### 2. KONFIGO_VAR_ Variable Integration

#### Variable Substitution
- **Test**: `var-integration.yaml`
- **Features**: External variable injection via environment
- **Environment Variables**:
  - `KONFIGO_VAR_KONFIGO_SERVICE_NAME=test-service`
  - `KONFIGO_VAR_KONFIGO_ENVIRONMENT=staging`
- **Result**: Variables are available in schema processing and concatenation

#### Combined Usage
- **Test**: `combined-env-vars.yaml`
- **Features**: Both `KONFIGO_KEY_` and `KONFIGO_VAR_` together
- **Result**: Both override types work simultaneously

### 3. Schema Integration

#### Basic Schema Validation
- **Test**: `env-with-schema.yaml`
- **Features**: Environment overrides with schema validation
- **Schema**: `env-friendly-schema.yaml` (no strict type validation)
- **Result**: Overrides work when not conflicting with type validation

#### Immutable Paths
- **Test**: `immutable-override.yaml`
- **Features**: Environment variables vs immutable path protection
- **Environment Variables**:
  - `KONFIGO_KEY_app.name=should-override-immutable` (blocked)
  - `KONFIGO_KEY_database.host=should-override-immutable` (blocked)
  - `KONFIGO_KEY_app.version=4.0.0` (allowed)
- **Result**: Immutable paths are protected, other paths can be overridden

### 4. Format Support

#### Multiple Output Formats
- **Tests**: `env-override.json`, `env-override.toml`, `env-override.env`
- **Features**: Environment overrides with different output formats
- **Result**: All formats (JSON, YAML, TOML, ENV) work correctly

### 5. Type Conversion and Edge Cases

#### Type Conversion
- **Test**: `type-conversion.yaml`
- **Features**: Environment string values in typed contexts
- **Environment Variables**:
  - `KONFIGO_KEY_database.ssl=false` (boolean as string)
  - `KONFIGO_KEY_features.cache=true` (boolean as string)
- **Result**: Values remain as strings when set via environment variables

#### Complex Key Names
- **Test**: `complex-keys.yaml`
- **Features**: Special characters in key names
- **Environment Variables**:
  - `KONFIGO_KEY_complex-key-name.with_underscore.and-dash=complex-value`
  - `KONFIGO_KEY_123numeric.start=numeric-key`
- **Result**: Hyphens, underscores, and numeric prefixes work correctly

## Important Behaviors and Limitations

### Type System Limitations
1. **Environment Variables Are Always Strings**: All `KONFIGO_KEY_` values are treated as strings, regardless of the original type in the configuration
2. **Schema Type Validation Conflicts**: If a schema expects an integer but an environment variable provides a string, validation will fail
3. **Boolean Conversion**: Environment boolean values become strings ("true", "false")

### Array Handling
1. **Array Index Override**: Using `KONFIGO_KEY_path.array.0=value` converts the array to an object with string keys
2. **No Array Append**: Cannot append new elements to arrays via environment variables
3. **Array Structure Change**: Original array structure is lost when any index is overridden

### Path Resolution
1. **Dot Notation**: Use dots to separate nested path segments
2. **New Path Creation**: Environment variables can create entirely new configuration branches
3. **Deep Nesting**: Arbitrarily deep nesting is supported

### Precedence Rules
1. **Environment > Files**: `KONFIGO_KEY_` variables always override file values
2. **Environment > Immutable**: Environment variables cannot override immutable paths
3. **Last Wins**: Among multiple `KONFIGO_KEY_` variables, the last one set wins

## Test Results

- **Total Tests**: 15
- **Passing Tests**: 15 (100%)
- **Failed Tests**: 0
- **Output Files**: All formats validated (JSON, YAML, TOML, ENV)

## Known Issues and Workarounds

### 1. Type Validation Conflicts
**Issue**: Environment variables are strings but schemas may expect other types
**Workaround**: Use schemas without strict type validation for paths that may be overridden by environment variables

### 2. Array Index Behavior
**Issue**: Array indices become object keys instead of modifying array elements
**Impact**: Changes array structure, may break consumers expecting arrays
**Workaround**: Avoid using environment variables to override array elements

### 3. Complex Type Values
**Issue**: Cannot set complex objects or arrays via environment variables
**Limitation**: Only scalar string values are supported
**Workaround**: Use file-based overrides for complex structures

## Usage Examples

### Basic Override
```bash
KONFIGO_KEY_app.port=8080 konfigo -s config.json -of output.yaml
```

### Nested Override
```bash
KONFIGO_KEY_database.connection.timeout=60 konfigo -s config.json -of output.yaml
```

### Variable Injection
```bash
KONFIGO_VAR_SERVICE_NAME=my-service konfigo -s config.json -S schema.yaml -of output.yaml
```

### Combined Usage
```bash
KONFIGO_KEY_app.debug=true \
KONFIGO_VAR_ENVIRONMENT=production \
konfigo -s config.json -S schema.yaml -of output.yaml
```

This test suite provides comprehensive coverage of Konfigo's environment variable integration capabilities and documents all important behaviors and limitations.
