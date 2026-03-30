# CP-001: Scope Prefix Model

## Purpose

Define logical organization using a single scope string, not hard namespaces.

## Inputs / Outputs

- Input: any API/CLI action requiring target location.
- Output: canonical scope string used across storage, auth, events, and drift.

## Public Interfaces

- `scope` field in all write/read APIs.
- `--scope` flag in CLI.

## Data Contracts

Format: `team/app/env/service`.

Validation:

- lowercase alnum + `-` in each segment
- segment count: 3..6
- no empty segments

Examples:

- `platform/payments/prod/api`
- `core/identity/dev/auth-service`

## Invariants

- Scope is case-sensitive canonical lowercase.
- Prefix matching is path-segment aware.

## Failure Modes

- `400 invalid_scope_format`
- `403 scope_access_denied`

## Acceptance Criteria

1. Invalid scope format is rejected with clear reason.
2. Prefix ACLs match only complete segments.
3. Same scope string is used in events and audit records.

## Out of Scope

- Tenant isolation by namespace.

## Dependencies

- CP-002 authz checks.
