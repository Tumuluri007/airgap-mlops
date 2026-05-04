#!/bin/sh
# cmba-verify init container.
#
# Reads the mounted ModelBinding, hashes the actual model file at the
# declared path, and exits non-zero on mismatch. Pod startup is blocked
# when the hashes do not match.
#
# Pure shell. Depends only on busybox plus yq, both available offline.
# Simple primitives are auditable; the script is short by design.

set -eu

MODEL_PATH="${CMBA_MODEL_PATH:-/models}"
BINDING_FILE="${CMBA_BINDING_FILE:-/etc/cmba/binding.yaml}"

EXPECTED_SHA=$(yq '.spec.modelArtifact.sha256' "$BINDING_FILE")
EXPECTED_PATH=$(yq '.spec.modelArtifact.path' "$BINDING_FILE")
ACTUAL_SHA=$(sha256sum "$EXPECTED_PATH" | awk '{print $1}')

if [ "$ACTUAL_SHA" != "$EXPECTED_SHA" ]; then
  echo "CMBA MISMATCH: expected=$EXPECTED_SHA actual=$ACTUAL_SHA"
  echo "Model artifact at $EXPECTED_PATH does not match ModelBinding declaration."
  exit 1
fi

echo "CMBA OK: model $EXPECTED_PATH verified against binding"
exit 0
