# Use Cases & Examples

This page demonstrates how to solve common, real-world configuration management problems using Konfigo.

## 1. Environment Promotion (Dev/Staging/Prod)

**The Goal:** Manage a base configuration and layer on environment-specific overrides for staging and production.

**The Setup:**
Create a directory with a base configuration and an override file for each environment.

**`configs/base.yml`**
```yaml
service:
  name: my-awesome-app
  port: 8080
database:
  host: localhost
  user: app_user
logging:
  level: debug
```

**`configs/production.yml`**
```yaml
database:
  host: prod-db.internal.net
logging:
  level: info
```

**The Command:**
The order of sources in the `-s` flag is critical. The last source specified wins in case of conflicts.

```bash
konfigo -s configs/base.yml,configs/production.yml
```

**The Result:**
```json{6,9}
{
  "database": {
    "host": "prod-db.internal.net",
    "user": "app_user"
  },
  "logging": {
    "level": "info"
  },
  "service": {
    "name": "my-awesome-app",
    "port": 8080
  }
}
```

**Explanation:**
The values for `database.host` and `logging.level` were overwritten by `production.yml` because it was the last source loaded. All other values from `base.yml` were preserved.

---

## 2. CI/CD Integration with Dynamic Tags & Secrets

**The Goal:** Build a configuration in a CI/CD pipeline that uses the Git commit tag for the Docker image and injects a database password from a secure environment variable.

**The Setup:**

**`configs/ci.yml`**
```yaml
# This file contains the base structure.
# The image tag will be supplied by a variable.
deployment:
  image: "my-registry.io/my-awesome-app:${RELEASE_VERSION}"
database:
  user: "ci_user"
```

**`schema.yml`**
```yaml
# The schema defines how to get the RELEASE_VERSION variable.
vars:
  - name: "RELEASE_VERSION"
    fromEnv: "CI_COMMIT_TAG" # Read from an env var set by the CI system
    defaultValue: "latest"
validate:
  - path: "database.password"
    rules:
      required: true
      minLength: 16
```

**The Command:**
In your CI/CD script, you would set the secure environment variables and run Konfigo.

```bash{1,2,5}
# These are provided by the CI/CD system's secret management and environment
export KONFIGO_KEY_database.password="a-very-secure-password-from-ci"
export CI_COMMIT_TAG="v1.2.3"

konfigo \
  -S schema.yml \
  -s configs/ci.yml
```

**The Result:**
```json
{
  "database": {
    "password": "a-very-secure-password-from-ci",
    "user": "ci_user"
  },
  "deployment": {
    "image": "my-registry.io/my-awesome-app:v1.2.3"
  }
}
```

**Explanation:**
- `KONFIGO_KEY_database.password` directly injected the secret into the configuration.
- The `vars` block in the schema read the `CI_COMMIT_TAG` environment variable.
- The `${RELEASE_VERSION}` placeholder was substituted with `v1.2.3` during processing.
