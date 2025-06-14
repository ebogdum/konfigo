# Konfigo Schema: The Engine of Configuration Processing

The Konfigo schema is a powerful YAML, JSON, or TOML file that you provide using the `-S` or `--schema` flag. It acts as the central engine for processing your merged configuration data, allowing you to define variables, generate new data, transform existing structures, and validate the final output.

A schema file can contain the following top-level keys:

*   [`apiVersion`](#apiversion) (Optional): Specifies the schema version.
*   [`inputSchema`](#inputschema--outputschema) (Optional): Defines a path to an external schema file for validating the raw input configuration *before* any processing.
*   [`outputSchema`](#inputschema--outputschema) (Optional): Defines a path to an external schema file for filtering the final processed configuration *after* all processing, ensuring only specified fields are included in the output.
*   [`immutable`](#immutable) (Optional): A list of dot-separated paths that should be treated as immutable during merges. Once a value is set for an immutable path from an earlier source, subsequent sources cannot override it. `KONFIGO_KEY_...` environment variables *can* still override immutable paths.
*   [`vars`](./variables.md) (Optional): Defines variables for substitution throughout the configuration and schema directives.
*   [`generators`](./generation.md) (Optional): Defines rules for generating new configuration data.
*   [`transform`](./transformation.md) (Optional): Defines rules for transforming existing configuration data (e.g., renaming keys, changing case).
*   [`validate`](./validation.md) (Optional): Defines rules for validating the processed configuration data.

## Core Concepts

### `apiVersion`

An optional string that can be used for schema versioning if you plan to evolve your schema definitions over time. Konfigo itself does not currently enforce specific versions but may do so in the future.

```yaml
apiVersion: "konfigo/v1alpha1"
```

### `inputSchema` & `outputSchema`

These directives allow you to define structural expectations for your configuration data at the beginning and end of the processing pipeline.

*   **`inputSchema`**:
    *   **Purpose**: Validates the structure of the configuration *after* all source files (`-s`) are merged but *before* any `vars`, `generators`, `transform`, or `validate` directives from the main schema are applied.
    *   **Fields**:
        *   `path` (Required): Path to an external schema file (JSON, YAML, or TOML format). This external schema's structure is compared against the merged input.
        *   `strict` (Optional, boolean, default: `false`): If `true`, the input configuration must *only* contain keys defined in the `inputSchema.path` file. Extra keys will result in an error. If `false`, extra keys are allowed.
    *   **Use Case**: Ensure that input configurations meet a basic structural contract before more complex processing begins.

*   **`outputSchema`**:
    *   **Purpose**: Filters the final configuration *after* all `vars`, `generators`, `transform`, and `validate` directives have been applied. Only the keys present in the `outputSchema.path` file will be included in the final output if `strict` is `false` (default). If `strict` is `true`, the processed configuration must exactly match the structure defined in `outputSchema.path`.
    *   **Fields**:
        *   `path` (Required): Path to an external schema file (JSON, YAML, or TOML format). This file acts as a template for the output.
        *   `strict` (Optional, boolean, default: `false`):
            *   If `false` (default): Only keys present in the `outputSchema.path` file will be included in the final output. Extra keys in the processed configuration are silently ignored. Keys in the `outputSchema.path` but missing from the processed configuration are also ignored.
            *   If `true`: The structure of the processed configuration must exactly match the structure defined in `outputSchema.path`. 
                *   Any key found in the processed configuration but not defined in `outputSchema.path` will result in an error.
                *   Any key defined in `outputSchema.path` but not found in the processed configuration will result in an error.
                *   If `outputSchema.path` defines a path as a map, but the processed configuration has a non-map type at that path, it will result in an error.
    *   **Use Case**: Produce a clean, well-defined output configuration. With `strict: false`, it removes any intermediate or temporary fields. With `strict: true`, it ensures the output conforms precisely to an expected contract.

**Example**:

```yaml
# Main schema.yml
inputSchema:
  path: "./schemas/expected_input_structure.json"
  strict: true
outputSchema:
  path: "./schemas/final_output_structure.json"

# ... other schema directives (vars, generators, etc.) ...
```

### `immutable`

A list of dot-separated configuration paths that should resist changes from subsequent configuration sources during the initial merge phase.

*   **Purpose**: Protect foundational or critical configuration values from being accidentally overridden by later, potentially more specific, configuration files.
*   **Behavior**:
    *   When Konfigo merges multiple source files (`-s`), if a key at an immutable path is set by an earlier source, any attempts by later sources to change that key's value will be ignored.
    *   **Exception**: `KONFIGO_KEY_...` environment variables *can* override values at immutable paths. This provides an escape hatch for essential runtime overrides.
*   **Use Case**: Defining global settings like `application.name` or `cluster.id` in a base configuration file and preventing environment-specific files from changing them.

**Example**:

```yaml
# schema.yml
immutable:
  - "service.name"
  - "logging.level"

# base.yml
service:
  name: "my-core-app"
logging:
  level: "INFO"

# env-specific.yml (attempted overrides will be ignored for immutable paths)
service:
  name: "my-specific-app" # This will be ignored
logging:
  level: "DEBUG" # This will be ignored
  format: "json" # This will be merged as it's not immutable
```

If `KONFIGO_KEY_service.name=my-runtime-app` is set, it *will* override `my-core-app`.

## Detailed Schema Sections

For in-depth information on each processing capability, refer to their dedicated pages:

*   **[Variables & Substitution](./variables.md)**: Learn how to define and use variables, including the powerful `konfigo_forEach` for batch processing.
*   **[Data Generation](./generation.md)**: Discover how to create new configuration values.
*   **[Data Transformation](./transformation.md)**: Explore ways to modify your configuration's structure and content.
*   **[Data Validation](./validation.md)**: Understand how to enforce rules and constraints on your configuration.

By combining these features, you can build robust and maintainable configuration management pipelines with Konfigo.
