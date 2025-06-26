# Quick Start: Your First Success in 5 Minutes

Welcome! Let's get you up and running with Konfigo immediately. By the end of this guide, you'll have successfully merged configuration files and understand how Konfigo works.

## Prerequisites

- Konfigo installed ([Installation Guide](./installation.md))
- Basic command line knowledge

## Step 1: Create Your First Configuration Files

Let's start with the simplest possible example. Create these two files:

::: code-group

```json [base.json]
{
  "app": {
    "name": "my-awesome-app",
    "port": 8080,
    "debug": false
  },
  "database": {
    "host": "localhost",
    "port": 5432
  }
}
```

```json [production.json]
{
  "app": {
    "port": 9090,
    "debug": false
  },
  "database": {
    "host": "prod-db.company.com",
    "ssl": true
  }
}
```

:::

::: details Click here if you prefer YAML format

```yaml
# base.yaml
app:
  name: "my-awesome-app"
  port: 8080
  debug: false
database:
  host: "localhost"
  port: 5432
```

```yaml
# production.yaml
app:
  port: 9090
  debug: false
database:
  host: "prod-db.company.com"
  ssl: true
```

:::

## Step 2: Your First Merge

Run this command to merge the files:

```bash
# Copy and paste this command:
konfigo -s base.json,production.json
```

::: tip Quick Copy
💡 **Pro tip**: Click the copy button in the top-right corner of code blocks to copy commands instantly!
:::

**Expected output**:
```json
{
  "app": {
    "name": "my-awesome-app",
    "port": 9090,
    "debug": false
  },
  "database": {
    "host": "prod-db.company.com",
    "port": 5432,
    "ssl": true
  }
}
```

**🎉 Congratulations!** You just merged two configuration files.

::: details 🔍 What Just Happened? (Click to expand)

Konfigo performed an intelligent merge:

1. **📖 Read both files** in the order you specified: `base.json` first, then `production.json`
2. **🏗️ Started with base.json** as the foundation
3. **🔄 Applied production.json** on top, overriding some values and adding new ones
4. **📤 Output the merged result** to your terminal

**Merge Logic**:
- ✅ `app.name` came from base.json (not overridden)
- 🔄 `app.port` was overridden by production.json (9090 instead of 8080)
- ✅ `database.port` came from base.json (not in production.json)
- ➕ `database.ssl` was added by production.json (new field)

This is **deep merging** - objects are combined intelligently rather than completely replaced.

:::

## Step 3: Save the Result

Save the merged configuration to a file:

```bash
# Save to JSON file
konfigo -s base.json,production.json -of final-config.json
```

::: code-group

```bash [Save as JSON]
konfigo -s base.json,production.json -of final-config.json
```

```bash [Save as YAML]
konfigo -s base.json,production.json -oy -of final-config.yaml
```

```bash [Save as TOML]
konfigo -s base.json,production.json -ot -of final-config.toml
```

:::

Check the result:
```bash
cat final-config.json
```

## Step 4: Try Different Formats

Konfigo works with multiple formats. Let's create a YAML override:

**`local-overrides.yaml`**:
```yaml
app:
  debug: true
  logLevel: "verbose"
```

Now merge all three:
```bash
konfigo -s base.json,production.json,local-overrides.yaml
```

**Output**:
```json
{
  "app": {
    "name": "my-awesome-app",
    "port": 9090,
    "debug": true,
    "logLevel": "verbose"
  },
  "database": {
    "host": "prod-db.company.com",
    "port": 5432,
    "ssl": true
  }
}
```

**What Just Happened?**
- Konfigo automatically detected the YAML format
- Applied the overrides in order: base → production → local
- The final `debug: true` came from the YAML file

## Step 5: Use Environment Variables

Override any value using environment variables:

```bash
KONFIGO_KEY_app.port=3000 konfigo -s base.json,production.json
```

Notice how `app.port` is now 3000, overriding both files.

**Environment variables always win** - they have the highest precedence.

## Step 6: Convert Formats

Output in different formats:

```bash
# Output as YAML
konfigo -s base.json,production.json -oy

# Output as TOML  
konfigo -s base.json,production.json -ot

# Save as YAML file
konfigo -s base.json,production.json -of config.yaml
```

## 🎯 Success Checklist

You've successfully:
- ✅ Merged JSON configuration files
- ✅ Mixed different formats (JSON + YAML)
- ✅ Used environment variable overrides
- ✅ Converted between formats
- ✅ Saved results to files

## Common Patterns You'll Use

### Environment-Specific Configs
```bash
# Development
konfigo -s base.json,dev.yaml -of dev-config.json

# Production
konfigo -s base.json,prod.yaml -of prod-config.json
```

### Runtime Overrides
```bash
# Override settings at runtime
KONFIGO_KEY_database.host=$DB_HOST konfigo -s config.yaml
```

### Format Conversion
```bash
# Convert legacy .env to modern YAML
konfigo -s legacy.env -oy -of modern.yaml
```

## What's Next?

Now that you understand the basics, here are your next steps:

### **Immediate Next Steps** (5-10 minutes)
- **[Basic Concepts](./concepts.md)** - Understand how Konfigo's processing pipeline works

### **Common Tasks** (30 minutes)
- **[User Guide](../guide/)** - Learn specific tasks like environment variables and validation

### **Advanced Power** (1-2 hours)  
- **[Schema Guide](../schema/)** - Unlock advanced features like variables, validation, and transformation

## Need Help?

- **Having issues?** Check the [Troubleshooting Guide](../reference/troubleshooting.md)
- **Want more examples?** Browse [Recipes & Examples](../guide/recipes.md)
- **Ready for schemas?** Start with [Schema Basics](../schema/)

**Great job getting started!** Konfigo has much more to offer - explore at your own pace.
