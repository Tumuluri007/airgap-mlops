# Policy Reference

Documentation for each of the 12 ML Kyverno ClusterPolicies in AGM.

## Air-Gap Group

### 01 — restrict-airgap-egress

- **Enforces**: every pod in `airgap.mlops/enforced=true` namespaces has a
  default-deny egress NetworkPolicy applied.
- **Blocks**: ML libraries that phone home (HuggingFace token checks,
  MLflow telemetry, Weights & Biases usage stats).
- **ML-specific**: standard K8s policies do not target ML namespaces with
  this strictness.
- **Air-gap operation**: pure admission-time check, no external lookups.

### 02 — require-local-registry-source

- **Enforces**: all container images come from `registry.airgap.local` or
  `mirror.airgap.local`.
- **Blocks**: pulls from docker.io, gcr.io, ghcr.io, quay.io, public.ecr.aws.
- **ML-specific**: serving images frequently inherit from public bases like
  `pytorch/pytorch:latest`.
- **Air-gap operation**: static allowlist match, no DNS lookups.

### 03 — block-internet-bound-init-containers

- **Enforces**: init container commands cannot include `curl`, `wget`,
  `pip install`, `conda install`, `npm install`, `apt-get update`.
- **Blocks**: data scientists adding `pip install scikit-learn==1.5` in
  init containers.
- **ML-specific**: very common ML anti-pattern.
- **Air-gap operation**: pattern match on PodSpec, no external dependencies.

### 04 — enforce-offline-trust-root

- **Enforces**: pods labeled `airgap.mlops/sig-verify=true` must mount the
  `airgap-trust-root` ConfigMap at `/etc/airgap/sigstore`.
- **Blocks**: pods configured for online Sigstore verification.
- **ML-specific**: model signing workflows default to online Rekor; this
  forces offline verification.

## Supply Chain Group

### 05 — verify-model-signature-offline

- **Enforces**: every container image in the `ml-serving` namespace has a
  valid Sigstore signature against the local trust root.
- **Blocks**: unsigned images, images signed against unknown identities.
- **ML-specific**: verifies both the container and any OCI-format model
  artifact mounted with it.
- **Air-gap operation**: `cosign verify --offline` against pre-staged bundles.

### 06 — require-mlbom-annotation

- **Enforces**: every InferenceService and Pod in `ml-serving` has the
  annotation `mlops.airgap/mlbom-digest=sha256:<...>`.
- **Blocks**: deployments with no ML-BOM reference.
- **ML-specific**: SBOM exists for software, ML-BOM extends to ML artifacts.

### 07 — require-slsa-provenance-attestation

- **Enforces**: container images carry an in-toto SLSA v1 provenance
  attestation that verifies offline.
- **Blocks**: images without provenance or with provenance that fails
  offline verification.
- **ML-specific**: distinguishes between model-training provenance (SLSA-ML
  extension) and standard build provenance.

### 08 — block-unsigned-model-images

- **Enforces**: ImagePullPolicy must be `IfNotPresent` for ml-serving images.
- **Blocks**: `ImagePullPolicy: Always` (would re-fetch and break air-gap).
- **ML-specific**: forces deployment-time verification, not pull-time.

## CMBA Group

### 09 — require-cmba-init-container

- **Enforces**: ml-serving pods include init container `cmba-verify` from
  the approved CMBA verifier image.
- **Blocks**: ML serving pods missing the init container.
- **ML-specific**: init container hashes the mounted model artifact and
  compares to the ModelBinding CRD before main starts.

### 10 — enforce-modelbinding-resource

- **Enforces**: every InferenceService has a paired ModelBinding CRD in
  the same namespace with the same name.
- **Blocks**: InferenceServices without a paired ModelBinding.
- **ML-specific**: this is the contract that links container identity to
  model identity.

### 11 — require-cmba-sentinel-sidecar

- **Enforces**: ml-serving pods include the cmba-sentinel sidecar.
- **Blocks**: pods that load models but do not continuously re-verify them.
- **ML-specific**: detects mid-flight model file replacement attacks.

## Governance Group

### 12 — require-ml-resource-limits

- **Enforces**: ml-serving pods declare GPU limits (`nvidia.com/gpu`) AND
  memory limits.
- **Blocks**: pods with unlimited or default resource requests.
- **ML-specific**: GPU resource specification is ML-specific; standard
  policies do not enforce GPU limits.

## Adding a New Policy

See [`CONTRIBUTING.md`](../CONTRIBUTING.md).
