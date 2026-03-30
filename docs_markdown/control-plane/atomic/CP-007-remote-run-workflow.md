# CP-007: Remote Run Workflow

## Purpose

Define running Konfigo from anywhere by pulling everything from control plane.

## Inputs / Outputs

- Input: bundle reference + scope + optional runtime overrides.
- Output: resolved output artifact(s) and run record.

## Public Interfaces

CLI flow:

1. `konfigo login`
2. `konfigo push schema ./schema.yml --name app-schema`
3. `konfigo push template ./template.yml --name app-template`
4. `konfigo push values ./values.prod.yml --env prod --scope platform/app/prod/api`
5. `konfigo bundle create --name release-2026-02-25 --pin`
6. `konfigo run --server https://konfigo.company --bundle release-2026-02-25`

API:

- `POST /v1/runs`
- `GET /v1/runs/{runId}`
- `GET /v1/runs/{runId}/artifact`

## Data Contracts

Run request:

```json
{
  "scope": "platform/app/prod/api",
  "bundleRef": "release-2026-02-25",
  "outputFormat": "yaml",
  "overrides": {}
}
```

## Invariants

- Run must resolve exact pinned bundle digest.
- Output must include `bundleRef` and `bundleDigest` metadata.
- Any override must be policy-allowed and provenance-marked.

## Failure Modes

- `404 bundle_not_found`
- `409 bundle_not_promoted_for_scope`
- `403 override_not_allowed`
- `422 connector_resolution_failed`

## Acceptance Criteria

1. Remote run succeeds with no local schema/template files.
2. Output artifact metadata contains bundle digest.
3. Denied override request returns policy reason.

## Out of Scope

- Interactive UI run editor.

## Dependencies

- CP-000 core model.
- CP-004 provenance.
