# CP-002: Authentication and Authorization

## Purpose

Define login and access control for humans, CI/CD, and services.

## Inputs / Outputs

- Input: OIDC assertions, service identity, requested scope/action.
- Output: short-lived access token and authorization decision.

## Public Interfaces

- `konfigo login` (browser/device)
- `POST /v1/auth/exchange-oidc`
- `POST /v1/auth/token`

## Data Contracts

Token claims required:

- `sub`
- `aud`
- `exp`
- `scope_prefixes[]`
- `actor_type` (`human|ci|service`)

## Invariants

- Username/password is disabled by default.
- Human auth requires OIDC SSO.
- CI/CD auth uses workload OIDC exchange.
- Access tokens are short-lived (default 15m).

## Failure Modes

- `401 invalid_identity_token`
- `403 action_not_allowed_for_scope`
- `403 protected_scope_requires_approval`

## Acceptance Criteria

1. Human CLI login works via OIDC browser/device flow.
2. GitHub/GitLab OIDC token can be exchanged for scoped token.
3. Expired token denies all actions.
4. Audit includes actor claims and decision reason.

## Out of Scope

- Legacy password auth.

## Dependencies

- CP-001 scope model.
