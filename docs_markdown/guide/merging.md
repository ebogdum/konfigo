# Merging Configurations

One of Konfigo's most powerful features is its ability to intelligently merge multiple configuration sources. This guide explains how merging works and how to use it effectively.

## Quick Start

Merge two configuration files:
```bash
konfigo -s base.yaml,prod.yaml -of result.json
```

## How Merging Works

Konfigo uses **deep merging** with clear precedence rules:

1. **Environment Variables** (`KONFIGO_KEY_*`) - Highest precedence
2. **Later Source Files** (rightmost in `-s` list)
3. **Earlier Source Files** (leftmost in `-s` list)
4. **Stdin Input** (lowest precedence)

### Example: Three-Way Merge

```bash
konfigo -s base.yaml,env.yaml,local.yaml
```

**Processing order**: `base.yaml` → `env.yaml` → `local.yaml` → Environment variables

## Merging Rules by Data Type

### Objects: Deep Merge
Objects are merged recursively, combining all keys:

```yaml
# base.yaml
app:
  name: "myapp"
  port: 8080
  features:
    auth: true

# prod.yaml  
app:
  port: 9090
  features:
    monitoring: true
```

**Result**:
```yaml
app:
  name: "myapp"        # from base
  port: 9090           # overridden by prod
  features:
    auth: true         # from base
    monitoring: true   # added by prod
```

### Scalars: Complete Replacement
Strings, numbers, and booleans are completely replaced:

```yaml
# base.yaml
database:
  host: "localhost"
  port: 5432
  ssl: false

# prod.yaml
database:
  host: "prod-db.com"
  ssl: true

# Result
database:
  host: "prod-db.com"  # replaced
  port: 5432           # preserved
  ssl: true            # replaced
```

### Arrays: Complete Replacement
Arrays are replaced entirely, not merged:

```yaml
# base.yaml
tags: ["app", "service"]

# override.yaml
tags: ["app", "production", "critical"]

# Result
tags: ["app", "production", "critical"]  # completely replaced
```

## Environment Variable Overrides

Override any configuration value using environment variables:

```bash
# Override nested keys using dot notation
export KONFIGO_KEY_database.host="override-host"
export KONFIGO_KEY_app.port=3000
export KONFIGO_KEY_features.auth=false

konfigo -s config.yaml
```

**Key format**: `KONFIGO_KEY_<path.to.key>=<value>`

Environment variables always have the highest precedence.

## Protected Paths (Immutable)

Use schemas to protect critical configuration paths:

```yaml
# schema.yaml
immutable:
  - "app.name"
  - "security.keys"
  - "database.credentials"

# These paths cannot be overridden by later sources
# (but environment variables can still override them)
```

## Practical Examples

### Multi-Environment Setup

```bash
# Base configuration + environment-specific overrides
konfigo -s base.yaml,environments/prod.yaml -of prod-config.json

# With environment variables for runtime customization
KONFIGO_KEY_app.port=8080 konfigo -s base.yaml,prod.yaml
```

### Configuration Layers

```bash
# Layer configurations: defaults → environment → local customization
konfigo -s defaults.json,env/production.yaml,local.toml
```

### Format Conversion During Merge

```bash
# Merge YAML and JSON, output as TOML
konfigo -s config.yaml,overrides.json -ot -of final.toml
```

## Best Practices

### 1. **Organize by Precedence**
Place more general configurations first, specific ones last:
```bash
konfigo -s defaults.yaml,environment.yaml,local.yaml
```

### 2. **Use Environment Variables for Runtime**
Perfect for containerized deployments:
```bash
# In Docker/Kubernetes
KONFIGO_KEY_database.host=$DB_HOST konfigo -s config.yaml
```

### 3. **Protect Critical Settings**
Use immutable paths for security-sensitive configuration:
```yaml
immutable:
  - "security.apiKeys"
  - "database.credentials"
```

