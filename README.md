# Air-Gapped MLOps Platform (AGM)

**The first open reference implementation for ML model lifecycle management in
network-isolated Kubernetes environments.**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
[![Helm Chart](https://img.shields.io/badge/helm-v0.1.0-orange)](helm/airgap-mlops/Chart.yaml)
[![Status](https://img.shields.io/badge/status-alpha-yellow)]()

AGM provides three integrated capabilities for production machine learning in
air-gapped Kubernetes clusters:

1. **Air-Gap Reference Architecture** — a four-plane design (Ingress, Governance,
   Execution, Evidence Ledger) for fully offline ML operations.
2. **ML Kyverno Policy Library** — twelve admission policies enforcing supply
   chain integrity, model verification, and air-gap discipline.
3. **Container-Model Binding Attestation (CMBA)** — a Kubernetes-native
   mechanism (CRD + admission webhook + init container + sentinel sidecar) that
   cryptographically binds container images to model artifacts and continuously
   re-verifies the binding at runtime.

## Quick Start

Install the entire platform on a local kind or minikube cluster:

```bash
helm install airgap-mlops oci://ghcr.io/tumuluri007/charts/airgap-mlops \
  --version 0.1.0 \
  -f helm/airgap-mlops/values-airgap.yaml
```

The full end-to-end walkthrough is in [`examples/e2e-walkthrough/`](examples/e2e-walkthrough/).
Expect under fifteen minutes from `helm install` to a working ML serving cluster.

## Why Air-Gapped MLOps?

Defense systems run inside SCIFs. Nuclear facility decommissioning robots operate
on air-gapped operational technology networks under NRC 10 CFR 73.54. Healthcare
ML inference for hospitals subject to strict patient data residency laws sits
behind isolated network perimeters.

Each of these environments shares a property that public MLOps documentation
does not address: there is no internet, and there will not be one. AGM is the
first published reference architecture for this case.

## Components

| Layer | What It Does | Where to Look |
|---|---|---|
| Reference Architecture | Four-plane design with explicit air-gap boundary | [`ARCHITECTURE.md`](ARCHITECTURE.md) |
| Kyverno Policy Library | 12 admission policies for ML in air-gapped clusters | [`policies/`](policies/) |
| CMBA | Container-Model Binding Attestation system | [`cmba/`](cmba/) |
| Two-Pipeline CI/CD | External build + internal deploy pattern | [`ci-cd/`](ci-cd/) |
| Helm Chart | One-command install for the entire platform | [`helm/airgap-mlops/`](helm/airgap-mlops/) |
| Transfer Bundler | Build, sign, and verify air-gap transfer bundles | [`transfer-bundler/`](transfer-bundler/) |
| Documentation | Architecture, policy reference, troubleshooting | [`docs/`](docs/) |
| Reproducibility | Reproduce paper benchmarks and walkthroughs | [`paper/`](paper/), [`benchmarks/`](benchmarks/) |

## Citing

If you use AGM in research or production, please cite:

```bibtex
@article{tumuluri2026agm,
  title   = {Trustworthy Machine Learning Operations in Air-Gapped Kubernetes Clusters:
             An Integrated Architecture for Data Governance, Supply Chain Security,
             and Runtime Verification},
  author  = {Tumuluri, Yashasvi and others},
  year    = {2026},
  note    = {Open-source reference implementation: https://github.com/Tumuluri007/airgap-mlops}
}
```

A `CITATION.cff` file is included so GitHub displays a "Cite this repository" button.

## License

Apache 2.0. See [LICENSE](LICENSE).

## Maintainer

[Yashasvi Tumuluri](https://github.com/Tumuluri007)
