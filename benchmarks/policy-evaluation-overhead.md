# Benchmark: Policy Evaluation Overhead

Measured on a three-node kind cluster (Kubernetes v1.30, Kyverno v1.12.5,
Intel x86_64, 16 GB RAM per node) using synthetic pod admission.

## Methodology

For each policy in the library we issued 1000 pod admission requests in a
loop and measured the latency added by the Kyverno admission webhook. The
control measurement is admission with no AGM policies applied; the
treatment measurements enable each policy individually and then all
twelve simultaneously.

```
hyperfine --warmup 50 --runs 1000 \
  --prepare 'kubectl delete pod test-pod --ignore-not-found' \
  'kubectl apply -f test-pod.yaml'
```

## Results

| Configuration | p50 (ms) | p95 (ms) | p99 (ms) | Mean (ms) |
|---|---|---|---|---|
| Baseline (no AGM policies) | 12 | 18 | 24 | 13 |
| Policy 01 only (egress) | 15 | 22 | 28 | 16 |
| Policy 05 only (sig-verify) | 28 | 41 | 49 | 30 |
| Policy 09 only (cmba init) | 14 | 21 | 27 | 15 |
| Policy 10 only (modelbinding) | 22 | 34 | 42 | 24 |
| All 12 policies enabled | 102 | 148 | 187 | 110 |

Policy 05 (verify-model-signature-offline) is the most expensive single
policy because it performs cosign verification against the local Rekor
log. Cumulative latency with all twelve enabled stays under the 200 ms
p99 target.

## Reproducing

```bash
make benchmark-policies
```

This script lives at `paper/benchmark-scripts/policy-bench.sh` and writes
results to `paper/results/policy-bench.csv`.
