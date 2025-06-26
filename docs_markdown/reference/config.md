# Configuration Options

Comprehensive reference for all Konfigo configuration options, environment variables, and behavioral settings.

## Command-Line Configuration

### Source File Discovery

Konfigo uses the following order to discover configuration files:

1. **Explicit paths** specified with `-s` flag
2. **Current directory** search for common config files
3. **Recursive discovery** with `-r` flag (searches subdirectories)
4. **Environment variable paths** from `KONFIGO_CONFIG_PATH`

### Default File Patterns

When no explicit sources are provided, Konfigo searches for:

```
config.{json,yaml,yml,toml,env}
app.{json,yaml,yml,toml,env}
konfigo.{json,yaml,yml,toml,env}
.env
```

### Case Sensitivity

| Flag | Behavior | Use Case |
|------|----------|----------|
| Default | Case-insensitive key matching | Most common usage |
| `-c` | Case-sensitive key matching | Strict configuration requirements |

**Example**:
```bash
# These keys are treated as the same (default)
{"Database": {...}} + {"database": {...}} = merged

# With -c flag, these remain separate
{"Database": {...}} + {"database": {...}} = both preserved
```

## Environment Variable Configuration

### Configuration Overrides

Override any configuration key using the `KONFIGO_KEY_` prefix:

```bash
# Override nested configuration
export KONFIGO_KEY_app.port=8080
export KONFIGO_KEY_database.host=prod-db.company.com
export KONFIGO_KEY_features.auth.enabled=true

# Override array elements (replaces entire array)
export KONFIGO_KEY_tags='["production", "critical"]'
```

### Schema Variables

Define variables for schema processing:

```bash
# High-priority schema variables
export KONFIGO_VAR_ENVIRONMENT=production
export KONFIGO_VAR_DATABASE_PASSWORD=secret123
export KONFIGO_VAR_API_ENDPOINT=https://api.company.com
```

### System Configuration

Control Konfigo's behavior:

```bash
# Logging level
export KONFIGO_LOG_LEVEL=DEBUG    # ERROR, WARN, INFO, DEBUG

# Default search paths
export KONFIGO_CONFIG_PATH=/etc/konfigo:/usr/local/konfigo

# Disable colored output
export NO_COLOR=1

# Force specific locale for parsing
export LC_ALL=C
```

## Schema Configuration

### Schema File Structure

```yaml
apiVersion: v1

# Schema metadata
metadata:
  name: "my-config-schema"
  version: "1.0.0"
  description: "Production configuration schema"

# Input validation (before processing)
inputSchema:
  path: "./schemas/input-validation.json"
  strict: false  # Allow extra keys not in schema

# Output filtering (after processing)
outputSchema:
  path: "./schemas/output-structure.json"
  strict: true   # Only include keys from schema

# Protected configuration paths
immutable:
  - "app.name"
  - "security.keys"
  - "database.credentials.username"

# Variable definitions
vars:
  - name: "ENVIRONMENT"
    value: "production"
    description: "Deployment environment"
  
  - name: "DATABASE_HOST"
    value: "prod-db.company.com"
    description: "Primary database hostname"

# Data generation
generators:
  - type: "uuid"
    targetPath: "app.instanceId"
    
  - type: "timestamp"
    targetPath: "metadata.generated"
    format: "2006-01-02T15:04:05Z07:00"

# Data transformation
transform:
  - type: "renameKey"
    from: "old.config.path"
    to: "new.config.path"
    
  - type: "setValue"
    path: "metadata.processed"
    value: true

# Data validation
validate:
  - path: "app.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
      
  - path: "environment"
    rules:
      required: true
      type: "string"
      enum: ["development", "staging", "production"]
```

## Processing Configuration

### Merge Behavior

#### Object Merging
- **Deep merge**: Objects are combined recursively
- **Key precedence**: Later sources override earlier sources
- **Null handling**: Explicit null values override previous values

#### Array Handling
- **Complete replacement**: Arrays are replaced entirely, not merged
- **No element merging**: Individual array elements are not combined

