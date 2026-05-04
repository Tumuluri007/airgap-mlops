# Architecture Deep Dive

This document expands the high-level architecture summary in `ARCHITECTURE.md`
with detailed explanations of each plane, the trust model, and failure modes.

## Why Four Planes

A pipeline diagram cannot capture work that happens continuously inside the
cluster (admission, runtime verification, audit). The four-plane structure
captures these concerns:

- **Ingress**: one-time work per artifact arrival.
- **Governance**: continuous policy enforcement at every admission.
- **Execution**: long-lived ML serving with continuous CMBA verification.
- **Evidence Ledger**: continuous append-only audit collection.

## Ingress Plane Details

### Manifest Verifier

Validates every file in the transfer bundle against the SHA-256 recorded in
`manifest.json`. Implemented as a one-shot Job that runs when a new bundle
appears in the drop-folder PVC.

### Local Container Registry

Either Harbor or Zot. Both are CNCF-graduated/sandbox projects with proven
offline operation. Harbor has more enterprise features (replication, content
trust, vulnerability scanning); Zot is lighter weight and OCI-native.

### Local Sigstore Stack

The local stack mirrors the public Sigstore service:
- **Rekor**: append-only transparency log, runs as a single deployment.
- **Fulcio**: certificate authority for keyless signing.
- **TUF root**: pre-staged at cluster provisioning time. Rotation requires
  generating new trust material on the connected side and shipping it through
  a transfer bundle.

## Governance Plane Details

### Kyverno in Offline Mode

Kyverno is configured with `--allowInsecureRegistry=false` (production) and
`--imageVerifyCacheEnabled=true`. The image-verify cache is pre-populated with
the signatures of every image in the local registry to avoid synchronous
Sigstore calls on every pod admission.

### ML Policy Library

See [`04-policy-reference.md`](04-policy-reference.md) for per-policy
documentation. The library is grouped into:
- Air-Gap (01-04): enforce that the air gap is observed.
- Supply Chain (05-08): enforce model and image provenance.
- CMBA (09-11): enforce the CMBA component contract.
- Governance (12): enforce ML resource governance.

## Execution Plane Details

### Model Registry

MLflow with embedded Postgres on a PVC. The artifact store points at a
PVC-backed local filesystem path; no S3 or GCS dependencies. Model promotion
between stages (Staging → Production) uses MLflow's stage-transition API
with Kyverno policy gating.

### KServe InferenceService

KServe controls the runtime, autoscaling, and traffic splitting. AGM does
not modify KServe; it adds the CMBA contract layer on top.

### Pod Internals

Every ml-serving pod contains exactly three containers:
1. `cmba-verify` (init): hashes the mounted model file before main starts.
2. Main model server: the actual inference container.
3. `cmba-sentinel` (sidecar): re-hashes the model file every recheck interval.

The init container shares a volume mount with the main container; both see
the same model file. The sentinel sidecar shares the same mount.

## Evidence Ledger Plane Details

### Append-Only Log

Loki collects events from every cluster component. A hash chain over the
last N events makes tampering detectable: each event records the SHA-256 of
the previous event entry. Validation tooling can scan the log forward and
flag any chain break.

### Audit DB

Postgres stores structured audit records that need relational querying
(e.g., "show me all binding mismatch events for ModelBinding X in the last
24 hours"). The PVC backing Postgres is WORM-mounted in production, with
log-structured append-only writes.

### Periodic Export

A cron job builds a tar archive of the last 24 hours of audit data and
writes it to a one-way data diode export folder. An external auditor on
the connected side picks up the archive. Nothing comes back through this
path.

## Trust Model Walkthrough

| Stage | Trusted Inputs | Verified | Failure Mode |
|---|---|---|---|
| Bundle arrival | Sigstore trust root | Bundle signature, file hashes | Bundle quarantined |
| Pod admission | Kyverno policies, ModelBinding CRD | All 12 policies, CMBA contract | Pod admission denied |
| Pod startup | ModelBinding spec | cmba-verify exit code | Pod startup blocked |
| Pod runtime | ModelBinding spec | cmba-sentinel re-hash every interval | Configured action (terminate, alert, crash) |
| Audit | Hash chain over events | Loki tail validation | Audit alert; export marked invalid |

## Threat Model

The threats AGM defends against:
- Tampered transfer bundles (bundle signature)
- Unsigned or untrusted container images (Policy 05, 07, 08)
- Pods that bypass CMBA (Policy 09, 10, 11)
- Model file replacement on shared volumes (cmba-sentinel)
- Audit log tampering (hash chain)
- Resource exhaustion attacks (Policy 12)

The threats AGM does NOT defend against:
- Cluster admin compromise (root in the cluster control plane)
- Hardware-level attacks (firmware, side-channel)
- Adversarial inputs to the deployed model (out of scope, requires runtime
  adversarial robustness testing)
