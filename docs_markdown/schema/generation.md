# Schema: Data Generation

Konfigo's schema allows you to define `generators` that create new data within your configuration. Generators construct values from existing configuration data and variables, enabling dynamic configuration composition.

## Available Generators

Konfigo provides several built-in generators:

- **`concat`**: Combines configuration data, variables, and literal text
- **`timestamp`**: Generates current timestamp in various formats
- **`random`**: Generates random values (integers, floats, strings, UUIDs)
- **`id`**: Generates various types of identifiers using alphanumeric characters

## `concat` Generator

The `concat` generator creates string values by combining configuration data, variable substitutions, and literal text. It operates in two phases: first replacing `{placeholder}` tokens with values from configuration paths, then resolving `${VARIABLE}` references.

### Structure

```yaml
generators:
  - type: "concat"
    targetPath: "path.to.new.key"
    format: "template with {placeholders} and ${VARIABLES}"
    sources:
      placeholder_name: "config.path.to.value"
```

### Fields

- **`type`** (Required): Must be `"concat"`
- **`targetPath`** (Required): Dot-separated path where the generated value will be placed
- **`format`** (Required): Template string with `{placeholders}` and `${VARIABLES}`
- **`sources`** (Required): Map of placeholder names to configuration paths

### Processing Order

1. Replace `{placeholder}` tokens using `sources` mapping
2. Resolve `${VARIABLE}` references from vars and environment
3. Insert result at `targetPath`

## `timestamp` Generator

The `timestamp` generator creates timestamp values in various formats.

### Structure

```yaml
generators:
  - type: "timestamp"
    targetPath: "path.to.timestamp"
    format: "rfc3339"  # Optional, defaults to "rfc3339"
```

### Fields

- **`type`** (Required): Must be `"timestamp"`
- **`targetPath`** (Required): Dot-separated path where the timestamp will be placed
- **`format`** (Optional): Timestamp format (defaults to "rfc3339")

### Supported Formats

- **`unix`**: Unix timestamp (seconds since epoch) - e.g., `1640995200`
- **`unixmilli`**: Unix timestamp in milliseconds - e.g., `1640995200000`
- **`rfc3339`**: RFC3339 format - e.g., `2021-12-31T16:00:00Z`
- **`iso8601`**: ISO8601 format - e.g., `2021-12-31T16:00:00Z`
- **Custom**: Any Go time format string - e.g., `2006-01-02 15:04:05`

### Example

```yaml
generators:
  - type: "timestamp"
    targetPath: "metadata.createdAt"
    format: "rfc3339"
  - type: "timestamp"
    targetPath: "metadata.buildTime"
    format: "2006-01-02 15:04:05"
```

## `random` Generator

The `random` generator creates random values in various formats.

### Structure

```yaml
generators:
  - type: "random"
    targetPath: "path.to.random.value"
    format: "string:16"  # Required
```

### Fields

- **`type`** (Required): Must be `"random"`
- **`targetPath`** (Required): Dot-separated path where the random value will be placed
- **`format`** (Required): Random value format specification

### Supported Formats

- **`int:min:max`**: Random integer between min and max (inclusive) - e.g., `int:1:100`
- **`float:min:max`**: Random float between min and max - e.g., `float:0.0:1.0`
- **`string:length`**: Random string using [a-zA-Z0-9] - e.g., `string:16`
- **`bytes:length`**: Random bytes as hex string - e.g., `bytes:8`
- **`uuid`**: UUID v4 format - e.g., `uuid`

### Examples

```yaml
generators:
  - type: "random"
    targetPath: "service.port"
    format: "int:8000:9000"
  - type: "random"
    targetPath: "session.id"
    format: "string:32"
  - type: "random"
    targetPath: "request.uuid"
    format: "uuid"
```

## `id` Generator

The `id` generator creates various types of identifiers using alphanumeric characters.

### Structure

```yaml
generators:
  - type: "id"
    targetPath: "path.to.id"
    format: "simple:8"  # Optional, defaults to "simple:8"
```

