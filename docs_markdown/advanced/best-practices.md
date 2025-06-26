# Best Practices

Following these best practices will help you build robust, maintainable, and secure configuration management workflows with Konfigo.

## Configuration Organization

### Directory Structure
Organize your configuration files in a logical hierarchy:

```
configs/
├── base/                    # Base configurations
│   ├── app.yaml
│   ├── database.yaml
│   └── services.yaml
├── environments/            # Environment-specific overrides
│   ├── development.yaml
│   ├── staging.yaml
│   └── production.yaml
├── schemas/                 # Schema definitions
│   ├── validation.yaml
│   ├── transformation.yaml
│   └── deployment.yaml
└── variables/               # Variable definitions
    ├── common.yaml
    ├── secrets.yaml.example  # Template for secrets
    └── batch-processing.yaml
```

### Naming Conventions
- **Use consistent naming**: Prefer kebab-case for files (`app-config.yaml`)
- **Include environment indicators**: `production-database.yaml`, `dev-services.yaml`
- **Use descriptive schema names**: `k8s-deployment-schema.yaml`, `validation-schema.yaml`
- **Version important configs**: `app-config-v2.yaml` for breaking changes

## Schema Design

### Modular Schema Approach
Break schemas into focused, reusable components:

```yaml
# base-schema.yaml - Common variables and transformations
vars:
  - name: "ENVIRONMENT"
    fromEnv: "DEPLOY_ENV"
    defaultValue: "development"
  - name: "VERSION"
    fromEnv: "APP_VERSION"
    defaultValue: "latest"

transform:
  - type: "setValue"
    path: "metadata.environment"
    value: "${ENVIRONMENT}"
  - type: "setValue"
    path: "metadata.version"
    value: "${VERSION}"

---
# validation-schema.yaml - Validation rules
validate:
  - path: "service.name"
    rules:
      required: true
      type: "string"
      minLength: 3
      regex: "^[a-z][a-z0-9-]*[a-z0-9]$"
  
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
```

### Progressive Validation Strategy
1. **Start simple**: Basic type and required field validation
2. **Add constraints**: Ranges, patterns, enums as requirements clarify
3. **Include business rules**: Domain-specific validation logic
4. **Document validation**: Clear error messages and examples

## Variable Management

### Hierarchical Variable Strategy
Organize variables by scope and precedence:

```yaml
# Global variables (lowest precedence)
global:
  APP_NAME: "my-service"
  VERSION: "1.0.0"
  LOG_FORMAT: "json"

# Environment-specific variables
vars:
  - name: "DATABASE_URL"
    fromEnv: "DATABASE_CONNECTION_STRING"
    defaultValue: "postgresql://localhost:5432/app"
  
  - name: "LOG_LEVEL"
    fromEnv: "LOG_LEVEL"  
    defaultValue: "info"
  
  - name: "FEATURE_FLAGS"
    fromPath: "features.enabled"
    defaultValue: []
```

### Secret Management Best Practices
- **Use environment variables**: Never commit secrets to configuration files
- **External secret systems**: Integrate with HashiCorp Vault, AWS Secrets Manager, etc.
- **Secret templates**: Provide `.example` files showing required secret structure
- **Validation**: Validate secret formats and requirements

```yaml
# secrets.yaml.example
vars:
  - name: "DATABASE_PASSWORD"
    fromEnv: "DB_PASSWORD"
    # defaultValue: "REQUIRED - Set DB_PASSWORD environment variable"
  
  - name: "API_KEY"
    fromEnv: "EXTERNAL_API_KEY"
    # Format: key-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Testing Strategies

### Configuration Testing
Implement comprehensive testing for your configuration pipeline:

```bash
#!/bin/bash
# test-configs.sh

set -e

# Test basic merging
test_merge() {
    echo "Testing configuration merging..."
    result=$(konfigo -s base.yaml,dev.yaml)
    if ! echo "$result" | grep -q "expected_setting"; then
        echo "ERROR: Missing expected setting in merge result"
        exit 1
    fi
    echo "✓ Merge test passed"
}

# Test schema validation
test_schema_validation() {
    echo "Testing schema validation..."
    if konfigo -s invalid-config.yaml -S strict-schema.yaml >/dev/null 2>&1; then
        echo "ERROR: Schema validation should have failed"
        exit 1
    fi
    echo "✓ Schema validation test passed"
}

# Test batch processing
test_batch_processing() {
    echo "Testing batch processing..."
    konfigo -s base.yaml -S schema.yaml -V batch-vars.yaml
    
    for env in dev staging prod; do
        if [[ ! -f "output/${env}-config.yaml" ]]; then
            echo "ERROR: Missing output file for $env"
            exit 1
        fi
    done
    echo "✓ Batch processing test passed"
}

# Run all tests
test_merge
test_schema_validation
test_batch_processing
echo "All tests passed!"
```

### CI/CD Integration
Integrate configuration validation into your deployment pipeline:

```yaml
# .github/workflows/config-validation.yml
name: Configuration Validation

