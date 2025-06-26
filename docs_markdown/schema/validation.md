# Schema: Data Validation

Konfigo's `validate` directives enforce rules and constraints on processed configuration data, ensuring integrity before final output. Validation executes after variables, generators, and transformations are complete.

## Validation Structure

```yaml
validate:
  - path: "config.path.to.validate"
    rules:
      required: true
      type: "string"
      minLength: 3
      regex: "^[a-zA-Z][a-zA-Z0-9_]*$"
```

### Fields

- **`path`** (Required): Dot-separated path to configuration value
- **`rules`** (Required): Object containing validation constraints

## Validation Rules

### Core Rules

**`required` (boolean, default: false)**
- When `true`, value must exist at path
- When `false` or missing, other rules skipped if path not found

**`type` (string)**
- Enforces specific data type
- Supported types: `"string"`, `"number"`, `"integer"`, `"boolean"`, `"slice"`, `"map"`

### String Rules

**`minLength` (integer)**
- Minimum string length

**`regex` (string)** 
- ECMA 262 JavaScript-style regular expression pattern

**`enum` (array of strings)**
- Value must match one of the provided options

### Numeric Rules

**`min` (number)**
- Minimum value (inclusive) for numbers and integers

**`max` (number)**
- Maximum value (inclusive) for numbers and integers

## Examples from Tests

### Basic Configuration Validation

**Input Configuration:**
```yaml
service:
  name: "user-service"
  port: 8080
  environment: "production"
database:
  host: "db-server"
  port: 5432
```

**Validation Schema:**
```yaml
validate:
  - path: "service.name"
    rules:
      required: true
      type: "string"
      minLength: 3
      regex: "^[a-zA-Z][a-zA-Z0-9_-]*$"
  
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
  
  - path: "service.environment"
    rules:
      type: "string"
      enum: ["development", "staging", "production"]
```

### Complex Nested Validation

**Schema with Multiple Constraints:**
```yaml
validate:
  # Database credentials validation
  - path: "database.credentials.username"
    rules:
      required: true
      type: "string"
      minLength: 3
      regex: "^[a-zA-Z][a-zA-Z0-9_]*$"
  
  - path: "database.credentials.password"
    rules:
      required: true
      type: "string"
      minLength: 8
      regex: "^[a-zA-Z0-9_@$!%*?&]{8,}$"
  
  # Array validation
  - path: "features.features"
    rules:
      type: "slice"
  
  # Multiple constraints on same field
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
  
  # Floating point validation
  - path: "timeouts.connect"
    rules:
      type: "number"
      min: 0.1
      max: 60.0
```

### Optional Field Validation

**Configuration:**
```yaml
api:
  endpoint: "https://api.example.com"
  timeout: 30
  # apiKey is optional
cache:
  enabled: true
  ttl: 300
```

**Schema:**
```yaml
validate:
  # Required fields
  - path: "api.endpoint"
    rules:
      required: true
      type: "string"
      regex: "^https?://"
  
  - path: "api.timeout"
    rules:
      required: true
      type: "number"
      min: 1
      max: 300
  
  # Optional field with validation when present
  - path: "api.apiKey"
    rules:
      required: false  # Optional
      type: "string"
      minLength: 32
      regex: "^[a-f0-9]{32}$"  # 32-char hex
  
  - path: "cache.enabled"
    rules:
      type: "boolean"
  
  - path: "cache.ttl"
    rules:
      type: "number"
      min: 60
      max: 3600
```

### Environment-Specific Validation

**Development vs Production Constraints:**
```yaml
# Different validation based on environment
vars:
  - name: "ENV"
    fromEnv: "ENVIRONMENT"
    defaultValue: "development"

validate:
  - path: "database.ssl"
    rules:
      required: true
      type: "boolean"
  
  # Stricter password rules in production
  - path: "database.password"
    rules:
      required: true
      type: "string"
      minLength: 12  # Longer in production
      regex: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{12,}$"
  
  # Debug mode validation
  - path: "app.debug"
    rules:
      type: "boolean"
```

## Validation Processing

### Order of Operations

1. **Path Resolution**: Check if configuration value exists at path
2. **Required Check**: If `required: true` and path missing, fail immediately
3. **Skip Missing**: If path missing and not required, skip other rules
4. **Type Validation**: Check data type if specified
5. **Constraint Validation**: Apply min/max, minLength, enum, regex rules

### Example Processing Flow

```yaml
# Configuration:
service:
  port: "8080"  # String instead of number

# Validation:
- path: "service.port"
  rules:
    required: true
    type: "number"
    min: 1024

# Processing:
# 1. Path "service.port" exists ✓
# 2. required: true and value exists ✓
# 3. type: "number" but value is string ✗ FAIL
# 4. min: 1024 not checked (type failed)
```

## Error Scenarios from Tests

### Type Mismatches

```yaml
# ERROR: Expected number, got string
service:
  port: "eight-thousand"

validate:
  - path: "service.port"
    rules:
      type: "number"
```

### Range Violations

```yaml
# ERROR: Value 99999 exceeds max 65535
service:
  port: 99999

validate:
  - path: "service.port"
    rules:
      type: "number"
      max: 65535
```

### Pattern Mismatches

```yaml
# ERROR: Invalid username format
user:
  name: "123invalid"

validate:
  - path: "user.name"
    rules:
      regex: "^[a-zA-Z][a-zA-Z0-9_]*$"
```

### Enum Violations

```yaml
# ERROR: "test" not in allowed values
environment: "test"

validate:
  - path: "environment"
    rules:
      enum: ["development", "staging", "production"]
```

