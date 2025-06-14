# Schema: Variables & Substitution

Konfigo's variable system is a cornerstone of its processing capabilities, allowing you to create dynamic and reusable configurations. Variables can be defined in multiple locations and are resolved using a strict precedence order. They are substituted into your configuration data and even within other schema directives using the `${VAR_NAME}` syntax.

## Variable Definition and Precedence

Variables can be defined in three main places, listed here from highest to lowest precedence:

1.  **Environment Variables (`KONFIGO_VAR_...`) (Highest Priority)**
    *   **Syntax**: `KONFIGO_VAR_VARNAME=value` (e.g., `export KONFIGO_VAR_API_KEY=secret123`)
    *   **Description**: Variables set in the environment using the `KONFIGO_VAR_` prefix. These override any other variable definitions.
    *   **Use Case**: Ideal for injecting secrets or highly dynamic, environment-specific values during runtime or in CI/CD pipelines.
    *   See [Environment Variables](../guide/environment-variables.md) for more details.

2.  **Variables File (`-V` or `--vars-file`)**
    *   **Syntax**: A separate YAML, JSON, or TOML file passed via the `-V` or `--vars-file` CLI flag.
    *   **Description**: This file can contain simple key-value pairs that define variables. It can also host the [`konfigo_forEach`](#batch-processing-with-konfigo_foreach) directive for batch processing.
    *   **Use Case**: Defining sets of variables for specific environments (e.g., `dev-vars.yml`, `prod-vars.yml`) or for controlling batch operations.
    *   **Example (`example-vars.yml` provided via `-V my-configs/example-vars.yml`):**
        ```yaml
        # Simple key-value pairs
        GREETING: "Hello from example-vars.yml"
        SERVICE_NAME: "my-awesome-app"
        REPLICA_COUNT: 3

        # Nested structures are also possible, though direct variable
        # substitution typically uses simple key-value pairs from the resolved map.
        DATABASE:
          HOST: "db.example.com"
          PORT: 5432
        ```
        In this example, `${GREETING}`, `${SERVICE_NAME}`, `${REPLICA_COUNT}`, `${DATABASE.HOST}`, and `${DATABASE.PORT}` (if flattened or accessed via path) would be available. Konfigo typically flattens these for direct substitution, or you might refer to nested values if your schema logic supports it (e.g., `fromPath` in schema `vars`).

3.  **Schema `vars` Block (Lowest Priority)**
    *   **Syntax**: The `vars:` block list within your main schema file (`-S`).
    *   **Description**: Defines the default set of variables, their sources (literal, from other environment variables, from other config paths), and fallback default values.
    *   **Use Case**: Establishing the baseline variable logic for your application.

## Defining Variables in the Schema (`vars:` block)

The `vars` block in your schema file is a list of variable definitions. Each definition is an object that specifies the variable's `name` and how its value should be determined.

### Common Fields for each Variable Definition:

*   `name` (Required, string): The name of the variable (e.g., `API_URL`). This is the name you'll use for substitution, like `${API_URL}`.

### Value Sources (choose one per definition):

*   `value` (string):
    *   **Description**: Defines a literal, static string value for the variable.
    *   **Example**:
        ```yaml
        vars:
          - name: "DEFAULT_REGION"
            value: "us-east-1"
        ```

*   `fromEnv` (string):
    *   **Description**: Sources the variable's value from a system environment variable (different from `KONFIGO_VAR_...`). This allows you to map existing system environment variables to Konfigo variables.
    *   **Example**:
        ```yaml
        vars:
          - name: "DOCKER_TAG"
            fromEnv: "CI_COMMIT_SHA" # Reads the value of the CI_COMMIT_SHA system env var
        ```

*   `fromPath` (string):
    *   **Description**: Sources the variable's value from another key within the *merged configuration data* (i.e., after all `-s` sources are merged, but before most schema processing like generators or transforms). The path is dot-separated.
    *   **Example**:
        ```yaml
        # Assuming merged config has: deployment: { namespace: "production" }
        vars:
          - name: "PRIMARY_NAMESPACE"
            fromPath: "deployment.namespace" # Value will be "production"
        ```

### Optional Fallback:

*   `defaultValue` (string):
    *   **Description**: Provides a fallback value if a variable defined using `fromEnv` or `fromPath` cannot be resolved (e.g., the environment variable is not set, or the path does not exist in the configuration).
    *   This is **not** used if `value` is specified.
    *   **Example**:
        ```yaml
        vars:
          - name: "RELEASE_VERSION"
            fromEnv: "CI_COMMIT_TAG"
            defaultValue: "latest" # If CI_COMMIT_TAG is not set, RELEASE_VERSION becomes "latest"
          - name: "OPTIONAL_SETTING"
            fromPath: "user.preferences.theme"
            defaultValue: "dark"
        ```

### Resolution Logic within Schema `vars`:

For each variable defined in the schema's `vars` block:
1.  If `value` is present, that's the variable's value.
2.  Else, if `fromEnv` is present, Konfigo attempts to read that system environment variable.
3.  Else, if `fromPath` is present, Konfigo attempts to read that path from the merged configuration.
4.  If the chosen source (`fromEnv` or `fromPath`) yields a value, that's used.
5.  If not, and `defaultValue` is present, that's used.
6.  If the source doesn't yield a value and no `defaultValue` is provided, Konfigo will error, as the variable cannot be resolved.

This resolved value from the schema `vars` block is then subject to being overridden by the `-V` file or `KONFIGO_VAR_...` environment variables as per the overall precedence rules.

## Variable Substitution

Once all variables are resolved, Konfigo performs substitution wherever `${VAR_NAME}` placeholders appear. This includes:

*   Values within your configuration data.
*   Certain fields within schema directives themselves (e.g., paths in `generators`, `transform`, `validate`, or even values in `setValue` transforms).

**Example of Substitution in Config:**

`config.yml`:
```yaml
server:
  url: "${API_HOST}:${API_PORT}/v1"
  timeout: "${DEFAULT_TIMEOUT}"
```

If `API_HOST=api.example.com`, `API_PORT=8443`, and `DEFAULT_TIMEOUT=30s` are resolved, the processed config will have:
```yaml
server:
  url: "api.example.com:8443/v1"
  timeout: "30s"
```

## Batch Processing with `konfigo_forEach`

Konfigo supports generating multiple output files from a single schema by iterating over sets of variables. This is a powerful feature for managing configurations across different environments, services, or any scenario requiring multiple variations of a base template.

This feature is activated by defining a `konfigo_forEach` block in the variables file specified with the `-V` or `--vars-file` flag.

### `konfigo_forEach` Structure

The `konfigo_forEach` block has the following structure within your `-V` variables file:

```yaml
# In your main variables file (e.g., -V loop-vars.yml)

# Optional: Global variables accessible to all iterations unless overridden by iteration-specific vars.
# These follow the standard variable precedence (KONFIGO_VAR_ > -V file global > schema vars).
GLOBAL_API_KEY: "default_global_key"
DEPLOYMENT_TIER: "general"

konfigo_forEach:
  # Specify the source of iteration data (choose ONE):
  items: # Option 1: Define variable sets directly as a list of maps.
    - SERVICE_NAME: "frontend"
      REPLICAS: 3
      PORT: 80
      DEPLOYMENT_TIER: "web" # Overrides DEPLOYMENT_TIER for this item
    - SERVICE_NAME: "backend-api"
      REPLICAS: 5
      PORT: 8080
      # DEPLOYMENT_TIER will be "general" (from global) for this item

  # itemFiles: # Option 2: List of external variable files (YAML, JSON, or TOML).
  #   # Paths are relative to the main variables file if not absolute.
  #   - "service-configs/frontend-vars.yml"
  #   - "service-configs/backend-vars.json"

  output:
    # Defines how output files are named and formatted for each iteration.
    # Placeholders:
    #   - `${VAR_NAME}`: Any variable from the current iteration's scope (iteration-specific, global, or schema-resolved).
    #   - `${ITEM_INDEX}`: The 0-based index of the current iteration.
    #   - `${ITEM_FILE_BASENAME}`: If using `itemFiles`, the basename of the current variable file (e.g., "frontend-vars" from "frontend-vars.yml").
    #                             This is an empty string if using `items`.
    filenamePattern: "dist/${SERVICE_NAME}/config-${ITEM_INDEX}.json" # Example: dist/frontend/config-0.json

    # Optional: Overrides the global output format (from -oX flags or filename extension from -of) for generated files.
    # Valid formats: "json", "yaml", "toml", "env".
    # If not set, format is inferred from filenamePattern's extension, or defaults to YAML if ambiguous.
    # format: "yaml"
```

### Key aspects of `konfigo_forEach`:

*   **Location**: Must be in the variables file supplied via `-V`.
*   **Iteration Data**:
    *   `items`: A list of maps, where each map represents a set of variables for one iteration.
    *   `itemFiles`: A list of paths to other variable files. Each file provides variables for one iteration. Paths are relative to the main `-V` file's directory if not absolute.
    *   You must use *either* `items` *or* `itemFiles`, not both.
*   **Global Variables**: Variables defined in the `-V` file *outside* the `konfigo_forEach` block are considered global. They are available to each iteration unless an iteration-specific variable (from `items` or an `itemFile`) has the same name.
*   **Output Configuration (`output`)**:
    *   `filenamePattern` (Required): A template for generating output filenames. It can use `${VAR_NAME}`, `${ITEM_INDEX}`, and `${ITEM_FILE_BASENAME}` placeholders.
        *   `${VAR_NAME}` resolution for filename patterns prioritizes:
            1.  Iteration-specific variables (from `items` or `itemFile`).
            2.  `KONFIGO_VAR_...` environment variables.
            3.  Simple `value` or `defaultValue` from the schema's `vars` block (does not resolve `fromEnv` or `fromPath` for filenames).
    *   `format` (Optional): Explicitly sets the output format for all generated files in the loop, overriding format inference from `filenamePattern`'s extension or global output flags.

### Variable Precedence in `konfigo_forEach` Mode

For each generated file during a `konfigo_forEach` loop, the variable resolution order is:

1.  **`KONFIGO_VAR_...` Environment Variables (Highest Priority)**
2.  **Current Iteration Variables**:
    *   If using `items`: Variables from the current item in the list.
    *   If using `itemFiles`: Variables loaded from the current item file.
3.  **Global Variables from `-V` file**: Variables defined in the main `-V` file *outside* the `konfigo_forEach` block.
4.  **Schema `vars` Block (Lowest Priority)**: Variables defined in the `vars:` section of your main schema file (`-S`).

### How It Works

1.  Konfigo loads the main schema (`-S`) and the primary variables file (`-V`).
2.  It detects the `konfigo_forEach` block within the `-V` file.
3.  It determines the iteration data (either `items` or `itemFiles`).
4.  For each iteration (each item in `items` or each file in `itemFiles`):
    a.  A deep copy of the base merged configuration (from `-s` sources) is created.
    b.  A unique set of variables is prepared for this iteration according to the precedence rules above. This includes `${ITEM_INDEX}` and `${ITEM_FILE_BASENAME}`.
    c.  The output filename is generated using `output.filenamePattern` and the current iteration's variables.
    d.  The main schema (`-S`) is processed against the copied configuration using this iteration's specific variable set. This includes all standard Konfigo steps: `vars` resolution (within the schema, if any are still relevant), `generators`, `transform`, variable substitution in config values, and `validate`.
    e.  The output directory for the generated file is created if it doesn't exist.
    f.  The resulting configuration is marshalled to the specified (or inferred) `output.format` and written to the generated filename.
5.  Once all iterations are complete, Konfigo exits. Normal output flags (`-of`, `-oj`, etc.) are ignored when `konfigo_forEach` is active, as output is fully controlled by the directive.

### Example Usage of `konfigo_forEach`

**Schema (`schema.yml`):**
```yaml
# schema.yml
vars:
  - name: "DEFAULT_TIMEOUT"
    value: "30s"
  - name: "LOG_LEVEL"
    defaultValue: "info"
config:
  serviceName: "${SERVICE_NAME}" # From iteration
  instanceCount: ${REPLICAS}    # From iteration
  apiPort: ${PORT}              # From iteration
  networkZone: "${ZONE}"        # From iteration or global
  timeout: "${DEFAULT_TIMEOUT}" # From schema vars
  logLevel: "${LOG_LEVEL}"      # From schema vars, potentially overridden
  globalSetting: "${GLOBAL_CONFIG_VAL}" # From -V file global
```

**Variables File (`loop-controller.yml` passed with `-V loop-controller.yml`):**
```yaml
# loop-controller.yml
GLOBAL_CONFIG_VAL: "shared-across-all"
ZONE: "default-zone" # Global, can be overridden by items

konfigo_forEach:
  items:
    - SERVICE_NAME: "user-service"
      REPLICAS: 2
      PORT: 8001
      LOG_LEVEL: "debug" # Overrides schema default for this item
    - SERVICE_NAME: "order-service"
      REPLICAS: 4
      PORT: 8002
      ZONE: "high-traffic-zone" # Overrides global ZONE for this item
  output:
    filenamePattern: "generated-configs/${ZONE}/${SERVICE_NAME}/app-config.v${ITEM_INDEX}.yml"
    # format: "yaml" # Optional, can be inferred from .yml in pattern
```

**Command:**
```bash
konfigo -s base-template.json -S schema.yml -V loop-controller.yml
```
(Assuming `base-template.json` is an empty JSON `{}` or contains foundational structure that doesn't conflict with schema-generated keys.)

**Expected Output Files:**

*   `generated-configs/default-zone/user-service/app-config.v0.yml`:
    ```yaml
    serviceName: user-service
    instanceCount: 2
    apiPort: 8001
    networkZone: default-zone
    timeout: 30s
    logLevel: debug
    globalSetting: shared-across-all
    ```
*   `generated-configs/high-traffic-zone/order-service/app-config.v1.yml`:
    ```yaml
    serviceName: order-service
    instanceCount: 4
    apiPort: 8002
    networkZone: high-traffic-zone
    timeout: 30s
    logLevel: info # Falls back to schema default as not set in item
    globalSetting: shared-across-all
    ```

This powerful combination of layered variable resolution and batch processing makes Konfigo highly adaptable for complex configuration scenarios.
