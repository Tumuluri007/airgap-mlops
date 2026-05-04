#!/bin/bash
# generate-modelbinding.sh — emit a ModelBinding YAML manifest given the
# container digest, the model SHA, and the in-pod model path.
#
# Usage:
#   generate-modelbinding.sh \
#     --container-digest <sha256:...> \
#     --model-sha <hex64> \
#     --model-path </models/...> \
#     --output <path>

set -euo pipefail

CONTAINER_DIGEST=""; MODEL_SHA=""; MODEL_PATH=""; OUTPUT=""

while [ $# -gt 0 ]; do
  case "$1" in
    --container-digest) CONTAINER_DIGEST="$2"; shift 2;;
    --model-sha)        MODEL_SHA="$2"; shift 2;;
    --model-path)       MODEL_PATH="$2"; shift 2;;
    --output)           OUTPUT="$2"; shift 2;;
    *) echo "Unknown flag: $1" >&2; exit 1;;
  esac
done

if [ -z "$CONTAINER_DIGEST" ] || [ -z "$MODEL_SHA" ] || [ -z "$MODEL_PATH" ] || [ -z "$OUTPUT" ]; then
  echo "All flags required: --container-digest --model-sha --model-path --output" >&2
  exit 1
fi

# Strip the algorithm prefix from the container digest if present, and
# concatenate with the model SHA to compute bindingHash.
DIGEST_HEX="${CONTAINER_DIGEST#sha256:}"
BINDING_HASH=$(printf '%s%s' "$DIGEST_HEX" "$MODEL_SHA" | sha256sum | awk '{print $1}')

cat > "$OUTPUT" <<EOF
apiVersion: cmba.airgap.mlops/v1alpha1
kind: ModelBinding
metadata:
  name: ml-model
  namespace: ml-serving
spec:
  containerImage: "registry.airgap.local/inference/ml-model@${CONTAINER_DIGEST}"
  modelArtifact:
    path: "${MODEL_PATH}"
    sha256: "${MODEL_SHA}"
    sizeBytes: 0
  bindingHash: "${BINDING_HASH}"
  verificationPolicy:
    initContainerRequired: true
    sentinelRequired: true
    recheckIntervalSeconds: 300
    onMismatch: terminate
EOF

echo "ModelBinding manifest written: $OUTPUT"
