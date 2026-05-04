# Evidence Ledger

Audit log architecture, format, and export for AGM.

## Layout

The Evidence Ledger plane has three components:

1. **Append-Only Event Log (Loki)**: collects all cluster events with a
   hash chain that makes tampering detectable.
2. **Local Audit DB (Postgres)**: stores structured records that need
   relational querying.
3. **Periodic Export Bundle Generator**: builds a tar archive of recent
   audit data and writes it to a one-way data diode export folder.

## Event Format

All AGM components write structured JSON events to `/var/log/cmba/events.log`:

```json
{
  "time": "2026-05-03T12:34:56Z",
  "reason": "CMBABindingDrift",
  "severity": "critical",
  "message": "model file hash does not match ModelBinding declaration",
  "expected": "5d41402abc4b2a76b9719d911017c592a37b4e5be2c6f4e0a9e1c4d4d9d4c8c1",
  "actual":   "deadbeef00000000000000000000000000000000000000000000000000000000",
  "path": "/models/iris.pkl",
  "podName": "iris-classifier-7d4c5b8f9-abcde",
  "namespace": "ml-serving",
  "modelBinding": "iris-classifier"
}
```

Loki ingests these events and adds a hash chain by computing
`SHA-256(prev_event_hash || current_event_json)` as a label on each entry.

## Audit DB Schema

Postgres tables (simplified):

```sql
CREATE TABLE bundle_imports (
    id BIGSERIAL PRIMARY KEY,
    bundle_sha256 CHAR(64) NOT NULL,
    bundle_size BIGINT NOT NULL,
    imported_at TIMESTAMP WITH TIME ZONE NOT NULL,
    imported_by TEXT NOT NULL,
    verification_result TEXT NOT NULL CHECK (verification_result IN ('pass', 'fail'))
);

CREATE TABLE binding_events (
    id BIGSERIAL PRIMARY KEY,
    namespace TEXT NOT NULL,
    binding_name TEXT NOT NULL,
    event_time TIMESTAMP WITH TIME ZONE NOT NULL,
    reason TEXT NOT NULL,
    severity TEXT NOT NULL,
    expected_sha CHAR(64),
    actual_sha CHAR(64),
    path TEXT
);
```

## Periodic Export

A CronJob runs every 24 hours and emits a bundle to the export folder:

```
audit-export-2026-05-03.tar
├── loki-events-2026-05-03.json   # raw event stream
├── audit-db-dump.sql              # Postgres dump
└── export-manifest.json           # SHA-256 of every file in the export
```

The export folder is mounted on a one-way data diode. An external auditor
on the connected side picks up the archive on a fixed schedule. Nothing
comes back through this path.

## Compliance Mapping

| Framework | Article | Evidence Ledger Output |
|---|---|---|
| EU AI Act 2024/1689 | Article 11 (Technical Documentation) | Bundle imports table; signed manifests |
| EU AI Act 2024/1689 | Article 17 (Quality Management System) | Hash-chained event log; export bundles |
| NIST AI RMF 1.0 | Measure function | Binding events table; benchmark CSVs |
| ISO/IEC 42001:2023 | A.6.2.5 (Documentation) | Periodic export bundles with manifest |

## Tamper Detection

Validation tooling reads the Loki event log forward and verifies the hash
chain at every entry:

```bash
agm-audit-validate /var/log/cmba/events.log
```

Returns non-zero if any chain break is detected. Operators can run this on
every export bundle before forwarding to an external auditor.
