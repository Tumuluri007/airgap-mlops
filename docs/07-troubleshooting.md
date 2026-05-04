# Troubleshooting

Common issues and how to diagnose them.

## CMBA Webhook Rejects All Pods

### Symptom

Every pod in `ml-serving` is rejected at admission with `CMBABindingViolation`.

### Common Causes

1. **Missing ModelBinding**: the pod references a binding that does not exist.
   ```bash
   kubectl get modelbindings.cmba.airgap.mlops -n ml-serving
   ```
2. **Image digest mismatch**: container image in pod does not match the
   ModelBinding spec.
3. **Trust root not mounted**: the webhook pod cannot find the trust root
   ConfigMap.
   ```bash
   kubectl describe pod -n cmba-system -l app=cmba-webhook
   ```

### Resolution

Check webhook logs:
```bash
kubectl logs -n cmba-system -l app=cmba-webhook --tail=100
```

The reject reason is in the `AdmissionResponse.Status.Message` field and
also logged.

## Init Container Loops with CMBA MISMATCH

### Symptom

Pod stuck in `Init:Error`, init container exits with `CMBA MISMATCH`.

### Diagnosis

```bash
kubectl logs <pod-name> -c cmba-verify
```

The log shows the expected SHA-256 (from ModelBinding) and the actual
SHA-256 (computed from the file). If they differ:

- The model file on the PVC has been replaced or corrupted.
- The ModelBinding was not regenerated for the current model artifact.

### Resolution

1. Recompute the model SHA-256: `sha256sum /path/to/model.pkl`.
2. Update the ModelBinding spec to match.
3. Re-sign the ModelBinding manifest.
4. Re-apply.

## Sentinel Sidecar Reports Frequent Drift

### Symptom

`ModelBinding.status.mismatchEvents` increments unexpectedly.

### Diagnosis

```bash
kubectl get modelbinding <name> -n ml-serving -o yaml
```

Look at `status.mismatchEvents` and `status.lastVerified`. If mismatches
correlate with model promotions or operational changes, the model file is
being modified out-of-band.

### Resolution

- Audit who has write access to the model PVC.
- Enable WORM storage on the PVC if available.
- Reduce write paths: use a CronJob to copy new models in, with a Kyverno
  policy requiring the CronJob's ServiceAccount.

## Helm Install Fails

### Symptom

`helm install airgap-mlops ./helm/airgap-mlops` returns errors.

### Common Causes

1. **CRDs not installed**: ModelBinding CRD must exist before any
   ModelBinding resource is created.
   ```bash
   kubectl apply -f cmba/crd/modelbinding-crd.yaml
   ```
2. **cert-manager missing**: the webhook needs cert-manager for TLS.
   ```bash
   kubectl apply -f https://github.com/cert-manager/cert-manager/releases/latest/download/cert-manager.yaml
   ```
3. **Kyverno not installed**: install Kyverno upstream chart first.
   ```bash
   helm repo add kyverno https://kyverno.github.io/kyverno/
   helm install kyverno kyverno/kyverno -n kyverno --create-namespace
   ```

## Bundle Verification Fails

### Symptom

`./transfer-bundler/verify-bundle.sh ./airgap-bundle.tar` returns non-zero.

### Diagnosis

Run with `-x` to see which check failed:
```bash
bash -x ./transfer-bundler/verify-bundle.sh ./airgap-bundle.tar
```

### Common Causes

- **Manifest hash mismatch**: file inside the bundle was modified after
  bundle.sh ran (corruption during transfer or tampering).
- **Cosign signature fails**: trust root mismatch; bundle was signed
  against a different identity.
- **`cosign` binary not found**: install cosign 2.4+ on the verifier host.

## Getting More Help

- File an issue: https://github.com/Tumuluri007/airgap-mlops/issues
- For security issues: see [`SECURITY.md`](../SECURITY.md).
