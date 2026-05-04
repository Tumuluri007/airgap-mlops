# Reproducibility

This directory contains the reproducibility artifacts referenced in the
AGM paper. Every benchmark in the paper can be reproduced with the
commands below.

## Environment

- Operating system: Ubuntu 22.04 LTS
- Kubernetes: kind v0.23.0 with Kubernetes v1.30.0
- Kyverno: v1.12.5
- KServe: v0.13.0
- MLflow: v2.15.0
- Helm: v3.14.0
- Cosign: v2.4.0
- yq: v4.40.5
- CPU: Intel x86_64, 16 GB RAM per node
- Number of nodes: 3 (kind cluster)

## Setup

```bash
# Install kind, helm, kubectl, cosign, yq.
make setup

# Create a three-node kind cluster.
kind create cluster --config paper/kind-config.yaml --name agm-bench

# Install AGM.
helm install airgap-mlops ./helm/airgap-mlops -f ./helm/airgap-mlops/values.yaml
```

## Benchmarks

```bash
# Policy evaluation overhead (Section 9.1, Table 9).
make benchmark-policies

# CMBA verification latency (Section 6.6, Table 5).
make benchmark-cmba

# Deployment time comparison (Section 9.1).
make benchmark-deploy
```

Output is written to `paper/results/` as CSV files named after each
benchmark.

## Defensive Demonstration Scenarios

```bash
# Scenario A: hot-swap on shared volume.
./paper/scenarios/A-hot-swap.sh

# Scenario B: mismatched model at admission.
./paper/scenarios/B-mismatched-admission.sh

# Scenario C: tampered binding signature.
./paper/scenarios/C-tampered-signature.sh
```

Each scenario script asserts the expected detection and exits non-zero
on unexpected behavior.

## Tear Down

```bash
kind delete cluster --name agm-bench
```
