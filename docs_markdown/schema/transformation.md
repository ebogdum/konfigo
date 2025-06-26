# Schema: Data Transformation

Konfigo's `transform` directives modify configuration structure and content after merging but before validation. Transformations process sequentially, with each operation receiving the output from the previous one.

## Processing Order

1. Variable substitution in transform directives
2. Execute transformations in list order
3. Each transformation modifies the configuration state
4. Results passed to validation (if defined)

## Transformation Types

### 1. `renameKey` - Move Configuration Keys

Moves a value from one path to another, creating nested structures as needed.

**Structure:**
```yaml
transform:
  - type: "renameKey"
    from: "old.path.to.key"
    to: "new.path.to.key"
```

**Fields:**
- **`type`** (Required): `"renameKey"`
- **`from`** (Required): Source path (must exist)
- **`to`** (Required): Destination path (created if needed)

**Example from Tests:**
```yaml
# Input config:
legacy:
  api_endpoint: "HTTP://OLD-DOMAIN.COM/api"

# Transform:
transform:
  - type: "renameKey"
    from: "legacy.api_endpoint"
    to: "service.url"

# Result:
service:
  url: "HTTP://OLD-DOMAIN.COM/api"
# legacy key removed
```

### 2. `changeCase` - Modify String Case

Converts string values to different case formats.

**Structure:**
```yaml
transform:
  - type: "changeCase"
    path: "path.to.string.value"
    case: "lower"  # upper, lower, snake, camel
```

**Fields:**
- **`type`** (Required): `"changeCase"`
- **`path`** (Required): Path to string value
- **`case`** (Required): Target case format

**Supported Cases:**
- `"upper"`: UPPERCASE
- `"lower"`: lowercase  
- `"snake"`: snake_case
- `"camel"`: camelCase

**Example from Tests:**
```yaml
# Input config:
service:
  url: "HTTP://OLD-DOMAIN.COM/api"

# Transform:
transform:
  - type: "changeCase"
    path: "service.url"
    case: "lower"

# Result:
service:
  url: "http://old-domain.com/api"
```

### 3. `addKeyPrefix` - Prefix Map Keys

Adds a prefix to all keys within a map object.

**Structure:**
```yaml
transform:
  - type: "addKeyPrefix"
    path: "path.to.map"
    prefix: "prefix_"
```

**Fields:**
- **`type`** (Required): `"addKeyPrefix"`
- **`path`** (Required): Path to map object
- **`prefix`** (Required): String to prepend to keys

**Example from Tests:**
```yaml
# Input config:
service:
  url: "http://old-domain.com/api"
  environment: "prod"

# Transform:
vars:
  - name: "ENV_PREFIX"
    value: "prod"
transform:
  - type: "addKeyPrefix"
    path: "service"
    prefix: "${ENV_PREFIX}_"

# Result:
service:
  prod_url: "http://old-domain.com/api"
  prod_environment: "prod"
```

### 4. `setValue` - Set Configuration Values

Sets any value at a specified path, with variable substitution support.

**Structure:**
```yaml
transform:
  - type: "setValue"
    path: "path.to.key"
    value: "any value type"
```

**Fields:**
- **`type`** (Required): `"setValue"`
- **`path`** (Required): Target path (created if needed)
- **`value`** (Required): Value to set (any type, strings get variable substitution)

**Examples from Tests:**

**String Value with Variables:**
```yaml
vars:
  - name: "ENV_PREFIX"
    value: "prod"
transform:
  - type: "setValue"
    path: "service.environment"
    value: "${ENV_PREFIX}"

# Result:
service:
  environment: "prod"
```

**Complex Object Value:**
```yaml
transform:
  - type: "setValue"
    path: "app.settings"
    value:
      enabled: true
      features: ["auth", "logging"]
      timeout: 30

# Result:
app:
  settings:
    enabled: true
    features: ["auth", "logging"]
    timeout: 30
```

## Combined Transformations

