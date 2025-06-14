# Quick Start

This guide will get you up and running with Konfigo in a few simple steps.

## Prerequisites

*   Konfigo installed (see [Installation](./installation.md)).

## 1. Create Configuration Files

Let's start with a couple of simple configuration files.

**`base.json`**:
```json
{
  "serviceName": "my-app",
  "logLevel": "info",
  "port": 8080
}
```

**`environment.yml`** (to override some base settings):
```yaml
logLevel: debug
port: 9090
featureFlags:
  newFeature: true
```

## 2. Merge Configurations

Use Konfigo to merge these files. By default, it outputs to YAML.

```bash
konfigo -s base.json,environment.yml
```

**Output:**
```yaml
featureFlags:
  newFeature: true
logLevel: debug
port: 9090
serviceName: my-app
```
You can see that `logLevel` and `port` from `environment.yml` have overridden the values from `base.json`, and `featureFlags` has been added.

## 3. Using a Simple Schema

Now, let's introduce a schema to perform some basic processing.

**`schema.yml`**:
```yaml
vars:
  - name: "ENVIRONMENT"
    value: "development"

config:
  # This structure will be implicitly used by Konfigo
  # to show how variables are substituted.
  # Actual output structure is based on merged data + schema operations.

transform:
  - type: "setValue"
    path: "deployment.environment"
    value: "${ENVIRONMENT}"
  - type: "renameKey"
    from: "serviceName"
    to: "application.name"

validate:
  - path: "application.name"
    rules:
      required: true
  - path: "port"
    rules:
      type: "integer"
      min: 1024
```

**Run Konfigo with the schema:**
```bash
konfigo -s base.json,environment.yml -S schema.yml
```

**Output:**
```yaml
application:
  name: my-app
deployment:
  environment: development
featureFlags:
  newFeature: true
logLevel: debug
port: 9090
```

**Explanation:**
*   The `ENVIRONMENT` variable was defined in the schema and used by the `setValue` transformer to add `deployment.environment`.
*   `serviceName` was renamed to `application.name`.
*   The `port` (9090) passed validation.

## 4. Output to a File in a Different Format

Let's output the result to a JSON file.

```bash
konfigo -s base.json,environment.yml -S schema.yml -of final_config.json
```

This will create `final_config.json` with the following content:
```json
{
  "application": {
    "name": "my-app"
  },
  "deployment": {
    "environment": "development"
  },
  "featureFlags": {
    "newFeature": true
  },
  "logLevel": "debug",
  "port": 9090
}
```

This quick start covered basic merging, schema usage (variables, transformation, validation), and output control. Explore the [User Guide](./guide/) and [Schema Guide](./schema/) for more advanced features and detailed explanations!