### 4. **Test Your Merges**
Always verify the merged result:
```bash
# Output to screen to verify before saving
konfigo -s base.yaml,prod.yaml
```

## Best Practices

### ✅ Do's

**Structure source files by purpose**
```bash
# Good: Clear hierarchy and responsibility
konfigo -s defaults.yaml,env/prod.yaml,local-overrides.yaml
```

**Use consistent naming patterns**
```yaml
# Good: Predictable structure
config/
  ├── base.yaml          # Common settings
  ├── env/
  │   ├── dev.yaml       # Development overrides
  │   ├── staging.yaml   # Staging overrides
  │   └── prod.yaml      # Production overrides
  └── local.yaml         # Developer-specific (git-ignored)
```

**Keep environment variables for runtime-only values**
```bash
# Good: Runtime secrets and deployment-specific values
KONFIGO_KEY_database.password=$DB_PASSWORD \
KONFIGO_KEY_app.instanceId=$INSTANCE_ID \
konfigo -s base.yaml,prod.yaml
```

**Test your merge logic**
```bash
# Good: Verify before deployment
konfigo -s base.yaml,prod.yaml --dry-run -v
```

### ❌ Don'ts

**Don't rely on file order for critical logic**
```bash
# Bad: Unclear which file wins
konfigo -s prod.yaml,base.yaml,override.yaml
```

**Don't mix configuration concerns**
```yaml
# Bad: Infrastructure mixed with application config
database:
  host: "prod-db.com"
kubernetes:
  replicas: 3
app:
  port: 8080
```

**Don't put secrets in source files**
```yaml
# Bad: Secrets in version control
database:
  password: "supersecret123"  # Use environment variables instead!
```

**Don't create deep nesting without purpose**
```yaml
# Bad: Unnecessarily complex structure
app:
  config:
    settings:
      database:
        connection:
          primary:
            host: "localhost"  # Too deep!

# Good: Reasonable nesting
database:
  host: "localhost"
```

### ⚠️ Common Pitfalls

::: warning Array Replacement
Arrays are completely replaced, not merged. If you need array merging, consider using objects with keys instead:

```yaml
# Instead of this:
services: ["web", "api", "worker"]

# Consider this:
services:
  web: { enabled: true }
  api: { enabled: true }
  worker: { enabled: false }
```
:::

::: warning Case Sensitivity
Configuration keys are case-sensitive by default. Be consistent:

```yaml
# These are different keys:
Database: { host: "..." }
database: { host: "..." }
```
:::

::: warning Environment Variable Precedence
Environment variables always win. Don't rely on file-based overrides for values that might be set in the environment:

```bash
# This file setting will be ignored:
# config.yaml: { port: 8080 }
KONFIGO_KEY_port=9090 konfigo -s config.yaml  # port will be 9090
```
:::

## Common Patterns

### Environment-Specific Deployment
```bash
# Development
konfigo -s base.yaml,env/dev.yaml -of dev-config.json

# Production  
konfigo -s base.yaml,env/prod.yaml -of prod-config.json
```

### Legacy System Integration
```bash
# Merge legacy .env with modern YAML
konfigo -s legacy.env,modern.yaml -of integrated.json
```

### Configuration Validation
```bash
# Merge and validate against schema
konfigo -s base.yaml,overrides.yaml -S validation.schema.yaml
```

## Troubleshooting

### Unexpected Overrides
Use verbose mode to see merge order:
```bash
konfigo -v -s file1.yaml,file2.yaml
```

### Array Merging Issues
Arrays are replaced, not merged. If you need array merging, use transformation schemas.

### Case Sensitivity
By default, merging is case-sensitive. Use `-c` flag for case-insensitive merging:
```bash
konfigo -c -s config1.yaml,config2.yaml
```

## Next Steps

- **[Environment Variables Guide](./environment-variables.md)** - Learn more about environment integration
- **[Schema Validation](../schema/validation.md)** - Add validation to your merging process
- **[Transformation](../schema/transformation.md)** - Modify data during merge
