# Batch Processing

Konfigo's batch processing feature using `konfigo_forEach` enables generating multiple configuration outputs from a single template. This is particularly powerful for creating deployment configurations, environment-specific files, or any scenario requiring multiple similar outputs with variations.

## Overview

Batch processing allows you to:
- Generate multiple configuration files from one template
- Iterate over data sets or files
- Use variables unique to each iteration
- Create structured output hierarchies
- Automate configuration generation for multiple environments/services

## `konfigo_forEach` Structure

The `konfigo_forEach` directive is defined in the variables file (`-V` flag):

```yaml
# variables.yml
konfigo_forEach:
  # Data source (choose one)
  items: [...]          # Inline array of objects
  itemFiles: [...]      # Array of file paths to load
  
  # Output configuration
  output:
    filenamePattern: "..."  # Template for output filenames
    format: "yaml"          # Optional: output format override

# Global variables available to all iterations
globalVar1: value1
globalVar2: value2
```

## Data Sources

### Inline Items
Define iteration data directly in the variables file:

```yaml
# batch-vars.yml
konfigo_forEach:
  items:
    - name: web-frontend
      image: nginx:1.21
      replicas: 2
      namespace: applications
    - name: api-backend  
      image: node:16-alpine
      replicas: 4
      namespace: backend
    - name: worker
      image: worker:latest
      replicas: 1
      namespace: processing
  output:
    filenamePattern: "deployments/${namespace}/${name}-deployment-${ITEM_INDEX}.yaml"

# Global variables
cluster: k8s-prod
```

### External Item Files
Load iteration data from separate files:

```yaml
# batch-vars.yml
konfigo_forEach:
  itemFiles:
    - "items/web-frontend.yml"
    - "items/api-backend.yml" 
    - "items/worker.yml"
  output:
    filenamePattern: "deployments/${name}-${ITEM_FILE_BASENAME}.yaml"
```

```yaml
# items/web-frontend.yml
name: web-frontend
image: nginx:1.21
replicas: 2
namespace: applications
```

```yaml
# items/api-backend.yml
name: api-backend
image: node:16-alpine  
replicas: 4
namespace: backend
environment: production
```

## Built-in Variables

Each iteration has access to special built-in variables:

- `${ITEM_INDEX}` - Zero-based iteration index (0, 1, 2, ...)
- `${ITEM_FILE_BASENAME}` - Filename without extension (for itemFiles mode)

## Real-World Example

Based on `test/batch/` test cases:

### Deployment Generation

**Base Configuration**:
```yaml
# base-config.yml
application:
  name: base-app
  version: 1.0.0
server:
  host: localhost
  port: 8080
database:
  host: localhost
  port: 5432
  ssl: false
```

**Schema Definition**:
```yaml
# deployment-schema.yml
vars:
  - name: "CLUSTER_NAME"
    value: "k8s-prod"
  - name: "NAMESPACE"
    defaultValue: "default"

transform:
  - type: "setValue"
    path: "deployment.name"
    value: "${SERVICE_NAME}-deployment"
  - type: "setValue"
    path: "deployment.replicas"
    value: "${REPLICAS}"
  - type: "setValue"
    path: "deployment.image"
    value: "${IMAGE_NAME}:${IMAGE_TAG}"
  - type: "setValue"
    path: "deployment.namespace"
    value: "${NAMESPACE}"
  - type: "setValue"
    path: "deployment.cluster"
    value: "${CLUSTER_NAME}"
```

**Batch Variables**:
```yaml
# deployments-batch.yml
konfigo_forEach:
  items:
    - SERVICE_NAME: web-frontend
      IMAGE_NAME: nginx
      IMAGE_TAG: "1.21"
      REPLICAS: "2"
      NAMESPACE: applications
    - SERVICE_NAME: api-backend
      IMAGE_NAME: node
      IMAGE_TAG: "16-alpine"
      REPLICAS: "4"
      NAMESPACE: backend
    - SERVICE_NAME: worker
      IMAGE_NAME: worker
      IMAGE_TAG: latest
      REPLICAS: "1"
      NAMESPACE: processing
  output:
    filenamePattern: "deployments/${NAMESPACE}/${SERVICE_NAME}-deployment-${ITEM_INDEX}.yaml"
    format: "yaml"

# Global variables
ENVIRONMENT: production
```

### Execution
```bash
konfigo -s base-config.yml -S deployment-schema.yml -V deployments-batch.yml
```

### Generated Files

**`deployments/applications/web-frontend-deployment-0.yaml`**:
```yaml
application:
  name: base-app
  version: 1.0.0
database:
  host: localhost
  port: 5432
  ssl: false
deployment:
  cluster: k8s-prod
  image: nginx:1.21
  name: web-frontend-deployment
  namespace: applications
  replicas: "2"
server:
  host: localhost
  port: 8080
```

**`deployments/backend/api-backend-deployment-1.yaml`**:
```yaml
application:
  name: base-app
  version: 1.0.0
database:
  host: localhost
  port: 5432
  ssl: false
deployment:
  cluster: k8s-prod
  image: node:16-alpine
  name: api-backend-deployment
  namespace: backend
  replicas: "4"
server:
  host: localhost
  port: 8080
```

## Advanced Patterns

