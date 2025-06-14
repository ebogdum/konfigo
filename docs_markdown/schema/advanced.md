# Schema: Advanced Features

This page covers advanced schema features for fine-grained control over your configuration.

## Immutable Fields

The `immutable` block lets you define a list of keys that, once set by an early-loading source, cannot be overwritten by a later source. This is useful for protecting critical, foundational values.

**`schema.yml`**
```yaml
immutable:
  - "service.name"
  - "database.port"
```

**Scenario:**
1.  `01-base.json` sets `database.port` to `5432`.
2.  `02-prod.yml` attempts to set `database.port` to `9999`.

Because `database.port` is immutable, the overwrite from `02-prod.yml` will be ignored, and the final value will remain `5432`.

::: danger
This rule is very strict and applies to all sources, including `KONFIGO_KEY_` environment variables.
:::

## Output Schema Filtering

The `outputSchema` block allows you to define a "mask" or "template" for your final output. Konfigo will produce a final configuration that **only** contains the keys present in your output schema file, discarding everything else.

This is useful for creating a clean, public-facing configuration from a larger, internal one that may contain secrets or intermediate values.

**`schema.yml`**
```yaml
outputSchema:
  path: "./public-config-shape.json"
```

**`public-config-shape.json`**
```json
{
  "service": {
    "name": "",
    "url": ""
  },
  "features": {
    "new_ui": false
  }
}
```

The final output from Konfigo will only contain the `service` and `features` blocks.
