# AGM Architecture

This document explains the AGM reference architecture in depth: the four-plane
structure inside the air-gapped cluster, the Transfer Preparation Zone on the
connected side, the trust model, the two-pipeline CI/CD pattern, and the failure
modes each plane handles.

## 1. The Air-Gap Boundary

AGM splits cleanly across an explicit air-gap boundary.

- **Connected side**: source repositories, public model hubs, and the public
  Sigstore transparency log feed into a **Transfer Preparation Zone**. The zone
  produces a single signed tar archive containing the container image, the model
  artifact, the CycloneDX ML-BOM, the Sigstore trust roots, and the Kyverno
  policy bundle.
- **Air-gap boundary**: physical transfer (removable media) or a one-way data
  diode. No bidirectional network.
- **Air-gapped side**: a self-contained Kubernetes cluster with no egress and
  no ingress beyond the controlled transfer path.

The signed transfer bundle is the only artifact that crosses the boundary. It
is a cryptographic abstraction that other organizations can adopt without
copying our specific tool choices.

## 2. The Four Planes

Inside the air-gapped cluster, AGM is organized as four horizontal planes.

### 2.1 Ingress Plane

The Ingress plane lands the transfer bundle and verifies it.

- **Manifest Verifier** — checks the bundle signature and the integrity of every
  artifact it contains.
- **Local Container Registry** — Harbor or Zot, both supporting offline operation.
- **Local Sigstore Stack** — local Rekor transparency log, local Fulcio CA, and
  a pre-staged TUF root. Pre-staging the trust roots before the cluster goes
  air-gapped is the key insight: after pre-staging, signature verification works
  offline indefinitely.

### 2.2 Governance Plane

The Governance plane runs admission control with no internet calls.

- **Kyverno Controller** — runs in offline mode with the image-verify cache
  pre-populated.
- **Sigstore Policy Controller** — configured for offline verification.
- **ML Policy Library** — 12 ClusterPolicies grouped into Air-Gap, Supply Chain,
  CMBA, and Governance categories. See [`policies/`](policies/) for full policies.

### 2.3 Execution Plane

The Execution plane runs the actual ML serving.

- **Local Model Registry** — MLflow with embedded Postgres on a PVC.
- **KServe Model Controller** — manages InferenceServices.
- **ML Serving Pod** — three containers per pod: a `cmba-verify` init container
  that hashes the mounted model file before startup, the model server itself,
  and a `cmba-sentinel` sidecar that re-hashes the model at a configurable
  interval.
- **ModelBinding CRD** — declares the cryptographic binding between container
  image and model artifact. See [`cmba/crd/modelbinding-crd.yaml`](cmba/crd/modelbinding-crd.yaml).

### 2.4 Evidence Ledger Plane

The Evidence Ledger plane keeps an append-only audit record.

- **Append-Only Event Log** — Loki with a hash chain so tampering is detectable.
- **Local Audit DB** — Postgres on a Write-Once-Read-Many (WORM) backed PVC.
- **Periodic Export Bundle Generator** — builds a tar archive that can leave
  the cluster through a one-way data diode for an external auditor. Nothing
  returns through that path.

The Evidence Ledger generates the documentation needed to demonstrate compliance
with EU AI Act Article 11 (technical documentation) and Article 17 (quality
management system).

## 3. The Two-Pipeline CI/CD Pattern

Standard CI/CD assumes the cluster can pull from a remote registry. In an
air-gapped environment that pull is impossible. AGM splits CI/CD across the
boundary:

### 3.1 External Pipeline

Runs on the connected internet (GitHub Actions, GitLab CI, Tekton). Builds the
container, signs it with Cosign, generates the ML-BOM with cdxgen, generates
the ModelBinding manifest, signs the manifest, and bundles everything into a
signed tar. See [`.github/workflows/airgap-bundle-build.yml`](.github/workflows/airgap-bundle-build.yml).

