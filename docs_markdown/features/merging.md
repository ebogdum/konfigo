# Configuration Merging

Konfigo merges multiple configuration sources using well-defined precedence rules. Understanding these rules is crucial for predictable configuration management.

## Merge Order and Precedence

Configuration sources are processed in this order:

1. **Sources in CLI order** - Processed left-to-right as listed in `-s`; later sources override earlier ones. Stdin (`-`) is merged at whatever position it appears.
2. **Environment Variables** (`KONFIGO_KEY_*`) - Merged last, giving them highest precedence

### Command Line Example
```bash
konfigo -s base.json,environment.yml,local.toml
```

Merge order: `base.json` → `environment.yml` → `local.toml` → Environment variables

## Basic Merging Rules

### Object Merging
Objects are merged recursively, combining keys from all sources:

```yaml
# base.json
{
  "service": {
    "name": "my-app",
    "port": 8080,
    "timeout": 30
  }
}

# override.json  
{
  "service": {
    "port": 9090,
    "environment": "production"
  }
}

# Result after: konfigo -s base.json,override.json
{
  "service": {
    "name": "my-app",        # from base
    "port": 9090,            # overridden by override
    "timeout": 30,           # from base
    "environment": "production" # added by override
  }
}
```

### Value Overwriting
Scalar values (strings, numbers, booleans) are completely replaced:

```yaml
# base.yml
database:
  host: localhost
  port: 5432
  ssl: false

# prod.yml
database:
  host: prod-db.example.com
  ssl: true

# Result: prod values override base values
database:
  host: prod-db.example.com  # overridden
  port: 5432                 # preserved from base
  ssl: true                  # overridden
```

### Array Replacement (Default)
By default, arrays are replaced entirely, not merged:

```yaml
# base.yml
features:
  - auth
  - logging

# override.yml
features:
  - auth
  - analytics
  - monitoring

# Result: array completely replaced
features:
  - auth
  - analytics
  - monitoring
```

### Array Union (`-m` flag)
With the `-m` flag, arrays are merged by union with deduplication:

```bash
konfigo -m -s base.yml,override.yml
```

```yaml
# Result: union of both arrays, duplicates removed
features:
  - auth
  - logging
  - analytics
  - monitoring
```

::: info Performance Note
For very large arrays (over 1,000 elements per side), deduplication is automatically skipped to avoid performance degradation. In this case, elements are simply appended.
:::

## Case Sensitivity

By default, Konfigo uses case-insensitive key matching. Use `-c` flag for case-sensitive mode.

### Case-Insensitive (Default)

::: tip Key Casing
In case-insensitive mode, when keys differ only in casing, the **last source's casing wins**. For example, if file A has `DatabaseHost` and file B has `databasehost`, the final key will be `databasehost`. Enable debug logging (`-d`) to see when key casing changes during merge.
:::

```bash
# These keys are treated as the same
konfigo -s config1.json,config2.json

# config1.json
{
  "Service": {
    "Name": "app"
  }
}

# config2.json
{
  "service": {
    "name": "my-app"
  }
}

# Result: keys merged despite case differences
{
  "service": {
    "name": "my-app"
  }
}
```

### Case-Sensitive Mode
```bash
# Use -c flag for case-sensitive merging
konfigo -s config1.json,config2.json -c

# Result: keys treated as different
{
  "Service": {
    "Name": "app"
  },
  "service": {
    "name": "my-app"
  }
}
```

## Environment Variable Integration

Environment variables with `KONFIGO_KEY_` prefix override any configuration file values:

```bash
# Set environment override
export KONFIGO_KEY_service.port=9999
export KONFIGO_KEY_database.ssl=true

konfigo -s base.json,override.json

# Environment variables take highest precedence
# service.port will be 9999 regardless of file contents
# database.ssl will be true regardless of file contents
```

## Immutable Paths

