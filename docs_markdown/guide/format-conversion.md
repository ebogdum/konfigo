# Format Conversion

Convert configuration files between JSON, YAML, TOML, and ENV formats with Konfigo.

## Overview

Konfigo can read any supported format and output to any other format, making it easy to:
- Convert legacy .env files to structured YAML
- Transform JSON to TOML for different tools
- Generate multiple formats from a single source

## Quick Reference

```bash
# Input → Output format mappings
konfigo -s config.yaml -oj    # YAML → JSON
konfigo -s config.json -oy    # JSON → YAML  
konfigo -s config.toml -oe    # TOML → ENV
konfigo -s .env -ot           # ENV → TOML

# Multiple outputs at once
konfigo -s config.yaml -oj -oy -ot -of base  # Creates base.json, base.yaml, base.toml
```

## Supported Formats

| Format | Extensions | Input Flag | Output Flag | Notes |
|--------|------------|------------|-------------|-------|
| **JSON** | `.json`, `.jsonc` | `-sj` | `-oj` | Comments supported in JSONC |
| **YAML** | `.yaml`, `.yml` | `-sy` | `-oy` | Full YAML 1.2 support |
| **TOML** | `.toml` | `-st` | `-ot` | TOML v1.0.0 compatible |
| **ENV** | `.env`, `.envrc` | `-se` | `-oe` | Key=value pairs |

## Detailed Examples

### ENV to YAML (Legacy Migration)

**Converting legacy .env files**:

```bash
# .env
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_SSL=true
APP_NAME=my-service
APP_DEBUG=false
FEATURE_AUTH_ENABLED=true
FEATURE_CACHE_TTL=3600
```

**Command**:
```bash
konfigo -s .env -oy -of config.yaml
```

**Result**:
```yaml
APP_DEBUG: false
APP_NAME: my-service
DATABASE_HOST: localhost
DATABASE_PORT: 5432
DATABASE_SSL: true
FEATURE_AUTH_ENABLED: true
FEATURE_CACHE_TTL: 3600
```

### JSON to YAML (API to Config)

**Converting API response to configuration**:

```json
// api-config.json
{
  "service": {
    "name": "user-service",
    "port": 8080,
    "timeout": 30
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "pool": {
      "min": 5,
      "max": 20
    }
  }
}
```

**Command**:
```bash
konfigo -s api-config.json -oy -of service.yaml
```

**Result**:
```yaml
database:
  host: localhost
  port: 5432
  pool:
    max: 20
    min: 5
service:
  name: user-service
  port: 8080
  timeout: 30
```

### YAML to TOML (Tool Integration)

**Converting for tools that prefer TOML**:

```yaml
# app.yaml
[package]
name = "my-app"
version = "1.0.0"

[dependencies]
database = "postgresql"
cache = "redis"

[server]
host = "0.0.0.0"
port = 8080
workers = 4
```

**Command**:
```bash
konfigo -s app.yaml -ot -of app.toml
```

### Multiple Format Output

**Generate configurations for different tools**:

```bash
# Create JSON for APIs, YAML for Kubernetes, TOML for Rust tools
konfigo -s base-config.yaml -oj -oy -ot -of deployment

# Creates:
# - deployment.json
# - deployment.yaml  
# - deployment.toml
```

## Format-Specific Options

### JSON Output Options

```bash
# Pretty-printed JSON (default)
konfigo -s config.yaml -oj

# Compact JSON (no whitespace)
konfigo -s config.yaml -oj --json-compact

# Custom indentation
konfigo -s config.yaml -oj --json-indent=4
```

### YAML Output Options

```bash
# Block style (default)
konfigo -s config.json -oy

# Flow style
konfigo -s config.json -oy --yaml-flow

# Custom indentation
konfigo -s config.json -oy --yaml-indent=4
```

### ENV Output Options

```bash
# Standard format
konfigo -s config.yaml -oe

# Quoted values
konfigo -s config.yaml -oe --env-quote

# Prefix for all keys
konfigo -s config.yaml -oe --env-prefix=APP_
```

## Advanced Conversion Patterns

### Nested Structure to Flat ENV

**Converting nested YAML to flat ENV**:

```yaml
# nested.yaml
database:
  primary:
    host: "db1.company.com"
    port: 5432
  replica:
    host: "db2.company.com"
    port: 5432
cache:
  redis:
    host: "redis.company.com"
    port: 6379
```

