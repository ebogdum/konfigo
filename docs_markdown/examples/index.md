# Examples and Tutorials

This section provides hands-on examples and step-by-step tutorials for common Konfigo scenarios.

## Quick Examples

### Basic Configuration Merging

**Files:**
```yaml
# base.yaml
app:
  name: "my-service"
  port: 8080
database:
  host: "localhost"
  timeout: 30

# production.yaml  
app:
  port: 9090
database:
  host: "prod-db.example.com"
  ssl: true
```

**Command:**
```bash
konfigo -s base.yaml,production.yaml
```

**Result:**
```yaml
app:
  name: "my-service"
  port: 9090                    # overridden
database:
  host: "prod-db.example.com"   # overridden
  timeout: 30                   # preserved
  ssl: true                     # added
```

### Format Conversion

```bash
# YAML to JSON
konfigo -s config.yaml -oj

# JSON to TOML  
konfigo -s config.json -ot

# Multiple formats
konfigo -s config.yaml -oj -ot -oe
```

### Environment Variable Override

```bash
# Override specific values
export KONFIGO_KEY_app.port=3000
export KONFIGO_KEY_database.host="override.example.com"

konfigo -s config.yaml
```

## Step-by-Step Tutorials

### Tutorial 1: Setting Up Multi-Environment Configuration

**Goal:** Create a configuration system that works across development, staging, and production environments.

**Step 1: Create base configuration**
```yaml
# configs/base.yaml
application:
  name: "web-api"
  version: "1.0.0"
  logging:
    level: "info"
    format: "json"

database:
  port: 5432
  timeout: 30
  pool:
    min: 5
    max: 20

features:
  auth: true
  monitoring: true
```

**Step 2: Create environment overrides**
```yaml
# configs/development.yaml
application:
  logging:
    level: "debug"
database:
  host: "localhost"
  pool:
    max: 5
features:
  debug: true

# configs/production.yaml
application:
  logging:
    level: "warn"
database:
  host: "prod-db.internal"
  ssl: true
  pool:
    min: 10
    max: 50
features:
  debug: false
```

**Step 3: Create validation schema**
```yaml
# schemas/validation.yaml
validate:
  - path: "application.name"
    rules:
      required: true
      type: "string"
      minLength: 3
  
  - path: "database.host"
    rules:
      required: true
      type: "string"
  
  - path: "database.port"
    rules:
      type: "number"
      min: 1024
      max: 65535
```

**Step 4: Generate configurations**
```bash
# Development
konfigo -s configs/base.yaml,configs/development.yaml \
       -S schemas/validation.yaml \
       -of dist/development.yaml

# Production  
konfigo -s configs/base.yaml,configs/production.yaml \
       -S schemas/validation.yaml \
       -of dist/production.yaml
```

### Tutorial 2: Schema-Driven Configuration Processing

**Goal:** Use schema processing to transform and validate configurations for Kubernetes deployment.

**Step 1: Create application configuration**
```yaml
# app-config.yaml
service:
  name: "web-api"
  image: "myregistry/web-api"
  port: 8080
  replicas: 3

resources:
  cpu: "100m"
  memory: "256Mi"
  
environment:
  - name: "NODE_ENV"
    value: "production"
  - name: "PORT"
    value: "8080"
```

**Step 2: Create Kubernetes transformation schema**
```yaml
# k8s-schema.yaml
vars:
  - name: "NAMESPACE"
    fromEnv: "K8S_NAMESPACE"
    defaultValue: "default"
  - name: "IMAGE_TAG"
    fromEnv: "IMAGE_TAG"
    defaultValue: "latest"

generators:
  - type: "concat"
    targetPath: "spec.template.spec.containers[0].image"
    format: "${service.image}:${IMAGE_TAG}"

transform:
  # Create Kubernetes Deployment structure
  - type: "setValue"
    path: "apiVersion"
    value: "apps/v1"
  - type: "setValue"
    path: "kind"
    value: "Deployment"
  - type: "setValue"
    path: "metadata.name"
    value: "${service.name}"
  - type: "setValue"
    path: "metadata.namespace"
    value: "${NAMESPACE}"
  
  # Configure deployment spec
  - type: "setValue"
    path: "spec.replicas"
    value: "${service.replicas}"
  - type: "setValue"
    path: "spec.selector.matchLabels.app"
    value: "${service.name}"
  
  # Configure pod template
  - type: "setValue"
    path: "spec.template.metadata.labels.app"
    value: "${service.name}"
  - type: "setValue"
    path: "spec.template.spec.containers[0].name"
    value: "${service.name}"
  - type: "setValue"
    path: "spec.template.spec.containers[0].ports[0].containerPort"
    value: "${service.port}"
  
  # Set resource limits
  - type: "setValue"
    path: "spec.template.spec.containers[0].resources.requests.cpu"
    value: "${resources.cpu}"
  - type: "setValue"
    path: "spec.template.spec.containers[0].resources.requests.memory"
    value: "${resources.memory}"

validate:
  - path: "spec.replicas"
    rules:
      type: "number"
      min: 1
      max: 10
  - path: "spec.template.spec.containers[0].image"
    rules:
      required: true
      type: "string"
      regex: "^[a-zA-Z0-9._/-]+:[a-zA-Z0-9._-]+$"
```

