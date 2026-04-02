# Konfigo - Merge, Transform, and Validate Configuration Files

<p align="center">
  <img src="konfigo_logo.png" alt="Konfigo - Configuration file merger, converter, and validator" width="200"/>
</p>

<p align="center">
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://github.com/ebogdum/konfigo/releases"><img src="https://img.shields.io/github/v/release/ebogdum/konfigo" alt="Latest Release"></a>
  <a href="https://goreportcard.com/report/github.com/ebogdum/konfigo"><img src="https://goreportcard.com/badge/github.com/ebogdum/konfigo" alt="Go Report Card"></a>
</p>

<p align="center">
  <strong>A fast, schema-driven CLI tool for merging, converting, validating, and generating configuration files across JSON, YAML, TOML, ENV, and INI formats.</strong>
</p>

---

Konfigo solves the problem of managing configuration across multiple files, formats, and environments. Instead of writing custom scripts to merge YAML files, convert JSON to TOML, or validate config values before deployment, Konfigo handles it all in a single command.

**Use cases:**
- Merge base + environment-specific configs for deployment pipelines
- Convert between JSON, YAML, TOML, and ENV formats
- Validate configuration values (types, ranges, patterns, required fields)
- Generate UUIDs, timestamps, and computed values at build time
- Batch-generate configs for multiple services or environments from a single template
- Override any config value via environment variables without editing files

## Quick Start

### Install

