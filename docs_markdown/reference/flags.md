# Command Flags Reference

This page provides comprehensive documentation for all Konfigo command-line flags and options.

## Command-Line Interface

### Global Options

| Flag | Long Form | Description | Default |
|------|-----------|-------------|---------|
| `-h` | `--help` | Show help message | - |
| `-v` | `--verbose` | Enable informational logging | false |
| `-d` | `--debug` | Enable debug logging | false |
| `-c` | `--case-sensitive` | Use case-sensitive key matching | false |
| `-r` | `--recursive` | Recursively search subdirectories | false |

### Source Input Options

| Flag | Long Form | Description | Notes |
|------|-----------|-------------|-------|
| `-s` | `--sources` | Comma-separated list of source files/directories | Required |
| `-sj` | `--source-json` | Force input parsing as JSON | For stdin or ambiguous files |
| `-sy` | `--source-yaml` | Force input parsing as YAML | For stdin or ambiguous files |
| `-st` | `--source-toml` | Force input parsing as TOML | For stdin or ambiguous files |
| `-se` | `--source-env` | Force input parsing as ENV | For stdin or ambiguous files |

### Schema Processing Options

| Flag | Long Form | Description | Notes |
|------|-----------|-------------|-------|
| `-S` | `--schema` | Path to schema file | Enables advanced processing |
| `-V` | `--vars-file` | Path to variables file | High-priority variable definitions |

### Output Format Options

| Flag | Long Form | Description | Notes |
|------|-----------|-------------|-------|
| `-oj` | `--output-json` | Output in JSON format | - |
| `-oy` | `--output-yaml` | Output in YAML format | - |
| `-ot` | `--output-toml` | Output in TOML format | - |
| `-oe` | `--output-env` | Output in ENV format | - |
| `-of` | `--output-file` | Write output to file | Extension determines format |

## Environment Variables

### Runtime Configuration Overrides

| Variable Pattern | Description | Example |
|------------------|-------------|---------|
| `KONFIGO_KEY_*` | Override any configuration key | `KONFIGO_KEY_app.port=8080` |
| `KONFIGO_VAR_*` | Define schema variables | `KONFIGO_VAR_DATABASE_HOST=prod-db.com` |

### Processing Control

| Variable | Description | Default |
|----------|-------------|---------|
| `KONFIGO_LOG_LEVEL` | Set logging level (ERROR, WARN, INFO, DEBUG) | ERROR |
| `KONFIGO_CONFIG_PATH` | Default search paths for config files | Current directory |

## Exit Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 0 | Success | Operation completed successfully |
| 1 | General Error | Invalid arguments or processing error |
| 2 | File Error | Source file not found or unreadable |
| 3 | Parse Error | Invalid syntax in source or schema files |
| 4 | Validation Error | Schema validation failed |
| 5 | Schema Error | Invalid or malformed schema |

## File Format Support

### Input Formats

| Format | Extensions | Parser | Notes |
|--------|------------|---------|-------|
| JSON | `.json`, `.jsonc` | Standard JSON with comment support | JSONC comments supported |
| YAML | `.yaml`, `.yml` | YAML 1.2 compliant | Full spec support |
| TOML | `.toml` | TOML v1.0.0 | Complete specification |
| ENV | `.env`, `.envrc` | Key=value pairs | Shell-style variables |

### Output Formats

All input formats are supported as output formats with automatic conversion between them.

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

### Large Configuration Files
- Use streaming parsers for files >100MB
- Consider splitting large schemas into modules
- Batch processing for multiple environments

### Memory Usage
- Konfigo loads all source files into memory
- Schema processing requires additional memory for transformations
- Monitor usage with large variable files

## Integration Patterns

### CI/CD Integration
```bash
# Typical pipeline usage
konfigo -s base.yaml,env/${ENVIRONMENT}.yaml -S schema.yaml -of config.json
```

### Container Integration
```dockerfile
# Multi-stage build pattern
FROM konfigo:latest as config-builder
COPY configs/ /configs/
RUN konfigo -s /configs/base.yaml,/configs/prod.yaml -of /tmp/final.json

FROM alpine:latest
COPY --from=config-builder /tmp/final.json /app/config.json
```

### Library Integration
Konfigo is designed as a CLI tool but can be integrated into build processes, deployment scripts, and configuration management workflows.

## Version Compatibility

### Schema Version Support
- `apiVersion: v1` - Current stable version
- `apiVersion: konfigo/v1alpha1` - Legacy format (deprecated)

### Breaking Changes
Major version updates may include breaking changes to:
- Schema format and processing
- Command-line flag syntax
- Output format structure

See release notes for migration guides and compatibility information.