### 3.2 Internal Pipeline

Runs inside the air-gapped cluster (Argo Workflows). Watches a drop-folder PVC
for the bundle, verifies the signature using the local Sigstore stack, loads
the container into the local registry, applies the ModelBinding CRD, and
records the import event in the Evidence Ledger. See
[`ci-cd/argo-workflows/airgap-deploy.yaml`](ci-cd/argo-workflows/airgap-deploy.yaml).

## 4. Trust Model

| Layer | Trusted | Verified | Untrusted |
|---|---|---|---|
| Build pipeline | Source repository (signed commits) | Container build | External dependencies |
| Transfer bundle | Bundle signature (Cosign) | Manifest hashes | Bundle origin |
| Internal registry | Local registry credentials | Image signatures | Image origin (until verified) |
| Pod admission | Kyverno controller | All policies (12) | Pod spec (until policies pass) |
| Pod runtime | cmba-verify init container | Mounted model file hash | Model file (until verified) |
| Continuous runtime | cmba-sentinel sidecar | Model file hash every recheck interval | Model file (between checks) |

## 5. Failure Modes

| Plane | Failure | Detection | Response |
|---|---|---|---|
| Ingress | Bundle signature fails | Manifest Verifier rejects bundle | Admin notified; bundle quarantined |
| Ingress | Local Sigstore unreachable | Verification timeout | Pod admission blocked by Policy 04 |
| Governance | Kyverno controller down | Webhook timeout | Admission denied (fail-closed) |
| Governance | Policy bundle missing | Policy not loaded | Pods that need policy enforcement blocked |
| Execution | Model file mismatch on startup | cmba-verify exits non-zero | Pod startup blocked |
| Execution | Model file mismatch at runtime | cmba-sentinel detects within recheck interval | Configured action: terminate, alert, or crash |
| Execution | Pod resource exhaustion | Policy 12 admission check | Pod admission denied if no GPU/memory limits |
| Evidence Ledger | Loki append-only chain breaks | Hash chain mismatch | Audit alert; export bundle marked invalid |

## 6. Comparison with Connected MLOps

| Concern | Connected MLOps | AGM (Air-Gapped) |
|---|---|---|
| Container image source | Public registry (docker.io, ghcr.io) | Local registry (registry.airgap.local) only |
| Model registry | Cloud-hosted (Vertex AI, SageMaker) | Local MLflow on PVC |
| Signature verification | Online Sigstore (Rekor public log) | Local Sigstore stack with pre-staged TUF root |
| Policy enforcement | Optional, often added late | Mandatory at admission, 12 ML-specific policies |
| Audit trail | Cloud-managed logs | Local Loki + Postgres on WORM PVC |
| CI/CD | Single pipeline pushes to cluster | Two-pipeline with signed transfer bundle |
| Observability | Cloud monitoring (Datadog, etc.) | Local Prometheus + Grafana + Loki |
| Cost model | Pay-per-use cloud services | Self-hosted on owned hardware |

## 7. Compliance Framework Alignment

| Framework | Article / Section | Where AGM Contributes |
|---|---|---|
| EU AI Act 2024/1689 | Article 9 (Risk Management) | CMBA threat model; defensive scenarios |
| EU AI Act 2024/1689 | Article 11 (Technical Documentation) | Evidence Ledger plane |
| EU AI Act 2024/1689 | Article 15 (Cybersecurity) | Kyverno policy library; CMBA |
| EU AI Act 2024/1689 | Article 17 (Quality Management) | Audit trail; reproducible benchmarks |
| NIST AI RMF 1.0 | Govern, Map, Measure, Manage | Policy library, threat model, benchmarks, runtime verification |
| ISO/IEC 42001:2023 | AI Management System | Four-plane operational structure |
| MITRE ATLAS v5.4.0 | AML.T0010 (Supply Chain Compromise) | CMBA defensive scenarios |
