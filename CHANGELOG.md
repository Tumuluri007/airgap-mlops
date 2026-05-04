# Changelog

All notable changes to AGM are documented in this file. Format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-03

### Added

- AGM reference architecture with four-plane design (Ingress, Governance,
  Execution, Evidence Ledger).
- Twelve Kyverno ClusterPolicies for ML workloads in air-gapped Kubernetes:
  - Air-Gap group: 01-restrict-airgap-egress, 02-require-local-registry-source,
    03-block-internet-bound-init-containers, 04-enforce-offline-trust-root.
  - Supply Chain group: 05-verify-model-signature-offline, 06-require-mlbom-annotation,
    07-require-slsa-provenance-attestation, 08-block-unsigned-model-images.
  - CMBA group: 09-require-cmba-init-container, 10-enforce-modelbinding-resource,
    11-require-cmba-sentinel-sidecar.
  - Governance group: 12-require-ml-resource-limits.
- Container-Model Binding Attestation (CMBA):
  - ModelBinding CRD (apiVersion `cmba.airgap.mlops/v1alpha1`).
  - Admission webhook (Go) with six validation checks.
  - Init container (`cmba-verify`) with deterministic file hashing.
  - Runtime sentinel sidecar (`cmba-sentinel`) with configurable recheck interval.
- Two-pipeline CI/CD pattern: external bundle build (GitHub Actions) and
  internal deployment (Argo Workflows).
- Helm chart for one-command deployment of the entire platform on kind, minikube,
  or production Kubernetes.
- Transfer bundler scripts for building, signing, and verifying air-gap transfer
  bundles.
- End-to-end walkthrough example with four numbered shell scripts.
- Reproducibility benchmarks for policy evaluation overhead and CMBA verification
  latency.
- Documentation: architecture deep dive, transfer preparation, policy reference,
  CMBA internals, evidence ledger, troubleshooting.

[Unreleased]: https://github.com/Tumuluri007/airgap-mlops/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Tumuluri007/airgap-mlops/releases/tag/v0.1.0
