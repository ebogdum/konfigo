# Control Plane Feature Specs

This folder contains implementation-level specifications for the Konfigo control plane roadmap items selected for the next major evolution:

- 4. Drift detection
- 13. GitOps and direct API mode
- 14. Terraform provider and Kubernetes operator
- 15. Events and webhooks
- 19. Query and search language
- 20. Disaster recovery (without operational drill workflows)

## Why This Set

These features make Konfigo a full control plane instead of only a local config processor:

- **Operate centrally** from any environment.
- **Enforce safety** through immutable bundles, policy gates, and approvals.
- **Integrate deeply** with CI/CD, Infrastructure-as-Code, and Kubernetes.
- **Observe and recover** with first-class events, drift visibility, and backup/restore.

## Shared Architecture Decisions

All specs in this folder assume the same baseline.

### 1. Logical Organization by Scope Prefix

No hard multi-tenant namespace model is required for this phase.

- Use a canonical `scope` string like `team/app/env/service`.
- Apply this same `scope` in all objects, access checks, storage keys, and events.
- Store remote values by provider prefix, for example:
  - Vault: `kv/konfigo/<scope>/...`
  - S3: `s3://konfigo-config/<scope>/...`

### 2. Auth and Identity

Username/password is not the primary path.

- Humans: OIDC SSO (browser or device code).
- CI/CD: workload identity federation (OIDC token exchange).
- Services: short-lived scoped tokens.
- Authorization: RBAC/ABAC on scope prefixes.

### 3. Immutable Versioned Objects

Core objects are immutable once published:

- `Schema@version`
- `Template@version`
- `Values@version`
- `Bundle@version` (the deployable resolved unit)

Promotion flows move references, not mutable content.

### 4. Event-First Control Plane

Every important state transition emits events:

- publish, promote, drift detected/resolved, rollback, restore

This is required by GitOps reconcile loops, operator behavior, notifications, and audit.

## Spec Files

- [04-drift-detection.md](./04-drift-detection.md)
- [13-gitops-and-api-mode.md](./13-gitops-and-api-mode.md)
- [14-terraform-provider-and-k8s-operator.md](./14-terraform-provider-and-k8s-operator.md)
- [15-events-and-webhooks.md](./15-events-and-webhooks.md)
- [19-query-and-search-language.md](./19-query-and-search-language.md)
- [20-disaster-recovery.md](./20-disaster-recovery.md)

## Suggested Delivery Sequence

1. Events and webhooks (`15`)
2. Query and search language (`19`)
3. Drift detection (`4`)
4. GitOps and API dual mode (`13`)
5. Terraform provider and Kubernetes operator (`14`)
6. Disaster recovery (`20`)

This order keeps foundational dependencies in front and reduces rework.

## Atomic Development Specs

For implementation planning and ticketing, use the atomic set:

- [atomic/README.md](./atomic/README.md)

The atomic specs split each concern into contained, testable units with:

- strict purpose and boundaries
- explicit interface contracts
- invariants and failure modes
- concrete acceptance criteria
