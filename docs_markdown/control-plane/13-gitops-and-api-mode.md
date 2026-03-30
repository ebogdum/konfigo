# 13: GitOps and Direct API Mode

## Overview

Konfigo supports two delivery modes that write into the same immutable control plane model:

- **GitOps mode**: desired state is declared in Git and reconciled.
- **API mode**: desired state is submitted directly via API/CLI.

Both modes must converge to identical internal objects (`Schema@version`, `Template@version`, `Values@version`, `Bundle@version`) and identical event/audit behavior.

## Goals

- Support teams that require PR-based workflows.
- Support automation that requires low-latency API writes.
- Avoid split-brain by using one shared validation and publication pipeline.

## Non-Goals

- Per-mode custom policy behavior.
- Mutable in-place editing of published objects.

## Source of Truth Model

- **GitOps mode** source of truth: tracked Git repo path(s).
- **API mode** source of truth: control plane object store.
- Effective runtime source of truth: published immutable bundle bound to scope.

## Shared Intent Pipeline

All writes pass through the same pipeline:

1. Parse and normalize submission.
2. Validate schema/template/value compatibility.
3. Enforce policy and access checks.
4. Create immutable versions.
5. Materialize `Bundle@version`.
6. Emit publish and audit events.

## GitOps Mode

### Repository Layout (Example)

```text
konfigo-state/
  scopes/
    platform/payments/dev/api/
      schema.yaml
      template.yaml
      values.yaml
      bundle.yaml
    platform/payments/prod/api/
      bundle.yaml
```

### Reconcile Behavior

- Reconciler watches branch (`main` by default).
- Signed commit verification optional or required by policy.
- Each accepted commit generates a reconciliation run and emits events.

### Drift with Git

- If control plane state diverges from declared repo revision, create `gitops.drift.detected` event.
- Optionally block further promotions until reconciled.

## API Mode

### API Endpoints

- `POST /v1/schemas`
- `POST /v1/templates`
- `POST /v1/values`
- `POST /v1/bundles`
- `POST /v1/promotions`

### CLI Flow

- `konfigo push schema ...`
- `konfigo push template ...`
- `konfigo push values ... --scope ...`
- `konfigo bundle create ... --pin`
- `konfigo promote ...`

## Mode Interoperability Rules

- A scope can be configured as:
  - `gitops_only`
  - `api_only`
  - `hybrid` (API writes allowed but reconciler still authoritative on protected refs)
- For `hybrid`, last writer is not enough; policy decides override authority.

## Conflict Resolution

When both modes target same scope:

1. Check scope mode policy.
2. Check actor role and approval status.
3. Compare object refs and version lineage.
4. Accept or reject with explicit conflict reason.

## Approval and Promotion Model

- Protected scopes (for example `*/prod/*`) require approvals.
- Approval records bind to exact bundle digest.
- Ref movement (`prod/stable`) allowed only after policy gate pass.

## Security

- Humans authenticate via OIDC SSO.
- CI/CD authenticates via OIDC token exchange.
- Scopes enforce prefix-based RBAC/ABAC.
- Signed artifacts optional in phase 1, required for protected scopes in phase 2.

## Audit Requirements

For each apply/reconcile action store:

- actor identity
- auth context (OIDC claims summary)
- source mode (`gitops` or `api`)
- object refs and digests
- policy decisions
- resulting event IDs

## API Examples

### Create Bundle from API

`POST /v1/bundles`

```json
{
  "scope": "platform/payments/prod/api",
  "schemaRef": "schema:payment:v7",
  "templateRef": "template:payment:v12",
  "valuesRef": "values:payments-prod:v9",
  "pinned": true
}
```

### Register GitOps Source

`POST /v1/gitops/sources`

```json
{
  "scopePrefix": "platform/payments/*",
  "repo": "git@github.com:org/konfigo-state.git",
  "branch": "main",
  "path": "scopes/platform/payments"
}
```

## Operational Metrics

- `gitops_reconcile_duration_seconds`
- `gitops_reconcile_failures_total`
- `api_publish_latency_seconds`
- `mode_conflict_total`

## Rollout Plan

### Phase 1

- API mode with immutable objects and promotion.
- GitOps source registration with manual reconcile trigger.

### Phase 2

- Continuous reconcile loop and protected ref policies.
- Hybrid mode conflict enforcement.

### Phase 3

- Signed commit and signed bundle enforcement.

## Testing Strategy

- Unit tests for mode policy and conflict checks.
- Integration tests for Git webhook -> reconcile -> publish pipeline.
- End-to-end tests for `gitops_only` and `hybrid` scope behavior.

## Open Decisions

- Single repo vs multi-repo support in phase 1.
- Required depth of commit signature verification on day one.
