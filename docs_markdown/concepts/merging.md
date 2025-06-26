# Configuration Merging

Konfigo's configuration merging is one of its core strengths, providing intelligent deep merging capabilities that respect data types and precedence rules.

## Merge Strategy

Konfigo uses **deep merging** with intelligent precedence handling:

- **Primitive Values**: Later sources override earlier ones
- **Arrays**: Later arrays completely replace earlier ones (configurable)
- **Objects**: Deep merge with later keys overriding earlier ones
- **Null Values**: Explicitly set nulls override previous values

## Basic Merging Example

**base.yaml**:
```yaml
app:
  name: "myapp"
  port: 8080
  features:
    auth: true
    cache: false
database:
  host: "localhost"
  port: 5432
```

**prod.yaml**:
```yaml
app:
  port: 9090
  features:
    cache: true
    monitoring: true
database:
  host: "prod-db.example.com"
  ssl: true
```

**Merged Result**:
```yaml
app:
  name: "myapp"           # from base
  port: 9090              # overridden by prod
  features:
    auth: true            # from base
    cache: true           # overridden by prod
    monitoring: true      # added by prod
database:
  host: "prod-db.example.com"  # overridden by prod
  port: 5432              # from base
  ssl: true               # added by prod
```

## Source Precedence

Configuration sources are merged in the order they are specified:

```bash
# Later sources override earlier ones
konfigo -s base.yaml,env.yaml,local.yaml
```

1. `base.yaml` (lowest precedence)
2. `env.yaml` (medium precedence)
3. `local.yaml` (highest precedence)

## Environment Variable Overrides

Environment variables can override any configuration key:

```bash
# Override nested configuration keys
export KONFIGO_KEY_database.host="override-host"
export KONFIGO_KEY_app.port=3000

konfigo -s config.yaml
```

Environment overrides have the highest precedence and are applied after all file merging.

## Immutable Paths

Schema can define immutable paths that cannot be overridden by subsequent sources:

```yaml
# schema.yaml
immutable:
  - "app.name"
  - "security.keys"
  - "database.credentials.username"
```

**Important**: Environment variables (`KONFIGO_KEY_*`) can still override immutable paths, as they represent explicit runtime configuration.

## Array Merging Strategies

### Replace Strategy (Default)
Arrays from later sources completely replace earlier arrays:

```yaml
# base.yaml
servers:
  - "server1"
  - "server2"

# prod.yaml  
servers:
  - "prod-server1"
  - "prod-server2"
  - "prod-server3"

# Result: ["prod-server1", "prod-server2", "prod-server3"]
```

### Append Strategy (Future Enhancement)
Future versions may support array appending strategies for specific use cases.

## Complex Merging Examples

### Multi-Environment Setup

```yaml
# base.yaml
app:
  name: "my-service"
  logging:
    level: "info"
    format: "json"
database:
  timeout: 30
  pool:
    min: 5
    max: 20

# staging.yaml
app:
  logging:
    level: "debug"
database:
  host: "staging-db.example.com"
  pool:
    max: 10

# Result combines both with staging overrides
app:
  name: "my-service"      # from base
  logging:
    level: "debug"        # from staging  
    format: "json"        # from base
database:
  timeout: 30             # from base
  host: "staging-db.example.com"  # from staging
  pool:
    min: 5                # from base
    max: 10               # from staging
```

### Null Value Handling

Explicit null values override existing values:

```yaml
# base.yaml
feature:
  cache:
    enabled: true
    ttl: 3600

# override.yaml
feature:
  cache:
    ttl: null  # Explicitly set to null

# Result
feature:
  cache:
    enabled: true
    ttl: null
```

## Best Practices

1. **Order Sources Thoughtfully**: Place most general configurations first, specific overrides last
2. **Use Immutable Paths**: Protect critical configuration that shouldn't be overridden
3. **Leverage Environment Variables**: Use `KONFIGO_KEY_*` for runtime-specific overrides
4. **Test Merging Logic**: Verify merge results with different source combinations
5. **Document Override Patterns**: Make it clear which files override which settings

## Troubleshooting Merging Issues

### Debug Merge Results
Use verbose logging to see the merge process:

```bash
konfigo -s base.yaml,prod.yaml -v
```

### Validate Sources
Check individual source files before merging:

```bash
# Test each source individually
konfigo -s base.yaml
konfigo -s prod.yaml
```

### Check for Type Conflicts
Ensure compatible data types across sources - Konfigo will report type mismatches clearly.
