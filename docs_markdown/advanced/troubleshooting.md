# Troubleshooting Guide

This guide helps you diagnose and resolve common issues when working with Konfigo.

## Common Issues and Solutions

### File and Path Issues

#### File Not Found
```
Error: failed to read file 'config.yaml'
```

**Solutions**:
- **Check file path**: Verify the file exists at the specified location
- **Use absolute paths**: Try absolute paths if relative paths fail
- **Check permissions**: Ensure the file is readable (`ls -la config.yaml`)
- **Verify current directory**: Confirm you're in the expected working directory

```bash
# Debug file access
ls -la config.yaml
pwd
konfigo -s "$(pwd)/config.yaml"  # Use absolute path
```

#### Permission Denied
```
Error: permission denied reading 'secret-config.yaml'
```

**Solutions**:
- **Fix file permissions**: `chmod 644 secret-config.yaml`
- **Check directory permissions**: Ensure parent directories are accessible
- **Run with appropriate user**: Switch to a user with proper access

### Parsing and Format Issues

#### YAML Parsing Error
```
Error: failed to parse file 'config.yaml'
  Caused by: yaml: line 5: found character that cannot start any token
```

**Solutions**:
- **Validate YAML syntax**: Use `yamllint config.yaml`
- **Check indentation**: YAML is whitespace-sensitive
- **Look for special characters**: Check for tabs, invisible characters
- **Verify encoding**: Ensure file is UTF-8 encoded

```bash
# Debug YAML issues
yamllint config.yaml
python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"
```

#### JSON Parsing Error
```
Error: failed to parse file 'config.json'
  Caused by: invalid character '}' after object key
```

**Solutions**:
- **Validate JSON syntax**: Use `jq . config.json`
- **Check for trailing commas**: JSON doesn't allow trailing commas
- **Verify quotes**: Ensure all strings use double quotes

```bash
# Debug JSON issues
jq . config.json
python3 -c "import json; json.load(open('config.json'))"
```

#### Format Detection Issues
```
Error: unable to determine file format for 'data.txt'
```

**Solutions**:
- **Use format flags**: Specify format explicitly (`-sj`, `-sy`, `-st`)
- **Check file extension**: Use standard extensions (`.json`, `.yaml`, `.toml`)
- **Verify content structure**: Ensure content matches expected format

```bash
# Force format detection
konfigo -s data.txt -sy  # Force YAML parsing
konfigo -s data.txt -sj  # Force JSON parsing
```

### Variable and Schema Issues

#### Variable Not Found
```
Error: variable 'DATABASE_URL' not found
```

**Solutions**:
- **Check variable definition**: Ensure variable is defined in schema or vars file
- **Verify environment variables**: Check if env vars are set (`env | grep KONFIGO_`)
- **Check spelling**: Verify variable names match exactly
- **Review precedence**: Understand variable resolution order

```bash
# Debug variable resolution
env | grep KONFIGO_VAR_
env | grep DATABASE_URL
konfigo -d -s config.yaml -S schema.yaml  # Debug mode shows variable resolution
```

#### Schema Validation Failed
```
Error: validation failed
  Caused by: path 'service.port' value 99999 exceeds maximum value 65535
```

**Solutions**:
- **Review validation rules**: Check schema validation constraints
- **Verify input data**: Ensure values match expected types and ranges
- **Check data source**: Verify which source provides the problematic value
- **Use debug mode**: See detailed validation process

```bash
# Debug validation
konfigo -d -s config.yaml -S schema.yaml
# Review the specific validation rule causing the issue
```

#### Generator/Transformer Failures
```
Error: schema processing failed
  Caused by: generator failed: source path 'service.hostname' not found
```

**Solutions**:
- **Verify source paths**: Ensure all referenced paths exist in merged configuration
- **Check processing order**: Generators run before transformations
- **Use debug output**: See intermediate processing results
- **Test without schema**: Verify base configuration first

```bash
# Debug schema processing step by step
konfigo -s config.yaml                    # Test merging only
konfigo -s config.yaml -S minimal-schema.yaml  # Test with simple schema
```

### Environment Variable Issues

#### Environment Override Not Working
```
# KONFIGO_KEY_service.port=9090 not overriding configuration
```

