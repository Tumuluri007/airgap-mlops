# CMBA Internals

How the Container-Model Binding Attestation system works under the hood.

## The Problem

In production Kubernetes, models are usually mounted into serving containers
from a PersistentVolume, an init container that pulls from object storage,
or a ConfigMap. The container image and the model artifact are not the same
object and they get updated through different paths.

A model file replaced on a shared volume after pod startup is not detected
by HTTP health checks. CMBA exists to close this gap.

## The ModelBinding CRD

CMBA introduces a new Kubernetes resource type:
`apiVersion: cmba.airgap.mlops/v1alpha1`, `kind: ModelBinding`.

The CRD declares a contract: this container image goes with this model
artifact, and here is the cryptographic proof. See
[`cmba/crd/modelbinding-crd.yaml`](../cmba/crd/modelbinding-crd.yaml) for
the full schema.

## The Three Controllers

### cmba-webhook (admission)

A Kubernetes ValidatingAdmissionWebhook intercepts pod creation in the
`ml-serving` namespace and performs six checks before admitting the pod:

1. Pod carries annotation `cmba.airgap.mlops/binding-name`.
2. The named ModelBinding exists in the pod's namespace.
3. Pod's container image matches `ModelBinding.spec.containerImage`.
4. Init container `cmba-verify` is present.
5. Sidecar `cmba-sentinel` is present (if `verificationPolicy.sentinelRequired`).
6. `ModelBinding.attestation.signatureBundle` verifies offline against the
   local trust root.

Any failed check rejects the pod with a clear reason.

### cmba-verify (init container)

Runs once at pod startup. Reads the mounted ModelBinding from a projected
volume, hashes the actual model file at the declared path, and exits
non-zero if the hashes do not match. Pod startup is blocked.

The init container is pure shell (busybox + yq). Simple primitives are
auditable; the entire script is short by design.

### cmba-sentinel (sidecar)

Runs alongside the model server. Every `recheckIntervalSeconds` (default
300) it re-hashes the mounted model file and compares to the
ModelBinding's expected hash.

On mismatch:
- Increments `ModelBinding.status.mismatchEvents`.
- Writes a structured event to `/var/log/cmba/events.log` (Loki picks it up).
- Sets a Kubernetes Event with reason `CMBABindingDrift`.
- Takes the configured action: `terminate` (SIGTERM main container),
  `alert` (event only), or `crash` (force pod restart).

## Air-Gap Behavior

- Webhook makes no external HTTPS calls. Trust root mounted as ConfigMap.
- Init container has no network requirements at all.
- Sentinel makes no external webhook calls (no Slack, no PagerDuty).

## Performance

- Admission webhook: 38 ms p99 on a three-node test cluster.
- Init container: under 1.5 s for models up to 5 GB.
- Sentinel: 3.8 MB resident memory, ~0.1% of one core averaged over 24 h.

## Defense in Depth

CMBA is one of three layers:
- **Build-time**: container image and model artifact are signed.
- **Admission-time**: ValidatingAdmissionWebhook enforces ModelBinding contract.
- **Runtime**: sentinel re-verifies at configurable intervals.

Bypassing CMBA requires compromising all three layers, plus the cluster's
admission controller chain. Defense-in-depth recommendations:
- Kyverno Policy 09 ensures cmba-verify is present.
- Kyverno Policy 11 ensures cmba-sentinel is present.
- ValidatingWebhookConfiguration uses `failurePolicy: Fail` (closed).
- Cluster admin RBAC restricted; webhook deletion requires a separate role.
