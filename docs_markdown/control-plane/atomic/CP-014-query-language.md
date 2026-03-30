# CP-014: Query Language

## Purpose

Provide one searchable language for objects, references, drift, and impact.

## Inputs / Outputs

- Input: query text.
- Output: paginated deterministic result set.

## Public Interfaces

- `POST /v1/query/execute`
- `POST /v1/query/where-used`
- `POST /v1/query/impact`
- `konfigo query "..."`

## Data Contracts

Supported operators (phase 1):

- `=`, `!=`, `<`, `>`, `<=`, `>=`
- `and`, `or`, `not`
- `in [..]`
- `like` (`%` wildcard)

## Invariants

- Query results filtered by scope authorization before return.
- Stable ordering for same query + cursor.

## Failure Modes

- `400 invalid_query_syntax`
- `403 query_scope_denied`

## Acceptance Criteria

1. `where_used schema:...` returns all direct bundle refs.
2. Scope ACL prevents cross-scope result leakage.
3. Large results paginate with cursor.

## Out of Scope

- Arbitrary SQL support.

## Dependencies

- CP-000 model indexes.
- CP-002 authz filters.
