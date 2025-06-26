# Frequently Asked Questions

Quick answers to the most common questions about Konfigo.

## Getting Started

### **Q: What exactly does Konfigo do?**
**A:** Konfigo merges, validates, and transforms configuration files. Think of it as a "configuration compiler" that takes multiple config sources (JSON, YAML, TOML, ENV) and produces a single, validated output in your preferred format.

### **Q: How is Konfigo different from other config tools?**
**A:** Konfigo uniquely combines:
- **Multi-format support** (JSON ↔ YAML ↔ TOML ↔ ENV)
- **Intelligent deep merging** with precedence rules
- **Schema-driven processing** (validation, transformation, generation)
- **Environment variable integration**
- **Batch output generation**

### **Q: Do I need to learn schemas to use Konfigo?**
**A:** No! You can start with simple merging and format conversion:
```bash
konfigo -s base.yaml,prod.yaml -of final.json
```
Schemas unlock advanced features when you're ready.

### **Q: Can I use Konfigo without installing anything?**
**A:** Yes! Konfigo is a single binary with no dependencies. Download and run:
```bash
curl -L https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64 -o konfigo
chmod +x konfigo
./konfigo -s config.yaml
```

## File Formats and Compatibility

### **Q: What file formats does Konfigo support?**
**A:** Input and output:
- **JSON** (`.json`, `.jsonc`) - Comments supported in JSONC
- **YAML** (`.yaml`, `.yml`) - Full YAML 1.2 support
- **TOML** (`.toml`) - TOML v1.0.0
- **ENV** (`.env`, `.envrc`) - Key=value pairs

### **Q: Can I mix different formats in one command?**
**A:** Absolutely! That's one of Konfigo's strengths:
```bash
konfigo -s base.yaml,override.json,local.toml -of final.json
```

### **Q: How does Konfigo handle invalid file formats?**
**A:** Konfigo provides detailed error messages and suggestions:
```bash
# Shows line numbers and specific syntax errors
konfigo -s invalid.yaml
# Error: YAML parsing failed at line 5: invalid indentation
```

### **Q: Can I convert between formats without merging?**
**A:** Yes! Use a single source file:
```bash
konfigo -s config.yaml -oj -of config.json  # YAML → JSON
konfigo -s data.json -oy -of data.yaml      # JSON → YAML
```

## Configuration Merging

### **Q: How does Konfigo decide which values to keep when merging?**
**A:** Konfigo uses clear precedence rules:
1. **Environment variables** (`KONFIGO_KEY_*`) - Highest
2. **Later source files** (rightmost in command)
3. **Earlier source files** (leftmost in command)

Example: `konfigo -s base.yaml,prod.yaml`
- `prod.yaml` overrides `base.yaml`
- Environment variables override both

### **Q: What happens to arrays when merging?**
**A:** Arrays are replaced completely, not merged:
```yaml
# base.yaml
tags: ["app", "service"]

# override.yaml  
tags: ["app", "production"]

# Result: ["app", "production"] (completely replaced)
```

### **Q: Can I protect certain configuration values from being overridden?**
**A:** Yes! Use immutable paths in schemas:
```yaml
immutable:
  - "app.name"
  - "security.keys"
```
Note: Environment variables can still override immutable paths.

### **Q: How do I merge all files in a directory?**
**A:** Use the recursive flag:
```bash
konfigo -r -s configs/        # All files in configs/
konfigo -r -s configs/ -of merged.yaml
```

## Environment Variables

### **Q: How do I override configuration with environment variables?**
**A:** Use `KONFIGO_KEY_` prefix with dot notation:
```bash
# Override nested values
export KONFIGO_KEY_database.host=prod-db.com
export KONFIGO_KEY_app.port=9090
konfigo -s config.yaml
```

### **Q: What's the difference between KONFIGO_KEY_ and KONFIGO_VAR_?**
**A:**
- **`KONFIGO_KEY_*`**: Directly overrides configuration values
- **`KONFIGO_VAR_*`**: Provides variables for schema substitution

```bash
# Direct override
KONFIGO_KEY_app.port=8080 konfigo -s config.yaml

# Schema variable  
KONFIGO_VAR_PORT=8080 konfigo -s config.yaml -S schema.yaml
# (schema must use ${PORT} somewhere)
```

### **Q: Can I use environment variables without schemas?**
**A:** Yes! `KONFIGO_KEY_*` variables work without schemas:
```bash
KONFIGO_KEY_database.host=localhost konfigo -s config.yaml
```

## Schemas

### **Q: When should I use schemas?**
**A:** Use schemas when you need:
- **Validation** (ensure configs are correct)
- **Variables** (environment-specific values)
- **Transformation** (modify data during processing)
- **Generation** (create UUIDs, timestamps, etc.)
- **Batch processing** (multiple outputs from templates)

### **Q: Are schemas required?**
**A:** No! Schemas are optional. You can use Konfigo for basic merging and format conversion without any schemas.

### **Q: Can I validate configurations without applying transformations?**
**A:** Yes! Use the validate-only flag:
```bash
konfigo -s config.yaml -S schema.yaml --validate-only
```

