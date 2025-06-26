# Error Messages Reference

Complete catalog of Konfigo error messages with explanations, causes, and solutions.

## File and Path Errors

### `ERROR: Source file not found: [filename]`

**Cause**: The specified source file doesn't exist or isn't accessible.

**Solutions**:
```bash
# Check file exists
ls -la config.yaml

# Use absolute path
konfigo -s /full/path/to/config.yaml

# Check current directory
pwd
```

### `ERROR: Permission denied reading file: [filename]`

**Cause**: Insufficient permissions to read the file.

**Solutions**:
```bash
# Check file permissions
ls -la config.yaml

# Fix permissions
chmod 644 config.yaml

# Run with appropriate user
sudo konfigo -s config.yaml
```

### `ERROR: Directory not found: [directory]`

**Cause**: Specified directory for recursive search doesn't exist.

**Solutions**:
```bash
# Check directory exists
ls -la configs/

# Create directory if needed
mkdir -p configs/

# Use correct path
konfigo -r -s ./configurations/
```

## Parsing Errors

### `ERROR: Invalid JSON syntax in file: [filename]`

**Cause**: Malformed JSON content.

**Common Issues**:
- Missing commas between elements
- Trailing commas
- Unquoted keys
- Unclosed brackets/braces

**Solutions**:
```bash
# Validate JSON syntax
jq . config.json

# Check for common issues
cat config.json | python -m json.tool
```

**Example Fix**:
```json
// ❌ Invalid
{
  "app": {
    "name": "myapp",  // No trailing comma before }
    "port": 8080,     // ❌ Trailing comma
  }
}

// ✅ Valid
{
  "app": {
    "name": "myapp",
    "port": 8080
  }
}
```

### `ERROR: Invalid YAML syntax in file: [filename]`

**Cause**: Malformed YAML content.

**Common Issues**:
- Inconsistent indentation
- Tab characters instead of spaces
- Missing colons
- Incorrect list formatting

**Solutions**:
```bash
# Validate YAML syntax
yamllint config.yaml

# Or use Python
python -c "import yaml; yaml.safe_load(open('config.yaml'))"
```

**Example Fix**:
```yaml
# ❌ Invalid (mixed indentation)
app:
  name: "myapp"
    port: 8080  # ❌ Wrong indentation

# ✅ Valid
app:
  name: "myapp"
  port: 8080
```

### `ERROR: Invalid TOML syntax in file: [filename]`

**Cause**: Malformed TOML content.

**Common Issues**:
- Unquoted strings with special characters
- Invalid datetime format
- Duplicate keys
- Incorrect array formatting

**Solutions**:
```bash
# Validate TOML syntax (if you have toml tool installed)
toml-test config.toml
```

### `ERROR: Invalid ENV format in file: [filename]`

**Cause**: Malformed environment variable file.

**Common Issues**:
- Spaces around equals sign
- Unquoted values with special characters
- Comments in wrong format

**Example Fix**:
```bash
# ❌ Invalid
APP_NAME = "myapp"    # ❌ Spaces around =
PORT=8080 # comment   # ❌ Space before comment

# ✅ Valid
APP_NAME="myapp"
PORT=8080             # Comment on separate line or no space
```

## Schema Errors

### `ERROR: Schema file is missing the required top-level field 'apiVersion'`

**Cause**: Schema file doesn't have the required `apiVersion` field.

**Solution**:
```yaml
# Add apiVersion to your schema
apiVersion: v1

vars:
  - name: "ENVIRONMENT"
    value: "production"
```

### `ERROR: Invalid schema version: [version]`

**Cause**: Unsupported or incorrect `apiVersion` value.

**Supported Versions**:
- `v1` (recommended)
- `konfigo/v1alpha1` (legacy, deprecated)

**Solution**:
```yaml
# Use supported version
apiVersion: v1
```

### `ERROR: Schema validation failed: [details]`

**Cause**: Schema structure doesn't match expected format.

**Common Issues**:
- Invalid field names
- Wrong data types
- Missing required fields

**Solution**: Check schema documentation and ensure proper structure.

## Variable and Substitution Errors

### `ERROR: Undefined variable: [variable_name]`

**Cause**: Schema references a variable that isn't defined.

**Solutions**:
```bash
# Define in schema
vars:
  - name: "UNDEFINED_VAR"
    value: "default_value"

# Or provide via environment
export KONFIGO_VAR_UNDEFINED_VAR="runtime_value"

# Or provide via variables file
konfigo -s config.yaml -S schema.yaml -V variables.yaml
```

### `ERROR: Circular variable dependency detected`

**Cause**: Variables reference each other in a loop.

**Example Problem**:
```yaml
vars:
  - name: "A"
    value: "${B}"
  - name: "B"
    value: "${A}"  # ❌ Circular reference
```

