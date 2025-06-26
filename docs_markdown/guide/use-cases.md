# Use Cases & Examples

This page demonstrates how to solve common, real-world configuration management problems using Konfigo. Each example includes the complete setup, commands, and explanations to help you implement similar solutions.

## 1. Environment Promotion (Dev/Staging/Prod)

**The Goal:** Manage a base configuration and layer on environment-specific overrides for staging and production.

**The Setup:**
Create a directory with a base configuration and an override file for each environment.

**`configs/base.yml`**
```yaml
service:
  name: my-awesome-app
  port: 8080
database:
  host: localhost
  user: app_user
logging:
  level: debug
```

**`configs/production.yml`**
```yaml
database:
  host: prod-db.internal.net
logging:
  level: info
```

**The Command:**
The order of sources in the `-s` flag is critical. The last source specified wins in case of conflicts.

```bash
konfigo -s configs/base.yml,configs/production.yml
```

**The Result:**
```json{6,9}
{
  "database": {
    "host": "prod-db.internal.net",
    "user": "app_user"
  },
  "logging": {
    "level": "info"
  },
  "service": {
    "name": "my-awesome-app",
    "port": 8080
  }
}
```

**Explanation:**
The values for `database.host` and `logging.level` were overwritten by `production.yml` because it was the last source loaded. All other values from `base.yml` were preserved.

---

## 2. CI/CD Integration with Dynamic Tags & Secrets

**The Goal:** Build a configuration in a CI/CD pipeline that uses the Git commit tag for the Docker image and injects a database password from a secure environment variable.

**The Setup:**

**`configs/ci.yml`**
```yaml
# This file contains the base structure.
# The image tag will be supplied by a variable.
deployment:
  image: "my-registry.io/my-awesome-app:${RELEASE_VERSION}"
database:
  user: "ci_user"
```

**`schema.yml`**
```yaml
# The schema defines how to get the RELEASE_VERSION variable.
vars:
  - name: "RELEASE_VERSION"
    fromEnv: "CI_COMMIT_TAG" # Read from an env var set by the CI system
    defaultValue: "latest"
validate:
  - path: "database.password"
    rules:
      required: true
      minLength: 16
```

**The Command:**
In your CI/CD script, you would set the secure environment variables and run Konfigo.

```bash{1,2,5}
# These are provided by the CI/CD system's secret management and environment
export KONFIGO_KEY_database.password="a-very-secure-password-from-ci"
export CI_COMMIT_TAG="v1.2.3"

konfigo \
  -S schema.yml \
  -s configs/ci.yml
```

**The Result:**
```json
{
  "database": {
    "password": "a-very-secure-password-from-ci",
    "user": "ci_user"
  },
  "deployment": {
    "image": "my-registry.io/my-awesome-app:v1.2.3"
  }
}
```

**Explanation:**
- `KONFIGO_KEY_database.password` directly injected the secret into the configuration.
- The `vars` block in the schema read the `CI_COMMIT_TAG` environment variable.
- The `${RELEASE_VERSION}` placeholder was substituted with `v1.2.3` during processing.

---

## 3. Microservices Configuration Management

**The Goal:** Generate individual configuration files for multiple microservices from a shared base configuration.

**The Setup:**

**`base-service.yaml`**
```yaml
app:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    port: 9090
database:
  timeout: 30
  pool:
    min: 5
    max: 20
```

**`services-schema.yaml`**
```yaml
vars:
  - name: "SERVICE_NAME"
    fromPath: "service.name"
  - name: "SERVICE_PORT"
    fromPath: "service.port"
  - name: "DATABASE_NAME"
    fromEnv: "DB_NAME"
    defaultValue: "${SERVICE_NAME}_db"

generators:
  - type: "concat"
    targetPath: "database.url"
    format: "postgresql://user:pass@localhost:5432/${DATABASE_NAME}"

transform:
  - type: "setValue"
    path: "app.name"
    value: "${SERVICE_NAME}"
  - type: "setValue"
    path: "app.port"
    value: "${SERVICE_PORT}"

validate:
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 3000
      max: 9000
```

**`services-batch.yaml`**
```yaml
konfigo_forEach:
  items:
    - service:
        name: "user-service"
        port: 3001
    - service:
        name: "order-service"  
        port: 3002
    - service:
        name: "payment-service"
        port: 3003
  
  output:
    filenamePattern: "configs/${SERVICE_NAME}-config.yaml"
    format: "yaml"
```

