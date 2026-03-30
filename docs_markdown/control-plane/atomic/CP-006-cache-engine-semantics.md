# CP-006: Cache Engine Semantics

## Purpose

Define strict caching behavior for remote connector data.

## Inputs / Outputs

- Input: connector read request.
- Output: cache hit/miss decision + value + freshness state.

## Public Interfaces

- `GET /v1/cache/entries/{key}` (internal/admin)
- `POST /v1/cache/invalidate`
- `konfigo cache stats|invalidate`

## Mandatory Rules

1. **Lease/TTL-aware cache**
- Never cache dynamic secrets beyond lease expiry.
- Effective TTL = `min(config_ttl, provider_ttl_or_lease)`.

2. **Stale-While-Revalidate (SWR) + watch invalidation**
- Return stale for bounded SWR window when enabled.
- Trigger async revalidation.
- Invalidate on provider watch signal where available.

## Provider Watch Hooks

- Consul: blocking query index advance
- etcd: watch revision event
- Redis: pub/sub channel event
- AWS AppConfig: version token change
- Optional hooks for Vault version checks

## Data Contracts

## Cache Key Contract

`<provider>:<scope>:<location>:<selector>:<auth_context_hash>`

## Cache Entry Contract

- `key`
- `valueHash`
- `insertedAt`
- `expiresAt`
- `staleUntil`
- `providerVersion`
- `leaseInfo`
- `lastRevalidateAt`

## Invariants

- Secrets are stored encrypted in cache at rest.
- No cache entry outlives lease expiration for dynamic secrets.
- SWR is disabled by default for protected scopes unless explicitly enabled.

## Failure Modes

- `500 cache_encryption_failure`
- `500 cache_state_corruption`
- `503 revalidation_failed` (with stale fallback if allowed)

## Acceptance Criteria

1. Dynamic Vault secret with 30s lease expires from cache <=30s.
2. SWR returns stale value and triggers background refresh.
3. Watch invalidation removes stale entry within configured latency budget.
4. Cache metrics expose hit/miss, stale returns, and revalidation failures.

## Out of Scope

- Multi-layer global cache topology tuning.

## Dependencies

- CP-005 connectors.
