# Schema: Data Generation

Konfigo's schema allows you to define `generators` that create new data within your configuration. Generators construct values from existing configuration data and variables, enabling dynamic configuration composition.

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
