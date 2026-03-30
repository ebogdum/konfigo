# 15: Events and Webhooks

## Overview

Konfigo requires a canonical event stream for auditability, automation, and integration.

This spec defines:

- internal event model
- event durability and ordering
- webhook delivery contract
- security and retry behavior

## Goals

- Emit every meaningful state transition as a structured event.
- Deliver events reliably to downstream systems.
- Provide verifiable webhook signatures.

## Non-Goals

- Full external message bus management in phase 1.
- Exactly-once semantics for external webhooks.

## Core Event Types

### Publish and Promotion

- `schema.published`
- `template.published`
- `values.published`
- `bundle.published`
- `bundle.promoted`
- `bundle.rollback.completed`

### Drift

- `drift.scan.started`
- `drift.scan.completed`
- `drift.finding.opened`
- `drift.finding.resolved`

### GitOps and API

- `gitops.reconcile.started`
- `gitops.reconcile.completed`
- `api.write.accepted`
- `api.write.rejected`

### Recovery

- `backup.completed`
- `restore.started`
- `restore.completed`

## Event Envelope

```json
{
  "eventId": "evt_01J...",
  "type": "bundle.published",
  "occurredAt": "2026-02-25T00:45:00Z",
  "scope": "platform/payments/prod/api",
  "actor": {
    "type": "service",
    "id": "ci/github-actions/payments-release"
  },
  "source": "control-plane",
  "version": "v1",
  "payload": {
    "bundleRef": "bundle:payments-prod:v22",
    "digest": "sha256:..."
  }
}
```

## Delivery Architecture

Use transactional outbox pattern:

1. Write domain state change and outbox row in one transaction.
2. Dispatcher publishes to internal bus and webhook queue.
3. Delivery workers handle retries and dead-letter routing.

## Ordering Guarantees

- Global ordering is not guaranteed.
- Ordering is guaranteed per `(scope, entity)` partition key.
- Consumers must treat event processing as idempotent.

## Webhook Model

### Endpoint Registration

`POST /v1/webhooks/endpoints`

```json
{
  "name": "slack-prod-alerts",
  "url": "https://hooks.example.com/konfigo",
  "eventTypes": ["drift.finding.opened", "bundle.promoted"],
  "scopePrefix": "platform/payments/prod/*"
}
```

### Signing

- Header: `X-Konfigo-Signature: v1=<hmac_sha256>`
- Header: `X-Konfigo-Timestamp: <unix_epoch_seconds>`
- Signature input: `<timestamp>.<raw_body>`

### Retry Policy

- Exponential backoff with jitter.
- Max attempts configurable (default 10).
- Final failure routes to dead-letter queue.

### Idempotency

- Include `eventId` header and in body.
- Receivers dedupe using `eventId`.

## Security

- TLS required.
- Optional mTLS for private destinations.
- Secret rotation support for webhook signing keys.
- Endpoint test and verification challenge flow.

## Observability

Required metrics:

- `events_emitted_total{type}`
- `webhook_deliveries_total{endpoint,status}`
- `webhook_delivery_latency_seconds`
- `webhook_retries_total`
- `webhook_dead_letter_total`

Required traces/log fields:

- `eventId`, `type`, `scope`, `endpointId`, `attempt`.

## API Endpoints

- `POST /v1/webhooks/endpoints`
- `GET /v1/webhooks/endpoints`
- `PATCH /v1/webhooks/endpoints/{id}`
- `DELETE /v1/webhooks/endpoints/{id}`
- `POST /v1/webhooks/endpoints/{id}/test`
- `GET /v1/events?scope=...&type=...`

## Retention

- Hot event query store: 30 to 90 days.
- Long-term archive (object storage): 1 year or policy-defined.

## Rollout Plan

### Phase 1

- Core event emission and internal query API.
- Basic webhook delivery with retries.

### Phase 2

- Dead-letter queue visibility and replay endpoint.
- Partition-aware ordering guarantees and throughput tuning.

### Phase 3

- Optional connectors for external bus sinks (SNS/EventBridge/PubSub/Event Grid).

## Testing Strategy

- Unit tests for signature generation/validation.
- Integration tests for retry and dead-letter paths.
- Load tests for burst event throughput.
- Chaos tests on webhook endpoint outages.

## Open Decisions

- Default event retention duration.
- Max webhook payload size and truncation behavior.