**Complete Example from Tests:**
```yaml
# Input config:
legacy:
  api_endpoint: "HTTP://OLD-DOMAIN.COM/api"
database:
  host: "db-server"

# Schema with transformations:
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

# Final result:
database:
  host: "db-server"
prod_service:
  prod_url: "http://old-domain.com/api"
  prod_environment: "prod"
```

## Advanced Transformation Patterns

### Conditional Value Setting

Use variables to set values conditionally:

```yaml
vars:
  - name: "ENVIRONMENT"
    fromEnv: "NODE_ENV"
    defaultValue: "development"
  - name: "DEBUG_MODE"
    value: "true"
    # Only in dev
transform:
  - type: "setValue"
    path: "app.debug"
    value: "${DEBUG_MODE}"
```

### Restructuring Legacy Configurations

Transform old configuration formats:

```yaml
transform:
  # Migrate old database config
  - type: "renameKey"
    from: "db.connection_string"
    to: "database.url"
  - type: "renameKey" 
    from: "db.max_connections"
    to: "database.pool.max"
  
  # Add new required fields
  - type: "setValue"
    path: "database.pool.min"
    value: 5
```

### Environment-Specific Transformations

Apply different transformations per environment:

```yaml
vars:
  - name: "ENV"
    fromEnv: "ENVIRONMENT"
    defaultValue: "dev"
transform:
  - type: "setValue"
    path: "app.environment"
    value: "${ENV}"
  - type: "addKeyPrefix"
    path: "database"
    prefix: "${ENV}_"
```

## Error Handling

Common transformation errors:

### Path Not Found
```yaml
transform:
  - type: "renameKey"
    from: "missing.path"  # ERROR: path doesn't exist
    to: "new.path"
```

### Type Mismatch
```yaml
transform:
  - type: "changeCase"
    path: "numeric.value"  # ERROR: value is not a string
    case: "lower"
```

### Invalid Case Format
```yaml
transform:
  - type: "changeCase"
    path: "string.value"
    case: "invalid"  # ERROR: unsupported case format
```

### Non-Map Prefix Target
```yaml
transform:
  - type: "addKeyPrefix"
    path: "string.value"  # ERROR: value is not a map
    prefix: "pre_"
```

## Best Practices

1. **Order Matters**: Plan transformation sequence carefully
2. **Path Validation**: Ensure source paths exist before renaming
3. **Variable Usage**: Leverage variables for dynamic prefixes and values
4. **Error Testing**: Test with invalid inputs to verify error handling
5. **Documentation**: Document complex transformation chains
    ```yaml
    transform:
      - type: "changeCase"
        path: "path.to.string.value" # Dot-separated path to the string value.
        case: "snake"                # Target case: "upper", "lower", "snake", "camel".
    ```
*   **Fields**:
    *   `type` (Required): `"changeCase"`
    *   `path` (Required, string): The dot-separated path to the string value to be modified.
    *   `case` (Required, string): The target case format. Supported values:
        *   `"upper"`: Converts to UPPERCASE.
        *   `"lower"`: Converts to lowercase.
        *   `"snake"`: Converts to snake_case.
        *   `"camel"`: Converts to camelCase (lower camel case).
*   **Behavior**:
    *   Retrieves the value from `path`.
    *   If the path is not found or the value is not a string, an error occurs.
    *   Converts the string to the specified `case`.
    *   Updates the value at `path` with the new cased string.
*   **Example**:
    ```yaml
    # Config before: { "apiSettings": { "RequestTimeout": "ThirtySeconds" } }
    transform:
      - type: "changeCase"
        path: "apiSettings.RequestTimeout"
        case: "snake"
    # Config after: { "apiSettings": { "RequestTimeout": "thirty_seconds" } }
    ```

### 3. `addKeyPrefix`

Adds a prefix to all keys within a map located at a specified path.

*   **Structure**:
    ```yaml
    transform:
      - type: "addKeyPrefix"
        path: "path.to.map.object" # Dot-separated path to the map.
        prefix: "my_prefix_"       # String to prepend to each key in the map.
    ```