**The Command:**
```bash
konfigo -s base-service.yaml -S services-schema.yaml -V services-batch.yaml
```

**The Result:**
Three files are generated:
- `configs/user-service-config.yaml`
- `configs/order-service-config.yaml`
- `configs/payment-service-config.yaml`

Each containing service-specific configurations with the correct ports and database URLs.

---

## 4. Kubernetes Deployment Generation

**The Goal:** Generate Kubernetes deployment manifests from application configuration.

**The Setup:**

**`app-config.yaml`**
```yaml
app:
  name: "web-api"
  image: "myregistry/web-api"
  version: "1.2.3"
  replicas: 3
  resources:
    requests:
      cpu: "100m"
      memory: "128Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  env:
    DATABASE_URL: "postgresql://db:5432/app"
    REDIS_URL: "redis://redis:6379"
```

**`k8s-deployment-schema.yaml`**
```yaml
vars:
  - name: "NAMESPACE"
    fromEnv: "K8S_NAMESPACE"
    defaultValue: "default"
  - name: "IMAGE_TAG"
    fromEnv: "IMAGE_TAG"
    defaultValue: "${app.version}"

generators:
  - type: "concat"
    targetPath: "spec.template.spec.containers[0].image"
    format: "${app.image}:${IMAGE_TAG}"

transform:
  - type: "setValue"
    path: "apiVersion"
    value: "apps/v1"
  - type: "setValue"
    path: "kind"
    value: "Deployment"
  - type: "setValue"
    path: "metadata.name"
    value: "${app.name}"
  - type: "setValue"
    path: "metadata.namespace"
    value: "${NAMESPACE}"
  - type: "setValue"
    path: "spec.replicas"
    value: "${app.replicas}"
  - type: "setValue"
    path: "spec.selector.matchLabels.app"
    value: "${app.name}"
  - type: "setValue"
    path: "spec.template.metadata.labels.app"
    value: "${app.name}"
  - type: "setValue"
    path: "spec.template.spec.containers[0].name"
    value: "${app.name}"
  - type: "setValue"
    path: "spec.template.spec.containers[0].resources"
    value: "${app.resources}"
  - type: "addKeyPrefix"
    path: "app.env"
    prefix: ""
  - type: "renameKey"
    from: "app.env"
    to: "spec.template.spec.containers[0].env"

validate:
  - path: "spec.replicas"
    rules:
      type: "number"
      min: 1
      max: 10
```

**The Command:**
```bash
export K8S_NAMESPACE="production"
export IMAGE_TAG="v1.2.3"

konfigo -s app-config.yaml -S k8s-deployment-schema.yaml -of deployment.yaml
```

---

## 5. Configuration Validation and Compliance

**The Goal:** Ensure all configurations meet security and compliance requirements.

**The Setup:**

**`security-schema.yaml`**
```yaml
validate:
  # Ensure HTTPS URLs only
  - path: "services.*.url"
    rules:
      type: "string"
      regex: "^https://"
  
  # Validate port ranges (no privileged ports)
  - path: "*.port"
    rules:
      type: "number"
      min: 1024
      max: 65535
  
  # Ensure no hardcoded passwords
  - path: "database.password"
    rules:
      type: "string"
      regex: "^\\$\\{[A-Z_]+\\}$"  # Must be a variable reference
  
  # Validate log levels
  - path: "logging.level"
    rules:
      type: "string"
      enum: ["error", "warn", "info", "debug"]
  
  # Ensure resource limits are set
  - path: "resources.limits.memory"
    rules:
      required: true
      type: "string"
      regex: "^\\d+[KMGT]i$"  # Kubernetes memory format
  
  # Validate environment-specific settings
  - path: "environment"
    rules:
      required: true
      enum: ["development", "staging", "production"]

# Security transformations
transform:
  # Remove debug settings in production
  - type: "setValue"
    path: "debug"
    value: false
    condition: '${ENVIRONMENT} == "production"'
  
  # Set secure defaults
  - type: "setValue"
    path: "security.tls.minVersion"
    value: "1.2"
```

