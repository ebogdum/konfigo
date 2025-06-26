# Advanced Schema Features

This section covers Konfigo's advanced schema capabilities for complex configuration management scenarios.

## Input and Output Schema Validation

### Input Schema Validation

Input schema validation ensures that your merged configuration meets expected structure requirements before processing begins.

```yaml
# input-requirements.yaml
type: "object"
required: ["service", "database"]
properties:
  service:
    type: "object"
    required: ["name", "port"]
    properties:
      name:
        type: "string"
        minLength: 3
      port:
        type: "number"
        minimum: 1024
  database:
    type: "object"
    required: ["host"]
    properties:
      host:
        type: "string"
      port:
        type: "number"
        minimum: 1
        maximum: 65535

# main-schema.yaml
inputSchema:
  path: "./input-requirements.yaml"
  strict: true  # Fail if extra properties are present

# ... rest of schema processing
```

**Benefits:**
- Catch configuration issues early in the pipeline
- Ensure required fields are present before processing
- Validate data types and constraints upfront
- Fail fast if configuration is malformed

### Output Schema Filtering

Output schema filtering controls what appears in the final configuration, useful for generating public APIs or removing sensitive data.

```yaml
# public-api-schema.yaml
type: "object"
properties:
  service:
    type: "object"
    properties:
      name:
        type: "string"
      port:
        type: "number"
      endpoints:
        type: "array"
  features:
    type: "object"

# main-schema.yaml
outputSchema:
  path: "./public-api-schema.yaml"
  strict: false  # Allow extra properties in the schema

# Any fields not in the output schema are removed from final result
```

**Use Cases:**
- Generate public configuration subsets
- Remove internal/debug information
- Create API-specific configuration views
- Ensure only approved fields are exposed

## Immutable Fields Protection

Immutable fields prevent critical configuration from being overridden by later sources or environment variables.

```yaml
# Protect critical configuration paths
immutable:
  - "service.name"                    # Service identity
  - "database.credentials.username"   # Security credentials
  - "security.certificates"          # Security configurations
  - "audit.settings"                 # Compliance settings

# Example: These values cannot be changed after initial merge
# base.yaml
service:
  name: "critical-service"
database:
  credentials:
    username: "app_user"

# override.yaml (these changes will be ignored)
service:
  name: "different-name"  # IGNORED - immutable
database:
  credentials:
    username: "hacker"    # IGNORED - immutable
    password: "new-pass"  # ALLOWED - not immutable
```

**Important Notes:**
- Environment variables (`KONFIGO_KEY_*`) can still override immutable paths (for operational flexibility)
- Immutable protection applies during file merging, not environment overrides
- Use for security-critical or compliance-required settings

## Complex Variable Resolution

### Path-Based Variables

Extract values from the configuration itself for use in other parts:

```yaml
vars:
  - name: "SERVICE_NAME"
    fromPath: "service.name"
    defaultValue: "unknown-service"
  
  - name: "FULL_VERSION"
    fromPath: "metadata.version"
    defaultValue: "dev"
  
  - name: "NAMESPACE"
    fromPath: "deployment.namespace"
    defaultValue: "default"

# Use extracted values in generators
generators:
  - type: "concat"
    targetPath: "deployment.image"
    format: "registry.example.com/${SERVICE_NAME}:${FULL_VERSION}"
```

### Conditional Variables

Variables can be set conditionally based on other configuration values:

```yaml
vars:
  - name: "LOG_LEVEL"
    fromEnv: "LOG_LEVEL"
    defaultValue: "info"
  
  - name: "DEBUG_MODE"
    value: "true"
    condition: '${LOG_LEVEL} == "debug"'
  
  - name: "MONITORING_ENABLED"
    value: "true"
    condition: '${ENVIRONMENT} != "development"'

transform:
  # Apply transformations conditionally
  - type: "setValue"
    path: "features.debug"
    value: "${DEBUG_MODE}"
    condition: '${DEBUG_MODE} == "true"'
```

## Advanced Transformations

### Conditional Transformations

Apply transformations only when certain conditions are met:

```yaml
transform:
  # Remove debug features in production
  - type: "deleteKey"
    path: "features.debug"
    condition: '${ENVIRONMENT} == "production"'
  
  # Set security headers in production
  - type: "setValue"
    path: "security.strict_transport_security"
    value: "max-age=31536000"
    condition: '${ENVIRONMENT} == "production"'
  
  # Enable verbose logging in development
  - type: "setValue"
    path: "logging.level"
    value: "debug"
    condition: '${ENVIRONMENT} == "development"'
```

### Nested Object Transformations

Work with complex nested structures:

