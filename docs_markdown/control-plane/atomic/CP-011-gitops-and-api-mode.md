# CP-011: GitOps and API Dual Mode

## Purpose

Support PR-driven and direct API workflows using one shared publication pipeline.

## Inputs / Outputs

- Input: Git commit state or direct API write request.
- Output: immutable objects + bundle publication or rejection.

## Public Interfaces

- `POST /v1/gitops/sources`
- `POST /v1/gitops/reconcile`
- `POST /v1/*` object publish APIs

## Data Contracts

Scope mode:

- `gitops_only`
- `api_only`
- `hybrid`

## Invariants

- Both modes run same validation/policy checks.
- Same input content leads to same object digest regardless of mode.

## Failure Modes

- `409 mode_conflict`
- `403 source_mode_not_allowed`

## Acceptance Criteria

1. GitOps and API publish produce identical bundle digest for same content.
2. `gitops_only` scope rejects direct API writes.
3. Hybrid conflict is deterministic and auditable.

## Out of Scope

- Multi-repo dependency graph optimizer.

## Dependencies

- CP-010 events.
- CP-002 authz.
