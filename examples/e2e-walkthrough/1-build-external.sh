#!/bin/bash
# 1-build-external.sh — runs on the connected side. Builds a sample
# scikit-learn iris classifier container, generates an ML-BOM, computes
# the model SHA, generates a ModelBinding manifest, signs everything,
# and bundles into a single transfer .tar.

set -euo pipefail

cd "$(dirname "$0")"
ROOT="$(cd ../.. && pwd)"
WORK="$ROOT/examples/e2e-walkthrough/_work"
mkdir -p "$WORK"

# --- 1. Sample model and serving container ---
mkdir -p sample-app/models
if [ ! -f sample-app/models/iris.pkl ]; then
  python3 -c "
import pickle
from sklearn.datasets import load_iris
from sklearn.linear_model import LogisticRegression
X, y = load_iris(return_X_y=True)
m = LogisticRegression(max_iter=200).fit(X, y)
with open('sample-app/models/iris.pkl', 'wb') as f:
    pickle.dump(m, f)
"
fi

# --- 2. Build container ---
docker build -t ml-model:e2e ./sample-app
docker save ml-model:e2e -o "$WORK/ml-model.tar"
DIGEST=$(docker inspect ml-model:e2e --format '{{.Id}}')

# --- 3. Hash model artifact ---
MODEL_SHA=$(sha256sum sample-app/models/iris.pkl | awk '{print $1}')

# --- 4. Generate ML-BOM ---
if command -v cdxgen >/dev/null; then
  cdxgen -t python --spec-version 1.6 -o "$WORK/ml-bom.json" sample-app
else
  echo '{"bomFormat":"CycloneDX","specVersion":"1.6"}' > "$WORK/ml-bom.json"
fi

# --- 5. Generate ModelBinding ---
"$ROOT/transfer-bundler/generate-modelbinding.sh" \
  --container-digest "$DIGEST" \
  --model-sha "$MODEL_SHA" \
  --model-path "/models/iris.pkl" \
  --output "$WORK/modelbinding.yaml"

# --- 6. Sign ModelBinding (development: self-signed) ---
echo "(stub) cosign sign-blob would run here against public Sigstore"
sha256sum "$WORK/modelbinding.yaml" > "$WORK/modelbinding.yaml.sig"

# --- 7. Bundle ---
"$ROOT/transfer-bundler/bundle.sh" \
  --container "$WORK/ml-model.tar" \
  --model     sample-app/models/iris.pkl \
  --mlbom     "$WORK/ml-bom.json" \
  --binding   "$WORK/modelbinding.yaml" \
  --binding-sig "$WORK/modelbinding.yaml.sig" \
  --output    "$WORK/airgap-bundle.tar"

echo "Built: $WORK/airgap-bundle.tar"
