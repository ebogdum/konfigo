# Best Practices

Proven patterns and recommendations for effective Konfigo usage in production environments.

## Configuration Architecture

### File Organization

#### ✅ Recommended Structure
```
configs/
├── base/
│   ├── app.yaml           # Core application settings
│   ├── database.yaml      # Database configuration
│   └── services.yaml      # Service discovery
├── environments/
│   ├── development.yaml   # Dev-specific overrides
│   ├── staging.yaml       # Staging overrides
│   └── production.yaml    # Production overrides
├── schemas/
│   ├── app.schema.yaml    # Validation schemas
│   └── variables.yaml     # Variable definitions
└── local/
    └── developer.yaml     # Local overrides (git-ignored)
```

#### ❌ Anti-patterns
```
# Don't mix concerns in single files
all-config.yaml          # Everything in one file
app-and-db-config.yaml   # Multiple concerns mixed

# Don't use unclear naming
config1.yaml, config2.yaml  # Unclear purpose
final.yaml, temp.yaml       # Temporary names in permanent use
```

### Naming Conventions

#### Configuration Files
```yaml
# ✅ Good: Clear, descriptive names
app-core.yaml
database-primary.yaml
microservice-auth.yaml
environment-production.yaml

# ❌ Bad: Unclear or inconsistent
config.yaml
settings.yml
app.toml
data.json  # Mixed formats without reason
```

#### Configuration Keys
```yaml
# ✅ Good: Consistent, hierarchical
app:
  name: "my-service"
  server:
    port: 8080
    timeout: 30s
  database:
    host: "localhost"
    pool:
      min: 5
      max: 20

# ❌ Bad: Inconsistent, flat
appName: "my-service"        # Mixed camelCase
server_port: 8080           # Mixed snake_case
databaseHost: "localhost"   # No clear hierarchy
```

## Security Best Practices

### Secret Management

#### ✅ Secure Approaches
```bash
# Use environment variables for secrets
export KONFIGO_VAR_DB_PASSWORD="$SECRET_PASSWORD"
export KONFIGO_VAR_API_KEY="$API_SECRET"

# Keep secrets out of version control
echo "local/*.yaml" >> .gitignore
echo "secrets/*.yaml" >> .gitignore

# Use external secret management
kubectl create secret generic app-secrets \
  --from-literal=db-password="$DB_PASSWORD"
```

#### ❌ Security Anti-patterns
```yaml
# ❌ Never commit secrets to version control
database:
  password: "supersecret123"  # Security risk!
  
api:
  key: "sk-abc123xyz789"      # Exposed secret!
```

### File Permissions
```bash
# ✅ Restrict access to sensitive configs
chmod 600 production.yaml
chmod 640 database-config.yaml  # Group read if needed

# ✅ Set proper directory permissions
chmod 755 configs/
chmod 700 secrets/
```

### Schema Security
```yaml
# ✅ Mark sensitive paths as immutable
immutable:
  - "security.apiKeys"
  - "database.credentials"
  - "encryption.keys"

# ✅ Validate sensitive data formats
validate:
  - path: "api.key"
    rules:
      required: true
      type: "string"
      regex: "^sk-[a-zA-Z0-9]{32}$"  # Validate format
```

## Performance Optimization

### Large Configuration Management

#### File Size Optimization
```bash
# ✅ Split large configs into logical modules
konfigo -s base.yaml,app-module.yaml,db-module.yaml

# ✅ Use recursive discovery for organized structures
konfigo -r -s ./configs/modules/

# ❌ Avoid single massive files
# Don't: everything-config.yaml (10MB+)
```

#### Memory Efficiency
```yaml
# ✅ Optimize schema processing order
apiVersion: v1

# Fast operations first
vars: [...]           # Variable definitions
validate: [...]       # Quick validation

# Heavy operations last
generators: [...]     # Data generation
transform: [...]      # Complex transformations
```

### Build Time Optimization

#### CI/CD Integration
```bash
# ✅ Cache intermediate results
konfigo -s base.yaml,env/${ENV}.yaml -of cache/merged-${ENV}.json

# ✅ Validate early in pipeline
konfigo --validate-only -s config.yaml -S schema.yaml

# ✅ Parallel processing for multiple environments
for env in dev staging prod; do
  konfigo -s base.yaml,env/${env}.yaml -of dist/${env}-config.json &
done
wait
```

## Development Workflow

### Local Development

#### Environment Setup
```bash
# ✅ Local override pattern
konfigo -s base.yaml,env/dev.yaml,local/$(whoami).yaml

# ✅ Developer-specific configs (git-ignored)
echo "local/" >> .gitignore
mkdir -p local/
cp env/dev.yaml local/$(whoami).yaml  # Starting template
```

#### Testing Configurations
```bash
# ✅ Validate before committing
konfigo --validate-only -s config.yaml -S schema.yaml

# ✅ Test with different environments
for env in dev staging prod; do
  echo "Testing $env environment..."
  konfigo -s base.yaml,env/${env}.yaml -S validation.schema.yaml
done

# ✅ Automated testing in CI
./validate_docs_examples.sh
```

### Version Control

#### Git Best Practices
```gitignore
# ✅ Proper .gitignore
local/              # Developer-specific configs
secrets/            # Secret files
*.local.yaml        # Local override files
.env.local          # Local environment files
dist/               # Generated configs
build/              # Build artifacts
```

#### Commit Strategies
```bash
# ✅ Separate config changes from code changes
git add configs/
git commit -m "config: update production database settings"

git add src/
git commit -m "feat: add new authentication service"
```

## Production Deployment

### Environment Management

