# Transfer Preparation

How to prepare an air-gap transfer bundle on the connected side.

## Bundle Layout

```
bundle.tar
├── manifest.json           # SHA-256 hashes of every file
├── ml-model.tar            # container image (docker save format)
├── model.<ext>             # model artifact
├── ml-bom.json             # CycloneDX 1.6 ML-BOM
├── modelbinding.yaml       # ModelBinding CRD manifest
└── modelbinding.yaml.sig   # detached cosign signature
```

## Bundle Generation

The simplest way is to use the GitHub Actions workflow:

```yaml
# Triggers on every tag matching v*
.github/workflows/airgap-bundle-build.yml
```

For manual generation:

```bash
./transfer-bundler/bundle.sh \
  --container ./build/ml-model.tar \
  --model ./models/iris.pkl \
  --mlbom ./build/ml-bom.json \
  --binding ./build/modelbinding.yaml \
  --binding-sig ./build/modelbinding.yaml.sig \
  --output ./out/airgap-bundle.tar
```

## Transfer Channels

AGM is agnostic about how the bundle physically crosses the air gap. The
boundary artifact (the signed tar) is what matters. Common transfer
channels:

- **Removable media**: USB drive, external SSD. Simple but operationally
  awkward for large bundles or frequent updates.
- **One-way data diode**: hardware appliance that physically allows data
  to flow only in one direction. Used in defense and critical infrastructure.
- **Cross-domain solution (CDS)**: software-and-policy-controlled transfer
  with automated content inspection. Used in classified environments.
- **Air-gapped relay cluster**: a staging cluster on the connected side that
  syncs to the production air-gapped cluster via controlled batch transfer.

## Trust Root Pre-Staging

Before the cluster goes air-gapped, you must pre-stage:

1. The Sigstore TUF root (`root.json`).
2. Fulcio root and intermediate certificates.
3. Rekor public key.

These are mounted into every verification pod as a ConfigMap named
`airgap-trust-root`. Updating the trust roots requires a controlled rotation
ceremony; see Sigstore's documentation on TUF root rotation.

## Transfer Audit

Every bundle import is recorded in the Evidence Ledger with:
- Bundle SHA-256
- Bundle size
- Generator identifier
- Import timestamp
- Importing user (Kubernetes ServiceAccount)
- Verification result (pass/fail)

These records are the basis for the export audit bundle that leaves the
cluster through the one-way data diode.
