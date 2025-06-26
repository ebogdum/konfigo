# User Guide: Common Tasks

Welcome to the Konfigo User Guide! This section is organized around what you want to accomplish, not just feature descriptions. Jump to any task to get step-by-step instructions.

## Quick Navigation

### 🚀 **Getting Things Done**
- **[Convert configuration formats](#convert-configuration-formats)** - JSON ↔ YAML ↔ TOML ↔ ENV
- **[Merge configurations from multiple sources](#merge-configurations)** - Combine files intelligently
- **[Use environment variables](#use-environment-variables)** - Runtime overrides and flexibility
- **[Validate configurations](#validate-configurations)** - Ensure correctness before deployment
- **[Generate multiple outputs](#generate-multiple-outputs)** - Batch processing for multiple environments

### 📚 **Reference & Advanced**
- **[CLI Reference](./cli-reference.md)** - Complete command-line options
- **[Recipes & Examples](./recipes.md)** - Real-world patterns and solutions
- **[Environment Variables](./environment-variables.md)** - Detailed environment integration

---

## Convert Configuration Formats

**Goal**: Transform configuration files between JSON, YAML, TOML, and ENV formats.

### Quick Examples

```bash
# YAML to JSON
konfigo -s config.yaml -oj -of config.json

# JSON to TOML  
konfigo -s config.json -ot -of config.toml

# ENV to YAML
konfigo -s .env -oy -of config.yaml

# Multiple formats at once
konfigo -s config.yaml -oj -ot -of base  # Creates base.json and base.toml
```

### Step-by-Step: Converting Legacy .env to Modern YAML

**Step 1**: Start with your .env file
```bash
# .env
DATABASE_HOST=localhost
DATABASE_PORT=5432
APP_DEBUG=true
FEATURE_AUTH=enabled
```

**Step 2**: Convert to YAML
```bash
konfigo -s .env -oy -of config.yaml
```

**Step 3**: Verify the result
```yaml
# config.yaml (generated)
APP_DEBUG: true
DATABASE_HOST: localhost
DATABASE_PORT: 5432
FEATURE_AUTH: enabled
```

### Advanced: Format-Specific Options

```bash
# Pretty-printed JSON
konfigo -s config.yaml -oj --json-indent=2

# Compact JSON
konfigo -s config.yaml -oj --json-compact

# YAML with specific style
konfigo -s config.json -oy --yaml-flow=false
```

### When to Use This
- **Legacy system modernization**: Convert old .env files to structured YAML
- **Tool integration**: Different tools prefer different formats
- **Team preferences**: Some teams prefer YAML, others JSON
- **Deployment requirements**: Kubernetes prefers YAML, APIs often use JSON

---

## Merge Configurations

**Goal**: Combine multiple configuration files into a single, unified configuration.

### Quick Examples

```bash
# Basic merge: base + environment
konfigo -s base.yaml,prod.yaml

# Multiple files with output
konfigo -s defaults.json,environment.yaml,local.toml -of final.json

# Recursive directory merge
konfigo -r -s configs/ -of merged.yaml
```

### Step-by-Step: Multi-Environment Setup

**Step 1**: Create your base configuration
```yaml
# base.yaml
app:
  name: "my-service"
  port: 8080
  timeout: 30
database:
  host: "localhost"
  port: 5432
  timeout: 10
```

**Step 2**: Create environment-specific overrides
```yaml
# prod.yaml
app:
  port: 9090
  timeout: 60
database:
  host: "prod-db.company.com"
  ssl: true
  pool_size: 20
```

**Step 3**: Merge for production
```bash
konfigo -s base.yaml,prod.yaml -of prod-config.json
```

**Step 4**: Verify the result
```json
{
  "app": {
    "name": "my-service",
    "port": 9090,
    "timeout": 60
  },
  "database": {
    "host": "prod-db.company.com",
    "port": 5432,
    "ssl": true,
    "pool_size": 20,
    "timeout": 10
  }
}
```

### Understanding Merge Order

```bash
# Left to right precedence (right wins)
konfigo -s base.yaml,env.yaml,local.yaml
#        ^lowest    ^medium   ^highest precedence
```

### Advanced Merging Patterns

```bash
# With environment variable overrides
KONFIGO_KEY_app.port=3000 konfigo -s base.yaml,prod.yaml

# Case-insensitive merging
konfigo -c -s config1.yaml,config2.yaml

# Verbose output to debug merging
konfigo -v -s file1.yaml,file2.yaml
```

### When to Use This
- **Multi-environment deployments**: dev, staging, prod configurations
- **Team collaboration**: Shared base + individual customizations
- **Microservices**: Common settings + service-specific overrides
- **Feature toggles**: Base config + feature-specific settings

---

## Use Environment Variables

**Goal**: Override configuration values at runtime using environment variables.

### Quick Examples

```bash
# Override any configuration key
KONFIGO_KEY_database.host=prod-db konfigo -s config.yaml

# Multiple overrides
KONFIGO_KEY_app.port=9000 KONFIGO_KEY_app.debug=false konfigo -s config.yaml

# Perfect for containers
docker run -e KONFIGO_KEY_database.url=$DATABASE_URL myapp
```

### Step-by-Step: Container Deployment

**Step 1**: Start with base configuration
```yaml
# app.yaml
app:
  port: 8080
  environment: "development"
database:
  host: "localhost"
  port: 5432
  name: "myapp"
```

**Step 2**: Use environment variables for runtime values
```bash
# Set production values via environment
export KONFIGO_KEY_app.environment="production"
export KONFIGO_KEY_database.host="prod-db.company.com"
export KONFIGO_KEY_database.port=3306
```

**Step 3**: Generate runtime configuration
```bash
konfigo -s app.yaml -of runtime-config.json
```

**Step 4**: Verify environment overrides worked
```json
{
  "app": {
    "port": 8080,
    "environment": "production"
  },
  "database": {
    "host": "prod-db.company.com",
    "port": 3306,
    "name": "myapp"
  }
}
```

### Key Format Rules

```bash
# Dot notation for nested keys
KONFIGO_KEY_database.connection.pool_size=10

# Array indices (if supported)
KONFIGO_KEY_servers.0.host=server1.com

# Complex nested structures
KONFIGO_KEY_features.auth.providers.oauth.client_id=abc123
```

### Integration Patterns

```bash
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - env:
        - name: KONFIGO_KEY_database.url
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

```bash
# Docker Compose
services:
  app:
    environment:
      - KONFIGO_KEY_app.port=8080
      - KONFIGO_KEY_database.host=db
```

### When to Use This
- **Containerized applications**: Docker, Kubernetes deployments
- **CI/CD pipelines**: Different values per environment
- **Secret management**: Passwords, API keys from secure stores
- **Testing**: Quick value overrides without file changes

---

## Validate Configurations

**Goal**: Ensure your configurations are correct before deployment using schemas.

### Quick Examples

```bash
# Validate against schema
konfigo -s config.yaml -S validation.schema.yaml

# Validate during merge
konfigo -s base.yaml,prod.yaml -S schema.yaml -of validated.json

# Input validation only
konfigo -s config.yaml -S schema.yaml --validate-only
```

### Step-by-Step: Configuration Validation

**Step 1**: Create a validation schema
```yaml
# validation.schema.yaml
inputSchema:
  strict: true
  
vars: []

transforms: []

validation:
  - path: "app.port"
    required: true
    type: "number"
    min: 1024
    max: 65535
  - path: "app.name"
    required: true
    type: "string"
    minLength: 3
  - path: "database.host"
    required: true
    type: "string"
    regex: "^[a-zA-Z0-9.-]+$"
```

**Step 2**: Test with valid configuration
```yaml
# valid-config.yaml
app:
  name: "my-service"
  port: 8080
database:
  host: "prod-db.company.com"
```

```bash
konfigo -s valid-config.yaml -S validation.schema.yaml
# ✅ Validation passes, outputs merged config
```

**Step 3**: Test with invalid configuration
```yaml
# invalid-config.yaml
app:
  name: "x"  # Too short
  port: 80   # Too low
database:
  host: "invalid host!"  # Invalid characters
```

```bash
konfigo -s invalid-config.yaml -S validation.schema.yaml
# ❌ Validation fails with detailed error messages
```

### Common Validation Rules

```yaml
validation:
  # Required field
  - path: "app.name"
    required: true
    
  # Type checking
  - path: "app.port"
    type: "number"
    
  # Value constraints
  - path: "app.replicas"
    type: "number"
    min: 1
    max: 100
    
  # String validation
  - path: "app.environment"
    type: "string"
    enum: ["dev", "staging", "prod"]
    
  # Pattern matching
  - path: "database.host"
    type: "string"
    regex: "^[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    
  # Length constraints
  - path: "app.description"
    type: "string"
    minLength: 10
    maxLength: 200
```

### When to Use This
- **Production deployments**: Prevent invalid configurations from causing outages
- **Team collaboration**: Ensure everyone follows configuration standards
- **CI/CD pipelines**: Automated validation before deployment
- **Complex configurations**: Catch errors in large, complex config files

---

## Generate Multiple Outputs

**Goal**: Create multiple configuration files from a single source for different environments or services.

### Quick Examples

```bash
# Multiple formats from one source
konfigo -s config.yaml -oj -ot -oy -of app  # Creates app.json, app.toml, app.yaml

# Batch processing with forEach
konfigo -s base.yaml -S schema.yaml -V environments.yaml

# Directory output
konfigo -s template.yaml -S schema.yaml -V services.yaml -od outputs/
```

### Step-by-Step: Multi-Environment Generation

**Step 1**: Create base template
```yaml
# base.yaml
app:
  name: "my-service"
  port: 8080
database:
  host: "${DATABASE_HOST}"
  port: "${DATABASE_PORT}"
environment: "${ENVIRONMENT}"
```

**Step 2**: Create variables for environments
```yaml
# environments.yaml
konfigo_forEach:
  - name: "dev"
    vars:
      DATABASE_HOST: "dev-db.company.com"
      DATABASE_PORT: 5432
      ENVIRONMENT: "development"
  - name: "staging"
    vars:
      DATABASE_HOST: "staging-db.company.com"
      DATABASE_PORT: 5432
      ENVIRONMENT: "staging"
  - name: "prod"
    vars:
      DATABASE_HOST: "prod-db.company.com"
      DATABASE_PORT: 3306
      ENVIRONMENT: "production"
```

**Step 3**: Create processing schema
```yaml
# schema.yaml
vars:
  - name: "DATABASE_HOST"
    required: true
  - name: "DATABASE_PORT"
    required: true
  - name: "ENVIRONMENT"
    required: true

transforms:
  - path: "database.host"
    setValue: "${DATABASE_HOST}"
  - path: "database.port"
    setValue: "${DATABASE_PORT}"
  - path: "environment"
    setValue: "${ENVIRONMENT}"
```

**Step 4**: Generate all environments
```bash
konfigo -s base.yaml -S schema.yaml -V environments.yaml
```

**Output**: Creates `dev.yaml`, `staging.yaml`, and `prod.yaml` with environment-specific values.

### Advanced: Service Generation

**For microservices deployment**:
```yaml
# services.yaml
konfigo_forEach:
  - name: "user-service"
    vars:
      SERVICE_NAME: "user-service"
      SERVICE_PORT: 8001
      DB_NAME: "users"
  - name: "order-service"
    vars:
      SERVICE_NAME: "order-service"
      SERVICE_PORT: 8002
      DB_NAME: "orders"
```

### When to Use This
- **Multi-environment deployment**: Generate configs for dev, staging, prod
- **Microservices**: Create similar configs for multiple services
- **A/B testing**: Generate configurations for different feature variations
- **Infrastructure as Code**: Generate Terraform or Kubernetes manifests

---

## Next Steps

### **Master the CLI**
- **[CLI Reference](./cli-reference.md)** - Complete command-line documentation
- **[Environment Variables Guide](./environment-variables.md)** - Deep dive into environment integration

### **Advanced Features**
- **[Schema Guide](../schema/)** - Unlock validation, transformation, and generation
- **[Recipes & Examples](./recipes.md)** - Real-world patterns and solutions

### **Get Help**
- **[Troubleshooting](../reference/troubleshooting.md)** - Common issues and solutions
- **[FAQ](../reference/faq.md)** - Frequently asked questions

The User Guide gives you practical solutions for real tasks. Start with what you need to accomplish, then explore advanced features as your needs grow!
        iii. **Generators**: Data generation rules are applied.
        iv. **Transformers**: Data transformation rules are applied.
        v.  **Global Variable Substitution**: `${VAR_NAME}` placeholders are substituted throughout the configuration.
        vi. **Validation**: The processed configuration is validated against the `validate` rules in the schema.
        vii. **Output Schema Filtering**: If `outputSchema` is defined, the configuration is filtered.
    d.  **(If Batch Mode with `konfigo_forEach` in `-V` file and Schema Provided)**:
        *   Steps c.ii through c.vii are performed for *each iteration* defined in `konfigo_forEach`, using a deep copy of the merged configuration from step 4.b and iteration-specific variables. Each iteration produces its own output file.
    e.  **Output**: The final configuration (or multiple configurations in batch mode) is written to the specified output file(s) or to stdout in the chosen format.

This guide, along with the [Schema Documentation](../schema/index.md), aims to provide you with all the information needed to master Konfigo.
