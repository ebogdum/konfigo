# CP-016: Governance, Approvals, and Audit

## Purpose

Define governance controls for protected scopes and environments.

## Inputs / Outputs

- Input: change request (publish, promote, override).
- Output: approved/rejected decision and immutable audit trail.

## Public Interfaces

- `POST /v1/approvals/requests`
- `POST /v1/approvals/requests/{id}/approve`
- `POST /v1/approvals/requests/{id}/reject`
- `GET /v1/audit/events`

## Data Contracts

Approval request fields:

- `requestId`
- `scope`
- `action`
- `targetRef`
- `requestedBy`
- `requiredApprovers`
- `status`

Audit event fields:

- `auditId`
- `occurredAt`
- `actor`
- `scope`
- `action`
- `decision`
- `reason`
- `relatedRefs[]`

## Invariants

- Protected scopes require approval before reference movement.
- Approver cannot self-approve when policy forbids it.
- Audit entries are immutable and append-only.

## Failure Modes

- `403 approval_required`
- `403 self_approval_not_allowed`
- `409 approval_request_expired`

## Acceptance Criteria

1. Promotion to `*/prod/*` without approval is blocked.
2. Approved request unblocks only exact referenced action.
3. Rejection reason is mandatory and queryable in audit trail.
4. Audit timeline reconstructs full decision chain for any promotion.

## Out of Scope

- External ticketing system workflow orchestration.

## Dependencies

- CP-002 authz.
- CP-003 protected refs.
- CP-010 events.
