# Quick Start

Get AGM running on a local kind or minikube cluster in under 15 minutes.

## Prerequisites

- Docker Desktop or Docker Engine
- kind 0.23+ or minikube 1.33+
- kubectl 1.28+
- helm 3.14+

## Steps

### 1. Create a Local Cluster

```bash
kind create cluster --name agm-test
# or
minikube start --kubernetes-version=v1.30.0
```

### 2. Install AGM

```bash
git clone https://github.com/Tumuluri007/airgap-mlops.git
cd airgap-mlops

helm install airgap-mlops ./helm/airgap-mlops \
  -f ./helm/airgap-mlops/values.yaml
```

### 3. Verify Installation

```bash
kubectl get pods -n cmba-system
kubectl get pods -n ml-serving
kubectl get clusterpolicies | grep -i ml
kubectl get crd modelbindings.cmba.airgap.mlops
```

You should see:
- `cmba-webhook` deployment running with 1+ replicas
- 12 ClusterPolicies prefixed with the policy IDs (01- through 12-)
- The `modelbindings.cmba.airgap.mlops` CRD registered

### 4. Run the Walkthrough

```bash
cd examples/e2e-walkthrough
./1-build-external.sh
./2-transfer.sh
./3-import-internal.sh
./4-deploy-and-verify.sh
```

### 5. Tear Down

```bash
kind delete cluster --name agm-test
```

## Troubleshooting

See [`07-troubleshooting.md`](07-troubleshooting.md).
