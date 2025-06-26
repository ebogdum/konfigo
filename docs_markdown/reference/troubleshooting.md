# Troubleshooting Guide

Systematic approach to diagnosing and solving common Konfigo issues. Start with the symptom you're experiencing.

## Quick Diagnostics

### 1. **Basic Health Check**
```bash
# Verify Konfigo is working
konfigo --version

# Test with minimal example
echo '{"test": true}' | konfigo -s -
```

### 2. **Enable Debug Output**
```bash
# See what Konfigo is doing
konfigo -v -s your-config.yaml

# Maximum detail
konfigo -d -s your-config.yaml
```

### 3. **Validate Without Processing**
```bash
# Check syntax only
konfigo --syntax-only -s config.yaml

# Test schema without applying
konfigo -s config.yaml -S schema.yaml --validate-only
```

---

## Common Issues by Symptom

### 🔥 **"Command not found: konfigo"**

**Cause**: Konfigo not installed or not in PATH

**Solutions**:
```bash
# Check if Konfigo is installed
which konfigo

# Add to PATH (if installed elsewhere)
export PATH=$PATH:/path/to/konfigo

# Reinstall if missing
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64 -o konfigo
chmod +x konfigo
sudo mv konfigo /usr/local/bin/
```

### 🔥 **"File not found" or "No such file"**

**Cause**: Source files don't exist or wrong paths

**Diagnosis**:
```bash
# Check if files exist
ls -la config.yaml
ls -la /full/path/to/config.yaml

# Verify current directory
pwd
```

**Solutions**:
```bash
# Use absolute paths
konfigo -s /full/path/to/config.yaml

# Check relative paths from current directory
konfigo -s ./configs/base.yaml

# List available files
ls configs/
```

### 🔥 **"YAML parsing failed" or "Invalid format"**

**Cause**: Syntax errors in configuration files

**Diagnosis**:
```bash
# Check YAML syntax
yamllint config.yaml

# Or use Python
python -c "import yaml; yaml.safe_load(open('config.yaml'))"

# Check JSON syntax
jq . config.json
```

**Solutions**:
```bash
# Fix YAML indentation (most common issue)
# Ensure consistent spacing (2 or 4 spaces, not tabs)

# Check for common YAML errors:
# - Inconsistent indentation
# - Tabs instead of spaces  
# - Missing quotes around special characters
# - Unescaped colons in values

# Use explicit format if detection fails
konfigo -sy -s problematic-file.conf  # Force YAML
```

### 🔥 **"Schema validation failed"**

**Cause**: Configuration doesn't meet schema requirements

**Diagnosis**:
```bash
# See detailed validation errors
konfigo -v -s config.yaml -S schema.yaml

# Check specific validation rules
konfigo -d -s config.yaml -S schema.yaml --validate-only
```

**Solutions**:
```bash
# Check required fields
# Ensure all paths in schema exist in config

# Verify data types
# Numbers should be numbers, not strings
# Booleans should be true/false, not "true"/"false"

# Check value constraints
# min/max ranges, enum values, regex patterns

# Example fixes:
# Wrong: port: "8080"
# Right: port: 8080

# Wrong: enabled: "true"  
# Right: enabled: true
```

### 🔥 **"Variable not found" or "Undefined variable"**

**Cause**: Schema references variables that aren't defined

**Diagnosis**:
```bash
# Check variable definitions
konfigo -v -s config.yaml -S schema.yaml -V vars.yaml

# List environment variables
env | grep KONFIGO_VAR
```

**Solutions**:
```bash
# Define missing variables in schema
vars:
  - name: "MISSING_VAR"
    value: "default_value"

# Or set environment variable
export KONFIGO_VAR_MISSING_VAR="value"

# Or provide variables file
konfigo -s config.yaml -S schema.yaml -V variables.yaml
```

### 🔥 **"Permission denied"**

**Cause**: File permission issues

**Diagnosis**:
```bash
# Check file permissions
ls -la config.yaml
ls -la output-directory/

# Check if output directory exists
ls -ld output-directory/
```

**Solutions**:
```bash
# Fix file permissions
chmod 644 config.yaml

# Fix directory permissions  
chmod 755 output-directory/

# Create output directory
mkdir -p output-directory/

# Use different output location
konfigo -s config.yaml -of ~/temp/output.json
```

### 🔥 **"Memory limit exceeded" or "Out of memory"**

**Cause**: Large configuration files or complex processing

**Solutions**:
```bash
# Increase memory limit
konfigo --memory-limit 2GB -s large-config.yaml

# Use streaming for large files
konfigo --stream -s huge-file.json

# Process in chunks
split -l 1000 large-config.yaml chunk_
for chunk in chunk_*; do
  konfigo -s "$chunk" -of "output_$chunk.yaml"
done
```

### 🔥 **"Schema processing timeout"**

