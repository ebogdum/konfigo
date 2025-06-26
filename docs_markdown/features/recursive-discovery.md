# Recursive Discovery

Konfigo can automatically discover and process configuration files within directory trees using the `-r` (recursive) flag. This feature is essential for managing complex projects with distributed configuration files.

## Basic Recursive Discovery

### Enable Recursive Mode
```bash
# Recursively discover files in directory
konfigo -s config/ -r

# Recursive discovery with specific file patterns
konfigo -s src/,config/ -r

# Combine with other sources
konfigo -s base.json,config/ -r
```

### File Discovery Rules

Konfigo automatically discovers files based on supported extensions:
- `.json` - JSON files
- `.yml`, `.yaml` - YAML files  
- `.toml` - TOML files
- `.env` - Environment files
- `.ini` - INI files

### Exclusion Patterns

Files and directories are automatically excluded:
- Hidden files (starting with `.`)
- Common build directories (`node_modules`, `dist`, `build`, `.git`)
- Backup files (ending with `~`, `.bak`)

## Directory Structure Examples

Based on `test/recursive-discovery/` test cases:

### Project Configuration
```
project/
├── base.json                 # Root configuration
├── environments/
│   ├── development.yml       # Development overrides
│   ├── staging.yml          # Staging overrides  
│   └── production.yml       # Production overrides
├── services/
│   ├── api/
│   │   ├── config.toml      # API service config
│   │   └── secrets.env      # API secrets
│   └── database/
│       ├── config.yml       # Database config
│       └── migrations.json  # Migration settings
└── features/
    ├── auth.yml            # Authentication config
    ├── logging.yml         # Logging config
    └── monitoring.toml     # Monitoring config
```

```bash
# Discover all configuration files
konfigo -s project/ -r

# Files processed in alphabetical order:
# project/base.json
# project/environments/development.yml
# project/environments/production.yml  
# project/environments/staging.yml
# project/features/auth.yml
# project/features/logging.yml
# project/features/monitoring.toml
# project/services/api/config.toml
# project/services/api/secrets.env
# project/services/database/config.yml
# project/services/database/migrations.json
```

### Microservices Configuration
```
services/
├── gateway/
│   ├── base.yml
│   ├── routes.json
│   └── environments/
│       ├── dev.yml
│       └── prod.yml  
├── auth-service/
│   ├── config.toml
│   ├── oauth.yml
│   └── providers/
│       ├── google.yml
│       └── github.yml
└── user-service/
    ├── config.yml
    ├── database.yml
    └── cache.yml
```

```bash
# Process all microservice configurations
konfigo -s services/ -r
```

## File Processing Order

Files are processed in **alphabetical order** by full path, ensuring deterministic merging:

```
config/
├── 01-base.yml        # Processed first
├── 02-database.yml    # Processed second  
├── 99-overrides.yml   # Processed last
└── modules/
    ├── auth.yml       # Processed after parent directory
    └── logging.yml    # Processed last
```

### Controlling Order with Naming

Use prefixes to control processing order:
```
config/
├── 00-defaults.yml    # Base defaults
├── 10-database.yml    # Database configuration
├── 20-services.yml    # Service configuration
├── 30-features.yml    # Feature toggles
└── 99-overrides.yml   # Final overrides
```

## Combining Sources

Recursive discovery can be combined with explicit file sources:

```bash
# Explicit base + recursive discovery  
konfigo -s base.json,config/ -r

# Multiple directories with recursive discovery
konfigo -s shared/,project-specific/ -r

# Mix explicit files and recursive directories
konfigo -s global.yml,services/,local-overrides.json -r
```

## Real-World Examples

### Kubernetes Configuration Management
```
k8s-configs/
├── base/
│   ├── namespace.yml
│   ├── secrets.yml
│   └── configmaps.yml
├── services/
│   ├── frontend/
│   │   ├── deployment.yml
│   │   ├── service.yml  
│   │   └── ingress.yml
│   └── backend/
│       ├── deployment.yml
│       ├── service.yml
│       └── database.yml
└── environments/
    ├── development/
    │   ├── resources.yml
    │   └── replicas.yml
    └── production/
        ├── resources.yml
        ├── replicas.yml
        └── monitoring.yml
```