**Solution**: Break the circular dependency by providing explicit values.

### `ERROR: Maximum variable substitution depth exceeded`

**Cause**: Too many nested variable substitutions (limit: 1000).

**Solution**: Simplify variable dependencies or check for unintended recursion.

## Validation Errors

### `Validation error: The field '[path]' is expected to be of type '[expected]', but got '[actual]'`

**Cause**: Data type doesn't match schema validation rules.

**Example**:
```yaml
# Schema expects number
validate:
  - path: "app.port"
    rules:
      type: "number"

# But config has string
app:
  port: "8080"  # ❌ String instead of number
```

**Solution**:
```yaml
# Fix config data type
app:
  port: 8080  # ✅ Number
```

### `Validation error: Required field '[path]' is missing`

**Cause**: Schema requires a field that isn't present in the configuration.

**Solution**:
```yaml
# Ensure required field is present
app:
  name: "myapp"  # ✅ Required field added
  port: 8080
```

### `Validation error: Value '[value]' is not in allowed enum: [options]`

**Cause**: Value doesn't match enumerated options.

**Example**:
```yaml
# Schema allows specific values
validate:
  - path: "environment"
    rules:
      enum: ["development", "staging", "production"]

# But config has invalid value
environment: "dev"  # ❌ Not in enum
```

**Solution**:
```yaml
# Use valid enum value
environment: "development"  # ✅ Valid
```

## Merge and Processing Errors

### `ERROR: Cannot merge incompatible types at path: [path]`

**Cause**: Attempting to merge different data types at the same path.

**Example**:
```json
// File 1
{"config": {"debug": true}}

// File 2  
{"config": {"debug": {"level": "info"}}}
// ❌ Can't merge boolean with object
```

**Solution**: Ensure compatible data types across all source files.

### `ERROR: Immutable path '[path]' cannot be overridden`

**Cause**: Attempting to override a path marked as immutable in schema.

**Note**: Environment variables (`KONFIGO_KEY_*`) can still override immutable paths.

**Solution**:
```bash
# Remove from immutable list in schema, or
# Use environment variable override
export KONFIGO_KEY_app.name="new-name"
```

## Environment Variable Errors

### `ERROR: Invalid environment variable format: [variable]`

**Cause**: Malformed `KONFIGO_KEY_*` or `KONFIGO_VAR_*` variable.

**Common Issues**:
- Missing underscore after prefix
- Invalid key path format
- Special characters in key names

**Examples**:
```bash
# ❌ Invalid
export KONFIGO_KEYapp.port=8080      # Missing underscore
export KONFIGO_KEY_app..port=8080    # Double dot
export KONFIGO_KEY_app port=8080     # Space in key

# ✅ Valid
export KONFIGO_KEY_app.port=8080
export KONFIGO_KEY_nested.config.value=test
```

## Output Errors

### `ERROR: Cannot write to output file: [filename]`

**Cause**: Insufficient permissions or invalid path for output file.

**Solutions**:
```bash
# Check directory permissions
ls -la $(dirname output.json)

# Create directory if needed
mkdir -p output/

# Check disk space
df -h
```

### `ERROR: Invalid output format for file extension: [extension]`

**Cause**: Unsupported file extension for output.

**Supported Extensions**: `.json`, `.yaml`, `.yml`, `.toml`, `.env`

**Solution**: Use supported extension or explicit format flag.

## Common Error Patterns and Quick Fixes

### "It worked before but now fails"

**Check**:
1. File permissions changed?
2. File moved or renamed?
3. Schema updated with new requirements?
4. Environment variables changed?

### "Merge results are unexpected"

**Debug with**:
```bash
# See merge order and precedence
konfigo -v -s file1.yaml,file2.yaml

# Check environment variable overrides
env | grep KONFIGO_KEY_
```

### "Schema processing fails"

**Common Fixes**:
1. Add `apiVersion: v1` to schema
2. Check variable definitions are complete
3. Verify validation rules match data types
4. Ensure all referenced paths exist

### "Output is empty or missing data"

**Check**:
1. Source files are being found and read
2. Merge precedence isn't causing unexpected overrides
3. Output schema isn't filtering out required data
4. Validation isn't rejecting the configuration

## Getting Help

When reporting errors, include:

1. **Full error message**
2. **Command used**
3. **Relevant configuration files** (sanitized)
4. **Environment variables** (non-sensitive)
5. **Konfigo version**: `konfigo --help` (shows usage)
6. **Operating system and version**

Use debug mode for detailed information:
```bash
konfigo -d -s your-config.yaml
```

This comprehensive error reference helps diagnose and resolve any issues you might encounter with Konfigo.
