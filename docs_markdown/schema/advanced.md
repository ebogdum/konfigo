# Schema: Advanced Features

Advanced schema features provide fine-grained control over configuration processing, including immutable fields, input/output schema validation, and complex processing patterns.

## Immutable Fields

The `immutable` directive protects configuration paths from being overwritten during merging. Once set by an earlier source, immutable paths cannot be modified by later sources.

### Structure

```yaml
immutable:
  - "path.to.protected.field"
  - "another.immutable.path"
```

### Behavior

1. **Source Order Matters**: First source to set an immutable path wins
2. **Complete Protection**: Later sources cannot override immutable values
3. **Environment Override**: `KONFIGO_KEY_*` environment variables **can** still override immutable paths
4. **Deep Path Protection**: Protects specific nested paths, not entire objects

### Examples from Tests

**Schema with Immutable Paths:**
```yaml
immutable:
  - "service.name"
  - "database.port"
  - "security.apiKey"

# Rest of schema...
generators:
  - type: "concat"
    targetPath: "service.url"
    format: "https://{name}.example.com"
    sources:
      name: "service.name"
```

**Source Files:**

**`01-base.yaml` (loaded first):**
```yaml
service:
  name: "user-service"
  version: "1.0.0"
database:
  host: "localhost"
  port: 5432
security:
  apiKey: "base-secret-key"
```

**`02-override.yaml` (loaded second):**
```yaml
service:
  name: "overridden-service"  # ← Will be IGNORED (immutable)
  version: "2.0.0"            # ← Will be applied (not immutable)
database:
  host: "prod-db.example.com"
  port: 9999                  # ← Will be IGNORED (immutable)
security:
  apiKey: "new-secret"        # ← Will be IGNORED (immutable)
  protocol: "https"           # ← Will be applied (not immutable)
```

**Final Result:**
```yaml
service:
  name: "user-service"        # ← Protected by immutable
  version: "2.0.0"            # ← Updated by second source
  url: "https://user-service.example.com"  # ← Generated using protected name
database:
  host: "prod-db.example.com"
  port: 5432                  # ← Protected by immutable
security:
  apiKey: "base-secret-key"   # ← Protected by immutable
  protocol: "https"           # ← Added by second source
```

### Environment Variable Override

Even immutable fields can be overridden by environment variables:

```bash
# This WILL override the immutable service.name
export KONFIGO_KEY_service__name=env-override-service

konfigo -s 01-base.yaml -s 02-override.yaml -S schema.yaml
```

**Result:**
```yaml
service:
  name: "env-override-service"  # ← Environment variable wins
  # ... rest unchanged
```

## Input Schema Validation

The `inputSchema` directive validates merged configuration structure **before** any schema processing (variables, generators, transformations).

### Structure

```yaml
inputSchema:
  path: "path/to/input-schema.yaml"
  strict: false  # Optional, default: false
```

### Fields

- **`path`** (Required): Path to external schema file defining expected input structure
- **`strict`** (Optional, default: `false`): 
  - `false`: Extra keys allowed in input
  - `true`: Input must contain **only** keys defined in schema

### Example

**Input Schema Definition** (`input-structure.yaml`):
```yaml
# Defines expected structure after merging sources
service:
  name: ""
  port: 0
database:
  host: ""
  credentials:
    username: ""
```

**Main Schema** (`schema.yaml`):
```yaml
inputSchema:
  path: "./input-structure.yaml"
  strict: false

# Continue with normal processing
vars:
  - name: "DB_PASSWORD"
    fromEnv: "DATABASE_PASSWORD"
    defaultValue: "default-pass"

generators:
  - type: "concat"
    targetPath: "database.connectionString"
    format: "postgresql://{user}:${DB_PASSWORD}@{host}:5432/app"
    sources:
      user: "database.credentials.username"
      host: "database.host"
```

**Valid Input** (passes validation):
```yaml
service:
  name: "api-service"
  port: 8080
  # Extra field allowed when strict: false
  environment: "production"
database:
  host: "db.example.com"
  credentials:
    username: "api_user"
```

**Invalid Input** (fails validation):
```yaml
service:
  # Missing required port field
  name: "api-service"
database:
  # Missing host field
  credentials:
    username: "api_user"
```

