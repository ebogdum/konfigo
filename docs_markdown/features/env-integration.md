# Environment Integration

Konfigo provides powerful environment variable integration for configuration management, supporting both direct key overrides and schema-level variable resolution.

## Environment Variable Types

Konfigo handles two distinct types of environment variables:

1. **Configuration Overrides** (`KONFIGO_KEY_*`) - Direct configuration value overrides
2. **Schema Variables** (`KONFIGO_VAR_*`) - Variables for schema processing

## Configuration Overrides (`KONFIGO_KEY_`)

Environment variables with the `KONFIGO_KEY_` prefix directly override configuration values at any depth using dot notation.

### Basic Usage
```bash
# Override top-level values
export KONFIGO_KEY_port=9090
export KONFIGO_KEY_debug=true

# Override nested values using dot notation
export KONFIGO_KEY_service.name=my-service
export KONFIGO_KEY_database.host=prod-db.example.com
export KONFIGO_KEY_logging.level=debug

konfigo -s base.json
```

### Deep Path Overrides
```bash
# Override deeply nested configuration
export KONFIGO_KEY_api.auth.oauth.client_id=new-client-id
export KONFIGO_KEY_features.payments.stripe.secret_key=sk_live_xxx

# These override any values in configuration files
konfigo -s config.yml
```

### Precedence Rules
`KONFIGO_KEY_` variables have the **highest precedence**, overriding:
- All configuration files
- Variables files (`-V`)
- Schema-defined values
- Even immutable paths defined in schema

```bash
# This will override immutable paths
export KONFIGO_KEY_service.name=override-service

# Even with immutable schema definition:
# immutable:
#   - "service.name"
```

## Schema Variables (`KONFIGO_VAR_`)

Environment variables with the `KONFIGO_VAR_` prefix provide values for schema variable resolution.

### Variable Resolution
```bash
# Set schema variables
export KONFIGO_VAR_ENVIRONMENT=production
export KONFIGO_VAR_DATABASE_PASSWORD=secure-password
export KONFIGO_VAR_API_KEY=secret-api-key

# Schema can reference these variables
konfigo -s config.yml -S schema.yml
```

### Schema Variable Usage
```yaml
# schema.yml
vars:
  - name: "DATABASE_HOST"
    fromEnv: "DB_HOST"
    defaultValue: "localhost"
  
  - name: "API_ENDPOINT"
    value: "https://api.${ENVIRONMENT}.example.com"

generators:
  - type: "concat"
    targetPath: "database.url"
    format: "postgresql://${DATABASE_HOST}:5432/mydb"

transform:
  - type: "setValue"
    path: "api.key"
    value: "${API_KEY}"
```

## Real-World Examples

Based on `test/env-integration/` test cases:

### Development Environment
```bash
# Development overrides
export KONFIGO_KEY_service.environment=development
export KONFIGO_KEY_database.host=localhost
export KONFIGO_KEY_logging.level=debug
export KONFIGO_KEY_features.debug_mode=true

# Schema variables for development
export KONFIGO_VAR_API_BASE_URL=http://localhost:3000
export KONFIGO_VAR_CACHE_SIZE=100

konfigo -s base.json -S dev-schema.yml
```

### Production Environment
```bash
# Production overrides (often set by orchestration)
export KONFIGO_KEY_service.environment=production
export KONFIGO_KEY_database.host=prod-cluster.internal
export KONFIGO_KEY_database.pool.max=50
export KONFIGO_KEY_logging.level=warn

# Sensitive schema variables
export KONFIGO_VAR_DATABASE_PASSWORD=$(vault kv get -field=password secret/db)
export KONFIGO_VAR_API_SECRET=$(vault kv get -field=secret secret/api)

konfigo -s base.json,prod-overrides.json -S prod-schema.yml
```

### Container Deployment
```bash
# Kubernetes/Docker environment variables
export KONFIGO_KEY_service.name=${SERVICE_NAME}
export KONFIGO_KEY_service.port=${PORT}
export KONFIGO_KEY_database.host=${DB_HOST}
export KONFIGO_KEY_database.port=${DB_PORT}

# Schema variables from secrets
export KONFIGO_VAR_JWT_SECRET=${JWT_SECRET}
export KONFIGO_VAR_ENCRYPTION_KEY=${ENCRYPTION_KEY}

# Generate final configuration
konfigo -s /config/base.yml -S /config/schema.yml > /app/config.json
```

## Type Conversion

Environment variables are strings by default. Konfigo performs automatic type conversion:

