# Schema: Data Transformation

Konfigo's `transform` directives allow you to modify the structure and content of your merged configuration data. Transformations are applied after initial merging and variable resolution (for the transform directives themselves) but typically before final validation.

Transformations are defined as a list under the `transform` key in your schema file. Each transformation object in the list defines a specific operation.

## Common Behavior

*   **Variable Substitution in Directives**: Before a transformation is applied, Konfigo will perform variable substitution (e.g., `${MY_VAR}`) on the string values within the transformation definition itself (like `path`, `from`, `to`, `prefix`, or `value` if it's a string).
*   **Order of Execution**: Transformations are executed in the order they appear in the `transform` list. The output of one transformation becomes the input for the next.

## Transformation Types

### 1. `renameKey`

Changes the name of a key at a specified path, effectively moving its value to a new key and deleting the old one.

*   **Structure**:
    ```yaml
    transform:
      - type: "renameKey"
        from: "old.path.to.key"  # Dot-separated path of the key to rename.
        to: "new.path.to.key"    # Dot-separated new path for the key.
    ```
*   **Fields**:
    *   `type` (Required): `"renameKey"`
    *   `from` (Required, string): The current dot-separated path to the key you want to rename.
    *   `to` (Required, string): The new dot-separated path where the value should be moved.
*   **Behavior**:
    *   Retrieves the value from the `from` path.
    *   If the `from` path is not found, an error occurs.
    *   Sets the retrieved value at the `to` path (creating intermediate maps if necessary).
    *   Deletes the key at the original `from` path.
*   **Example**:
    ```yaml
    # Config before: { "user": { "name": "Alice", "id": 123 } }
    transform:
      - type: "renameKey"
        from: "user.name"
        to: "user.fullName"
    # Config after: { "user": { "fullName": "Alice", "id": 123 } }
    ```

### 2. `changeCase`

Modifies the case of a string value at a specified path.

*   **Structure**:
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
