#!/bin/bash
# 3-import-internal.sh — runs on the air-gapped side. Triggers the Argo
# Workflow that verifies the bundle, imports the container into the local
# registry, applies the ModelBinding CRD, and applies the InferenceService.

set -euo pipefail

NAMESPACE="airgap-system"

if ! kubectl get -n "$NAMESPACE" workflowtemplates airgap-deploy >/dev/null 2>&1; then
  kubectl apply -n "$NAMESPACE" -f ../../ci-cd/argo-workflows/airgap-deploy.yaml
fi

# Submit a workflow run from the template.
WORKFLOW_NAME=$(kubectl create -n "$NAMESPACE" -f - <<EOF
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: airgap-deploy-e2e-
spec:
  entrypoint: deploy
  workflowTemplateRef:
    name: airgap-deploy
EOF
)

echo "Submitted: $WORKFLOW_NAME"
echo "Watching for completion..."
kubectl wait -n "$NAMESPACE" --for=condition=Completed --timeout=10m "$WORKFLOW_NAME"
echo "Workflow completed successfully"