```bash
# String values
export KONFIGO_KEY_service.name=my-app

# Boolean values (case-insensitive)
export KONFIGO_KEY_debug=true
export KONFIGO_KEY_ssl.enabled=false

# Numeric values
export KONFIGO_KEY_port=8080
export KONFIGO_KEY_timeout=30.5

# Arrays (comma-separated)
export KONFIGO_KEY_allowed_hosts=localhost,127.0.0.1,example.com
```

Result:
```json
{
  "service": {
    "name": "my-app"
  },
  "debug": true,
  "ssl": {
    "enabled": false
  },
  "port": 8080,
  "timeout": 30.5,
  "allowed_hosts": ["localhost", "127.0.0.1", "example.com"]
}
```

## Complex Path Handling

### Array Index Access
```bash
# Override specific array elements
export KONFIGO_KEY_servers.0.host=server1.example.com
export KONFIGO_KEY_servers.1.host=server2.example.com

# Add new array elements
export KONFIGO_KEY_allowed_origins.0=https://app.example.com
export KONFIGO_KEY_allowed_origins.1=https://admin.example.com
```

### Special Characters in Keys
```bash
# Keys with hyphens or special characters
export KONFIGO_KEY_api-gateway.timeout=30
export KONFIGO_KEY_oauth2.client-id=my-client
```

## Variable Resolution Priority

For schema variables, resolution follows this priority order:

1. **`KONFIGO_VAR_*` Environment Variables** (Highest)
2. **Variables File** (`-V` flag)
3. **Schema `vars` Definitions**
4. **Default Values**

```yaml
# schema.yml
vars:
  - name: "DATABASE_HOST"
    fromEnv: "DB_HOST"        # Looks for KONFIGO_VAR_DB_HOST
    defaultValue: "localhost"  # Used if not found

# Resolution order:
# 1. KONFIGO_VAR_DB_HOST (if set)
# 2. DB_HOST environment variable (if set)
# 3. "localhost" default value
```

## Error Handling

### Invalid Paths
```bash
export KONFIGO_KEY_invalid..path=value
# Warning: Invalid path 'invalid..path' ignored
```

### Type Conversion Errors
```bash
export KONFIGO_KEY_port=invalid-number
# Warning: Could not convert 'invalid-number' to number for port, using string
```

## Best Practices

### Development
1. **Local Overrides**: Use `KONFIGO_KEY_` for quick development changes
2. **Environment Files**: Create `.env` files for team consistency
3. **Documentation**: Document expected environment variables

### Production
1. **Secrets Management**: Use `KONFIGO_VAR_` for sensitive schema variables
2. **Orchestration**: Set `KONFIGO_KEY_` values in deployment scripts
3. **Validation**: Verify environment variables before deployment

### Security
1. **Sensitive Data**: Use schema variables for secrets, not direct overrides
2. **Variable Scoping**: Limit variable access to necessary components
3. **Audit Trail**: Log which environment variables are used

## Integration Examples

### Docker Compose
```yaml
# docker-compose.yml
services:
  app:
    environment:
      - KONFIGO_KEY_service.environment=production
      - KONFIGO_KEY_database.host=db
      - KONFIGO_VAR_DATABASE_PASSWORD=secret123
    command: |
      sh -c "
        konfigo -s /app/config.yml -S /app/schema.yml > /app/final-config.json &&
        /app/start-server
      "
```

### Kubernetes
```yaml
# deployment.yaml
env:
  - name: KONFIGO_KEY_service.name
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
  - name: KONFIGO_KEY_service.namespace
    valueFrom:
      fieldRef:
        fieldPath: metadata.namespace
  - name: KONFIGO_VAR_DATABASE_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-secret
        key: password
```

### CI/CD Pipeline
```bash
#!/bin/bash
# Build-time configuration
export KONFIGO_KEY_build.timestamp=$(date -u +%Y-%m-%dT%H:%M:%SZ)
export KONFIGO_KEY_build.commit=${GITHUB_SHA}
export KONFIGO_KEY_build.branch=${GITHUB_REF##*/}

# Environment-specific variables
if [[ "$ENVIRONMENT" == "production" ]]; then
  export KONFIGO_VAR_API_BASE_URL=https://api.example.com
  export KONFIGO_VAR_CDN_URL=https://cdn.example.com
else
  export KONFIGO_VAR_API_BASE_URL=https://api-staging.example.com
  export KONFIGO_VAR_CDN_URL=https://cdn-staging.example.com
fi

# Generate configuration
konfigo -s base.yml,${ENVIRONMENT}.yml -S schema.yml > dist/config.json
```

## Test Coverage

Environment integration is thoroughly tested in `test/env-integration/`:
- Direct key overrides with `KONFIGO_KEY_`
- Schema variable resolution with `KONFIGO_VAR_`
- Type conversion verification
- Complex path handling
- Error condition testing
- Integration with other features