#### Multi-Environment Strategy
```bash
# ✅ Consistent deployment pattern
deploy() {
  local env=$1
  konfigo -s base.yaml,env/${env}.yaml \
          -S schemas/production.schema.yaml \
          -of dist/${env}-final.json
  
  # Validate deployment config
  validate-config dist/${env}-final.json
  
  # Deploy to environment
  deploy-to-${env} dist/${env}-final.json
}
```

#### Blue-Green Deployments
```bash
# ✅ Generate configs for both environments
konfigo -s base.yaml,env/blue.yaml -of configs/blue-config.json
konfigo -s base.yaml,env/green.yaml -of configs/green-config.json

# Validate both configs before switching
validate-config configs/blue-config.json
validate-config configs/green-config.json
```

### Monitoring and Observability

#### Configuration Tracking
```yaml
# ✅ Add metadata to generated configs
generators:
  - type: "timestamp"
    targetPath: "metadata.generated"
    
  - type: "setValue"
    targetPath: "metadata.source"
    value: "${KONFIGO_BUILD_SOURCE}"
    
  - type: "setValue"
    targetPath: "metadata.version"
    value: "${GIT_COMMIT_SHA}"
```

#### Health Checks
```bash
# ✅ Validate configuration after deployment
curl -f http://service/health/config || exit 1

# ✅ Monitor for configuration drift
konfigo -s current-config.json -S validation.schema.yaml --validate-only
```

## Schema Design Patterns

### Progressive Enhancement

#### Layered Validation
```yaml
# ✅ Start simple, add complexity gradually
apiVersion: v1

# Level 1: Basic requirements
validate:
  - path: "app.name"
    rules:
      required: true
      type: "string"

# Level 2: Add format validation
  - path: "app.version"
    rules:
      required: true
      type: "string"
      regex: "^\\d+\\.\\d+\\.\\d+$"  # Semantic versioning

# Level 3: Business logic validation
  - path: "app.port"
    rules:
      required: true
      type: "number"
      min: 1024
      max: 65535
```

#### Modular Schemas
```bash
# ✅ Split complex schemas into modules
konfigo -s config.yaml \
        -S schemas/base-validation.yaml \
        -S schemas/security-validation.yaml \
        -S schemas/performance-validation.yaml
```

### Variable Management

#### Hierarchical Variables
```yaml
# ✅ Organize variables by scope
vars:
  # Global settings
  - name: "ENVIRONMENT"
    value: "production"
  - name: "REGION"
    value: "us-west-2"
    
  # Application-specific
  - name: "APP_LOG_LEVEL"
    value: "info"
  - name: "APP_TIMEOUT"
    value: "30s"
    
  # Service-specific
  - name: "DB_POOL_SIZE"
    value: "20"
  - name: "CACHE_TTL"
    value: "3600"
```

## Error Handling and Recovery

### Graceful Degradation

#### Default Values Strategy
```yaml
# ✅ Provide sensible defaults
vars:
  - name: "TIMEOUT"
    value: "30s"           # Safe default
    description: "Request timeout"
    
  - name: "RETRY_COUNT"
    value: "3"             # Conservative default
    description: "Number of retry attempts"

# ✅ Validate with fallbacks
validate:
  - path: "app.timeout"
    rules:
      type: "string"
      default: "${TIMEOUT}"  # Use variable default
```

#### Error Recovery
```bash
# ✅ Robust deployment script
deploy_config() {
  local backup_config="config-backup-$(date +%s).json"
  
  # Backup current config
  cp current-config.json "$backup_config"
  
  # Generate new config
  if ! konfigo -s base.yaml,env/prod.yaml -of new-config.json; then
    echo "Config generation failed, keeping current config"
    return 1
  fi
  
  # Validate new config
  if ! validate-config new-config.json; then
    echo "Config validation failed, rolling back"
    return 1
  fi
  
  # Deploy new config
  cp new-config.json current-config.json
  echo "Config deployed successfully, backup saved as $backup_config"
}
```

## Team Collaboration

### Documentation

#### Self-Documenting Configs
```yaml
# ✅ Document complex configurations
app:
  # Port for HTTP server (must be > 1024 for non-root)
  port: 8080
  
  # Timeout in seconds for external API calls
  # Increase for slower networks, decrease for faster failover
  timeout: 30
  
  # Feature flags - enable/disable functionality
  features:
    auth: true      # OAuth2 authentication
    cache: true     # Redis caching layer
    metrics: false  # Prometheus metrics (disabled in dev)
```

#### Schema Documentation
```yaml
# ✅ Document schema intentions
apiVersion: v1
metadata:
  name: "production-config-schema"
  description: |
    Production configuration schema for microservice deployment.
    Enforces security requirements and performance constraints.
    
    Contact: devops@company.com
    Last updated: 2025-06-27

vars:
  - name: "DATABASE_HOST"
    value: "prod-db.company.com"
    description: |
      Primary database hostname. Points to load balancer in production.
      Override with KONFIGO_VAR_DATABASE_HOST for different environments.
```

### Code Review Process

#### Configuration Review Checklist
- [ ] No secrets in committed files
- [ ] Schema validation passes
- [ ] All environments tested
- [ ] Performance impact assessed
- [ ] Documentation updated
- [ ] Backward compatibility verified

#### Automated Checks
```bash
# ✅ Pre-commit hook
#!/bin/bash
# Check for secrets
if grep -r "password\|secret\|key" configs/ --exclude-dir=local; then
  echo "Potential secrets found in configs"
  exit 1
fi

# Validate all schemas
for schema in schemas/*.yaml; do
  konfigo --validate-only -S "$schema" || exit 1
done
```

This comprehensive guide provides battle-tested practices for using Konfigo effectively in production environments, from development through deployment and maintenance.
