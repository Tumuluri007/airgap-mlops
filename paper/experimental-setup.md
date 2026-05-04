# Experimental Setup

Detailed parameters for the experiments reported in the paper.

## Hardware

- Single workstation with 3 kind nodes
- CPU: Intel Core i7-12700K, 12 cores
- RAM: 64 GB DDR4
- Storage: NVMe SSD, ext4
- Network: localhost (kind cluster, no external traffic)

## Software Versions

| Component | Version |
|---|---|
| Linux kernel | 6.5.0 |
| Docker | 24.0.7 |
| kind | 0.23.0 |
| Kubernetes | 1.30.0 |
| Kyverno | 1.12.5 |
| KServe | 0.13.0 |
| MLflow | 2.15.0 |
| Cosign | 2.4.0 |
| Sigstore Rekor | 1.3.6 |
| Sigstore Fulcio | 1.5.0 |
| CycloneDX cdxgen | 10.9.0 |
| Argo Workflows | 3.5.6 |
| Loki | 3.0.0 |
| cert-manager | 1.15.0 |

## Benchmark Workload Profiles

### Policy Evaluation Overhead

- 1000 sequential pod admission requests per policy configuration
- 50-request warm-up phase before each measurement run
- Pod spec: minimal scikit-learn serving pod with valid CMBA annotations,
  cmba-verify init container, cmba-sentinel sidecar
- Each measurement run repeated 3 times; median reported

### CMBA Verification Latency

- 1000 pod admissions for webhook latency
- 100 init container runs per model size for hashing time
- 24-hour continuous run for sidecar footprint, recheckIntervalSeconds=300
- Model files generated as random data of the specified size

### Deployment Time Comparison

- Three configurations: connected baseline, air-gapped two-pipeline,
  air-gapped with all 12 policies + CMBA active
- Each measurement: 5 sequential model deployments, median reported
- Model: scikit-learn iris classifier (~12 KB)
- Container: registry.airgap.local/inference/sklearn-server:e2e

### Catch Rate Scenarios

- Each scenario run 50 times consecutively
- Detection time measured from event injection to alert in Evidence Ledger
- Pod state transitions observed via `kubectl get events --watch`

## Statistical Notes

- Latency percentiles reported as p50, p95, p99 unless otherwise stated
- Mean and standard deviation also recorded in raw CSVs
- All timing measurements use Go's `time.Now().UnixNano()` or
  `kubectl apply -o json | jq .metadata.creationTimestamp`