**The Command:**
```bash
# Validate development config
export ENVIRONMENT="development"
konfigo -s configs/dev-config.yaml -S security-schema.yaml

# Validate production config (will fail if insecure)
export ENVIRONMENT="production"
konfigo -s configs/prod-config.yaml -S security-schema.yaml
```

---

## 6. Legacy Configuration Migration

**The Goal:** Transform legacy configuration formats to modern structure.

**The Setup:**

**`legacy-config.yaml`**
```yaml
# Old flat structure
app_name: "old-app"
app_port: 8080
db_host: "localhost"
db_port: 5432
db_user: "user"
log_level: "INFO"
feature_auth_enabled: true
feature_cache_ttl: 3600
```

**`migration-schema.yaml`**
```yaml
transform:
  # Restructure into nested objects
  - type: "renameKey"
    from: "app_name"
    to: "application.name"
  - type: "renameKey"
    from: "app_port"
    to: "application.port"
  
  - type: "renameKey"
    from: "db_host"
    to: "database.host"
  - type: "renameKey"
    from: "db_port"
    to: "database.port"
  - type: "renameKey"
    from: "db_user"
    to: "database.user"
  
  - type: "renameKey"
    from: "log_level"
    to: "logging.level"
  - type: "changeCase"
    path: "logging.level"
    case: "lower"
  
  # Group feature flags
  - type: "renameKey"
    from: "feature_auth_enabled"
    to: "features.auth.enabled"
  - type: "renameKey"
    from: "feature_cache_ttl"
    to: "features.cache.ttl"

generators:
  # Generate new structured values
  - type: "concat"
    targetPath: "database.connectionString"
    format: "postgresql://{user}@{host}:{port}/app"
    sources:
      user: "database.user"
      host: "database.host"
      port: "database.port"

validate:
  # Ensure migration was successful
  - path: "application.name"
    rules:
      required: true
      type: "string"
  - path: "database.connectionString"
    rules:
      required: true
      type: "string"
      regex: "^postgresql://"
```

**The Command:**
```bash
konfigo -s legacy-config.yaml -S migration-schema.yaml -of modern-config.yaml
```

**The Result:**
```yaml
application:
  name: "old-app"
  port: 8080
database:
  host: "localhost"
  port: 5432
  user: "user"
  connectionString: "postgresql://user@localhost:5432/app"
logging:
  level: "info"
features:
  auth:
    enabled: true
  cache:
    ttl: 3600
```

---

## 7. Multi-Environment CI/CD Pipeline

**The Goal:** Create a complete CI/CD pipeline that generates environment-specific configurations.

**The Setup:**

**`.github/workflows/deploy.yml`**
```yaml
name: Deploy Application

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  validate-configs:
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
          for env in dev staging prod; do
            echo "Validating $env environment..."
            export ENVIRONMENT=$env
            konfigo -s configs/base.yaml,configs/$env.yaml \
                   -S schemas/validation.yaml \
                   -V variables/$env.yaml
          done

  deploy-dev:
    needs: validate-configs
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: development
    steps:
      - uses: actions/checkout@v3
      
      - name: Generate dev configuration
        env:
          DATABASE_PASSWORD: ${{ secrets.DEV_DB_PASSWORD }}
          API_KEY: ${{ secrets.DEV_API_KEY }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          konfigo -s configs/base.yaml,configs/dev.yaml \
                 -S schemas/k8s-deployment.yaml \
                 -V variables/dev.yaml \
                 -of k8s/dev-deployment.yaml
      
      - name: Deploy to development
        run: |
          kubectl apply -f k8s/dev-deployment.yaml

  deploy-prod:
    needs: deploy-dev
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v3
      
      - name: Generate prod configuration
        env:
          DATABASE_PASSWORD: ${{ secrets.PROD_DB_PASSWORD }}
          API_KEY: ${{ secrets.PROD_API_KEY }}
          IMAGE_TAG: ${{ github.sha }}
          REPLICAS: 5
        run: |
          konfigo -s configs/base.yaml,configs/prod.yaml \
                 -S schemas/k8s-deployment.yaml \
                 -V variables/prod.yaml \
                 -of k8s/prod-deployment.yaml
      
      - name: Deploy to production
        run: |
          kubectl apply -f k8s/prod-deployment.yaml
```

---

## 8. Configuration Testing and Validation

**The Goal:** Implement comprehensive testing for configuration management.

