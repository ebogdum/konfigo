# CP-010: Events and Webhooks

## Purpose

Emit and deliver reliable control-plane events for automation and audit.

## Inputs / Outputs

- Input: internal state transitions.
- Output: event stream + webhook deliveries.

## Public Interfaces

- `GET /v1/events`
- `POST /v1/webhooks/endpoints`
- `POST /v1/webhooks/endpoints/{id}/test`

## Data Contracts

Event envelope:

- `eventId`
- `type`
- `scope`
- `occurredAt`
- `actor`
- `payload`

Webhook signature headers:

- `X-Konfigo-Timestamp`
- `X-Konfigo-Signature`

## Invariants

- Event IDs are globally unique.
- Webhooks are retried with exponential backoff.
- Dead-letter queue stores exhausted deliveries.

## Failure Modes

- `503 webhook_delivery_failed`

## Acceptance Criteria

1. Event emitted for publish/promote/drift/restore paths.
2. Webhook signatures verify with shared secret.
3. Failed deliveries reach dead-letter queue after max attempts.

## Out of Scope

- Exactly-once external delivery guarantee.

## Dependencies

- CP-000 core model.
