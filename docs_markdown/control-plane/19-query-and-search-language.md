# 19: Query and Search Language

## Overview

Konfigo needs a native query language to answer operational questions quickly:

- Where is this key used?
- Which bundles reference this schema version?
- What changed for a specific scope over time?

This feature provides a consistent query interface across objects, versions, references, and drift findings.

## Goals

- Provide one query syntax for CLI and API.
- Support both key/value search and graph-style dependency queries.
- Enable fast impact analysis before promotions.

## Non-Goals

- Replacing SQL or external analytics warehouses.
- Arbitrary unbounded joins in phase 1.

## Query Model

Two query modes:

1. **Filter mode**: object lists with predicates.
2. **Graph mode**: relationship traversal (`where-used`, `depends-on`).

## Domain Entities

- Schema
- Template
- Values
- Bundle
- Policy
- Promotion
- Drift finding
- Event

## Example Queries

### Filter Queries

- `bundle where scope = "platform/payments/prod/api" and status = "published"`
- `values where scope like "platform/payments/prod/%" and key = "DB_HOST"`
- `drift where severity in ["critical","high"] and status = "open"`

### Graph Queries

- `where_used schema:payment:v7`
- `depends_on bundle:payments-prod:v22`
- `impact bundle:payments-prod:v22 -> scope_prefix("platform/payments/prod/*")`

## CLI

- `konfigo query "bundle where scope = 'platform/payments/prod/api'"`
- `konfigo where-used schema:payment:v7`
- `konfigo impact --bundle bundle:payments-prod:v22`

## API

- `POST /v1/query/execute`
- `POST /v1/query/where-used`
- `POST /v1/query/impact`

### Execute Query Example

`POST /v1/query/execute`

```json
{
  "query": "bundle where scope like 'platform/payments/prod/%' and pinned = true",
  "limit": 200,
  "cursor": null
}
```

## Language Design

### Syntax Rules (Phase 1)

- Case-insensitive keywords.
- String literals in single quotes.
- Comparison operators: `=`, `!=`, `<`, `>`, `<=`, `>=`.
- Boolean operators: `and`, `or`, `not`.
- List operator: `in [ ... ]`.
- Pattern operator: `like` with `%` wildcard.

### Reserved Functions

- `scope_prefix('team/app/%')`
- `age_lt('24h')`
- `ref_eq('schema:payment:v7')`

## Indexing Strategy

Primary indexes:

- `scope`
- `objectType`
- `version`
- `createdAt`
- `status`

Secondary inverted index:

- key path tokens
- text fields for name/description

Graph index:

- adjacency lists for references between object versions.

## Result Shape

Every result row contains:

- `entityType`
- `ref`
- `scope`
- `version`
- `summary`
- `links` (optional related refs)

## Security and Access Control

- Query engine applies scope-prefix authorization filter before returning any row.
- Unauthorized rows are excluded, not masked.
- Audit log stores original query and result count.

## Performance Targets

- P95 filter query latency < 300 ms for typical scope-filtered queries.
- P95 graph traversal latency < 800 ms for one-hop `where-used`.
- Support pagination via cursor for large result sets.

## Rollout Plan

### Phase 1

- Filter mode on bundle/values/drift/event entities.
- Basic `where-used` on schema/template refs.

### Phase 2

- Impact analysis and transitive dependency traversal.
- Saved queries.

### Phase 3

- Query alerts (run query on schedule and emit event on match).

## Testing Strategy

- Parser and validator unit tests.
- Query planner benchmark tests.
- Authorization leak tests (cross-scope data isolation).
- Snapshot tests for deterministic result ordering.

## Open Decisions

- Whether to expose the same syntax to UI advanced search as-is.
- Whether to allow user-defined macros in phase 2.
