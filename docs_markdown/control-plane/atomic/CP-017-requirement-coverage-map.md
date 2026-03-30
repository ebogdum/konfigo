# CP-017: Requirement Coverage Map

## Purpose

Map requested control-plane requirements to atomic specs for implementation tracking.

## Inputs / Outputs

- Input: requirement statements and feature requests.
- Output: deterministic mapping to atomic spec IDs.

## Public Interfaces

- This file is consumed by planning/review process; no runtime API surface.

## Data Contracts

- Requirement item -> one or more `CP-XXX` atomic spec references.
- Every mapped item must reference at least one acceptance-criteria-bearing spec.

## Invariants

- No requested requirement remains unmapped.
- Mapping references only existing files.

## Failure Modes

- stale mapping after spec renames
- missing mapping for newly added requirements

## Acceptance Criteria

1. Every requirement block has at least one spec reference.
2. All references in this file resolve to existing atomic files.
3. Mapping is updated before scope changes are implemented.

## Coverage Matrix

### Remote Backends and Connectors

- Vault KV1/KV2/logical/dynamic/transit -> `CP-005`
- AWS SSM/Secrets Manager/AppConfig/S3/DynamoDB -> `CP-005`
- GCP Secret Manager/GCS/Firestore/Runtime Config -> `CP-005`
- Azure Key Vault/App Configuration/Blob/Cosmos -> `CP-005`
- SOPS + key backends -> `CP-005`
- Consul/etcd/Redis + watch hooks -> `CP-005`, `CP-006`
- Extras (K8s, HTTP ETag, Git ref/file) -> `CP-005`

### Caching Rules

- Lease/TTL-aware cache hard rule -> `CP-006`
- SWR + watch invalidation hard rule -> `CP-006`

### Control Plane Capabilities

- Store/version schema/template/values/files/policies/bundles -> `CP-000`
- Remote run from server -> `CP-007`
- Pinning and immutable bindings -> `CP-003`
- Value provenance per resolved path -> `CP-004`
- Governance (RBAC, approvals, audit, protected env) -> `CP-002`, `CP-016`

### Core Model

- Project/Environment -> `CP-000`
- Schema@version/Template@version/Values@version/Bundle@version -> `CP-000`
- Policy model -> `CP-000`
- RunRecord -> `CP-004`

### Pinning and Safety

- Immutable versions -> `CP-003`
- Protected refs with approvals -> `CP-003`, `CP-016`
- Locked template-schema bindings -> `CP-003`
- Optional signed bundles -> `CP-003`, `CP-008`

### CLI Flow

- `konfigo login` -> `CP-002`, `CP-007`
- `konfigo push ...` -> `CP-007`
- `konfigo bundle create --pin` -> `CP-003`, `CP-007`
- `konfigo run --server ... --bundle ...` -> `CP-007`

### Offline Mode

- Cached bundle execution -> `CP-008`
- Signature verification for cached artifacts -> `CP-008`

## Notes

If any requirement changes, update this map first, then update only affected atomic specs.

## Out of Scope

- Defining runtime feature behavior directly (this file maps, not specifies).

## Dependencies

- CP-000 through CP-016.
