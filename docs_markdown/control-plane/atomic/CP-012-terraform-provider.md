# CP-012: Terraform Provider

## Purpose

Provide Infrastructure-as-Code management for Konfigo control-plane objects.

## Inputs / Outputs

- Input: Terraform plan/apply operations.
- Output: converged Konfigo object state and references.

## Public Interfaces

Terraform resources:

- `konfigo_schema`
- `konfigo_template`
- `konfigo_values`
- `konfigo_bundle`
- `konfigo_policy`
- `konfigo_promotion`

## Data Contracts

Resource IDs must include logical identity and immutable version reference.

## Invariants

- Provider never mutates immutable versions in place.
- Update creates new version and moves reference only if allowed.

## Failure Modes

- `409 immutable_update_attempt`
- `403 promotion_not_approved`

## Acceptance Criteria

1. `terraform apply` is idempotent with unchanged inputs.
2. Diff output clearly indicates ref movement vs new immutable version.
3. Provider surfaces policy-denial reason from API response.

## Out of Scope

- Terraform-managed connector secret material.

## Dependencies

- CP-000 core model.