```yaml
transform:
  # Transform all service URLs to HTTPS
  - type: "transformValues"
    path: "services.*.url"
    transformation:
      type: "regex"
      pattern: "^http://"
      replacement: "https://"
  
  # Add prefix to all environment variables
  - type: "addKeyPrefix"
    path: "environment.variables"
    prefix: "APP_"
  
  # Normalize all service names
  - type: "transformValues"
    path: "services.*.name"
    transformation:
      type: "changeCase"
      case: "lower"
```

## Complex Validation Rules

### Cross-Field Validation

Validate relationships between different configuration fields:

```yaml
validate:
  # Ensure database port matches database type
  - path: "database"
    rules:
      custom: |
        if database.type == "postgresql":
          assert database.port in [5432, 5433, 5434]
        elif database.type == "mysql":
          assert database.port in [3306, 3307]
  
  # Validate resource constraints
  - path: "resources"
    rules:
      custom: |
        memory_mb = parse_memory(resources.limits.memory)
        cpu_cores = parse_cpu(resources.limits.cpu)
        assert memory_mb / cpu_cores >= 512  # At least 512MB per CPU core
```

### Environment-Specific Validation

Different validation rules for different environments:

```yaml
validate:
  # Production-specific validation
  - path: "security.tls.enabled"
    rules:
      required: true
      type: "boolean"
      value: true
    condition: '${ENVIRONMENT} == "production"'
  
  # Development allows insecure settings
  - path: "security.tls.enabled"
    rules:
      type: "boolean"
    condition: '${ENVIRONMENT} == "development"'
  
  # Staging requires specific resource limits
  - path: "resources.limits.memory"
    rules:
      required: true
      regex: "^[1-9][0-9]*[GM]i$"
    condition: '${ENVIRONMENT} == "staging"'
```

## Advanced Generators

### Template-Based Generation

Generate complex structures using templates:

```yaml
generators:
  - type: "template"
    targetPath: "kubernetes.deployment"
    template: |
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: ${SERVICE_NAME}
        namespace: ${NAMESPACE}
        labels:
          app: ${SERVICE_NAME}
          version: ${VERSION}
      spec:
        replicas: ${REPLICAS}
        selector:
          matchLabels:
            app: ${SERVICE_NAME}
        template:
          metadata:
            labels:
              app: ${SERVICE_NAME}
              version: ${VERSION}
          spec:
            containers:
            - name: ${SERVICE_NAME}
              image: ${IMAGE}
              ports:
              - containerPort: ${PORT}
              env:
              {{- range $key, $value := .environment }}
              - name: {{ $key }}
                value: "{{ $value }}"
              {{- end }}
```

### Dynamic List Generation

Generate lists based on configuration data:

```yaml
generators:
  - type: "generateList"
    targetPath: "services.endpoints"
    template:
      - name: "health"
        path: "/health"
        port: "${service.port}"
      - name: "metrics"
        path: "/metrics"
        port: "${service.monitoring.port}"
      - name: "ready"
        path: "/ready"
        port: "${service.port}"
    condition: '${service.type} == "web"'
```

## Schema Composition and Inheritance

### Schema Inheritance

Build complex schemas from simpler base schemas:

```yaml
# base-schema.yaml
extends: []
vars:
  - name: "ENVIRONMENT"
    fromEnv: "NODE_ENV"
    defaultValue: "development"

validate:
  - path: "service.name"
    rules:
      required: true
      type: "string"

# web-service-schema.yaml
extends: ["base-schema.yaml"]
vars:
  - name: "PORT"
    fromEnv: "PORT"
    defaultValue: "3000"

validate:
  - path: "service.port"
    rules:
      required: true
      type: "number"
      min: 3000
      max: 8000

# final-schema.yaml
extends: ["web-service-schema.yaml"]
transform:
  - type: "setValue"
    path: "metadata.processed_by"
    value: "konfigo"
```

### Modular Schema Components

Break schemas into reusable components:

```yaml
# components/validation.yaml
validate:
  - path: "service.name"
    rules:
      required: true
      minLength: 3
      regex: "^[a-z][a-z0-9-]*$"

# components/security.yaml
validate:
  - path: "security.tls.enabled"
    rules:
      type: "boolean"
      value: true

transform:
  - type: "setValue"
    path: "security.headers.strict_transport_security"
    value: "max-age=31536000"

# main-schema.yaml
includes:
  - "components/validation.yaml"
  - "components/security.yaml"

# Additional schema-specific processing
generators:
  - type: "concat"
    targetPath: "service.url"
    format: "https://${service.name}.${DOMAIN}"
```

These advanced features enable sophisticated configuration management workflows that can handle complex enterprise requirements while maintaining clarity and maintainability.
