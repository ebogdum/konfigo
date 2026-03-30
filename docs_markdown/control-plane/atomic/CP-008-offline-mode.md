# CP-008: Offline Mode

## Purpose

Allow execution when control plane/network is unavailable, using verified cached bundles.

## Inputs / Outputs

- Input: offline run request + local cache.
- Output: run artifact or explicit offline policy denial.

## Public Interfaces

- `konfigo bundle cache pull --bundle <ref>`
- `konfigo run --offline --bundle <ref>`

## Data Contracts

Offline cache package must include:

- bundle manifest
- bundle digest
- signature (if required)
- dependency refs and digests
- fetched timestamp
- expiry policy metadata

## Invariants

- Offline execution allowed only for non-expired cache package.
- Signature verification required when scope policy enforces signed bundles.
- Offline mode records `mode=offline` in RunRecord.

## Failure Modes

- `404 offline_bundle_not_cached`
- `422 offline_bundle_expired`
- `422 offline_signature_verification_failed`
- `403 offline_not_allowed_for_scope`

## Acceptance Criteria

1. Cached bundle runs with network disabled.
2. Expired cache package is rejected.
3. Tampered bundle cache fails signature/digest verification.
4. RunRecord marks mode and reason for offline execution.

## Out of Scope

- Offline promotion or write operations.

## Dependencies

- CP-003 pinning/signature rules.
- CP-004 run records.
