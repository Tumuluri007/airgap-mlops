#!/bin/bash
# 4-deploy-and-verify.sh — send a sample inference request and inspect
# the Evidence Ledger entries.

set -euo pipefail

NAMESPACE="ml-serving"
ISVC_NAME="iris-classifier"

echo "Waiting for InferenceService to be ready..."
kubectl wait -n "$NAMESPACE" --for=condition=Ready --timeout=5m \
  "inferenceservice/$ISVC_NAME"

INGRESS=$(kubectl get -n "$NAMESPACE" "inferenceservice/$ISVC_NAME" \
  -o jsonpath='{.status.url}')

echo "Sending inference request to $INGRESS"
curl -sS "$INGRESS/v1/models/$ISVC_NAME:predict" \
  -H "Content-Type: application/json" \
  -d '{"instances":[[5.1, 3.5, 1.4, 0.2]]}'
echo

echo "Inspecting Evidence Ledger entries..."
kubectl logs -n airgap-system -l app=loki --tail=50 \
  | grep -E "AGMBundle|CMBABinding" || true

echo
echo "ModelBinding status:"
kubectl get -n "$NAMESPACE" modelbindings.cmba.airgap.mlops "$ISVC_NAME" \
  -o yaml | yq '.status'