**Cause**: Complex schema or infinite loops

**Solutions**:
```bash
# Increase timeout
konfigo --timeout 60s -s config.yaml -S complex-schema.yaml

# Simplify schema temporarily
# Remove complex transformations to isolate issue

# Check for circular references
# Ensure variables don't reference themselves
```

---

## Advanced Debugging

### Environment Variable Issues

**Problem**: Environment overrides not working

**Diagnosis**:
```bash
# Check environment variables
env | grep KONFIGO_KEY
env | grep KONFIGO_VAR

# Test with simple override
KONFIGO_KEY_test=value konfigo -s config.yaml
```

**Common fixes**:
```bash
# Correct format for nested keys
export KONFIGO_KEY_database.host=localhost  # ✅
export KONFIGO_KEY_database_host=localhost  # ❌

# Escape special characters
export KONFIGO_KEY_app.url="https://example.com"  # ✅

# Check variable precedence
# Environment variables have highest precedence
```

### Merging Issues

**Problem**: Files not merging as expected

**Diagnosis**:
```bash
# Check merge order
konfigo -v -s file1.yaml,file2.yaml

# Test individual files
konfigo -s file1.yaml
konfigo -s file2.yaml
```

**Common fixes**:
```bash
# Check merge order (later files override earlier)
konfigo -s base.yaml,override.yaml  # override wins

# Verify file contents
# Ensure files contain expected data

# Check for case sensitivity issues
konfigo -c -s config1.yaml,config2.yaml  # case-insensitive
```

### Performance Issues

**Problem**: Konfigo running slowly

**Diagnosis**:
```bash
# Time the operation
time konfigo -s config.yaml -S schema.yaml

# Profile with verbose output
konfigo -v -s config.yaml -S schema.yaml

# Check file sizes
ls -lh config.yaml schema.yaml
```

**Solutions**:
```bash
# Enable parallel processing
konfigo --parallel 4 -s configs/*

# Optimize schema
# Remove unnecessary transformations
# Simplify complex validations

# Use streaming for large files
konfigo --stream -s large-config.json
```

---

## Error Message Reference

### Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| `0` | Success | None needed |
| `1` | General error | Check command syntax |
| `2` | Invalid arguments | Review command line flags |
| `3` | File not found | Verify file paths |
| `4` | Parse error | Check file syntax |
| `5` | Schema validation failed | Review validation rules |
| `6` | Write error | Check output permissions |

### Common Error Patterns

```bash
# "cannot read file"
# → Check file exists and permissions

# "invalid YAML"  
# → Run yamllint or check indentation

# "validation failed: required field missing"
# → Add missing field or make it optional

# "variable 'X' not found"
# → Define variable in schema or environment

# "output directory does not exist"
# → Create directory or use different path
```

---

## Getting Help

### Self-Diagnosis Steps

1. **Check the basics**: File exists, permissions, syntax
2. **Use verbose mode**: `konfigo -v` shows processing steps
3. **Isolate the problem**: Test with minimal config
4. **Check environment**: Variables, PATH, permissions
5. **Review documentation**: Ensure correct usage

### Debug Information to Collect

When reporting issues, include:

```bash
# Konfigo version
konfigo --version

# Command that failed
konfigo -v -s config.yaml -S schema.yaml

# File contents (sanitized)
cat config.yaml
cat schema.yaml

# Environment variables (sanitized)
env | grep KONFIGO

# Operating system
uname -a
```

### Community Resources

- **GitHub Issues**: [Report bugs and feature requests](https://github.com/ebogdum/konfigo/issues)
- **Documentation**: [Browse all guides](../index.md)
- **Examples**: [Real-world patterns](../guide/recipes.md)

### Before Reporting a Bug

1. **Update to latest version**: `konfigo --version`
2. **Check existing issues**: Search GitHub issues
3. **Create minimal reproduction**: Simplest possible example
4. **Include debug output**: Use `-d` flag for details

---

## Prevention Tips

### **Always validate configurations**
```bash
# Test before deployment
konfigo -s config.yaml -S schema.yaml --validate-only
```

### **Use version control**
```bash
# Track configuration changes
git add configs/ schemas/
git commit -m "Update production config"
```

### **Test in development**
```bash
# Test merging logic before production
konfigo -s base.yaml,dev.yaml -S schema.yaml
```

### **Document your patterns**
```bash
# Create README for your config setup
echo "# Configuration Guide" > configs/README.md
echo "Run: konfigo -s base.yaml,\$ENV.yaml" >> configs/README.md
```

### **Use schemas for validation**
```yaml
# Always validate critical fields
validation:
  - path: "database.host"
    required: true
  - path: "app.port"
    type: "number"
    min: 1024
```

Most Konfigo issues are configuration or syntax related. Use verbose mode, check your file syntax, and validate your schemas to prevent common problems!
