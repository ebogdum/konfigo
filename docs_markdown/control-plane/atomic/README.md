# Konfigo Control Plane: Atomic Feature Specs

This folder contains implementation-ready, atomic feature specs.

Each file is intentionally narrow and follows the same structure:

1. Feature ID
2. Purpose
3. Inputs / Outputs
4. Public Interfaces (API/CLI/events)
5. Data Contracts
6. Invariants
7. Failure Modes
8. Acceptance Criteria (testable)
9. Out of Scope
10. Dependencies

## Atomic Spec List

### Core Foundation

- [CP-000-core-object-model.md](./CP-000-core-object-model.md)
- [CP-001-scope-prefix-model.md](./CP-001-scope-prefix-model.md)
- [CP-002-authn-authz.md](./CP-002-authn-authz.md)
- [CP-003-pinning-and-immutability.md](./CP-003-pinning-and-immutability.md)
- [CP-004-provenance-and-run-records.md](./CP-004-provenance-and-run-records.md)

### Remote Connectors and Cache

- [CP-005-provider-connector-matrix.md](./CP-005-provider-connector-matrix.md)
- [CP-006-cache-engine-semantics.md](./CP-006-cache-engine-semantics.md)
- [CP-007-remote-run-workflow.md](./CP-007-remote-run-workflow.md)
- [CP-008-offline-mode.md](./CP-008-offline-mode.md)

### Selected Roadmap Features

- [CP-009-drift-detection.md](./CP-009-drift-detection.md)
- [CP-010-events-and-webhooks.md](./CP-010-events-and-webhooks.md)
- [CP-011-gitops-and-api-mode.md](./CP-011-gitops-and-api-mode.md)
- [CP-012-terraform-provider.md](./CP-012-terraform-provider.md)
- [CP-013-kubernetes-operator.md](./CP-013-kubernetes-operator.md)
- [CP-014-query-language.md](./CP-014-query-language.md)
- [CP-015-disaster-recovery.md](./CP-015-disaster-recovery.md)
- [CP-016-governance-approvals-and-audit.md](./CP-016-governance-approvals-and-audit.md)
- [CP-017-requirement-coverage-map.md](./CP-017-requirement-coverage-map.md)

## How To Use For Development

- Treat each file as an independent work unit.
- Implement acceptance criteria exactly as written.
- Do not expand scope inside the same ticket.
- Add new atomic files instead of overloading existing ones.