### Fields

- **`type`** (Required): Must be `"id"`
- **`targetPath`** (Required): Dot-separated path where the ID will be placed
- **`format`** (Optional): ID format specification (defaults to "simple:8")

### Supported Formats

- **`simple:length`**: Random ID using [a-zA-Z0-9] - e.g., `simple:12`
- **`prefix:prefix:length`**: ID with prefix + random chars - e.g., `prefix:user_:8`
- **`numeric:length`**: Numeric ID using [0-9] - e.g., `numeric:6`
- **`alpha:length`**: Alphabetic ID using [a-zA-Z] - e.g., `alpha:10`
- **`sequential`**: Sequential counter-based ID (starts from 1)
- **`timestamp`**: Timestamp + random suffix - e.g., `1640995200abcd`

### Examples

```yaml
generators:
  - type: "id"
    targetPath: "user.id"
    format: "prefix:usr_:8"     # Results in: usr_A9Kx2mP1
  - type: "id"
    targetPath: "session.counter"
    format: "sequential"        # Results in: 1, 2, 3, ...
  - type: "id"
    targetPath: "trace.id"
    format: "timestamp"         # Results in: 1640995200A9Kx
```

## Examples from Tests

### Basic Generation

**Input Configuration:**
```yaml
service:
  name: "data-processor"
  instanceId: "instance-007"
  port: 8080
region: "us-west-2"
```

**Schema:**
```yaml
vars:
  - name: "APP_VERSION"
    value: "1.2.3"
generators:
  - type: "concat"
    targetPath: "service.identifier"
    format: "Service: {name} (ID: {id}) running in {region}"
    sources:
      name: "service.name"
      id: "service.instanceId"
      region: "region"
```

**Result:**
```yaml
service:
  identifier: "Service: data-processor (ID: instance-007) running in us-west-2"
  # ... other fields preserved
```

### Multiple Generators with Variables

**Schema:**
```yaml
vars:
  - name: "APP_VERSION"
    value: "1.2.3"
  - name: "DOMAIN"
    value: "example.com"
generators:
  - type: "concat"
    targetPath: "service.url"
    format: "https://{service}.${DOMAIN}:{port}"
    sources:
      service: "service.name"
      port: "service.port"
  - type: "concat"
    targetPath: "database.connectionString"
    format: "postgresql://{host}:{port}/{db}"
    sources:
      host: "database.host"
      port: "database.port"
      db: "database.name"
  - type: "concat"
    targetPath: "service.fullIdentifier"
    format: "{name}-{version} - ${APP_VERSION} ({env})"
    sources:
      name: "service.name"
      version: "service.version"
      env: "environment"
```

**Result:**
```yaml
service:
  url: "https://data-processor.example.com:8080"
  fullIdentifier: "data-processor-1.2.3 - 1.2.3 (production)"
database:
  connectionString: "postgresql://db-server:5432/app_db"
```

### Variables-Only Generation

Generators can use only variables without configuration sources:

**Schema:**
```yaml
vars:
  - name: "APP_VERSION"
    value: "1.2.3"
  - name: "DOMAIN"
    value: "example.com"
generators:
  - type: "concat"
    targetPath: "metadata.buildInfo"
    format: "Built version ${APP_VERSION} for ${DOMAIN}"
    sources: {}  # No configuration sources needed
```

## Error Handling

Common errors and their solutions:

### Missing Source Path
```yaml
# ERROR: 'service.missing' not found in configuration
sources:
  name: "service.missing"  # This path doesn't exist
```

### Empty Format
```yaml
# ERROR: format cannot be empty
format: ""
```

### No Sources When Needed
```yaml
# ERROR: Missing sources for placeholder {name}
format: "Service: {name}"
sources: {}  # Missing required source mapping
```

## Advanced Patterns

### Cascading Generation

Generators can reference previously generated values:

