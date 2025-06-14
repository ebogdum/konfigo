# Schema: Data Generation

Konfigo's schema allows you to define `generators` that create new data within your configuration. This is useful for constructing values based on existing configuration or variables.

Generators are defined as a list under the `generators` key in your schema file. Each generator must specify a `type` and the necessary parameters for that type.

## `concat` Generator

The `concat` generator is currently the primary generator type. It constructs a new string value by concatenating other string values, which can be sourced from existing configuration paths or resolved variables.

### Structure:

```yaml
generators:
  - type: "concat"
    targetPath: "path.to.new.key"  # Dot-separated path where the generated string will be placed.
    format: "Value1: {placeholder1}, Value2: {placeholder2}, Var: ${MY_VARIABLE}" # String template.
    sources: # Map of placeholders in 'format' to their source paths in the config.
      placeholder1: "path.to.source.value1"
      placeholder2: "another.config.value"
```

### Fields:

*   `type` (Required, string): Must be `"concat"`.
*   `targetPath` (Required, string): A dot-separated path specifying where the newly generated string value should be inserted into the configuration. If the path doesn't exist, it will be created. If it exists, its value will be overwritten.
*   `format` (Required, string): A template string that defines the structure of the generated output.
    *   It can contain literal text.
    *   It can contain placeholders in the format `{placeholder_name}`. These placeholders will be replaced by values from the `sources` map.
    *   It can also contain standard Konfigo variables like `${MY_VARIABLE}`, which will be resolved *after* the `sources` placeholders are processed.
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
