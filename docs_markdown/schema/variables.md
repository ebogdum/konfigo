# Schema: Variables & Substitution

Konfigo's variable system enables dynamic, reusable configurations through layered variable resolution and `${VAR_NAME}` substitution. Variables can be defined in multiple locations with strict precedence ordering.

## Variable Precedence (Highest to Lowest)

1. **Environment Variables** (`KONFIGO_VAR_*`)
2. **Variables File** (`-V` or `--vars-file`)  
3. **Schema vars Block** (in schema file)

## Schema Variables Block

Define variables in your schema file using the `vars` key:

```yaml
vars:
  - name: "VARIABLE_NAME"
    value: "literal value"
    # OR use one of these sources:
    # fromEnv: "SYSTEM_ENV_VAR"
    # fromPath: "config.path.to.value"
    # defaultValue: "fallback"  # Used with fromEnv/fromPath
```

### Variable Definition Fields

**`name`** (Required): Variable name for `${VAR_NAME}` substitution

**Value Sources** (choose one):
- **`value`**: Literal string value
- **`fromEnv`**: Read from system environment variable
- **`fromPath`**: Read from merged configuration path

**`defaultValue`** (Optional): Fallback when `fromEnv`/`fromPath` fails

## Examples from Tests

### Basic Variable Definition

**Schema:**
```yaml
vars:
  # Literal value
  - name: "API_HOST"
    value: "api.example.com"
  
  # From environment with default
  - name: "API_PORT"
    fromEnv: "SERVICE_PORT"
    defaultValue: "8080"
    
  # From configuration path
  - name: "TARGET_NAMESPACE"
    fromPath: "deployment.namespace"
    
  # Environment with fallback
  - name: "DATABASE_PASSWORD"
    fromEnv: "DB_PASS"
    defaultValue: "default-password"
```

**Configuration Usage:**
```yaml
config:
  api:
    host: "${API_HOST}"
    port: "${API_PORT}"
    endpoint: "${API_HOST}:${API_PORT}/api/v1"
  database:
    connectionString: "postgres://user:${DATABASE_PASSWORD}@${API_HOST}:5432/db"
  settings:
    namespace: "${TARGET_NAMESPACE}"
```

### Input Configuration

**`base-config.yaml`:**
```yaml
deployment:
  namespace: "production"
  replicas: 3
database:
  host: "db-server"
  port: 5432
```

### Variable Resolution Process

**Without Environment Variables:**
```yaml
# Resolved variables:
# API_HOST = "api.example.com" (from value)
# API_PORT = "8080" (fromEnv failed, used defaultValue)
# TARGET_NAMESPACE = "production" (fromPath successful)
# DATABASE_PASSWORD = "default-password" (fromEnv failed, used defaultValue)

# Final result:
config:
  api:
    host: "api.example.com"
    port: "8080"
    endpoint: "api.example.com:8080/api/v1"
  database:
    connectionString: "postgres://user:default-password@api.example.com:5432/db"
  settings:
    namespace: "production"
```

**With Environment Variables:**
```bash
export SERVICE_PORT=9000
export DB_PASS=secure-password
export KONFIGO_VAR_API_HOST=prod-api.example.com
```

```yaml
# Resolved variables:
# API_HOST = "prod-api.example.com" (KONFIGO_VAR_* highest precedence)
# API_PORT = "9000" (fromEnv successful)
# TARGET_NAMESPACE = "production" (fromPath successful)
# DATABASE_PASSWORD = "secure-password" (fromEnv successful)

# Final result:
config:
  api:
    host: "prod-api.example.com"
    port: "9000"
    endpoint: "prod-api.example.com:9000/api/v1"
  database:
    connectionString: "postgres://user:secure-password@prod-api.example.com:5432/db"
  settings:
    namespace: "production"
```

## Variables File Override

**External Variables File** (`variables-basic.yaml`):
```yaml
# Overrides schema-defined variables
API_HOST: "external-api.example.com"
NESTED_VAR: "from-external-file"
NEW_VAR: "only-in-external"
```

**Command:**
```bash
konfigo -s base-config.yaml -S schema-basic.yaml -V variables-basic.yaml
```

**Result:** Variables from `-V` file override schema `vars` with same names.

## Variable Substitution Contexts

Variables are substituted in multiple contexts:

### 1. Configuration Values
```yaml
# In merged configuration
database:
  url: "postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}"
  timeout: "${CONNECTION_TIMEOUT}"
```