Download a pre-built binary from [Releases](https://github.com/ebogdum/konfigo/releases) (Linux, macOS, Windows, FreeBSD, OpenBSD, NetBSD - amd64/arm64), or build from source:

```bash
go install github.com/ebogdum/konfigo/cmd/konfigo@latest
```

### Merge configuration files

```bash
konfigo -s base.yaml,production.yaml -of config.json
```

Later files override earlier files. Nested objects are deep-merged.

### Convert between formats

```bash
konfigo -s config.yaml -oj        # YAML to JSON (stdout)
konfigo -s config.json -oy        # JSON to YAML
konfigo -s config.toml -oe        # TOML to ENV
konfigo -s legacy.ini -of out.yaml  # INI to YAML
```

### Override values from environment variables

```bash
export KONFIGO_KEY_database.host=prod-db.example.com
export KONFIGO_KEY_database.port=5432
konfigo -s config.yaml -of config.json
```

Environment variables always take highest precedence over file sources.

## Features

### Multi-Format Configuration Merging

Merge any combination of JSON, YAML, TOML, ENV, and INI files with deterministic precedence:

```bash
# Later sources override earlier ones; objects deep-merge, arrays replace
konfigo -s defaults.yaml,environment.json,secrets.env

# Recursive directory discovery
konfigo -s configs/ -r -of merged.yaml

# Case-sensitive key matching
konfigo -s config1.yaml,config2.yaml -c

# Array union merge (deduplicated instead of replaced)
konfigo -s base.yaml,override.yaml -m
```

### Schema-Driven Processing

Apply a schema to validate, transform, and generate configuration values:

```yaml
# schema.yaml
vars:
  - name: ENVIRONMENT
    fromEnv: NODE_ENV
    defaultValue: development
  - name: VERSION
    value: "2.0.0"

generators:
  - type: concat
    targetPath: service.url
    format: "https://{host}:${PORT}"
    sources:
      host: service.host

  - type: timestamp
    targetPath: metadata.buildTime
    format: rfc3339

  - type: random
    targetPath: session.secret
    format: uuid

transform:
  - type: renameKey
    from: legacy.dbHost
    to: database.host
  - type: changeCase
    path: service.name
    case: snake
  - type: setValue
    path: app.version
    value: "${VERSION}"

validate:
  - path: database.port
    rules:
      required: true
      type: number
      min: 1
      max: 65535
  - path: service.name
    rules:
      type: string
      minLength: 3
      regex: "^[a-z][a-z0-9-]*$"
  - path: environment
    rules:
      enum: [development, staging, production]

immutable:
  - database.credentials
  - security.apiKey
```

```bash
konfigo -s config.yaml -S schema.yaml -V vars/production.yaml -of deploy.json
```

### Variable Substitution

Three-tier variable precedence: environment (`KONFIGO_VAR_*`) > variables file (`-V`) > schema defaults.

```yaml
# config.yaml
database:
  url: "postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}"
  pool_size: "${POOL_SIZE}"
```

```bash
export KONFIGO_VAR_DB_HOST=prod-db.example.com
konfigo -s config.yaml -V vars.yaml -oj
```

### Batch Processing

Generate multiple config files from a single template using `forEach`:

```yaml
# batch-vars.yaml
CLUSTER: k8s-prod

forEach:
  items:
    - SERVICE_NAME: api
      PORT: "8080"
      REPLICAS: "3"
    - SERVICE_NAME: worker
      PORT: "8081"
      REPLICAS: "5"
  output:
    filenamePattern: "deploy/${SERVICE_NAME}/config-${ITEM_INDEX}.yaml"
```

```bash
konfigo -s base.yaml -S schema.yaml -V batch-vars.yaml
# Creates: deploy/api/config-0.yaml, deploy/worker/config-1.yaml
```

### Input and Output Schema Validation

Validate structure before processing, filter sensitive data from output:

```yaml
# schema.yaml
inputSchema:
  path: "../schemas/required-structure.json"
  strict: true   # reject unexpected keys

outputSchema:
  path: "../schemas/public-api.json"
  strict: false  # include only defined keys, drop extras
```

### Immutable Path Protection

Protect critical config values from being overwritten by later sources, transformers, or generators. Child paths are automatically protected:

```yaml
immutable:
  - database.credentials   # also protects database.credentials.username, etc.
  - security.apiKey
```

## CLI Reference

```
konfigo [flags] -s <sources>
```

| Flag | Description |
|------|-------------|
| `-s <paths>` | Comma-separated source files/directories. Use `-` for stdin |
| `-r` | Recursive directory discovery |
| `-c` | Case-sensitive key matching (default: case-insensitive) |
| `-m` | Merge arrays by union with deduplication instead of replacing |
| `-S, --schema <path>` | Schema file for processing (JSON, YAML, or TOML) |
| `-V, --vars-file <path>` | Variables file (high-priority variables + forEach) |
| `-sj / -sy / -st / -se` | Force input format: JSON / YAML / TOML / ENV |
| `-of <path>` | Output file (format from extension, or use with `-oX`) |
| `-oj / -oy / -ot / -oe` | Output format: JSON / YAML / TOML / ENV |
| `-v` | Verbose logging (INFO) |
| `-d` | Debug logging (DEBUG + INFO) |
| `-h` | Show help |

### Environment Variables

| Pattern | Purpose | Example |
|---------|---------|---------|
| `KONFIGO_KEY_*` | Override any config value (highest precedence) | `KONFIGO_KEY_database.host=prod-db` |
| `KONFIGO_VAR_*` | Set schema variables (highest variable precedence) | `KONFIGO_VAR_ENV=production` |

## Supported Formats

| Format | Input | Output | Extensions |
|--------|-------|--------|------------|
| JSON | Yes | Yes | `.json` |
| YAML | Yes | Yes | `.yaml`, `.yml` |
| TOML | Yes | Yes | `.toml` |
| ENV | Yes | Yes | `.env` |
| INI | Yes | No | `.ini` |

## Processing Pipeline

```
Source Files -> Parse -> Merge -> [Schema Processing] -> Output
                                       |
                         Input Schema Validation
                         Variable Resolution
                         Generator Execution
                         Transformer Execution
                         Variable Substitution
                         Validator Execution
                         Output Schema Filtering
```

## Real-World Examples

### CI/CD: Environment-Specific Deployment Config

```bash
# Generate production config from layered sources
konfigo -s config/base.yaml,config/production.yaml \
  -S schemas/deploy.yaml \
  -V vars/production.yaml \
  -of deploy/config.json
```

### Docker: Build-Time Configuration

```dockerfile
FROM golang:1.22 AS config
COPY configs/ /configs/
COPY schemas/ /schemas/
RUN konfigo -s /configs/ -r -S /schemas/app.yaml -of /app-config.json

FROM alpine:latest
COPY --from=config /app-config.json /etc/app/config.json
```

### Kubernetes: Generate Per-Service Configs

```bash
konfigo -s k8s/base.yaml -S k8s/schema.yaml -V k8s/services-batch.yaml
# Generates one config file per service defined in forEach items
```

### Config Validation in Pre-Commit Hooks

```bash
konfigo -s config/ -r -S schemas/validation.yaml > /dev/null
# Exit code 0 = valid, 1 = errors found
```

### Format Migration

```bash
# Convert legacy INI to YAML
konfigo -s legacy.ini -of modern.yaml

# Convert .env to JSON for application consumption
konfigo -s .env -of config.json
```

## Documentation

Full documentation: **[ebogdum.github.io/konfigo](https://ebogdum.github.io/konfigo/)**

- [User Guide](https://ebogdum.github.io/konfigo/guide/) - Task-oriented walkthroughs
- [Schema Guide](https://ebogdum.github.io/konfigo/schema/) - Variables, generators, transformers, validators
- [CLI Reference](https://ebogdum.github.io/konfigo/guide/cli-reference) - Complete flag documentation
- [Batch Processing](https://ebogdum.github.io/konfigo/features/batch-processing) - forEach and multi-output generation

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Konfigo is licensed under the [MIT License](./LICENSE).