**Solutions**:
- **Check variable format**: Use correct `KONFIGO_KEY_` prefix
- **Verify path syntax**: Use dot notation for nested paths
- **Check case sensitivity**: Use `-c` flag if needed
- **Export variables**: Ensure variables are exported

```bash
# Debug environment variables
export KONFIGO_KEY_service.port=9090
env | grep KONFIGO_KEY_
konfigo -v -s config.yaml  # Verbose mode shows environment processing
```

#### Variable Substitution Not Working
```
# ${ENVIRONMENT} not being replaced in output
```

**Solutions**:
- **Check variable definition**: Ensure variable is properly defined
- **Verify schema processing**: Substitution happens during schema processing
- **Use debug mode**: See variable resolution process
- **Check syntax**: Ensure `${VAR_NAME}` format is correct

### Performance and Memory Issues

#### Slow Processing
**Solutions**:
- **Profile with debug mode**: Identify bottlenecks
- **Reduce complexity**: Simplify schemas or split large files
- **Use parallel processing**: Default behavior, but verify
- **Check file sizes**: Consider streaming for very large files

```bash
# Profile processing time
time konfigo -d -s large-config.yaml -S complex-schema.yaml
```

#### Memory Issues
```
Error: out of memory processing large configuration
```

**Solutions**:
- **Use streaming**: Process large files via stdin
- **Split configurations**: Break large files into smaller pieces
- **Increase memory limits**: In containerized environments
- **Optimize schemas**: Reduce complex transformations

```bash
# Use streaming for large files
cat large-config.yaml | konfigo -sy -S schema.yaml
```

## Debug Techniques

### Step-by-Step Debugging

1. **Test basic merging**:
```bash
konfigo -s base.yaml,env.yaml
```

2. **Add schema processing**:
```bash
konfigo -s base.yaml,env.yaml -S schema.yaml
```

3. **Add variables**:
```bash
konfigo -s base.yaml,env.yaml -S schema.yaml -V vars.yaml
```

### Debug Mode Usage

Enable detailed logging to understand processing flow:

```bash
# Enable debug logging
konfigo -d -s config.yaml -S schema.yaml

# Enable info logging (less verbose)
konfigo -v -s config.yaml -S schema.yaml
```

### Component Validation

Validate individual components:

```bash
# Test YAML syntax
yamllint config.yaml

# Test JSON syntax
jq . config.json

# Test TOML syntax
toml-test config.toml

# Check environment variables
env | grep KONFIGO_
```

### Output Intermediate Results

Save intermediate processing results for analysis:

```bash
# Save merged configuration
konfigo -s base.yaml,env.yaml -oy > merged.yaml

# Save after schema processing
konfigo -s merged.yaml -S schema.yaml -oy > processed.yaml

# Compare differences
diff -u merged.yaml processed.yaml
```

## Error Message Reference

### Common Error Patterns

| Error Pattern | Likely Cause | Solution |
|---------------|--------------|----------|
| `failed to read file` | File access issue | Check path, permissions |
| `failed to parse` | Syntax error | Validate file format |
| `variable not found` | Missing variable | Check definition, environment |
| `validation failed` | Schema constraint | Review validation rules |
| `path not found` | Missing configuration | Check merge results |
| `type mismatch` | Incompatible types | Verify data types |

### Exit Codes

- **0**: Success
- **1**: General error (file not found, parsing failed, etc.)
- **2**: Validation error
- **3**: Schema processing error
- **4**: Environment/variable error

## Getting Help

### Built-in Help
```bash
konfigo -h  # Show help message
```

### Verbose Output
```bash
konfigo -v  # Show processing information
konfigo -d  # Show detailed debug information
```

### Community Resources
- **GitHub Issues**: Report bugs and get help
- **Documentation**: Comprehensive guides and examples
- **Discussions**: Community Q&A and best practices

### Creating Minimal Reproduction Cases

When seeking help, create minimal examples:

```yaml
# minimal-config.yaml
service:
  name: test-app
  port: 8080

# minimal-schema.yaml  
validate:
  - path: "service.port"
    rules:
      type: "number"
      max: 65535
```

```bash
# Minimal command that reproduces the issue
konfigo -s minimal-config.yaml -S minimal-schema.yaml
```

Include:
- Exact command used
- Full error message
- Input files (minimal versions)
- Expected vs actual behavior
- Environment details (OS, Konfigo version)