**Step 3: Generate Kubernetes deployment**
```bash
export K8S_NAMESPACE="production"
export IMAGE_TAG="v1.2.3"

konfigo -s app-config.yaml \
       -S k8s-schema.yaml \
       -of k8s-deployment.yaml
```

### Tutorial 3: Batch Configuration Generation

**Goal:** Generate multiple service configurations from a template.

**Step 1: Create service template**
```yaml
# template.yaml
service:
  logging:
    level: "info"
    format: "json"
database:
  timeout: 30
  pool:
    min: 5
features:
  monitoring: true
```

**Step 2: Create batch processing schema**
```yaml
# batch-schema.yaml
vars:
  - name: "SERVICE_NAME"
    fromPath: "metadata.name"
  - name: "SERVICE_PORT"
    fromPath: "metadata.port"
  - name: "DATABASE_NAME"
    fromPath: "metadata.database"

generators:
  - type: "concat"
    targetPath: "service.name"
    format: "${SERVICE_NAME}"
  - type: "concat"
    targetPath: "service.port"
    format: "${SERVICE_PORT}"
  - type: "concat"
    targetPath: "database.name"
    format: "${DATABASE_NAME}"
  - type: "concat"
    targetPath: "database.url"
    format: "postgresql://user:pass@db:5432/${DATABASE_NAME}"

validate:
  - path: "service.port"
    rules:
      type: "number"
      min: 3000
      max: 9000
```

**Step 3: Create batch variables file**
```yaml
# batch-vars.yaml
konfigo_forEach:
  items:
    - metadata:
        name: "user-service"
        port: 3001
        database: "users"
    - metadata:
        name: "order-service"
        port: 3002
        database: "orders"
    - metadata:
        name: "payment-service"
        port: 3003
        database: "payments"

  output:
    filenamePattern: "configs/${SERVICE_NAME}.yaml"
    format: "yaml"
```

**Step 4: Generate all service configurations**
```bash
konfigo -s template.yaml \
       -S batch-schema.yaml \
       -V batch-vars.yaml

# Results in:
# configs/user-service.yaml
# configs/order-service.yaml  
# configs/payment-service.yaml
```

## Common Patterns

### Pattern 1: Environment-Specific Secrets

**Problem:** Need to inject different secrets per environment without hardcoding them.

**Solution:**
```bash
# Development
export KONFIGO_VAR_DB_PASSWORD="dev-password"
export KONFIGO_VAR_API_KEY="dev-api-key"

# Production (from secret management)
export KONFIGO_VAR_DB_PASSWORD="$(vault kv get -field=password secret/prod/db)"
export KONFIGO_VAR_API_KEY="$(vault kv get -field=key secret/prod/api)"

# Same command for all environments
konfigo -s base.yaml,env-specific.yaml -S schema-with-secrets.yaml
```

### Pattern 2: Configuration Validation

**Problem:** Need to ensure configurations meet compliance requirements.

**Solution:**
```yaml
# compliance-schema.yaml
validate:
  # Ensure HTTPS only
  - path: "services.*.url"
    rules:
      type: "string"
      regex: "^https://"
  
  # No privileged ports
  - path: "*.port"
    rules:
      type: "number"
      min: 1024
  
  # Required security settings
  - path: "security.tls.enabled"
    rules:
      required: true
      type: "boolean"
      value: true
```

### Pattern 3: Legacy Migration

**Problem:** Need to migrate from old configuration format to new one.

**Solution:**
```yaml
# migration-schema.yaml
transform:
  # Rename old keys to new structure
  - type: "renameKey"
    from: "app_name"
    to: "application.name"
  - type: "renameKey"
    from: "db_host"
    to: "database.host"
  
  # Update format
  - type: "changeCase"
    path: "environment"
    case: "lower"
  
  # Generate new required fields
  - type: "setValue"
    path: "metadata.migrated"
    value: true
```

These examples demonstrate the progression from simple configuration merging to complex, schema-driven processing workflows. Each pattern can be adapted and combined to meet your specific requirements.
