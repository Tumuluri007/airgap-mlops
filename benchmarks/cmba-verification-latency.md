# Benchmark: CMBA Verification Latency

Measured on a three-node kind cluster (Kubernetes v1.30, Intel x86_64,
16 GB RAM per node).

## Methodology

We instrument the admission webhook, init container, and sentinel
sidecar to record verification latency. Pod admission latency is
measured end-to-end from `kubectl apply` to admission response.

## Admission Webhook Latency

| Percentile | Latency (ms) |
|---|---|
| p50 | 22 |
| p90 | 31 |
| p95 | 34 |
| p99 | 38 |

Target: under 50 ms p99. Result: 38 ms p99.

## Init Container Hashing Time

Time to hash the mounted model file in the cmba-verify init container.

| Model Size | Time (ms) |
|---|---|
| 10 MB | 38 |
| 100 MB | 87 |
| 500 MB | 165 |
| 1 GB | 1140 |
| 5 GB | 1380 |

For models up to 5 GB, init container overhead added to pod boot stays
under 1.5 seconds.

## Sentinel Sidecar Footprint

Measured over a 24-hour run with `recheckIntervalSeconds=300`.

| Metric | Value |
|---|---|
| Steady-state RSS | 3.8 MB |
| Average CPU | ~0.1% of one core |
| Re-verification time, 500 MB model | 165 ms |

## Reproducing

```bash
make benchmark-cmba
```

Results written to `paper/results/cmba-bench.csv`.
