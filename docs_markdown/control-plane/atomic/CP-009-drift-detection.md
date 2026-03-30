# CP-009: Drift Detection

## Purpose

Detect divergence between desired bundle state and actual provider/runtime state.

## Inputs / Outputs

- Input: scope + bundleRef + provider set.
- Output: drift report with typed findings and severity.

## Public Interfaces

- `POST /v1/drift/scans`
- `GET /v1/drift/scans/{id}`
- `konfigo drift scan --scope ... --bundle ...`

## Data Contracts

Drift types:

- `missing`
- `extra`
- `value_mismatch`
- `schema_violation`
- `policy_violation`
- `stale_version`

## Invariants

- Drift scan always references a concrete bundle digest.
- Findings must include path + source metadata.

## Failure Modes

- `404 bundle_not_found`
- `422 actual_state_fetch_failed`

## Acceptance Criteria

1. Scan returns deterministic findings for same inputs.
2. Severity classification matches policy mapping.
3. `drift.finding.opened` and `drift.finding.resolved` events are emitted.

## Out of Scope

- Automatic remediation.

## Dependencies

- CP-005 connectors.
- CP-010 events.
