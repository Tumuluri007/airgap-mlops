# airgap-mlops Helm Chart

One-command deployment of the entire AGM platform: Kyverno + 12 ML
policies, CMBA admission webhook + sentinel + init container, ModelBinding
CRD, local container registry, local Sigstore stack, MLflow model
registry, KServe, and the Evidence Ledger.

## Install

### Developer (kind/minikube)

```bash
helm install airgap-mlops ./helm/airgap-mlops \
  -f ./helm/airgap-mlops/values.yaml
```

### Air-Gapped Production

```bash
helm install airgap-mlops oci://registry.airgap.local/charts/airgap-mlops \
  --version 0.1.0 \
  -f ./helm/airgap-mlops/values-airgap.yaml
```

## What Gets Installed

| Template | What It Deploys |
|---|---|
| `namespace.yaml` | `ml-serving`, `ml-training`, `airgap-system`, `cmba-system` namespaces with `airgap.mlops/enforced=true` label |
| `kyverno-installation.yaml` | Kyverno controller in offline mode |
| `policies-configmap.yaml` | 12 ClusterPolicies wrapped as a single applyable manifest |
| `cmba-webhook.yaml` | CMBA admission webhook + ValidatingWebhookConfiguration |
| `cmba-sentinel-rbac.yaml` | RBAC for the sentinel sidecar to update ModelBinding status |
| `modelbinding-crd.yaml` | The ModelBinding CRD |
| `local-registry.yaml` | Harbor or Zot for local container images |
| `local-sigstore.yaml` | Local Rekor + Fulcio + TUF root server |
| `mlflow.yaml` | Local MLflow with embedded Postgres |
| `kserve.yaml` | KServe controller |
| `evidence-ledger.yaml` | Loki + Postgres on WORM-backed PVC |
| `network-policies.yaml` | Default-deny egress for all `ml-*` namespaces |

## Values Reference

See `values.yaml` for the full set of overridable values, and
`values-airgap.yaml` for the production air-gapped configuration.

## Uninstall

```bash
helm uninstall airgap-mlops
```

## Tested Versions

- Kubernetes 1.28+
- Kyverno 1.12+
- KServe 0.13+
- MLflow 2.15+
- Helm 3.14+
