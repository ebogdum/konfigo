# CP-003: Pinning and Immutability

## Purpose

Prevent mutable drift by enforcing immutable versions and protected reference movement.

## Inputs / Outputs

- Input: bundle create/promotion request.
- Output: pinned immutable bundle or explicit rejection.

## Public Interfaces

- `konfigo bundle create --pin`
- `POST /v1/bundles`
- `POST /v1/promotions`

## Data Contracts

`lockedBindings` example:

```json
{
  "templateRef": "template:app:v12",
  "schemaRef": "schema:app:v7",
  "mode": "hard"
}
```

Protected refs example:

- `platform/payments/prod/stable`

## Invariants

- `hard` locked binding cannot be overridden.
- Protected refs move only via approval policy.
- Optional bundle signature verification controlled by policy.

## Failure Modes

- `409 locked_binding_violation`
- `403 protected_ref_update_denied`
- `422 signature_verification_failed`

## Acceptance Criteria

1. Attempt to pair `template:v12` with non-allowed schema is rejected.
2. Protected ref update without approval is rejected.
3. Signed bundle verification works when policy requires signatures.

## Out of Scope

- Signature key management UX.

## Dependencies

- CP-002 authz.
- CP-000 object model.
