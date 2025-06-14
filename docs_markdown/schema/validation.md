# Schema: Data Validation

Konfigo's `validate` directives allow you to enforce rules and constraints on your processed configuration data, ensuring its integrity and correctness before final output. Validation occurs after variables are substituted, generators are run, and transformations are applied.

Validation rules are defined within the `validate` key in your schema file. This key holds a list of validation groups, where each group targets a specific path in the configuration and applies a set of rules to the value at that path.

## Structure of a Validation Group

```yaml
validate:
  - path: "path.to.value.to.validate"  # Dot-separated path to the configuration key.
    rules:                             # Object containing the validation rules for this path.
      required: true
      type: "string"                 # Expected data type.
      minLength: 5
      enum: ["active", "inactive", "pending"]
      regex: "^[a-zA-Z0-9_-]+$"
      # For numbers:
      # type: "number" # or "integer"
      # min: 0
      # max: 100
```

### Fields:

*   `path` (Required, string): A dot-separated path to the value in the configuration that this set of rules should validate.
*   `rules` (Required, object): An object containing one or more validation rules.

## Validation Rules

All rules are optional within the `rules` object. If a value at the specified `path` is not found, only the `required` rule is checked. If `required` is `false` or not set, and the path is not found, other validation rules for that path are skipped.

*   `required` (boolean, default: `false`):
    *   If `true`, the value at `path` must exist in the configuration. If it's missing, validation fails.
    *   **Example**: `required: true`

*   `type` (string):
    *   Specifies the expected data type of the value.
    *   Supported types:
        *   `"string"`
        *   `"number"` (matches floating-point or integer numbers; JSON numbers are typically `float64`)
        *   `"integer"` (specifically checks if a number is a whole number, e.g., `10.0` is a valid integer, but `10.5` is not)
        *   `"boolean"`
        *   `"array"` (Note: Konfigo currently uses `reflect.TypeOf(val).Kind().String()` which might return `slice` for arrays/lists from JSON/YAML)
        *   `"map"` (Note: Konfigo currently uses `reflect.TypeOf(val).Kind().String()` which might return `map` for objects from JSON/YAML)
    *   If the actual type does not match the expected type, validation fails.
    *   **Example**: `type: "number"`

*   `min` (number):
    *   For values of `type: "number"` or `"integer"`. The value must be greater than or equal to `min`.
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
