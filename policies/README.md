# ML Kyverno Policy Library

Twelve Kyverno ClusterPolicies designed for ML workloads in air-gapped
Kubernetes clusters. As of submission, the official Kyverno policy library at
[kyverno.io/policies](https://kyverno.io/policies) contains zero policies for
machine learning workloads.

## Layout

| Group | Policy IDs | Purpose |
|---|---|---|
| `airgap/` | 01-04 | Enforce air-gap discipline (egress, registries, init containers, trust roots) |
| `supply-chain/` | 05-08 | Enforce supply chain integrity (signing, ML-BOM, SLSA, image policy) |
| `cmba/` | 09-11 | Enforce CMBA components (init container, ModelBinding, sentinel sidecar) |
| `governance/` | 12 | Enforce ML resource governance (GPU and memory limits) |

## Applying Policies

```bash
kubectl apply -f policies/airgap/
kubectl apply -f policies/supply-chain/
kubectl apply -f policies/cmba/
kubectl apply -f policies/governance/
```

Or all at once:

```bash
kubectl apply -R -f policies/
```

## Testing

Each policy ships with chainsaw test cases. Run with:

```bash
make policy-test
```
