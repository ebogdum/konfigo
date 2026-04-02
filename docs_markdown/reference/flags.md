# Command Flags Reference

This page provides comprehensive documentation for all Konfigo command-line flags and options.

## Command-Line Interface

### Global Options

| Flag | Description | Default |
|------|-------------|---------|
| `-h` | Show help message | - |
| `-v` | Enable informational (INFO) logging | false |
| `-d` | Enable debug (DEBUG + INFO) logging, overrides `-v` | false |
| `-c` | Use case-sensitive key matching | false (case-insensitive) |
| `-m` | Merge arrays by union with deduplication | false (arrays replaced) |
| `-r` | Recursively search subdirectories | false |

### Source Input Options

| Flag | Description | Notes |
|------|-------------|-------|
| `-s` | Comma-separated list of source files/directories | Required. Use `-` for stdin |
| `-sj` | Force input parsing as JSON | Required for stdin |
| `-sy` | Force input parsing as YAML | Required for stdin |
| `-st` | Force input parsing as TOML | Required for stdin |
| `-se` | Force input parsing as ENV | Required for stdin |

### Schema Processing Options

| Flag | Long Form | Description | Notes |
|------|-----------|-------------|-------|
| `-S` | `--schema` | Path to schema file | JSON, YAML, or TOML only |
| `-V` | `--vars-file` | Path to variables file | High-priority variable definitions |

### Output Format Options

| Flag | Description | Notes |
|------|-------------|-------|
| `-oj` | Output in JSON format | - |
| `-oy` | Output in YAML format | Default when no format specified |
| `-ot` | Output in TOML format | - |
| `-oe` | Output in ENV format | - |
| `-of` | Write output to file | Extension determines format |

## Environment Variables

### Runtime Configuration Overrides

| Variable Pattern | Description | Example |
|------------------|-------------|---------|
| `KONFIGO_KEY_*` | Override any configuration key (highest precedence, but subject to immutable paths) | `KONFIGO_KEY_app.port=8080` |
| `KONFIGO_VAR_*` | Define schema variables (highest variable precedence) | `KONFIGO_VAR_DATABASE_HOST=prod-db.com` |

## Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Operation completed successfully |
| 1 | Error | Any error (parsing, validation, file I/O, schema processing, etc.) |

## File Format Support

### Input Formats

| Format | Extensions | Notes |
|--------|------------|-------|
| JSON | `.json` | Standard JSON (no comments) |
| YAML | `.yaml`, `.yml` | Single-document YAML |
| TOML | `.toml` | Full TOML v1.0.0 support |
| ENV | `.env` | `KEY=value` pairs, `#` comments, dot notation for nesting |
| INI | `.ini` | Sections become nested maps, input only |

### Output Formats

| Format | Notes |
|--------|-------|
| JSON | Pretty-printed with 2-space indentation |
| YAML | 2-space indentation |
| TOML | Standard TOML with sections |
| ENV | Flattened to `UPPERCASE_UNDERSCORE=value`, sorted keys, auto-quoting |

INI is input-only; it cannot be used as an output format.

## Error Handling

### Common Error Patterns

- **File not found**: Check file paths and permissions
- **Parse errors**: Validate syntax with format-specific tools
- **Schema errors**: Verify schema structure and required fields
- **Validation failures**: Check data types and constraints

### Debug Information

Use `-v` or `-d` flags to get detailed information about:
- File discovery and loading
- Merge order and precedence
- Schema processing steps
- Variable substitution
- Validation results

## Performance Considerations

### File Size Limits
- Configuration files are limited to **50 MiB** each; larger files are rejected with an error
- Schema files are limited to **10 MiB**
- Consider splitting large configurations into multiple files and merging

### Memory Usage
- Konfigo loads all source files into memory
- Schema processing requires additional memory for transformations
- Array union merging (`-m`) with very large arrays (>1,000 elements per side) automatically skips deduplication to avoid O(n*m) performance degradation

## Integration Patterns

### CI/CD Integration
```bash
# Typical pipeline usage
konfigo -s base.yaml,env/${ENVIRONMENT}.yaml -S schema.yaml -of config.json
```

### Container Integration
```dockerfile
# Multi-stage build pattern
FROM golang:1.22 AS config-builder
COPY . /src
WORKDIR /src
RUN go build -o /usr/local/bin/konfigo ./cmd/konfigo
COPY configs/ /configs/
RUN konfigo -s /configs/base.yaml,/configs/prod.yaml -of /tmp/final.json

FROM alpine:latest
COPY --from=config-builder /tmp/final.json /app/config.json
```

### Library Integration
Konfigo is designed as a CLI tool but can be integrated into build processes, deployment scripts, and configuration management workflows.

## Schema Version

The `apiVersion` field in schema files is informational only. Konfigo does not enforce or validate this field. Common values used in the community:
- `konfigo/v1alpha1`