*   **Fields**:
    *   `type` (Required): `"addKeyPrefix"`
    *   `path` (Required, string): The dot-separated path to the map object whose keys will be prefixed.
    *   `prefix` (Required, string): The string to prepend to each key within the map at `path`.
*   **Behavior**:
    *   Retrieves the value from `path`.
    *   If the path is not found or the value is not a map, an error occurs.
    *   Creates a new map where each key from the original map is prepended with `prefix`.
    *   Updates the value at `path` with this new map.
*   **Example**:
    ```yaml
    # Config before: { "settings": { "timeout": 30, "retries": 3 } }
    transform:
      - type: "addKeyPrefix"
        path: "settings"
        prefix: "http_"
    # Config after: { "settings": { "http_timeout": 30, "http_retries": 3 } }
    ```

### 4. `setValue`

Sets a specific value at a given path, potentially overwriting an existing value or creating the path if it doesn't exist. The value can be of any valid YAML/JSON type (string, number, boolean, list, map).

*   **Structure**:
    ```yaml
    transform:
      - type: "setValue"
        path: "path.to.target.key" # Dot-separated path where the value will be set.
        value: "New Static Value"   # The value to set. Can be any type.
                                    # If a string, it undergoes variable substitution.
    ```
    ```yaml
    transform:
      - type: "setValue"
        path: "feature.flags"
        value:
          newToggle: true
          betaFeature: false
    ```
*   **Fields**:
    *   `type` (Required): `"setValue"`
    *   `path` (Required, string): The dot-separated path where the `value` will be set.
    *   `value` (Required, any): The value to set at the `path`. This can be a simple literal (string, number, boolean) or a complex nested structure (map, list). If `value` is a string, it will undergo `${VAR_NAME}` substitution before being set.
*   **Behavior**:
    *   Sets the provided `value` at the specified `path`.
    *   If the path or parts of it do not exist, they are created.
    *   If a value already exists at `path`, it is overwritten.
*   **Example (String Value with Variable Substitution)**:
    ```yaml
    # Schema vars:
    # vars:
    #   - name: "ADMIN_EMAIL"
    #     value: "admin@example.com"
    # Config before: { }
    transform:
      - type: "setValue"
        path: "contact.admin"
        value: "Email: ${ADMIN_EMAIL}"
    # Config after: { "contact": { "admin": "Email: admin@example.com" } }
    ```
*   **Example (Complex Value)**:
    ```yaml
    # Config before: { "app": { "name": "MyApp" } }
    transform:
      - type: "setValue"
        path: "app.settings.notifications"
        value:
          enabled: true
          channels: ["email", "sms"]
    # Config after:
    # {
    #   "app": {
    #     "name": "MyApp",
    #     "settings": {
    #       "notifications": {
    #         "enabled": true,
    #         "channels": ["email", "sms"]
    #       }
    #     }
    #   }
    # }
    ```

## Example Combining Transformations

```yaml
# schema.yml
vars:
  - name: "ENV_PREFIX"
    value: "prod"
transform:
  - type: "renameKey"
    from: "legacy.api_endpoint"
    to: "service.url"
  - type: "changeCase"
    path: "service.url" # Assuming it's a string like "HTTP://OLD-DOMAIN.COM"
    case: "lower"
  - type: "setValue"
    path: "service.environment"
    value: "${ENV_PREFIX}"
  - type: "addKeyPrefix"
    path: "service" # Assuming service map now exists
    prefix: "${ENV_PREFIX}_"

# Initial config (merged from sources):
# {
#   "legacy": { "api_endpoint": "HTTP://OLD-DOMAIN.COM/api" },
#   "other_setting": 123
# }

# Expected final config (YAML):
# other_setting: 123
# prod_service:
#   environment: prod
#   url: http://old-domain.com/api
```

This sequence first renames `legacy.api_endpoint` to `service.url`. Then, it changes the case of the string at `service.url`. After that, it sets `service.environment` using a variable. Finally, it prefixes all keys within the `service` map (which now contains `url` and `environment`) with `prod_`.
