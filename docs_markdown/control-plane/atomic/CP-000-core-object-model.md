# CP-000: Core Object Model

## Purpose

Define immutable control-plane entities and reference format.

## Inputs / Outputs

- Input: publish requests for schema/template/values/files/policies/bundles.
- Output: immutable versioned objects and references.

## Public Interfaces

- `POST /v1/schemas`
- `POST /v1/templates`
- `POST /v1/values`
- `POST /v1/files`
- `POST /v1/policies`
- `POST /v1/bundles`

## Data Contracts

- `Schema@version`: `schema:<name>:<version>`
- `Template@version`: `template:<name>:<version>`
- `Values@version`: `values:<scope>:<version>`
- `Bundle@version`: `bundle:<scope>:<version>`
- `RunRecord`: immutable execution record tied to `bundleDigest`

`Bundle` required fields:

- `scope`
- `schemaRef`
- `templateRef`
- `valuesRef`
- `policyRefs[]`
- `digest`
- `createdAt`

## Invariants

- Published versions are immutable.
- Bundle refs must point to immutable versions only.
- Bundle digest must be deterministic for same inputs.

## Failure Modes

- `400 invalid_ref_format`
- `409 version_already_exists`
- `422 incompatible_bundle_inputs`

## Acceptance Criteria

1. Re-publishing same ref returns `409`.
2. Creating bundle with mutable ref is rejected.
3. Same inputs produce same bundle digest.
4. `GET` returns exact object without mutation.

## Out of Scope

- Promotion policies.
- Connector sync.

## Dependencies

- CP-001 scope model.
- CP-003 immutability policy.
