# End-to-End Walkthrough

Four numbered shell scripts that take you from external bundle build to
internal model serving in under thirty minutes on a developer laptop.

## Prerequisites

- Docker
- kind or minikube
- helm 3.14+
- kubectl
- cosign 2.4+
- yq 4.x

## Steps

```bash
# 1. Build a signed transfer bundle on the connected side.
./1-build-external.sh

# 2. Simulate physical transfer (just moves the file into the drop folder).
./2-transfer.sh

# 3. Run the internal pipeline: verify, import, apply.
./3-import-internal.sh

# 4. Send a sample inference request and inspect the Evidence Ledger.
./4-deploy-and-verify.sh
```

## What You Should See

After step 4, you should see:

- An InferenceService named `iris-classifier` running in the `ml-serving` namespace.
- A successful inference response from the model.
- Three event entries in the Evidence Ledger:
  - `AGMBundleVerified` — bundle signature verified on import.
  - `AGMBundleImported` — bundle artifacts loaded into local registry.
  - `CMBABindingVerified` — sentinel sidecar's first successful re-verification.

## Cleanup

```bash
kind delete cluster --name agm-test
```
