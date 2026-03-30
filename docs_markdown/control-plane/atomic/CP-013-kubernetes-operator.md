# CP-013: Kubernetes Operator

## Purpose

Materialize pinned Konfigo bundles into Kubernetes runtime resources.

## Inputs / Outputs

- Input: `KonfigoBundleRef` CR.
- Output: updated ConfigMap/Secret + status conditions.

## Public Interfaces

CRD:

- `KonfigoBundleRef`

Fields:

- `spec.scope`
- `spec.bundleRef`
- `spec.target.kind` (`ConfigMap|Secret`)
- `spec.target.name`

## Data Contracts

`KonfigoBundleRef.status` fields:

- `observedBundleRef`
- `observedBundleDigest`
- `conditions[]`
- `lastSuccessAt`

## Invariants

- Operator always applies exact requested bundle ref.
- On failure, last known good target is preserved.

## Failure Modes

- `Degraded` condition with reason:
  - auth failure
  - bundle fetch failure
  - verification failure

## Acceptance Criteria

1. CR apply reconciles target resource from bundle.
2. Bundle ref change updates target deterministically.
3. Reconcile failure sets status and emits event.

## Out of Scope

- Runtime request/response enforcement.

## Dependencies

- CP-007 remote run APIs.
- CP-010 events.