```yaml
generators:
  - type: "concat"
    targetPath: "temp.namespace"
    format: "{app}-{env}"
    sources:
      app: "metadata.appName"
      env: "environment"
  - type: "concat"
    targetPath: "kubernetes.namespace"
    format: "k8s-{namespace}"
    sources:
      namespace: "temp.namespace"  # Reference previous generation
```

### Deep Path Creation

Generators automatically create nested paths:

```yaml
generators:
  - type: "concat"
    targetPath: "deeply.nested.generated.value"
    format: "Generated at {path}"
    sources:
      path: "metadata.location"
```

Creates the entire path structure if it doesn't exist.

## Best Practices

1. **Clear Naming**: Use descriptive placeholder names in `sources`
2. **Path Validation**: Ensure source paths exist in your configuration
3. **Variable Order**: Define variables before generators that use them
4. **Avoid Conflicts**: Don't overwrite critical configuration paths
5. **Test Generation**: Verify generated values match expectations
*   `sources` (Required, map): A map where:
    *   Keys are the `placeholder_name` strings used in the `format` string (without the curly braces).
    *   Values are dot-separated paths to existing values within the *current state of the configuration* (after merges, but typically before most other schema processing for the current generator pass).

### How it Works:

1.  For each `placeholder` defined in `sources`:
    a.  Konfigo retrieves the value from the configuration at the specified `path`.
    b.  If a source path is not found, Konfigo will return an error.
2.  The `format` string is processed:
    a.  Each `{placeholder_name}` is replaced with the corresponding value retrieved from `sources`.
3.  The resulting string (after `sources` substitution) is then processed for standard Konfigo variable substitution (e.g., `${MY_VARIABLE}`).
4.  The final, fully resolved string is set at the `targetPath` in the configuration.

### Example:

**Schema (`schema.yml`):**
```yaml
vars:
  - name: "APP_VERSION"
    value: "1.2.3"
config: # Assume this is the state of the config before this generator runs
  service:
    name: "data-processor"
    instanceId: "instance-007"
  region: "us-west-2"

generators:
  - type: "concat"
    targetPath: "service.identifier"
    format: "Service: {name} (ID: {id}) running in {region_val} - Version: ${APP_VERSION}"
    sources:
      name: "service.name"
      id: "service.instanceId"
      region_val: "region"
  - type: "concat"
    targetPath: "service.url"
    format: "https://${FQDN_VAR}" # Using only a global variable
    sources: {} # No local sources needed if format string only uses global vars
```

**Variables (e.g., from `-V vars.yml` or environment):**
```yaml
# vars.yml
FQDN_VAR: "myapp.example.com"
```

**Processing Steps:**

1.  **First Generator (`service.identifier`):**
    *   `sources`:
        *   `name` -> `service.name` -> "data-processor"
        *   `id` -> `service.instanceId` -> "instance-007"
        *   `region_val` -> `region` -> "us-west-2"
    *   `format` after `sources` substitution: `"Service: data-processor (ID: instance-007) running in us-west-2 - Version: ${APP_VERSION}"`
    *   `format` after `${APP_VERSION}` substitution: `"Service: data-processor (ID: instance-007) running in us-west-2 - Version: 1.2.3"`
    *   `config.service.identifier` becomes `"Service: data-processor (ID: instance-007) running in us-west-2 - Version: 1.2.3"`

2.  **Second Generator (`service.url`):**
    *   `sources`: (empty)
    *   `format` after `sources` substitution (no change): `"https://${FQDN_VAR}"`
    *   `format` after `${FQDN_VAR}` substitution: `"https://myapp.example.com"`
    *   `config.service.url` becomes `"https://myapp.example.com"`

**Resulting Configuration (snippet):**
```yaml
service:
  name: "data-processor"
  instanceId: "instance-007"
  identifier: "Service: data-processor (ID: instance-007) running in us-west-2 - Version: 1.2.3"
  url: "https://myapp.example.com"
region: "us-west-2"
# ... other config ...
```

Generators are applied in the order they are defined in the `generators` list. The output of one generator can potentially be used as a source for a subsequent generator if its `targetPath` is referenced in the later generator's `sources`.
