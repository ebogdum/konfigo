# Environment Variables

Konfigo leverages environment variables for two primary purposes: directly setting configuration values and providing variables for substitution within schemas. This offers a flexible way to manage environment-specific settings, secrets, or dynamic parameters without modifying your configuration files or schemas.

## Overriding Configuration Values (`KONFIGO_KEY_...`)

You can directly set or override any configuration value by using environment variables prefixed with `KONFIGO_KEY_`. The part of the variable name following this prefix is treated as a dot-separated path to the desired key within your configuration structure.

*   **Syntax**: `KONFIGO_KEY_path.to.your.key=value`
*   **Precedence**: These overrides have the **highest precedence** over all other configuration sources, including files specified with `-s`, and even immutable paths defined in a schema.
*   **Use Case**: Ideal for injecting sensitive data (like API keys or database passwords) from a secure environment (e.g., CI/CD pipeline secrets) or for making quick, temporary changes without altering files.

### Examples

Given a base configuration file `config.yml`:

```yaml
server:
  host: localhost
  port: 8080
database:
  url: postgres://user:pass@host/db
```

1.  **Override server port**:
    ```bash
    export KONFIGO_KEY_server.port=9000
    konfigo -s config.yml
    ```
    Output (YAML):
    ```yaml
    server:
      host: localhost
      port: 9000 # Overridden
    database:
      url: postgres://user:pass@host/db
    ```

2.  **Override a nested database URL**:
    ```bash
    export KONFIGO_KEY_database.url="mysql://new_user:new_pass@new_host/new_db"
    konfigo -s config.yml
    ```
    Output (YAML):
    ```yaml
    server:
      host: localhost
      port: 8080
    database:
      url: "mysql://new_user:new_pass@new_host/new_db" # Overridden
    ```

3.  **Set a new key**:
    ```bash
    export KONFIGO_KEY_featureFlags.newFeature=true
    konfigo -s config.yml
    ```
    Output (YAML):
    ```yaml
    server:
      host: localhost
      port: 8080
    database:
      url: postgres://user:pass@host/db
    featureFlags:
      newFeature: true # New key added
    ```

## Supplying Substitution Variables (`KONFIGO_VAR_...`)

Environment variables prefixed with `KONFIGO_VAR_` are used to supply values for variable substitution (e.g., `${MY_VARIABLE}`) within your configuration files and schema directives.

*   **Syntax**: `KONFIGO_VAR_VARNAME=value`
*   **Precedence**: These variables have the **highest priority** in the variable resolution order:
    1.  `KONFIGO_VAR_...` environment variables.
    2.  Variables from the file specified by the `-V` or `--vars-file` flag.
    3.  Variables defined in the schema's `vars:` block.
*   **Use Case**: Perfect for providing dynamic values that change between environments (e.g., API endpoints, release versions, resource limits) to be used in schema processing or directly in config templates.

### Examples

Given a schema `schema.yml`:

```yaml
vars:
  - name: "GREETING"
    defaultValue: "Hello from schema"
  - name: "TARGET_ENV"
    fromEnv: "DEPLOY_ENV" # Will try to read DEPLOY_ENV from system
    defaultValue: "development"
config:
  message: "${GREETING}, running in ${TARGET_ENV}!"
  serviceUrl: "${API_ENDPOINT}"
```

And a base configuration `config.json`:
```json
{
  "someKey": "someValue"
}
```

1.  **Override `GREETING` and provide `API_ENDPOINT` via `KONFIGO_VAR_`**:
    ```bash
    export KONFIGO_VAR_GREETING="Hi from environment"
    export KONFIGO_VAR_API_ENDPOINT="https://api.production.example.com"
    konfigo -s config.json -S schema.yml
    ```
    Output (YAML):
    ```yaml
    message: "Hi from environment, running in development!" # GREETING overridden, TARGET_ENV uses schema default
    serviceUrl: "https://api.production.example.com" # API_ENDPOINT provided
    someKey: someValue
    ```

2.  **Let `TARGET_ENV` be sourced from `DEPLOY_ENV` (which is also set)**:
    ```bash
    export DEPLOY_ENV="staging" # System environment variable
    export KONFIGO_VAR_API_ENDPOINT="https://api.staging.example.com"
    konfigo -s config.json -S schema.yml
    ```
    Output (YAML):
    ```yaml
    message: "Hello from schema, running in staging!" # GREETING uses schema default, TARGET_ENV uses DEPLOY_ENV
    serviceUrl: "https://api.staging.example.com"
    someKey: someValue
    ```

3.  **`KONFIGO_VAR_` takes precedence over `-V` file**:

    `vars.yml` (passed with `-V vars.yml`):
    ```yaml
    GREETING: "Hello from vars file"
    API_ENDPOINT: "https://api.varsfile.example.com"
    ```

    Command:
    ```bash
    export KONFIGO_VAR_GREETING="Top priority greeting"
    konfigo -s config.json -S schema.yml -V vars.yml
    ```
    Output (YAML):
    ```yaml
    message: "Top priority greeting, running in development!" # GREETING from KONFIGO_VAR_
    serviceUrl: "https://api.varsfile.example.com" # API_ENDPOINT from vars.yml (as no KONFIGO_VAR_API_ENDPOINT)
    someKey: someValue
    ```

By effectively using these two types of environment variables, you can create highly adaptable and secure configuration workflows with Konfigo.