### Environment Matrix Generation
```yaml
# environment-matrix.yml
konfigo_forEach:
  items:
    - environment: development
      database_host: dev-db.internal
      api_url: https://api-dev.example.com
      replicas: 1
      debug: true
    - environment: staging
      database_host: staging-db.internal
      api_url: https://api-staging.example.com
      replicas: 2
      debug: false
    - environment: production
      database_host: prod-db.internal
      api_url: https://api.example.com
      replicas: 5
      debug: false
  output:
    filenamePattern: "configs/${environment}/app-config.json"
    format: "json"
```

### Service Configuration Generation
```yaml
# services-batch.yml  
konfigo_forEach:
  itemFiles:
    - "services/user-service.yml"
    - "services/order-service.yml"
    - "services/payment-service.yml"
    - "services/notification-service.yml"
  output:
    filenamePattern: "k8s/${service_name}/${ITEM_FILE_BASENAME}-config.yml"

# Global service defaults
default_replicas: 3
default_namespace: services
monitoring_enabled: true
```

### Multi-Environment Deployment
```yaml
# multi-env-batch.yml
konfigo_forEach:
  items:
    - env: dev
      cluster: dev-cluster
      namespace: development
      replicas: 1
      image_tag: latest
      resources:
        cpu: 100m
        memory: 256Mi
    - env: staging
      cluster: staging-cluster
      namespace: staging
      replicas: 2
      image_tag: v1.2.3
      resources:
        cpu: 200m
        memory: 512Mi
    - env: prod
      cluster: prod-cluster
      namespace: production
      replicas: 5
      image_tag: v1.2.3
      resources:
        cpu: 500m
        memory: 1Gi
  output:
    filenamePattern: "environments/${env}/${cluster}/deployment.yaml"
```

## Variable Resolution

### Variable Priority (per iteration)
1. **Item-specific variables** (highest precedence)
2. **Global variables** from variables file
3. **Schema variables** 
4. **Environment variables** (`KONFIGO_VAR_*`)

### Variable Substitution in Patterns
Filename patterns support full variable substitution:

```yaml
output:
  filenamePattern: "clusters/${cluster}/namespaces/${namespace}/${service}-${env}-${ITEM_INDEX}.${format}"
  format: "yaml"

# With variables:
# cluster: k8s-prod
# namespace: backend  
# service: api
# env: production
# format: yaml

# Generates: clusters/k8s-prod/namespaces/backend/api-production-1.yaml
```

## Format Control

### Automatic Format Detection
Format is determined by file extension in `filenamePattern`:

```yaml
output:
  filenamePattern: "configs/${env}/app.json"    # JSON output
  filenamePattern: "configs/${env}/app.yaml"    # YAML output
  filenamePattern: "configs/${env}/app.toml"    # TOML output
```

### Explicit Format Override
```yaml
output:
  filenamePattern: "configs/${env}/app-config"  # No extension
  format: "json"                                # Explicit format
```

## Error Handling

### Invalid Batch Configuration
```yaml
konfigo_forEach:
  items: [...]
  itemFiles: [...]  # Error: cannot have both items and itemFiles
```

### Missing Required Fields
```yaml
konfigo_forEach:
  items: [...]
  # Error: output.filenamePattern is required
```

### File Processing Errors
```bash
# Item file not found
konfigo_forEach:
  itemFiles:
    - "missing-file.yml"  # Error: file not found

# Invalid item file format
konfigo_forEach:
  itemFiles:
    - "invalid.yml"       # Warning: skipped due to parse error
```

## Best Practices

### Organization
1. **Separate Concerns**: Keep templates, data, and schemas separate
2. **Descriptive Naming**: Use clear variable names and file patterns
3. **Version Control**: Track templates and data separately
4. **Documentation**: Document variable requirements and outputs

### Performance
1. **Minimize Items**: Use only necessary iteration data
2. **Efficient Patterns**: Avoid deeply nested output directories
3. **Batch Size**: Consider memory usage with large item sets
4. **Parallel Processing**: Konfigo processes iterations efficiently

### Maintenance
1. **Schema Validation**: Use schemas to validate item structure
2. **Default Values**: Provide sensible defaults for optional variables
3. **Error Handling**: Test with invalid data scenarios
4. **Output Cleanup**: Clean up old generated files when patterns change

## Integration Examples

### CI/CD Pipeline
```bash
#!/bin/bash
# Generate deployment configurations for all environments
konfigo -s base.yml -S k8s-schema.yml -V environments-batch.yml

# Deploy each generated configuration
for env_dir in output/environments/*/; do
  env=$(basename "$env_dir")
  kubectl apply -f "$env_dir" --context="$env-cluster"
done
```

### Terraform Configuration
```yaml
# terraform-batch.yml
konfigo_forEach:
  items:
    - region: us-east-1
      instance_type: t3.medium
      availability_zones: ["us-east-1a", "us-east-1b"]
    - region: us-west-2
      instance_type: t3.large
      availability_zones: ["us-west-2a", "us-west-2b", "us-west-2c"]
  output:
    filenamePattern: "terraform/${region}/main.tf.json"
    format: "json"
```

### Docker Compose Services
```yaml
# services-batch.yml
konfigo_forEach:
  itemFiles:
    - "services/web.yml"
    - "services/api.yml"
    - "services/worker.yml"
    - "services/db.yml"
  output:
    filenamePattern: "docker/${ITEM_FILE_BASENAME}/docker-compose.yml"
```

## Test Coverage

Batch processing is thoroughly tested in `test/batch/`:
- Multiple iteration sources (items vs itemFiles)
- Variable resolution and substitution
- Filename pattern processing
- Format detection and override
- Error condition handling
- Integration with schema processing
- Complex real-world scenarios
