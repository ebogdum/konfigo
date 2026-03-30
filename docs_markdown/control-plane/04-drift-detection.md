# 04: Drift Detection

## Overview

Drift detection compares **desired state** (pinned Konfigo bundle and policy) with **actual state** (runtime/config provider values) for a specific scope.

This feature answers:

- What changed outside approved Konfigo workflows?
- Which differences are dangerous vs informational?
- What is the recommended remediation path?

## Goals

- Detect configuration drift continuously and on demand.
- Classify drift by type and severity.
- Emit structured drift events for automation.
- Provide actionable remediation commands and APIs.

## Non-Goals

- Automatic self-healing in initial version.
- Runtime payload enforcement (HTTP/gRPC stream checks).
- Full policy engine redesign.

## Drift Types

Drift result entries must use one of these classes:

- `missing`: expected key/path absent in actual state.
- `extra`: unexpected key/path exists in actual state.
- `value_mismatch`: same key exists but value differs.
- `schema_violation`: actual value no longer conforms to schema.
- `policy_violation`: value exists but violates policy constraints.
- `stale_version`: runtime uses outdated bundle reference.

## Scope Model

Every drift run is bound to a `scope` prefix:

- Example: `platform/payments/prod/api`
- Vault source path: `kv/konfigo/platform/payments/prod/api/...`
- S3 source path: `s3://konfigo-config/platform/payments/prod/api/...`

## Detection Modes

### 1. On-Demand Scan

Triggered by user/CI:

- CLI: `konfigo drift scan --scope platform/payments/prod/api --bundle release-2026-02-25`
- API: `POST /v1/drift/scans`

### 2. Scheduled Scan

Configured cadence per scope:

- every 5m, 15m, 1h, etc.
- Optional quiet windows.

### 3. Change-Triggered Scan

Triggered by source change signals:

- Consul watch
- etcd watch
- Redis pub/sub signal
- Vault secret version update

## Data Inputs

A scan builds two normalized snapshots:

1. **Desired snapshot**
- Resolved from `Bundle@version`
- Includes schema and policy context

2. **Actual snapshot**
- Pulled from configured providers
- Includes metadata (version, lease, etag, updatedAt)

## Normalization Rules

- Convert all sources to canonical map form.
- Apply deterministic ordering for object keys.
- Normalize number/bool/string format where schema allows.
- Preserve provider metadata separately from value payload.

## Detection Algorithm

1. Load bundle, schema, policy for scope.
2. Resolve desired values to canonical representation.
3. Fetch actual values from configured providers.
4. Run structural diff (`missing`, `extra`, `value_mismatch`).
5. Validate actual against schema and policy.
6. Mark stale bundle or stale template/schema refs.
7. Assign severity score per finding.
8. Persist report and emit events.

## Severity Model

- `critical`: policy violation, schema violation on required key.
- `high`: missing required key, stale version in protected scope.
- `medium`: value mismatch on managed key.
- `low`: extra key not on denylist.
- `info`: cosmetic non-impacting differences.

## Suppression and Ignore Rules

Support scoped suppressions:

- Temporary suppression (`expiresAt`).
- Path-based ignore patterns.
- Provider-metadata-only ignore toggles.

Example ignore rule:

```yaml
ignore:
  - scope: "platform/payments/prod/api"
    path: "metadata.lastRotated"
    reason: "managed by external rotator"
    expiresAt: "2026-03-15T00:00:00Z"
```

## API Contract

### Create Scan

`POST /v1/drift/scans`

```json
{
  "scope": "platform/payments/prod/api",
  "bundleRef": "release-2026-02-25",
  "providers": ["vault", "s3"],
  "mode": "on_demand"
}
```

### Get Scan Result

`GET /v1/drift/scans/{scanId}`

Response includes:

- summary counts by drift type and severity
- detailed findings
- recommended remediation actions

### List Active Drift

`GET /v1/drift/findings?scope=platform/payments/prod/api&status=open`

## CLI UX

- `konfigo drift scan --scope <scope> --bundle <bundleRef>`
- `konfigo drift list --scope <scope>`
- `konfigo drift show <scanId>`
- `konfigo drift resolve <findingId> --reason "accepted external change"`

## Event Emission

Emit at least:

- `drift.scan.started`
- `drift.scan.completed`
- `drift.finding.opened`
- `drift.finding.resolved`

Event payload includes `scope`, `bundleRef`, `findingType`, `severity`, and source metadata.

## Remediation Patterns

- Re-apply pinned bundle to target provider.
- Promote updated bundle through approval flow.
- Add justified suppression when external system is authoritative.

## Security

- Scan permissions require read access to scope and provider credentials.
- Findings are access-controlled by scope prefix.
- Sensitive values are redacted in reports, but hash fingerprints are retained.

## Metrics

- `drift_scan_duration_seconds`
- `drift_findings_total{type,severity}`
- `drift_open_findings`
- `drift_time_to_resolution_seconds`

## Rollout Plan

### Phase 1

- On-demand scans for Vault + S3.
- Structural diff only.

### Phase 2

- Scheduled scans and schema/policy drift classes.

### Phase 3

- Source-triggered scans (watch integrations).
- Suppression lifecycle automation.

## Testing Strategy

- Unit tests for diff classifiers and normalizer.
- Provider integration tests with fixture data.
- End-to-end tests on protected scope with stale version scenarios.
- Load tests for high-cardinality scopes.

## Open Decisions

- Maximum supported scan interval granularity.
- Default retention period for scan reports.
- Whether to include optional auto-remediation for low-severity drift in a later version.
