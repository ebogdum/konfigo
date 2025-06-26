# Format Conversion

Konfigo supports reading and writing configuration files in multiple formats: JSON, YAML, TOML, ENV, and INI. The tool can automatically detect input formats and convert between any supported formats.

## Supported Formats

| Format | Extensions | Input | Output | Auto-detect |
|--------|------------|-------|--------|-------------|
| JSON   | `.json`    | ✅    | ✅     | ✅          |
| YAML   | `.yml`, `.yaml` | ✅ | ✅     | ✅          |
| TOML   | `.toml`    | ✅    | ✅     | ✅          |
| ENV    | `.env`     | ✅    | ✅     | ❌*         |
| INI    | `.ini`     | ✅    | ✅     | ✅          |

*ENV format requires explicit format flag when reading from stdin

## Input Format Control

### Automatic Detection
Konfigo automatically detects input format based on file extensions:

```bash
# Auto-detected as JSON
konfigo -s config.json

# Auto-detected as YAML  
konfigo -s config.yml

# Auto-detected as TOML
konfigo -s config.toml
```

### Explicit Format Override
Force specific input format using format flags:

```bash
# Force JSON parsing
konfigo -s config.txt -sj

# Force YAML parsing
konfigo -s config.txt -sy

# Force TOML parsing  
konfigo -s config.txt -st

# Force ENV parsing
konfigo -s config.txt -se
```

### Stdin Format Requirements
When reading from stdin, format must be explicitly specified:

```bash
# JSON from stdin
cat config.json | konfigo -sj

# YAML from stdin
cat config.yml | konfigo -sy

# TOML from stdin
cat config.toml | konfigo -st

# ENV from stdin
cat config.env | konfigo -se
```

## Output Format Control

### Single Format Output
Specify output format using format flags:

```bash
# Output as JSON
konfigo -s config.yml -oj

# Output as YAML (default)
konfigo -s config.json -oy

# Output as TOML
konfigo -s config.yml -ot

# Output as ENV
konfigo -s config.json -oe
```

### File Output with Format Detection
Use `-of` to write to file with automatic format detection:

```bash
# Format determined by extension
konfigo -s config.yml -of output.json  # JSON output
konfigo -s config.json -of output.toml # TOML output
konfigo -s config.toml -of output.yml  # YAML output
```

### Multiple Output Formats
Generate multiple output formats simultaneously:

```bash
# Output to stdout in both JSON and YAML
konfigo -s config.toml -oj -oy

# File output with additional stdout format
konfigo -s config.yml -of output.json -oy
```

## Format Conversion Examples

Based on test cases from `test/format-conversion/`:

### JSON to YAML
```bash
# Input: config.json
{
  "service": {
    "name": "my-app",
    "port": 8080
  },
  "database": {
    "host": "localhost",
    "port": 5432
  }
}

# Convert to YAML
konfigo -s config.json -oy

# Output:
service:
  name: my-app
  port: 8080
database:
  host: localhost
  port: 5432
```

### YAML to TOML
```bash
# Input: config.yml
service:
  name: my-app
  port: 8080
features:
  - auth
  - logging

# Convert to TOML
konfigo -s config.yml -ot

# Output:
[service]
name = "my-app"
port = 8080
features = ["auth", "logging"]
```

### Multiple Formats to ENV
```bash
# Merge JSON and YAML, output as ENV
konfigo -s base.json,override.yml -oe

# Output:
SERVICE_NAME=my-app
SERVICE_PORT=9090
DATABASE_HOST=prod-db.example.com
DATABASE_SSL=true
```

## Format-Specific Features

### JSON
- Preserves exact numeric types
- Supports nested objects and arrays
- UTF-8 encoding support

### YAML
- Human-readable format
- Supports comments
- Multi-document files (first document used)
- Type inference for scalars

### TOML
- Strongly typed
- Clear section hierarchy
- Good for configuration files
- Date/time support

### ENV
- Key-value pairs only
- Nested objects flattened with underscores
- All values treated as strings
- Uppercase key normalization

### INI
- Section-based structure
- Simple key-value pairs
- Comments supported
- Case-insensitive keys

## Error Handling

### Invalid Format
```bash
# Malformed JSON
konfigo -s invalid.json
# Error: failed to parse file: invalid character '}' looking for beginning of object key string

# Mismatched format override
echo "invalid: json" | konfigo -sj
# Error: failed to parse stdin: invalid character ':' after top-level value
```

### Unsupported Conversions
All supported formats can convert to any other supported format. Konfigo handles type coercion automatically:

- Numbers preserved where possible
- Booleans converted to strings in ENV format
- Complex structures flattened for ENV output
- Arrays serialized as comma-separated values in ENV

## Best Practices

1. **Auto-detection**: Let Konfigo detect formats when possible for cleaner commands
2. **Explicit stdin**: Always specify format when reading from stdin
3. **File extensions**: Use standard extensions for automatic format detection
4. **ENV limitations**: Be aware that ENV format has structural limitations
5. **Type preservation**: Use JSON/YAML/TOML for preserving complex data types

## Test Coverage

Format conversion is tested comprehensively in `test/format-conversion/`:
- Cross-format conversion matrix
- Type preservation verification
- Error condition handling
- Edge cases and special characters