The schema can define paths as immutable, preventing later sources from overriding earlier values. Immutable protection also extends to all child paths — marking `service` as immutable also protects `service.name`, `service.port`, etc.

```yaml
# schema.yml
immutable:
  - "service.name"
  - "database.host"

# base.json (loaded first)
{
  "service": {
    "name": "core-service",
    "port": 8080
  }
}

# override.json (loaded second)
{
  "service": {
    "name": "different-service",  # This will be IGNORED
    "port": 9090                  # This will override
  }
}

# Result: immutable paths protected
{
  "service": {
    "name": "core-service",  # Preserved from base
    "port": 9090             # Overridden normally
  }
}
```

**Important**: `KONFIGO_KEY_` environment variables are also subject to immutable path protection. Once a value is set at an immutable path, no source can override it.

## Real-World Example

Based on `test/merging/` test cases:

```bash
# Base configuration
# base-config.json
{
  "application": {
    "name": "my-app",
    "version": "1.0.0", 
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "ssl": false
  },
  "logging": {
    "level": "info",
    "format": "text"
  }
}

# Production overrides  
# override-prod.json
{
  "application": {
    "port": 9090,
    "environment": "production"
  },
  "database": {
    "host": "prod-db.example.com",
    "ssl": true,
    "pool": {
      "min": 5,
      "max": 20
    }
  },
  "logging": {
    "level": "warn",
    "format": "json"
  },
  "secrets": {
    "api_key": "prod-key-123"
  }
}

# Merge command
konfigo -s base-config.json,override-prod.json

# Final result
{
  "application": {
    "name": "my-app",              # from base
    "version": "1.0.0",            # from base
    "port": 9090,                  # overridden
    "environment": "production"     # added
  },
  "database": {
    "host": "prod-db.example.com", # overridden
    "port": 5432,                  # from base
    "ssl": true,                   # overridden
    "pool": {                      # added
      "min": 5,
      "max": 20
    }
  },
  "logging": {
    "level": "warn",               # overridden
    "format": "json",              # overridden
    "output": "stdout"             # from base
  },
  "features": {
    "auth": true,                  # from base
    "cache": true,                 # overridden
    "monitoring": true             # added
  },
  "secrets": {                     # added entirely
    "api_key": "prod-key-123"
  }
}
```

## Stdin Integration

Stdin data is merged after all file sources:

```bash
# Files merged first, then stdin
echo '{"service": {"debug": true}}' | konfigo -s base.json,prod.json -sj

# Stdin has higher precedence than files
# but lower than environment variables
```

## Error Handling

### File Not Found
```bash
konfigo -s base.json,missing.json
# Error: failed to stat path missing.json: no such file or directory
```

### Parse Errors
```bash
konfigo -s base.json,invalid.json
# Warning: Skipping file invalid.json due to parse error: invalid character '}'
# Processing continues with remaining files
```

## Best Practices

1. **Predictable Order**: List files from most general to most specific
2. **Environment Overrides**: Use `KONFIGO_KEY_` for runtime overrides
3. **Immutable Protection**: Define critical paths as immutable in schema
4. **Error Tolerance**: Konfigo continues processing when individual files fail
5. **Case Consistency**: Use consistent key casing or explicit case-sensitive mode

## Advanced Patterns

### Multi-Environment Setup
```bash
# Base → Environment → Local overrides
konfigo -s base.json,environments/${ENV}.json,local.json
```

### Configuration Layers
```bash
# System → Application → User → Runtime
konfigo -s system.conf,app.toml,~/.config/app.yml,runtime.env
```

### Development Workflow
```bash
# Committed → Generated → Local development
konfigo -s config.yml,secrets.yml,dev-overrides.yml
```

## Test Coverage

Merging functionality is comprehensively tested in `test/merging/`:
- Multi-format merging
- Precedence rule verification  
- Case sensitivity modes
- Error condition handling
- Complex nested object merging
- Array replacement behavior