```bash
# Generate environment-specific configuration
konfigo -s k8s-configs/base/,k8s-configs/services/,k8s-configs/environments/production/ -r
```

### Application Configuration
```
config/
├── defaults/
│   ├── logging.yml
│   ├── database.yml  
│   └── security.yml
├── features/
│   ├── payments/
│   │   ├── stripe.yml
│   │   └── paypal.yml
│   ├── auth/
│   │   ├── oauth.yml
│   │   └── saml.yml
│   └── notifications/
│       ├── email.yml
│       └── sms.yml
├── environments/
│   ├── local.yml
│   ├── staging.yml
│   └── production.yml
└── overrides/
    └── local-dev.yml
```

```bash
# Development environment
konfigo -s config/defaults/,config/features/,config/environments/local.yml,config/overrides/ -r

# Production environment  
konfigo -s config/defaults/,config/features/,config/environments/production.yml -r
```

## Error Handling

### Directory Not Found
```bash
konfigo -s missing-directory/ -r
# Error: failed to stat path missing-directory/: no such file or directory
```

### Permission Denied
```bash
konfigo -s /root/config/ -r
# Warning: Cannot access /root/config/secret.yml: permission denied
# Processing continues with accessible files
```

### Parse Errors
```bash
konfigo -s config/ -r
# Warning: Skipping file config/invalid.yml due to parse error: invalid YAML syntax
# Processing continues with remaining files
```

## Performance Considerations

### Large Directory Trees
- Recursive discovery scans entire directory trees
- Consider using specific subdirectories for better performance
- Use explicit file lists for large projects when possible

### File System Optimization
```bash
# More specific paths for better performance
konfigo -s config/core/,config/env/production/ -r

# Instead of scanning everything
konfigo -s config/ -r
```

## Best Practices

### Directory Organization
1. **Hierarchical Structure**: Organize configs by logical groupings
2. **Naming Conventions**: Use prefixes to control merge order
3. **Environment Separation**: Keep environment-specific configs in subdirectories
4. **Feature Isolation**: Group related configurations together

### File Naming
1. **Descriptive Names**: Use clear, descriptive filenames
2. **Order Prefixes**: Use numeric prefixes when order matters
3. **Consistent Extensions**: Use appropriate file extensions for auto-detection
4. **Avoid Conflicts**: Don't use conflicting names across directories

### Project Structure Examples

#### Clean Architecture
```
config/
├── 00-base/           # Foundation configuration
├── 10-infrastructure/ # Database, cache, queues
├── 20-services/      # Business logic services  
├── 30-features/      # Feature-specific config
├── 40-integrations/  # External service config
└── 99-environment/   # Environment overrides
```

#### Domain-Driven Design
```
config/
├── shared/           # Cross-cutting concerns
├── user-domain/      # User management
├── order-domain/     # Order processing
├── payment-domain/   # Payment handling
└── notification-domain/ # Notifications
```

#### Microservices
```
config/
├── infrastructure/   # Shared infrastructure
├── gateway/         # API gateway config
├── services/        # Individual service configs
│   ├── user-service/
│   ├── order-service/
│   └── payment-service/
└── environments/    # Environment-specific
    ├── development/
    ├── staging/
    └── production/
```

## Integration with Other Features

### Schema Processing
```bash
# Recursive discovery with schema processing
konfigo -s config/ -r -S schema.yml
```

### Environment Variables
```bash
# Combine recursive discovery with environment overrides
export KONFIGO_KEY_environment=production
konfigo -s config/ -r
```

### Output Generation
```bash
# Process recursively and output to file
konfigo -s config/ -r -of final-config.json
```

## Test Coverage

Recursive discovery is tested comprehensively in `test/recursive-discovery/`:
- Multiple directory levels
- Mixed file formats
- File processing order verification
- Error condition handling
- Integration with merging logic
- Performance with large directory trees
