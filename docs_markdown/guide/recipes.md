# Recipes & Examples

Real-world patterns and solutions for common Konfigo use cases. Each recipe includes the problem, complete solution, and explanation.

## Quick Navigation

### 🏗️ **Deployment Patterns**
- [Multi-Environment Deployment](#multi-environment-deployment)
- [Container & Kubernetes](#container--kubernetes)
- [CI/CD Integration](#cicd-integration)

### 🔧 **Development Workflows**
- [Team Configuration Management](#team-configuration-management)  
- [Local Development Setup](#local-development-setup)
- [Legacy System Modernization](#legacy-system-modernization)

### 🚀 **Advanced Patterns**
- [Microservices Configuration](#microservices-configuration)
- [Feature Flag Management](#feature-flag-management)
- [Configuration Validation Pipeline](#configuration-validation-pipeline)

---

## Multi-Environment Deployment

**Problem**: Deploy the same application to dev, staging, and production with environment-specific settings.

### Solution Structure

```
configs/
├── base.yaml           # Common settings
├── environments/
│   ├── dev.yaml       # Development overrides
│   ├── staging.yaml   # Staging overrides
│   └── prod.yaml      # Production overrides
└── schema.yaml        # Validation rules
```

### Files

**`base.yaml`** - Common configuration:
```yaml
app:
  name: "user-service"
  version: "1.0.0"
  timeout: 30
  features:
    auth: true
    logging: true
    
database:
  port: 5432
  pool_size: 10
  timeout: 5
  
monitoring:
  enabled: true
  metrics_port: 9090
```

**`environments/dev.yaml`** - Development overrides:
```yaml
app:
  debug: true
  log_level: "debug"
  
database:
  host: "dev-db.company.com"
  pool_size: 5
  
monitoring:
  sampling_rate: 1.0
```

**`environments/prod.yaml`** - Production overrides:
```yaml
app:
  debug: false
  log_level: "info"
  replicas: 3
  
database:
  host: "prod-db.company.com"
  pool_size: 50
  ssl: true
  
monitoring:
  sampling_rate: 0.1
  alerts_enabled: true
```

**`schema.yaml`** - Validation:
```yaml
validation:
  - path: "app.name"
    required: true
    type: "string"
  - path: "database.host"
    required: true
    type: "string"
  - path: "app.replicas"
    type: "number"
    min: 1
    max: 10
```

### Deployment Commands

```bash
# Development deployment
konfigo -s base.yaml,environments/dev.yaml -S schema.yaml -of dev-config.json

# Staging deployment  
konfigo -s base.yaml,environments/staging.yaml -S schema.yaml -of staging-config.json

# Production deployment
konfigo -s base.yaml,environments/prod.yaml -S schema.yaml -of prod-config.json
```

### Automation Script

```bash
#!/bin/bash
# deploy-config.sh

ENVIRONMENT=${1:-dev}
BASE_CONFIG="base.yaml"
ENV_CONFIG="environments/${ENVIRONMENT}.yaml"
SCHEMA="schema.yaml"
OUTPUT="deploy/${ENVIRONMENT}-config.json"

if [[ ! -f "$ENV_CONFIG" ]]; then
  echo "Error: Environment config $ENV_CONFIG not found"
  exit 1
fi

echo "Generating configuration for $ENVIRONMENT..."
konfigo -s "$BASE_CONFIG,$ENV_CONFIG" -S "$SCHEMA" -of "$OUTPUT"

if [[ $? -eq 0 ]]; then
  echo "✅ Configuration generated: $OUTPUT"
else
  echo "❌ Configuration generation failed"
  exit 1
fi
```

**Usage**:
```bash
./deploy-config.sh dev
./deploy-config.sh staging  
./deploy-config.sh prod
```

---

## Container & Kubernetes

**Problem**: Deploy containerized applications with runtime configuration overrides.

### Docker Deployment

**`docker-compose.yml`**:
```yaml
version: '3.8'
services:
  app:
    image: myapp:latest
    environment:
      # Override config at runtime
      - KONFIGO_KEY_database.host=${DB_HOST:-localhost}
      - KONFIGO_KEY_database.password=${DB_PASSWORD}
      - KONFIGO_KEY_app.port=${APP_PORT:-8080}
      - KONFIGO_KEY_app.environment=${ENVIRONMENT:-development}
    volumes:
      - ./configs:/app/configs
    command: |
      sh -c "
        konfigo -s /app/configs/base.yaml -of /app/runtime-config.json &&
        exec myapp --config /app/runtime-config.json
      "
```

**Environment file (`.env`)**:
```bash
# .env
DB_HOST=prod-db.company.com
DB_PASSWORD=secure_password
APP_PORT=9090
ENVIRONMENT=production
```

**Run with**:
```bash
docker-compose up
```

### Kubernetes Deployment

**`kubernetes/configmap.yaml`**:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-base-config
data:
  base.yaml: |
    app:
      name: "user-service"
      port: 8080
      timeout: 30
    database:
      port: 5432
      pool_size: 10
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  database-password: <base64-encoded-password>
```

**`kubernetes/deployment.yaml`**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: app
        image: user-service:latest
        env:
        # Runtime configuration overrides
        - name: KONFIGO_KEY_database.host
          value: "prod-db.company.com"
        - name: KONFIGO_KEY_database.password
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-password
        - name: KONFIGO_KEY_app.environment
          value: "production"
        volumeMounts:
        - name: base-config
          mountPath: /app/configs
        command:
        - /bin/sh
        - -c
        - |
          konfigo -s /app/configs/base.yaml -of /app/runtime-config.json
          exec /app/user-service --config /app/runtime-config.json
      volumes:
      - name: base-config
        configMap:
          name: app-base-config
```

---

## Team Configuration Management

**Problem**: Manage shared team configurations while allowing individual customization.

### Solution Structure

```
team-configs/
├── shared/
│   ├── base.yaml          # Team defaults
│   ├── tools.yaml         # Shared tool configs
│   └── standards.yaml     # Team standards
├── environments/
│   ├── dev.yaml          # Development settings
│   └── test.yaml         # Testing settings
├── personal/
│   └── .gitignore        # Don't commit personal configs
└── scripts/
    └── setup-dev.sh      # Setup script
```

### Files

**`shared/base.yaml`** - Team defaults:
```yaml
development:
  database:
    host: "team-dev-db.company.com"
    port: 5432
  api:
    base_url: "https://api-dev.company.com"
    timeout: 30
    
tools:
  editor:
    tab_size: 2
    format_on_save: true
  linter:
    max_line_length: 100
    enforce_style: true
```

**`personal/alice.yaml`** - Individual overrides:
```yaml
development:
  database:
    host: "localhost"  # Local database for development
  api:
    timeout: 60        # Longer timeout for debugging
    
tools:
  editor:
    tab_size: 4        # Personal preference
```

**Setup script (`scripts/setup-dev.sh`)**:
```bash
#!/bin/bash

USER=${1:-$(whoami)}
PERSONAL_CONFIG="personal/${USER}.yaml"
DEV_CONFIG="environments/dev.yaml"

echo "Setting up development environment for $USER..."

# Create personal config if it doesn't exist
if [[ ! -f "$PERSONAL_CONFIG" ]]; then
  echo "Creating personal config: $PERSONAL_CONFIG"
  cp personal/template.yaml "$PERSONAL_CONFIG"
fi

# Generate development configuration
konfigo \
  -s shared/base.yaml,shared/tools.yaml,"$DEV_CONFIG","$PERSONAL_CONFIG" \
  -of ".dev-config.json"

echo "✅ Development configuration ready: .dev-config.json"
echo "💡 Edit $PERSONAL_CONFIG to customize your settings"
```

**Usage**:
```bash
# Setup for specific team member
./scripts/setup-dev.sh alice

# Setup for current user
./scripts/setup-dev.sh
```

---

## Microservices Configuration

**Problem**: Generate similar configurations for multiple microservices with service-specific customization.

### Solution: Template + Variables

**`template.yaml`** - Service template:
```yaml
service:
  name: "${SERVICE_NAME}"
  port: "${SERVICE_PORT}"
  version: "1.0.0"
  
database:
  host: "${DB_HOST}"
  database: "${DB_NAME}"
  pool_size: "${DB_POOL_SIZE}"
  
monitoring:
  enabled: true
  port: "${METRICS_PORT}"
  service_name: "${SERVICE_NAME}"
  
resources:
  memory: "${MEMORY_LIMIT}"
  cpu: "${CPU_LIMIT}"
```

**`services.yaml`** - Service definitions:
```yaml
konfigo_forEach:
  - name: "user-service"
    vars:
      SERVICE_NAME: "user-service"
      SERVICE_PORT: 8001
      DB_NAME: "users"
      DB_POOL_SIZE: 20
      METRICS_PORT: 9001
      MEMORY_LIMIT: "512Mi"
      CPU_LIMIT: "500m"
      
  - name: "order-service"
    vars:
      SERVICE_NAME: "order-service"
      SERVICE_PORT: 8002
      DB_NAME: "orders"
      DB_POOL_SIZE: 30
      METRICS_PORT: 9002
      MEMORY_LIMIT: "1Gi"
      CPU_LIMIT: "1000m"
      
  - name: "payment-service"
    vars:
      SERVICE_NAME: "payment-service"
      SERVICE_PORT: 8003
      DB_NAME: "payments"
      DB_POOL_SIZE: 15
      METRICS_PORT: 9003
      MEMORY_LIMIT: "256Mi"
      CPU_LIMIT: "250m"
```

**`schema.yaml`** - Processing rules:
```yaml
vars:
  - name: "DB_HOST"
    value: "postgres.company.com"  # Common database host

transforms:
  - path: "service.name"
    setValue: "${SERVICE_NAME}"
  - path: "service.port"
    setValue: "${SERVICE_PORT}"
  - path: "database.host"
    setValue: "${DB_HOST}"
  - path: "database.database"
    setValue: "${DB_NAME}"

validation:
  - path: "service.port"
    type: "number"
    min: 8000
    max: 9000
  - path: "service.name"
    type: "string"
    regex: "^[a-z-]+$"
```

**Generate all service configs**:
```bash
konfigo -s template.yaml -S schema.yaml -V services.yaml -od configs/

# Creates:
# configs/user-service.yaml
# configs/order-service.yaml  
# configs/payment-service.yaml
```

**Kubernetes deployment generation**:
```bash
# Generate Kubernetes manifests
for config in configs/*.yaml; do
  service=$(basename "$config" .yaml)
  envsubst < k8s-template.yaml > "k8s/${service}-deployment.yaml"
done
```

---

## Feature Flag Management

**Problem**: Manage feature flags across environments with easy toggle capabilities.

### Solution Structure

**`features/base.yaml`** - Default feature states:
```yaml
features:
  auth:
    oauth_login: true
    two_factor: false
    social_login: true
    
  payments:
    stripe_integration: true
    paypal_integration: false
    crypto_payments: false
    
  ui:
    new_dashboard: false
    dark_mode: true
    analytics_widget: true
    
  experimental:
    ai_recommendations: false
    performance_mode: false
```

**`features/environments/prod.yaml`** - Production overrides:
```yaml
features:
  auth:
    two_factor: true     # Enable 2FA in production
    
  payments:
    paypal_integration: true  # Enable PayPal in prod
    
  ui:
    new_dashboard: true  # New dashboard ready for prod
    
  experimental:
    # Keep experimental features disabled in prod
```

**`features/rollout.yaml`** - Gradual rollout configuration:
```yaml
konfigo_forEach:
  - name: "canary"
    vars:
      ROLLOUT_PERCENTAGE: 5
      NEW_DASHBOARD: true
      AI_RECOMMENDATIONS: true
      
  - name: "beta"
    vars:
      ROLLOUT_PERCENTAGE: 25
      NEW_DASHBOARD: true
      AI_RECOMMENDATIONS: false
      
  - name: "stable"
    vars:
      ROLLOUT_PERCENTAGE: 100
      NEW_DASHBOARD: false
      AI_RECOMMENDATIONS: false
```

**`features/schema.yaml`** - Feature processing:
```yaml
transforms:
  - path: "features.ui.new_dashboard"
    setValue: "${NEW_DASHBOARD}"
  - path: "features.experimental.ai_recommendations"
    setValue: "${AI_RECOMMENDATIONS}"
  - path: "rollout.percentage"
    setValue: "${ROLLOUT_PERCENTAGE}"

validation:
  - path: "rollout.percentage"
    type: "number"
    min: 0
    max: 100
```

**Usage**:
```bash
# Generate standard environment configs
konfigo -s features/base.yaml,features/environments/prod.yaml -of prod-features.json

# Generate rollout configurations
konfigo -s features/base.yaml -S features/schema.yaml -V features/rollout.yaml -od rollout/

# Override specific features via environment
KONFIGO_KEY_features.experimental.ai_recommendations=true \
  konfigo -s prod-features.json -of runtime-features.json
```

---

## Configuration Validation Pipeline

**Problem**: Ensure all configurations are valid before deployment through automated validation.

### CI/CD Pipeline Integration

**`.github/workflows/config-validation.yml`**:
```yaml
name: Configuration Validation

on:
  pull_request:
    paths:
      - 'configs/**'
      - 'schemas/**'

jobs:
  validate-configs:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Download Konfigo
      run: |
        curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64 -o konfigo
        chmod +x konfigo
        sudo mv konfigo /usr/local/bin/
    
    - name: Validate Development Configs
      run: |
        ./scripts/validate-environment.sh dev
        
    - name: Validate Staging Configs
      run: |
        ./scripts/validate-environment.sh staging
        
    - name: Validate Production Configs
      run: |
        ./scripts/validate-environment.sh prod
        
    - name: Generate Test Configs
      run: |
        ./scripts/generate-all-configs.sh --validate-only
```

**Validation script (`scripts/validate-environment.sh`)**:
```bash
#!/bin/bash
set -e

ENVIRONMENT=$1
if [[ -z "$ENVIRONMENT" ]]; then
  echo "Usage: $0 <environment>"
  exit 1
fi

echo "🔍 Validating $ENVIRONMENT configuration..."

BASE_CONFIG="configs/base.yaml"
ENV_CONFIG="configs/environments/${ENVIRONMENT}.yaml"
SCHEMA="schemas/validation.yaml"

# Check if files exist
for file in "$BASE_CONFIG" "$ENV_CONFIG" "$SCHEMA"; do
  if [[ ! -f "$file" ]]; then
    echo "❌ Missing file: $file"
    exit 1
  fi
done

# Validate configuration
if konfigo -s "$BASE_CONFIG,$ENV_CONFIG" -S "$SCHEMA" --validate-only; then
  echo "✅ $ENVIRONMENT configuration is valid"
else
  echo "❌ $ENVIRONMENT configuration validation failed"
  exit 1
fi

# Test merge output
echo "📋 Testing merge output..."
OUTPUT=$(mktemp)
konfigo -s "$BASE_CONFIG,$ENV_CONFIG" -S "$SCHEMA" -of "$OUTPUT"

# Verify required fields exist
REQUIRED_FIELDS=("app.name" "database.host" "app.port")
for field in "${REQUIRED_FIELDS[@]}"; do
  if ! jq -e ".$field" "$OUTPUT" > /dev/null; then
    echo "❌ Missing required field: $field"
    cat "$OUTPUT"
    rm "$OUTPUT"
    exit 1
  fi
done

rm "$OUTPUT"
echo "✅ All required fields present"
```

**Comprehensive validation schema (`schemas/validation.yaml`)**:
```yaml
# Input validation
inputSchema:
  strict: true

# Core validation rules
validation:
  # Application configuration
  - path: "app.name"
    required: true
    type: "string"
    regex: "^[a-z][a-z0-9-]*$"
    
  - path: "app.port"
    required: true
    type: "number"
    min: 1024
    max: 65535
    
  - path: "app.environment"
    required: true
    type: "string"
    enum: ["development", "staging", "production"]
    
  # Database configuration
  - path: "database.host"
    required: true
    type: "string"
    minLength: 3
    
  - path: "database.port"
    required: true
    type: "number"
    min: 1
    max: 65535
    
  # Security validation
  - path: "database.ssl"
    required: true
    type: "boolean"
    
  - path: "security.api_keys"
    required: false
    type: "array"
    minItems: 1

# Post-processing validation
transforms:
  - path: "metadata.validated_at"
    setValue: "${TIMESTAMP}"
  - path: "metadata.validation_version"
    setValue: "1.0"

# Immutable production settings
immutable:
  - "app.name"
  - "security"
```

### Local Development Integration

**Pre-commit hook (`.git/hooks/pre-commit`)**:
```bash
#!/bin/bash

echo "🔍 Validating configurations before commit..."

# Find all changed config files
CHANGED_CONFIGS=$(git diff --cached --name-only | grep -E '\.(yaml|yml|json|toml)$' | grep configs/ || true)

if [[ -z "$CHANGED_CONFIGS" ]]; then
  echo "ℹ️  No configuration files changed"
  exit 0
fi

# Validate each changed config
for config in $CHANGED_CONFIGS; do
  echo "Validating $config..."
  
  # Determine environment from path
  if [[ "$config" =~ environments/([^/]+)\.yaml$ ]]; then
    ENV="${BASH_REMATCH[1]}"
    if ! ./scripts/validate-environment.sh "$ENV"; then
      echo "❌ Configuration validation failed for $config"
      exit 1
    fi
  fi
done

echo "✅ All configurations valid"
```

---

## Best Practices Summary

### 1. **Structure Your Configs**
- Use consistent directory structures
- Separate base configs from environment overrides
- Keep schemas and validation rules in version control

### 2. **Automate Everything**
- Use scripts for common operations
- Integrate validation into CI/CD pipelines
- Set up pre-commit hooks for immediate feedback

### 3. **Security First**
- Never commit secrets to version control
- Use environment variables for sensitive data
- Validate configurations before deployment

### 4. **Documentation**
- Document your configuration patterns
- Provide examples for team members
- Maintain README files for complex setups

### 5. **Testing**
- Test configuration merges in development
- Validate against schemas regularly
- Use dry-run modes when available

These recipes provide proven patterns for real-world Konfigo usage. Adapt them to your specific needs and build upon these foundations!
