# CP-005: Provider Connector Matrix

## Purpose

Define supported remote connectors and exact feature surface for each.

## Inputs / Outputs

- Input: provider config + scope prefix + key selectors.
- Output: normalized key/value documents + source metadata.

## Public Interfaces

- `POST /v1/connectors`
- `GET /v1/connectors`
- `POST /v1/connectors/{id}/test`
- `konfigo connector add|test|list`

## Data Contracts

## Connector Matrix (Phase Targets)

### Vault

- KV v2 reads
- KV v1 reads
- generic logical path reads
- dynamic secret engines: database/aws/gcp/azure
- lease metadata capture
- optional Transit decrypt reads

### AWS

- SSM Parameter Store: `GetParameter`, `GetParametersByPath`
- Secrets Manager reads
- AppConfig reads
- S3 object config reads
- optional DynamoDB document source

### GCP

- Secret Manager reads
- GCS object config reads
- optional Firestore document source
- optional/deprecated Runtime Config adapter

### Azure

- Key Vault Secrets reads
- Azure App Configuration (keys, labels, feature flags)
- Blob object config reads
- optional Cosmos DB document source

### SOPS

- Encrypted YAML/JSON/TOML/ENV files
- key backends: age, PGP, AWS KMS, GCP KMS, Azure Key Vault
- local and remote file source support

### Consul

- KV key reads
- KV prefix reads
- blocking query/watch support

### etcd

- v3 key reads
- v3 prefix reads
- watch support

### Redis

- string key reads
- hash field reads
- optional RedisJSON document reads
- optional pub/sub invalidation hooks

### Nice-to-Have Extras

- Kubernetes ConfigMap/Secret source
- HTTP(S) config endpoint source with ETag behavior
- Git ref/file source

### Normalized Output Contract

Each fetched item must provide:

- `path`
- `value` (or redacted marker)
- `source.provider`
- `source.location`
- `source.version` (if available)
- `source.etag_or_revision` (if available)
- `source.lease` (if available)
- `fetchedAt`

## Invariants

- Connector never returns provider-native shape to downstream pipeline.
- All connectors output normalized contract.
- Connector errors are typed and include provider + location.

## Failure Modes

- `400 invalid_connector_config`
- `401 connector_auth_failed`
- `404 connector_path_not_found`
- `429 connector_rate_limited`
- `503 connector_upstream_unavailable`

## Acceptance Criteria

1. Each listed phase-1 provider has at least one integration test.
2. Test endpoint validates credentials and minimal read path.
3. Normalized output contract matches schema for all providers.
4. Unsupported optional provider features return explicit capability error.

## Out of Scope

- Write-back to remote providers.

## Dependencies

- CP-006 cache semantics.
- CP-002 auth.
