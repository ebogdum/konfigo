# Input and Output Formats

Konfigo provides comprehensive support for multiple configuration formats, enabling seamless conversion and processing across different file types.

## Supported Formats

| Format | Extensions | Input | Output | Features |
|--------|------------|-------|--------|----------|
| **JSON** | `.json` | ✅ | ✅ | Precise typing, compact, widely supported |
| **YAML** | `.yaml`, `.yml` | ✅ | ✅ | Human-readable, comments, multi-document |
| **TOML** | `.toml` | ✅ | ✅ | Configuration-focused, strongly typed |
| **ENV** | `.env` | ✅ | ✅ | Environment variables, simple key-value |
| **INI** | `.ini` | ✅ | ❌ | Legacy support, sections |

## Format Detection

Konfigo automatically detects input formats using multiple methods:

### 1. File Extension (Primary)
```bash
konfigo -s config.json    # Detected as JSON
konfigo -s app.yaml       # Detected as YAML
konfigo -s settings.toml  # Detected as TOML
```

### 2. Content Analysis (Fallback)
When file extensions are missing or ambiguous, Konfigo analyzes content structure:

```bash
konfigo -s config          # Auto-detects based on content
konfigo -s data.txt        # Analyzes content structure
```

### 3. Format Override (Explicit)
Force specific format parsing:

```bash
konfigo -s config.txt -sj  # Force JSON parsing
konfigo -s data -sy        # Force YAML parsing
konfigo -s settings -st    # Force TOML parsing
```

## Format-Specific Features

### JSON
- **Strengths**: Precise typing, compact, universal support
- **Use Cases**: APIs, data exchange, compact storage
- **Example**:
```json
{
  "app": {
    "name": "my-service",
    "port": 8080,
    "features": ["auth", "cache"]
  }
}
```

### YAML
- **Strengths**: Human-readable, supports comments, multi-document
- **Use Cases**: Configuration files, documentation, complex structures
- **Example**:
```yaml
# Application configuration
app:
  name: my-service
  port: 8080
  features:
    - auth
    - cache
  # Development settings
  debug: true
```

### TOML
- **Strengths**: Configuration-focused, strongly typed, readable
- **Use Cases**: Application configuration, settings files
- **Example**:
```toml
[app]
name = "my-service"
port = 8080
features = ["auth", "cache"]

[database]
host = "localhost"
port = 5432
```

### ENV
- **Strengths**: Simple key-value, environment integration
- **Use Cases**: Environment variables, Docker configs, CI/CD
- **Example**:
```env
APP_NAME=my-service
APP_PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
```

### INI (Input Only)
- **Strengths**: Legacy support, simple sections
- **Use Cases**: Legacy applications, simple configurations
- **Example**:
```ini
[app]
name = my-service
port = 8080

[database]
host = localhost
port = 5432
```

## Multiple Output Formats

Generate output in multiple formats simultaneously:

### Stdout Multiple Formats
```bash
# Output JSON and YAML to stdout
konfigo -s config.yaml -oj -oy

# Output all supported formats
konfigo -s config.yaml -oj -oy -ot -oe
```

### Files and Stdout Combined
```bash
# Save to file and display on stdout
konfigo -s config.yaml -of final.json -oy

# Multiple file formats
konfigo -s config.yaml -of app.json -of settings.toml -oy
```

## Stdin Processing

Process configurations from standard input with format specification:

### Basic Stdin Usage
```bash
# From file redirection
konfigo -sj < config.json

# From pipe
cat config.yaml | konfigo -sy
```

### Advanced Stdin Examples
```bash
# From command output
kubectl get configmap app-config -o yaml | konfigo -sy -S schema.yaml

# From curl/API
curl -s https://api.example.com/config | konfigo -sj -of local.yaml

# From heredoc
konfigo -sy -S schema.yaml <<EOF
app:
  name: test-app
  port: 8080
  features:
    auth: true
EOF
```

### Stdin with Schema Processing
```bash
# Process stdin with schema and variables
echo '{"env": "dev"}' | konfigo -sj -S schema.yaml -V vars.yaml
```

## Format Conversion Examples

### JSON to YAML
```bash
konfigo -s config.json -oy
```

### YAML to TOML
```bash
konfigo -s app.yaml -ot
```

### Multiple Sources to Single Format
```bash
konfigo -s base.yaml,prod.json,local.toml -oj
```

### ENV to Structured Format
```bash
konfigo -s .env -oy
# Converts flat key-value to nested YAML structure
```

## Format-Specific Considerations

### Data Type Preservation
- **JSON**: Maintains precise typing (numbers, booleans, null)
- **YAML**: Supports rich typing with automatic inference
- **TOML**: Strong typing with explicit type specification
- **ENV**: String-based with optional type conversion

### Comments and Documentation
- **YAML**: Full comment support preserved in processing
- **TOML**: Comment support with preservation
- **JSON**: No native comment support
- **ENV**: Comment support with `#` prefix

### Complex Data Structures
- **JSON/YAML**: Full support for nested objects and arrays
- **TOML**: Good support with table syntax
- **ENV**: Flat key-value with path notation support
- **INI**: Basic sections, limited nesting

## Best Practices

1. **Choose Format by Use Case**:
   - YAML for human-edited configuration
   - JSON for APIs and data exchange
   - TOML for application settings
   - ENV for environment-specific values

2. **Consistent Format Strategy**:
   - Use consistent formats within projects
   - Document format choices for team members

3. **Leverage Auto-Detection**:
   - Use clear file extensions
   - Let Konfigo detect formats automatically

4. **Test Format Conversions**:
   - Verify data integrity across format conversions
   - Check for type preservation requirements

5. **Handle Stdin Gracefully**:
   - Always specify format for stdin input
   - Validate piped input before processing
