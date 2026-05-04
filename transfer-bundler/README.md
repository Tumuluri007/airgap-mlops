# Transfer Bundler

Shell scripts for building, signing, and verifying air-gap transfer bundles.

## Files

- `bundle.sh` — build a signed transfer bundle from a container image, a
  model artifact, an ML-BOM, and a signed ModelBinding manifest.
- `verify-bundle.sh` — verify a bundle on the internal (air-gapped) side.
  Checks manifest hashes and cosign signature against the pre-staged
  trust root.
- `generate-modelbinding.sh` — emit a ModelBinding YAML manifest with the
  bindingHash computed from container digest and model SHA.
- `manifest-schema.json` — JSON Schema for the bundle manifest format.

## Bundle Format

Every transfer bundle is a `.tar` archive with this layout:

```
bundle.tar
├── manifest.json           # bundle-level manifest (SHA-256 of every file)
├── ml-model.tar            # container image saved as docker-archive
├── model.<ext>             # model artifact (pickle, onnx, safetensors, etc.)
├── ml-bom.json             # CycloneDX 1.6 ML-BOM
├── modelbinding.yaml       # ModelBinding CRD manifest
└── modelbinding.yaml.sig   # detached cosign signature for the binding
```

## End-to-End

```bash
# External side (connected): build the bundle.
./bundle.sh \
  --container ./ml-model.tar \
  --model ./models/iris.pkl \
  --mlbom ./ml-bom.json \
  --binding ./modelbinding.yaml \
  --binding-sig ./modelbinding.yaml.sig \
  --output ./airgap-bundle.tar

# Physical transfer to the air-gapped cluster (out of scope of this script).

# Internal side (air-gapped): verify the bundle.
./verify-bundle.sh ./airgap-bundle.tar /etc/airgap/sigstore
```