**Command with transformation**:
```bash
konfigo -s nested.yaml -S flatten.schema.yaml -oe
```

**Schema for flattening**:
```yaml
# flatten.schema.yaml
transforms:
  - path: "*"
    flattenKeys: true
    separator: "_"
```

**Result**:
```bash
DATABASE_PRIMARY_HOST=db1.company.com
DATABASE_PRIMARY_PORT=5432
DATABASE_REPLICA_HOST=db2.company.com
DATABASE_REPLICA_PORT=5432
CACHE_REDIS_HOST=redis.company.com
CACHE_REDIS_PORT=6379
```

### Format Conversion with Validation

**Ensure data integrity during conversion**:

```bash
# Convert with validation
konfigo -s legacy.env -S validation.schema.yaml -oy -of validated.yaml
```

```yaml
# validation.schema.yaml
validation:
  - path: "DATABASE_HOST"
    required: true
    type: "string"
  - path: "DATABASE_PORT"
    required: true
    type: "number"
    min: 1
    max: 65535
```

### Batch Format Conversion

**Convert multiple files at once**:

```bash
# Convert all YAML files in a directory to JSON
for file in configs/*.yaml; do
  name=$(basename "$file" .yaml)
  konfigo -s "$file" -oj -of "json/$name.json"
done

# Or using find
find configs/ -name "*.yaml" -exec bash -c '
  konfigo -s "$1" -oj -of "json/$(basename "$1" .yaml).json"
' _ {} \;
```

## Common Use Cases

### 1. **Legacy System Modernization**
```bash
# Convert old .env files to structured YAML
konfigo -s legacy/.env -oy -of modern/config.yaml
```

### 2. **Multi-Tool Deployment**
```bash
# Generate configs for different deployment tools
konfigo -s app.yaml -oj -ot -of deploy/app  # Kubernetes (YAML), Terraform (JSON), Rust tools (TOML)
```

### 3. **API Integration**
```bash
# Convert API responses to local config format
curl -s api.example.com/config | konfigo -sj -oy -of local-config.yaml
```

### 4. **Configuration Standardization**
```bash
# Standardize team configurations to YAML
konfigo -s team-configs/*.{json,toml,env} -oy -of standardized.yaml
```

## Troubleshooting

### Format Detection Issues

If Konfigo can't detect the format automatically:

```bash
# Explicitly specify input format
konfigo -sj -s config.txt     # Treat as JSON
konfigo -sy -s data.conf      # Treat as YAML
konfigo -se -s variables.txt  # Treat as ENV
```

### Invalid Character Handling

**ENV format limitations**:
- Keys cannot contain spaces or special characters
- Use transformation to clean keys:

```yaml
# clean-keys.schema.yaml
transforms:
  - path: "*"
    renameKey:
      pattern: "[^A-Z0-9_]"
      replacement: "_"
```

### Large File Performance

**For large configuration files**:
```bash
# Use streaming for large files
konfigo -s large-config.json -oy --stream -of output.yaml

# Process in chunks
split -l 1000 large.json chunk_
for chunk in chunk_*; do
  konfigo -s "$chunk" -oy -of "output_$chunk.yaml"
done
```

## Best Practices

### 1. **Validate After Conversion**
Always verify the converted output:
```bash
# Convert and validate
konfigo -s source.env -S validation.schema.yaml -oy -of target.yaml
```

### 2. **Use Consistent Naming**
Establish patterns for converted files:
```bash
# Environment-specific naming
konfigo -s base.yaml -oj -of "configs/$(date +%Y%m%d)-config.json"
```

### 3. **Preserve Metadata**
Add conversion metadata:
```yaml
# Add to schema
generation:
  - path: "_metadata.converted_from"
    setValue: "legacy.env"
  - path: "_metadata.converted_at"
    setValue: "${TIMESTAMP}"
```

### 4. **Test Conversions**
Always test converted configurations:
```bash
# Test the converted config
konfigo -s converted.yaml --validate-only -S test.schema.yaml
```

## Next Steps

- **[Merging Configurations](./merging.md)** - Combine multiple files during conversion
- **[Environment Variables](./environment-variables.md)** - Add runtime values during conversion
- **[Schema Transformation](../schema/transformation.md)** - Advanced data transformation during conversion

Format conversion is often the first step in configuration management. Master this, then explore how to combine it with merging and validation for powerful configuration workflows!