### **Q: How do I create my first schema?**
**A:** Start simple with variables:
```yaml
# my-schema.yaml
vars:
  - name: "DATABASE_HOST"
    value: "localhost"

transforms:
  - path: "database.host"
    setValue: "${DATABASE_HOST}"
```

## Performance and Limits

### **Q: How large can my configuration files be?**
**A:** Konfigo handles large files well. For files >100MB, use streaming:
```bash
konfigo --stream -s large-config.json
```

### **Q: Can I process configurations in parallel?**
**A:** Yes! Use the parallel flag:
```bash
konfigo --parallel 4 -r -s configs/
```

### **Q: Does Konfigo have any dependencies?**
**A:** No! Konfigo is a single binary with zero dependencies. Just download and run.

## Security

### **Q: How should I handle secrets in configurations?**
**A:** Best practices:
1. **Never commit secrets** to version control
2. **Use environment variables** for sensitive data:
   ```bash
   KONFIGO_KEY_database.password=$SECRET_PASSWORD konfigo -s config.yaml
   ```
3. **Use external secret management** and inject at runtime
4. **Use immutable paths** to protect critical settings

### **Q: Can I exclude sensitive data from output?**
**A:** Yes! Use output schemas to filter sensitive fields:
```yaml
outputSchema:
  exclude:
    - "database.password"
    - "api.secret_key"
```

## Integration and Deployment

### **Q: How do I use Konfigo in Docker?**
**A:** Common pattern:
```dockerfile
FROM alpine
RUN wget https://github.com/ebogdum/konfigo/releases/latest/download/konfigo-linux-amd64 -O /usr/local/bin/konfigo && chmod +x /usr/local/bin/konfigo
COPY configs/ /configs/
CMD ["sh", "-c", "konfigo -s /configs/base.yaml -of /app/config.json && exec myapp"]
```

### **Q: How do I integrate with Kubernetes?**
**A:** Use ConfigMaps and environment variables:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  base.yaml: |
    app:
      name: "my-service"
---
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: KONFIGO_KEY_database.host
          value: "prod-db"
        command: ["konfigo", "-s", "/etc/config/base.yaml", "-of", "/tmp/config.json"]
```

### **Q: Can I use Konfigo in CI/CD pipelines?**
**A:** Absolutely! Common CI/CD pattern:
```bash
# Validate configurations
konfigo -s base.yaml,prod.yaml -S schema.yaml --validate-only

# Generate environment-specific configs
konfigo -s base.yaml,prod.yaml -S schema.yaml -of prod-config.json

# Deploy with generated config
kubectl apply -f deployment.yaml
```

## Error Handling

### **Q: What should I do if Konfigo fails with "file not found"?**
**A:** Check these common issues:
1. **File paths** - Use absolute paths or verify current directory
2. **File permissions** - Ensure files are readable
3. **File names** - Check for typos

```bash
# Debug with verbose output
konfigo -v -s config.yaml
```

### **Q: How do I debug schema validation errors?**
**A:** Use debug mode to see detailed validation information:
```bash
konfigo -d -s config.yaml -S schema.yaml --validate-only
```

### **Q: What exit codes does Konfigo use?**
**A:**
- `0` - Success
- `1` - General error
- `3` - File not found
- `4` - Parse error
- `5` - Schema validation failed

## Common Use Cases

### **Q: How do I manage configurations for multiple environments?**
**A:** Use the layer pattern:
```bash
# Development
konfigo -s base.yaml,environments/dev.yaml -of dev-config.json

# Production
konfigo -s base.yaml,environments/prod.yaml -of prod-config.json
```

### **Q: How do I generate configurations for microservices?**
**A:** Use batch processing with schemas:
```yaml
# services.yaml
konfigo_forEach:
  - name: "user-service"
    vars: {SERVICE_PORT: 8001, DB_NAME: "users"}
  - name: "order-service"  
    vars: {SERVICE_PORT: 8002, DB_NAME: "orders"}
```

```bash
konfigo -s template.yaml -S schema.yaml -V services.yaml
```

### **Q: How do I convert legacy .env files to modern YAML?**
**A:** Simple format conversion:
```bash
konfigo -s legacy.env -oy -of modern.yaml
```

## Getting Help

### **Q: Where can I find more examples?**
**A:** Check out:
- **[Recipes & Examples](../guide/recipes.md)** - Real-world patterns
- **[User Guide](../guide/)** - Task-oriented guides
- **[Schema Guide](../schema/)** - Advanced features

### **Q: How do I report bugs or request features?**
**A:** 
- **GitHub Issues**: [github.com/ebogdum/konfigo/issues](https://github.com/ebogdum/konfigo/issues)
- **Bug reports**: Include Konfigo version and minimal reproduction steps
- **Feature requests**: Describe your use case and desired behavior

### **Q: Can I contribute to Konfigo?**
**A:** Yes! Check the project's contribution guidelines on GitHub. Documentation improvements, bug fixes, and feature contributions are welcome.

---

Still have questions? Check the [Troubleshooting Guide](./troubleshooting.md) or browse the [complete documentation](../index.md).