## Output Schema Filtering

The `outputSchema` directive filters final configuration to include only specified keys, creating clean, controlled output.

### Structure

```yaml
outputSchema:
  path: "path/to/output-schema.yaml"
  strict: false  # Optional, default: false
```

### Fields

- **`path`** (Required): Path to external schema file defining output structure
- **`strict`** (Optional, default: `false`):
  - `false`: Include only keys present in output schema, ignore extras
  - `true`: Final output must exactly match output schema structure

### Example

**Output Schema Definition** (`public-api.yaml`):
```yaml
# Define what should be included in final output
api:
  endpoint: ""
  version: ""
database:
  host: ""
# Note: credentials intentionally excluded
```

**Main Schema** (`schema.yaml`):
```yaml
outputSchema:
  path: "./public-api.yaml"
  strict: false

# Processing that creates more data than we want to expose
vars:
  - name: "API_VERSION"
    value: "v1.2.3"

generators:
  - type: "concat"
    targetPath: "api.endpoint"
    format: "https://{host}/api/{version}"
    sources:
      host: "service.host"
      version: "api.version"
  
  # Generate internal fields that won't appear in output
  - type: "concat"
    targetPath: "internal.buildId"
    format: "build-{timestamp}"
    sources:
      timestamp: "build.timestamp"
```

**Full Processed Configuration** (before output filtering):
```yaml
api:
  endpoint: "https://api.example.com/api/v1.2.3"
  version: "v1.2.3"
  host: "api.example.com"  # Extra field
database:
  host: "db.example.com"
  credentials:
    username: "secret_user"  # Sensitive data
    password: "secret_pass"
service:
  host: "api.example.com"
internal:
  buildId: "build-2024-06-26T10:30:00Z"
  debug: true
```

**Final Filtered Output:**
```yaml
api:
  endpoint: "https://api.example.com/api/v1.2.3"
  version: "v1.2.3"
  # api.host excluded (not in output schema)
database:
  host: "db.example.com"
  # credentials excluded (not in output schema)
# service and internal sections completely excluded
```

## Strict Mode Behavior

When `strict: true` is used with input or output schemas:

### Input Schema Strict Mode

```yaml
inputSchema:
  path: "./strict-input.yaml"
  strict: true
```

**Strict Input Schema** (`strict-input.yaml`):
```yaml
service:
  name: ""
  port: 0
```

**Valid Input:**
```yaml
service:
  name: "api"
  port: 8080
# No extra keys allowed
```

**Invalid Input:**
```yaml
service:
  name: "api"
  port: 8080
  environment: "prod"  # ← ERROR: Extra key not in schema
```

### Output Schema Strict Mode

```yaml
outputSchema:
  path: "./strict-output.yaml"
  strict: true
```

Requires processed configuration to exactly match output schema structure - no missing keys, no extra keys.

## Combined Advanced Features

**Complete Advanced Schema Example:**
```yaml
# Protect critical configuration
immutable:
  - "service.name"
  - "database.credentials.username"

# Validate input structure
inputSchema:
  path: "./schemas/input-validation.yaml"
  strict: false

# Define clean output
outputSchema:
  path: "./schemas/public-output.yaml"
  strict: false

# Standard processing
vars:
  - name: "ENVIRONMENT"
    fromEnv: "NODE_ENV"
    defaultValue: "development"

generators:
  - type: "concat"
    targetPath: "service.url"
    format: "https://{name}-${ENVIRONMENT}.example.com"
    sources:
      name: "service.name"  # Protected by immutable

transform:
  - type: "setValue"
    path: "metadata.processed"
    value: true

validate:
  - path: "service.url"
    rules:
      required: true
      type: "string"
      regex: "^https://"
```

This provides comprehensive control over configuration processing with validation, protection, and filtering.

## Best Practices

1. **Immutable Core Values**: Protect essential configuration that shouldn't change
2. **Input Validation**: Verify merged configuration meets expectations before processing
3. **Output Filtering**: Create clean, secure output by excluding internal/sensitive data
4. **Gradual Adoption**: Start with basic features and add advanced controls as needed
5. **Schema Evolution**: Use separate schema files for input/output validation to enable independent evolution
