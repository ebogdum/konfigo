# 14: Terraform Provider and Kubernetes Operator

## Overview

Konfigo should integrate natively with infrastructure workflows:

- **Terraform provider** for declarative control plane object management.
- **Kubernetes operator** for bundle-driven runtime materialization.

Both integrations must use pinned immutable object references and emit standard control plane events.

## Goals

- Allow platform teams to manage Konfigo resources in Terraform.
- Allow Kubernetes workloads to consume pinned Konfigo bundles.
- Keep integrations deterministic and auditable.

## Non-Goals

- Full runtime traffic enforcement in operator.
- Complex custom resource templating in provider v1.

## Terraform Provider

### Provider Configuration

```hcl
provider "konfigo" {
  endpoint = "https://konfigo.company"
  auth_mode = "oidc_workload"
  scope     = "platform/payments/prod/api"
}
```

### Initial Resource Set

- `konfigo_schema`
- `konfigo_template`
- `konfigo_values`
- `konfigo_bundle`
- `konfigo_promotion`
- `konfigo_policy`

### Data Sources

- `konfigo_bundle`
- `konfigo_resolved_config`
- `konfigo_drift_findings`

### Idempotency Rules

- Published object versions are immutable.
- Update operations create new versions and update references.
- Resource IDs include stable logical name + concrete version.

### Terraform Example

```hcl
resource "konfigo_schema" "payment" {
  name    = "payment"
  version = "v7"
  content = file("schema.yaml")
}

resource "konfigo_template" "payment" {
  name    = "payment"
  version = "v12"
  content = file("template.yaml")
}

resource "konfigo_values" "prod" {
  scope   = "platform/payments/prod/api"
  version = "v9"
  content = file("values.prod.yaml")
}

resource "konfigo_bundle" "release" {
  scope        = "platform/payments/prod/api"
  schema_ref   = konfigo_schema.payment.ref
  template_ref = konfigo_template.payment.ref
  values_ref   = konfigo_values.prod.ref
  pinned       = true
}
```

## Kubernetes Operator

### CRDs

#### KonfigoBundleRef

Specifies desired control plane bundle and local projection target.

```yaml
apiVersion: konfigo.io/v1alpha1
kind: KonfigoBundleRef
metadata:
  name: payments-api-config
spec:
  scope: "platform/payments/prod/api"
  bundleRef: "bundle:payments-prod:v22"
  target:
    kind: ConfigMap
    name: payments-api-konfigo
```

#### KonfigoPolicyBinding (Optional v1beta)

Allows namespace/workload-local constraints over imported bundle keys.

### Operator Reconcile Loop

1. Read `KonfigoBundleRef`.
2. Fetch bundle metadata and resolved config.
3. Verify digest/signature if enforced.
4. Write/update target `ConfigMap`/`Secret`.
5. Update status with observed bundle version.
6. Emit reconcile events.

### Update Strategy

- Default: rolling update by changing mounted config and annotating deployment.
- Optional: strict mode requiring explicit approval annotation for `prod`.

### Failure Behavior

- If pull fails, preserve last known good target.
- Mark CR status as degraded with reason.
- Emit `operator.reconcile.failed` event.

## Security Model

- Operator uses workload identity and scope-limited permissions.
- Secret values can be written only to Kubernetes `Secret` targets.
- Terraform provider supports short-lived auth only.

## Drift and Integration

- Operator-reported observed version feeds drift detection.
- Terraform data sources expose active drift findings.

## API Requirements for Integrations

- `GET /v1/bundles/{ref}`
- `GET /v1/bundles/{ref}/resolved`
- `POST /v1/events` (internal or authenticated ingestion)
- `GET /v1/drift/findings`

## Metrics

### Provider

- `tf_konfigo_apply_duration_seconds`
- `tf_konfigo_publish_failures_total`

### Operator

- `konfigo_operator_reconcile_duration_seconds`
- `konfigo_operator_reconcile_failures_total`
- `konfigo_operator_last_success_timestamp`

## Rollout Plan

### Phase 1

- Terraform: schema/template/values/bundle resources.
- Operator: `KonfigoBundleRef` -> ConfigMap target only.

### Phase 2

- Terraform promotion and policy resources.
- Operator Secret target support and signature verification.

### Phase 3

- Advanced policy binding and promotion awareness in operator.

## Testing Strategy

- Contract tests between provider and control plane API.
- Acceptance tests using Terraform plugin framework.
- Operator e2e in KinD with rollout and failure scenarios.

## Open Decisions

- Whether to support multiple bundle refs merged into one target in phase 1.
- Whether operator should trigger deployment rollouts directly or by annotation only.
