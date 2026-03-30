# 20: Disaster Recovery

## Overview

This spec defines recovery capabilities for Konfigo control plane data and availability:

- encrypted backups
- cross-region replication
- point-in-time restore

This version explicitly excludes scheduled game-day drill automation.

## Goals

- Protect all critical control plane state.
- Meet defined RPO and RTO objectives.
- Restore safely without violating immutability or audit requirements.

## Non-Goals

- Automated periodic drill orchestration.
- Multi-primary write topology in phase 1.

## Recovery Objectives

Target defaults (adjust per environment tier):

- **RPO**: <= 15 minutes
- **RTO**: <= 60 minutes

Protected scopes may have stricter targets.

## Data to Protect

- Immutable object store (`Schema`, `Template`, `Values`, `Bundle` versions)
- Reference mappings (promoted refs per scope)
- Policy definitions and approval records
- Event/outbox history
- Audit logs
- Integration metadata (webhooks, git sources, provider connections)

## Backup Strategy

### Full + Incremental

- Daily full snapshot.
- Frequent incremental backups (for example every 5 minutes).

### Encryption

- Encrypt in transit and at rest.
- Use cloud KMS (AWS KMS, GCP KMS, Azure Key Vault) or Vault transit.
- Rotate encryption keys with key-version metadata in backup manifest.

### Integrity

- Generate checksum manifest for every backup artifact.
- Sign manifests for tamper evidence.

## Replication Strategy

### Cross-Region Replication

- Primary region handles writes.
- Secondary region receives replicated backup stream and optional warm standby data.
- Replication lag exposed as metric and alert.

### Failover Modes

- Cold restore (backup only).
- Warm standby (replicated datastore and object store metadata).

## Restore Workflows

### 1. Point-in-Time Restore

- Select target timestamp.
- Restore metadata and object store references.
- Verify bundle digests and policy link integrity.
- Rebuild query indexes.

### 2. Scope-Targeted Restore

- Restore only specific `scope` prefix content.
- Preserve global event timeline with restore boundary markers.

### 3. Full Region Restore

- Restore all data into new environment.
- Rebind external integrations (webhooks, provider credentials).
- Promote recovered endpoint to active.

## Restore Safety Checks

Before finalize:

- object digest validation
- referential integrity checks
- policy and approval linkage checks
- event stream continuity checks

## API and CLI

### API

- `POST /v1/recovery/backups/create`
- `GET /v1/recovery/backups`
- `POST /v1/recovery/restores`
- `GET /v1/recovery/restores/{id}`

### CLI

- `konfigo recovery backup create`
- `konfigo recovery backup list`
- `konfigo recovery restore start --timestamp ...`
- `konfigo recovery restore status <restoreId>`

## Event Integration

Emit:

- `backup.completed`
- `backup.failed`
- `restore.started`
- `restore.completed`
- `restore.failed`

## Security

- Recovery actions require elevated role and explicit reason.
- Separate break-glass policy for production restore.
- Backup data access is tightly scoped and audited.

## Observability

Key metrics:

- `backup_duration_seconds`
- `backup_last_success_timestamp`
- `replication_lag_seconds`
- `restore_duration_seconds`
- `restore_failures_total`

Alerts:

- missed backup window
- replication lag threshold breach
- failed restore operation

## Rollout Plan

### Phase 1

- Full/incremental encrypted backups.
- Manual restore execution via API and CLI.

### Phase 2

- Cross-region replication and warm standby option.
- Scope-targeted restore.

### Phase 3

- Harden restore validations and automated post-restore health verification.

## Testing Strategy

- Unit tests for backup manifest integrity checks.
- Integration tests for backup and restore pipelines.
- Scale tests for large object/version counts.
- Failure injection tests for partial restore and lag scenarios.

## Open Decisions

- Backup retention tiers by environment class.
- Warm standby cost constraints vs desired RTO for production.