### 2. Schema Directives

**In Generators:**
```yaml
generators:
  - type: "concat"
    targetPath: "service.url"
    format: "https://{service}.${DOMAIN}:{port}"
    sources:
      service: "service.name"
      port: "service.port"
```

**In Transformations:**
```yaml
transform:
  - type: "setValue"
    path: "app.environment"
    value: "${ENVIRONMENT}"
  - type: "addKeyPrefix"
    path: "database"
    prefix: "${ENV_PREFIX}_"
```

**In Validation:**
```yaml
validate:
  - path: "app.version"
    rules:
      regex: "^${VERSION_PATTERN}$"
```

## Advanced Variable Patterns

### Conditional Variables with Defaults

```yaml
vars:
  # Environment-specific configuration
  - name: "LOG_LEVEL"
    fromEnv: "APP_LOG_LEVEL"
    defaultValue: "info"
  
  - name: "DEBUG_MODE"
    fromEnv: "DEBUG"
    defaultValue: "false"
  
  # Database configuration
  - name: "DB_SSL_MODE"
    fromEnv: "DATABASE_SSL"
    defaultValue: "disable"
  
  # Feature flags
  - name: "FEATURE_AUTH"
    fromPath: "features.authentication.enabled"
    defaultValue: "true"
```

### Complex Path Resolution

```yaml
vars:
  # Extract nested values
  - name: "REPLICA_COUNT"
    fromPath: "deployment.replicas"
    defaultValue: "2"
  
  - name: "SERVICE_VERSION"
    fromPath: "metadata.labels.version"
    defaultValue: "latest"
  
  # Combine with configuration usage
config:
  deployment:
    replicas: "${REPLICA_COUNT}"
    image: "myapp:${SERVICE_VERSION}"
```

### Variable Cascading

Variables can reference configuration values that contain other variables:

```yaml
# Input configuration
app:
  name: "user-service"
  version: "v1.2.3"

# Schema variables
vars:
  - name: "APP_NAME"
    fromPath: "app.name"
  - name: "APP_VERSION"
    fromPath: "app.version"

# Generator using both variables
generators:
  - type: "concat"
    targetPath: "deployment.image"
    format: "${APP_NAME}:${APP_VERSION}"
    sources: {}  # No config sources, only variables
```

## Error Handling

### Missing Required Variables

```yaml
vars:
  - name: "REQUIRED_VAR"
    fromEnv: "MISSING_ENV_VAR"
    # No defaultValue - will fail if env var not set
```

**Error:** `Failed to resolve variable REQUIRED_VAR: environment variable MISSING_ENV_VAR not found`

### Invalid Path References

```yaml
vars:
  - name: "MISSING_PATH"
    fromPath: "nonexistent.config.path"
    # No defaultValue - will fail if path not found
```

**Error:** `Failed to resolve variable MISSING_PATH: path nonexistent.config.path not found`

### Circular References

```yaml
# Avoid circular variable references
config:
  value1: "${VAR2}"
  value2: "${VAR1}"

vars:
  - name: "VAR1"
    fromPath: "value1"
  - name: "VAR2"
    fromPath: "value2"
```

## Variable Testing Strategies

### 1. Test Variable Precedence

```bash
# Test 1: Schema variables only
konfigo -s config.yaml -S schema.yaml

# Test 2: External variables override
konfigo -s config.yaml -S schema.yaml -V vars.yaml

# Test 3: Environment variables override all
KONFIGO_VAR_API_HOST=env-override konfigo -s config.yaml -S schema.yaml -V vars.yaml
```

### 2. Test Default Fallbacks

```bash
# Ensure defaults work when environment variables missing
unset SERVICE_PORT
konfigo -s config.yaml -S schema.yaml
```

### 3. Test Path Resolution

```bash
# Verify configuration paths are correctly resolved
konfigo -s config-with-target-paths.yaml -S schema-with-frompath.yaml
```

## Best Practices

1. **Use Descriptive Names**: Variable names should clearly indicate their purpose
2. **Provide Defaults**: Always include `defaultValue` for optional variables
3. **Document Sources**: Comment variable definitions to explain their sources
4. **Test All Paths**: Verify variables work with and without external sources
5. **Environment Separation**: Use `-V` files for environment-specific variables
6. **Security**: Use `KONFIGO_VAR_*` environment variables for secrets
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