on:
  pull_request:
    paths:
      - 'configs/**'
      - 'schemas/**'
      - 'variables/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Konfigo
        run: |
          curl -L -o konfigo https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64
          chmod +x konfigo
          sudo mv konfigo /usr/local/bin/
      
      - name: Validate all environments
        run: |
          for env in development staging production; do
            echo "Validating $env configuration..."
            konfigo -s configs/base.yaml,configs/environments/$env.yaml \
                   -S schemas/validation.yaml \
                   -V variables/common.yaml
          done
      
      - name: Test batch processing
        run: |
          konfigo -s configs/base.yaml \
                 -S schemas/deployment.yaml \
                 -V variables/services.yaml
          
          # Verify outputs were generated
          ls -la output/
```

## Security Considerations

### Access Control
- **File permissions**: Restrict read access to configuration files
- **Repository access**: Limit who can modify configuration repositories
- **Approval workflows**: Require reviews for configuration changes
- **Audit logging**: Track configuration changes and access

### Secret Security
- **Never commit secrets**: Use `.gitignore` for secret files
- **Rotate regularly**: Implement secret rotation policies
- **Validate formats**: Ensure secrets match expected patterns
- **Monitor usage**: Track secret access and unusual patterns

```yaml
# Security-focused validation schema
validate:
  # Ensure no hardcoded secrets
  - path: "database.password"
    rules:
      type: "string"
      regex: "^\\$\\{[A-Z_]+\\}$"  # Must be a variable reference
  
  # Validate API key format
  - path: "api.key"
    rules:
      type: "string"
      regex: "^key-[a-f0-9]{32}$"
  
  # Ensure HTTPS URLs only
  - path: "services.*.url"
    rules:
      type: "string"
      regex: "^https://"
```

## Documentation Standards

### Schema Documentation
Document your schemas thoroughly:

```yaml
# deployment-schema.yaml
# Purpose: Generates Kubernetes deployment manifests for microservices
# Usage: konfigo -s base.yaml -S deployment-schema.yaml -V service-vars.yaml
# Requirements: 
#   - DOCKER_REGISTRY environment variable
#   - SERVICE_NAME in variables file

apiVersion: "konfigo/v1alpha1"

vars:
  # Docker registry for image pulling
  # Example: registry.example.com/myapp
  - name: "REGISTRY"
    fromEnv: "DOCKER_REGISTRY"
    defaultValue: "localhost:5000"
  
  # Service replica count
  # Range: 1-10 for non-production, 2-20 for production
  - name: "REPLICAS"
    fromEnv: "SERVICE_REPLICAS"
    defaultValue: "2"

# Generate deployment-specific values
generators:
  - type: "concat"
    targetPath: "spec.template.spec.containers[0].image"
    format: "${REGISTRY}/${SERVICE_NAME}:${VERSION}"
    sources:
      SERVICE_NAME: "metadata.name"
      VERSION: "metadata.version"
```

### README Documentation
Create comprehensive README files:

```markdown
# Configuration Management

## Overview
This directory contains configuration files for the MyApp service using Konfigo.

## Structure
- `base/`: Base configuration files
- `environments/`: Environment-specific overrides
- `schemas/`: Konfigo schema definitions
- `variables/`: Variable definitions for different contexts

## Usage

### Local Development
```bash
konfigo -s base/app.yaml,environments/development.yaml -of local-config.yaml
```

### Production Deployment
```bash
export DATABASE_PASSWORD="$(vault kv get -field=password secret/myapp/db)"
export API_KEY="$(vault kv get -field=key secret/myapp/external-api)"

konfigo -s base/app.yaml,environments/production.yaml \
       -S schemas/production.yaml \
       -V variables/production.yaml \
       -of production-config.yaml
```

## Required Environment Variables
- `DATABASE_PASSWORD`: Database connection password
- `API_KEY`: External API authentication key
- `DOCKER_REGISTRY`: Container registry URL

## Testing
Run `./scripts/test-configs.sh` to validate all configurations.
```

## Performance Optimization

### Large Configuration Files
- **Split large files**: Break monolithic configs into focused files
- **Use batch processing**: Process multiple outputs efficiently
- **Cache schemas**: Reuse schema definitions across environments
- **Minimize variable resolution**: Avoid complex path-based variables in hot paths

### Memory Management
- **Stream processing**: Use stdin for very large configurations
- **Parallel processing**: Process multiple environments concurrently
- **Resource limits**: Set appropriate memory limits in containers

## Error Handling and Recovery

### Graceful Degradation
```yaml
# Provide sensible defaults for optional features
vars:
  - name: "CACHE_ENABLED"
    fromEnv: "ENABLE_CACHE"
    defaultValue: "false"
  
  - name: "FEATURE_FLAG_SERVICE"
    fromEnv: "FEATURE_SERVICE_URL"
    defaultValue: ""  # Empty means use local defaults

# Validate critical vs optional settings
validate:
  # Critical - must be present
  - path: "database.url"
    rules:
      required: true
      type: "string"
  
  # Optional - can be empty
  - path: "features.external_service"
    rules:
      type: "string"
      # No required: true
```

### Monitoring and Alerting
- **Configuration drift**: Monitor for unexpected configuration changes
- **Validation failures**: Alert on schema validation failures
- **Performance metrics**: Track configuration processing time
- **Success rates**: Monitor batch processing success rates

Following these best practices will help you build robust, secure, and maintainable configuration management systems with Konfigo.
