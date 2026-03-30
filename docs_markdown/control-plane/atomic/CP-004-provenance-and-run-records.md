# CP-004: Provenance and Run Records

## Purpose

Make every resolved value and run explainable.

## Inputs / Outputs

- Input: bundle execution request.
- Output: resolved config + provenance map + run record.

## Public Interfaces

- `konfigo run --server ... --bundle ... --explain`
- `GET /v1/runs/{runId}`
- `GET /v1/runs/{runId}/provenance`

## Data Contracts

Per-value provenance fields:

- `path`
- `resolvedValueHash`
- `sourceType` (`bundle_values|connector|override|default`)
- `sourceRef` (e.g. `values:...:v9`, `vault:kv/...#42`)
- `resolvedAt`

RunRecord fields:

- `runId`
- `scope`
- `bundleRef`
- `bundleDigest`
- `actor`
- `result`
- `startedAt/completedAt`

## Invariants

- Every output path has provenance entry unless intentionally omitted by policy.
- Sensitive values are redacted; hashes are retained.

## Failure Modes

- `500 provenance_capture_failed`
- `403 provenance_access_denied`

## Acceptance Criteria

1. `--explain` returns deterministic provenance entries.
2. Run record always includes bundle digest and actor.
3. Redaction policy masks secret values but preserves source metadata.

## Out of Scope

- Full GUI explain visualizations.

## Dependencies

- CP-000 object model.
- CP-005 connectors.
