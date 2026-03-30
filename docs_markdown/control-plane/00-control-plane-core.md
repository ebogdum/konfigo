# 00: Control Plane Core

## Overview

This document defines the core model and behavior of the Konfigo control plane.

It turns Konfigo from a local-only processor into a centralized system that can:

- store and version schemas, templates, values, files, and policies
- resolve and execute bundle-based runs remotely
- enforce immutable pinning and protected promotions
- provide provenance and auditability

## Goals

- Centralize configuration assets and execution intent.
- Keep object versions immutable and reproducible.
- Support controlled promotion between environments.
- Make all decisions auditable.

## Non-Goals

- Runtime traffic filtering or API gateway policy enforcement.
- Dynamic policy language redesign in this phase.

## Canonical Scope Model

Use a logical `scope` prefix format:

- `team/app/env/service`
- Example: `platform/payments/prod/api`

`scope` is used consistently for:

- storage paths
- access control checks
- event partitioning
- drift scans
- promotion policies

## Core Object Model

### Project

Represents top-level ownership and lifecycle boundary.

Fields:

- `projectId`
- `name`
- `description`
- `ownerGroup`
- `createdAt`

### Environment

Represents deployment stage constraints (`dev`, `stage`, `prod`, custom).

Fields:

- `environmentId`
- `projectId`
- `name`
- `isProtected`
- `promotionPolicyRef`

### Schema@version

Immutable schema specification.

Fields:

- `schemaRef` (`schema:<name>:<version>`)
- `contentDigest`
- `content`
- `createdBy`
- `createdAt`

### Template@version

Immutable template content.

Fields:

- `templateRef` (`template:<name>:<version>`)
- `contentDigest`
- `content`
- `createdBy`
- `createdAt`

### Values@version

Immutable values set bound to scope.

Fields:

- `valuesRef` (`values:<scope>:<version>`)
- `scope`
- `contentDigest`
- `content`
- `createdBy`
- `createdAt`

### FileAsset@version

Optional immutable file blobs used by bundles.

Fields:

- `fileRef`
- `path`
- `contentDigest`
- `contentLocation`

### Policy@version

Immutable policy package.

Fields:

- `policyRef`
- `rules`
- `enforcementMode`

### Bundle@version

Deployable immutable unit.

Fields:

- `bundleRef` (`bundle:<scope>:<version>`)
- `scope`
- `schemaRef`
- `templateRef`
- `valuesRef`
- `fileRefs[]`
- `policyRefs[]`
- `lockedBindings`
- `digest`
- `signature` (optional or required by scope policy)
- `createdAt`, `createdBy`

### RunRecord

Execution record linked to bundle digest.

Fields:

- `runId`
- `scope`
- `bundleRef`
- `bundleDigest`
- `mode` (`local`, `remote`, `offline`)
- `actor`
- `startedAt`, `completedAt`
- `result` (`success`, `failure`, `partial`)
- `outputArtifacts[]`

## Lifecycle and State Transitions

1. Publish immutable objects (`Schema`, `Template`, `Values`, optional `FileAsset`, `Policy`).
2. Create pinned immutable `Bundle` from exact references.
3. Promote bundle reference across environment refs (`dev/stable`, `prod/stable`).
4. Execute runs against pinned bundle.
5. Track drift and emit findings.

## API Surface (Core)

- `POST /v1/projects`
- `POST /v1/environments`
- `POST /v1/schemas`
- `POST /v1/templates`
- `POST /v1/values`
- `POST /v1/files`
- `POST /v1/policies`
- `POST /v1/bundles`
- `POST /v1/promotions`
- `POST /v1/runs`
- `GET /v1/runs/{runId}`

## Required Invariants

- Published versions are immutable.
- `Bundle` references immutable versions only.
- Protected references cannot move without policy approval.
- Every run references exact bundle digest.

## Promotion Model

Reference aliases are mutable pointers with guardrails.

Examples:

- `platform/payments/dev/stable -> bundle:...:v21`
- `platform/payments/prod/stable -> bundle:...:v20`

Promotion updates pointer only after:

- approval checks
- policy checks
- compatibility checks

## Configuration Resolution Order

For bundle execution:

1. Bundle pinned references
2. Runtime overrides permitted by policy
3. External provider values (if configured)
4. Fallback defaults

Any override must be reflected in provenance output.

## CLI (Control Plane Core)

- `konfigo login`
- `konfigo project create ...`
- `konfigo env create ...`
- `konfigo push schema ...`
- `konfigo push template ...`
- `konfigo push values ...`
- `konfigo push file ...`
- `konfigo push policy ...`
- `konfigo bundle create --pin ...`
- `konfigo promote ...`
- `konfigo run --server ... --bundle ...`

## Security Baseline

- OIDC-based human auth
- workload identity for CI/CD
- short-lived scoped tokens
- scope-prefix RBAC/ABAC
- full audit logging

## Data Retention

- immutable objects: retained by policy (default long-term)
- run records: hot 90 days, archive configurable
- promotion history: long-term

## Observability

- object publish counts and latency
- promotion success/failure counts
- run success/failure by scope
- policy denial counts

## Rollout Plan

### Phase 1

- Core immutable model and bundle creation
- Run record persistence

### Phase 2

- Promotion refs with approval gates
- file assets and policy bundles

### Phase 3

- advanced compatibility checks and signed bundle enforcement by policy

