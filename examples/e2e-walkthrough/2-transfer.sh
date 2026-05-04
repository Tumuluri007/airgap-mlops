#!/bin/bash
# 2-transfer.sh — simulate physical transfer by copying the bundle into
# the drop folder PVC mount on the air-gapped cluster.
#
# In production: this is a USB drive, a one-way data diode, or a
# write-once cross-domain solution.

set -euo pipefail

cd "$(dirname "$0")"
WORK="_work"

if [ ! -f "$WORK/airgap-bundle.tar" ]; then
  echo "Bundle not found. Run 1-build-external.sh first." >&2
  exit 1
fi

# Copy into the drop folder PVC inside the air-gapped cluster.
DROP_NAMESPACE="airgap-system"
DROP_POD=$(kubectl get pod -n "$DROP_NAMESPACE" -l app=drop-folder -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || true)

if [ -z "$DROP_POD" ]; then
  echo "drop-folder pod not found; saving bundle to local _work dir for next step"
  cp "$WORK/airgap-bundle.tar" "$WORK/transferred.tar"
else
  kubectl cp -n "$DROP_NAMESPACE" "$WORK/airgap-bundle.tar" "$DROP_POD:/drop/airgap-bundle.tar"
  echo "Bundle copied to $DROP_NAMESPACE/$DROP_POD:/drop/airgap-bundle.tar"
fi