## Advanced Validation Patterns

### Multi-Field Dependencies

While Konfigo doesn't support cross-field validation directly, you can use generators and transformations to prepare validation:

```yaml
# Generate computed field for validation
generators:
  - type: "concat"
    targetPath: "validation.dbConnection"
    format: "{host}:{port}"
    sources:
      host: "database.host"
      port: "database.port"

validate:
  - path: "validation.dbConnection"
    rules:
      regex: "^[a-zA-Z0-9.-]+:[0-9]+$"
```

### Conditional Validation with Variables

```yaml
vars:
  - name: "REQUIRE_SSL"
    fromEnv: "PROD_MODE"
    defaultValue: "false"

# Set required SSL field based on environment
transform:
  - type: "setValue"
    path: "database.sslRequired"
    value: "${REQUIRE_SSL}"

validate:
  - path: "database.sslRequired"
    rules:
      type: "boolean"
  
  # Only validate SSL cert if required
  - path: "database.sslCert"
    rules:
      required: false  # Handle conditionally in app logic
      type: "string"
```

## Best Practices

1. **Start Simple**: Begin with basic required/type validation
2. **Gradual Enhancement**: Add constraints incrementally
3. **Test Edge Cases**: Validate with invalid data to verify error handling
4. **Document Constraints**: Explain validation rules in configuration templates
5. **Environment Awareness**: Consider different validation needs per environment
6. **Error Messages**: Validation errors include path and rule details
    *   **Example**: `min: 0`

*   `max` (number):
    *   For values of `type: "number"` or `"integer"`. The value must be less than or equal to `max`.
    *   **Example**: `max: 100`

*   `minLength` (integer):
    *   For values of `type: "string"`. The string's length must be greater than or equal to `minLength`.
    *   **Example**: `minLength: 3`

*   `enum` (list of strings):
    *   For values of `type: "string"`. The string value must be one of the values present in the `enum` list.
    *   **Example**: `enum: ["production", "staging", "development"]`

*   `regex` (string):
    *   For values of `type: "string"`. The string value must match the provided ECMA 262 (JavaScript-style) regular expression.
    *   **Example**: `regex: "^\\d{3}-\\d{2}-\\d{4}$"` (for a US SSN format)

## How Validation Works

1.  Konfigo iterates through each validation group defined in the `validate` list.
2.  For each group, it attempts to retrieve the value from the configuration at the specified `path`.
3.  **Existence Check**:
    *   If the `rules.required` is `true` and the value is not found, validation fails immediately for this group.
    *   If the value is not found and `rules.required` is `false` (or not set), this validation group is skipped, and Konfigo moves to the next one.
4.  **Rule Application**: If the value is found, Konfigo applies all specified rules in the `rules` object to it.
    *   **Type Check**: If `type` is specified, it's checked first. If it fails, an error is reported.
    *   **Other Rules**: Subsequent rules (`min`, `max`, `minLength`, `enum`, `regex`) are checked. These rules generally assume the type check (if specified) has passed or that the value is of a compatible type (e.g., `min` expects a number).
5.  If any rule fails, Konfigo stops processing and reports a validation error, typically indicating the path, the problematic value, and the rule that failed.
6.  If all validation groups pass, the configuration is considered valid according to the schema.

## Examples

### Example 1: Basic Service Configuration Validation

```yaml
# schema.yml
validate:
  - path: "service.name"
    rules:
      required: true
      type: "string"
      minLength: 3
  - path: "service.port"
    rules:
      required: true
      type: "integer"
      min: 1024
      max: 65535
  - path: "service.environment"
    rules:
      type: "string"
      enum: ["dev", "staging", "prod"]
  - path: "service.apiKey"
    rules:
      required: false # Optional API key
      type: "string"
      regex: "^[a-f0-9]{32}$" # Example: 32-char hex string
```

**Valid Configuration (`config.yml`):**
```yaml
service:
  name: "user-auth"
  port: 8080
  environment: "prod"
  apiKey: "abcdef0123456789abcdef0123456789"
```

**Invalid Configuration (and why):**
```yaml
service:
  name: "db" # Fails service.name minLength: 3
  # port is missing - Fails service.port required: true
  environment: "testing" # Fails service.environment enum
```

### Example 2: Validating Nested Structures and Optional Fields

```yaml
# schema.yml
validate:
  - path: "database.host"
    rules:
      required: true
      type: "string"
  - path: "database.port"
    rules:
      required: true
      type: "integer"
      min: 1
      max: 65535
  - path: "database.credentials.username"
    rules:
      required: true
      type: "string"
  - path: "database.credentials.password" # Password is required
    rules:
      required: true
      type: "string"
      minLength: 8
  - path: "featureFlags.betaEnabled" # Optional boolean
    rules:
      type: "boolean"
  - path: "timeouts.read"
    rules:
      type: "number"
      min: 0.5 # e.g., 0.5 seconds
      max: 60
```

**Valid Configuration (`config.json`):**
```json
{
  "database": {
    "host": "db.example.com",
    "port": 5432,
    "credentials": {
      "username": "admin",
      "password": "complex_password123"
    }
  },
  "featureFlags": {
    "betaEnabled": true
  },
  "timeouts": {
    "read": 30
  }
}
```

If `timeouts.read` was `"fast"`, it would fail the `type: "number"` check. If `database.port` was `99999`, it would fail the `max: 65535` check.

By defining comprehensive validation rules, you can significantly increase the reliability and robustness of your application's configuration.