#### Primitive Values
- **Direct replacement**: Strings, numbers, booleans are replaced completely

### Precedence Order

1. **Environment variables** (`KONFIGO_KEY_*`) - Highest precedence
2. **Later source files** (rightmost in `-s` list)
3. **Earlier source files** (leftmost in `-s` list)
4. **Stdin input** - Lowest precedence

### Schema Processing Order

1. **Input validation** (if `inputSchema` specified)
2. **Variable substitution** (from `vars` and `KONFIGO_VAR_*`)
3. **Data generation** (from `generators`)
4. **Data transformation** (from `transform`)
5. **Data validation** (from `validate`)
6. **Output filtering** (if `outputSchema` specified)

## Output Configuration

### Format Detection

Output format is determined by:

1. **Explicit format flags** (`-oj`, `-oy`, `-ot`, `-oe`)
2. **Output file extension** (when using `-of`)
3. **Default to JSON** (when outputting to stdout)

### Format-Specific Options

#### JSON Output
- **Pretty printing**: Enabled by default
- **Compact mode**: Not available (use external tools like `jq -c`)

#### YAML Output
- **Indentation**: 2 spaces (standard)
- **Flow style**: Block style for readability
- **Comments**: Not preserved in conversion

#### TOML Output
- **Table organization**: Automatic based on structure
- **Array formatting**: Multi-line for readability

#### ENV Output
- **Key flattening**: Nested keys use dot notation
- **Value quoting**: Automatic based on content
- **Comment preservation**: Not supported

## Performance Configuration

### Memory Management

```bash
# For large configurations
export KONFIGO_MAX_MEMORY=2G        # Not implemented yet
export KONFIGO_STREAMING_THRESHOLD=100M  # Not implemented yet
```

### Processing Limits

- **Maximum source files**: No hard limit (limited by system resources)
- **Maximum file size**: No hard limit (memory dependent)
- **Maximum nesting depth**: 100 levels (prevents infinite recursion)
- **Maximum variable substitutions**: 1000 per key (prevents infinite loops)

## Security Configuration

### Sensitive Data Handling

- **Never log sensitive values**: Use `-d` flag carefully in production
- **Environment variable security**: Ensure `KONFIGO_KEY_*` and `KONFIGO_VAR_*` are properly secured
- **File permissions**: Konfigo respects system file permissions
- **Output sanitization**: Sensitive data is not filtered automatically

### Best Practices

1. **Use environment variables for secrets**:
   ```bash
   export KONFIGO_VAR_DB_PASSWORD=$SECRET_PASSWORD
   ```

2. **Restrict file permissions**:
   ```bash
   chmod 600 sensitive-config.yaml
   ```

3. **Validate input sources**:
   ```bash
   # Use schemas to ensure expected structure
   konfigo -s config.yaml -S validation-schema.yaml
   ```

## Debugging Configuration

### Logging Levels

| Level | Description | Output |
|-------|-------------|--------|
| ERROR | Errors only | Critical failures |
| WARN | Warnings and errors | Potential issues |
| INFO | Informational messages | Processing steps |
| DEBUG | Detailed debug info | Everything |

### Debug Output Examples

```bash
# Basic verbose output
konfigo -v -s config.yaml

# Full debug information
konfigo -d -s config1.yaml,config2.yaml -S schema.yaml

# Specific debug areas
export KONFIGO_DEBUG_MERGE=true
export KONFIGO_DEBUG_SCHEMA=true
export KONFIGO_DEBUG_VALIDATION=true
```

### Common Debug Scenarios

- **Merge conflicts**: Use `-v` to see merge order and precedence
- **Schema issues**: Use `-d` to see variable substitution and validation steps
- **File discovery**: Use `-v` to see which files are found and loaded
- **Environment overrides**: Use `-d` to see which `KONFIGO_KEY_*` variables are applied

This configuration reference provides complete control over Konfigo's behavior for any use case, from simple file merging to complex enterprise configuration management.