**The Setup:**

**`test/test-configs.sh`**
```bash
#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Testing Konfigo configurations..."

# Test 1: Basic configuration merging
test_basic_merge() {
    echo -n "Testing basic merge... "
    result=$(konfigo -s configs/base.yaml,configs/dev.yaml)
    if echo "$result" | grep -q "name.*my-app"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        echo "Expected 'my-app' in merged configuration"
        exit 1
    fi
}

# Test 2: Schema validation
test_schema_validation() {
    echo -n "Testing schema validation... "
    # This should pass
    konfigo -s test/fixtures/valid-config.yaml -S schemas/validation.yaml >/dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        echo "Valid configuration failed validation"
        exit 1
    fi
    
    # This should fail
    if konfigo -s test/fixtures/invalid-config.yaml -S schemas/validation.yaml >/dev/null 2>&1; then
        echo -e "${RED}✗${NC}"
        echo "Invalid configuration passed validation"
        exit 1
    fi
}

# Test 3: Environment variable integration
test_env_vars() {
    echo -n "Testing environment variables... "
    export KONFIGO_KEY_app.environment="test"
    export KONFIGO_VAR_VERSION="test-version"
    
    result=$(konfigo -s configs/base.yaml -S schemas/basic.yaml)
    if echo "$result" | grep -q "test" && echo "$result" | grep -q "test-version"; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        echo "Environment variables not properly applied"
        exit 1
    fi
    
    unset KONFIGO_KEY_app.environment
    unset KONFIGO_VAR_VERSION
}

# Test 4: Batch processing
test_batch_processing() {
    echo -n "Testing batch processing... "
    
    # Clean up any existing outputs
    rm -rf test/output/
    mkdir -p test/output/
    
    cd test/output/
    konfigo -s ../fixtures/base.yaml -S ../fixtures/batch-schema.yaml -V ../fixtures/batch-vars.yaml
    
    # Check that all expected files were created
    expected_files=("service1-config.yaml" "service2-config.yaml" "service3-config.yaml")
    for file in "${expected_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            echo -e "${RED}✗${NC}"
            echo "Expected output file $file not found"
            exit 1
        fi
    done
    
    cd "$PROJECT_ROOT"
    echo -e "${GREEN}✓${NC}"
}

# Test 5: Format conversion
test_format_conversion() {
    echo -n "Testing format conversion... "
    
    # YAML to JSON
    json_result=$(konfigo -s configs/base.yaml -oj)
    if ! echo "$json_result" | jq . >/dev/null 2>&1; then
        echo -e "${RED}✗${NC}"
        echo "YAML to JSON conversion failed"
        exit 1
    fi
    
    # JSON to TOML
    toml_result=$(konfigo -s test/fixtures/sample.json -ot)
    if [[ -z "$toml_result" ]]; then
        echo -e "${RED}✗${NC}"
        echo "JSON to TOML conversion failed"
        exit 1
    fi
    
    echo -e "${GREEN}✓${NC}"
}

# Run all tests
echo "Running configuration tests..."
test_basic_merge
test_schema_validation
test_env_vars
test_batch_processing
test_format_conversion

echo -e "\n${GREEN}All tests passed!${NC} 🎉"
```

---

## Best Practices from Use Cases

Based on these real-world examples, here are key best practices:

### 1. **Structure Your Configuration Hierarchy**
- Use a clear base → environment → local override pattern
- Keep environment-specific changes minimal and focused
- Document the merge order and precedence rules

### 2. **Leverage Schema Validation Early**
- Implement validation rules that catch common mistakes
- Use enum validation for controlled values
- Validate security requirements (HTTPS, port ranges, etc.)

### 3. **Automate Configuration Testing**
- Test configuration generation in CI/CD pipelines
- Validate all environments before deployment
- Use batch processing to test multiple scenarios

### 4. **Secure Secret Management**
- Use environment variables for secrets, never commit them
- Implement validation for secret formats
- Use external secret management systems in production

### 5. **Plan for Migration and Evolution**
- Use transformation schemas to migrate legacy configurations
- Version your schemas and configuration structure
- Test migration paths thoroughly

These examples demonstrate Konfigo's flexibility in handling complex, real-world configuration management scenarios. Each pattern can be adapted and combined to meet your specific requirements.
