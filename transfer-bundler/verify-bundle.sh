#!/bin/bash
# verify-bundle.sh — verify an air-gap transfer bundle on the internal side.
#
# Usage: verify-bundle.sh <bundle.tar> [<trust-root-dir>]
#
# Steps:
#   1. Extract the bundle into a temp directory.
#   2. Verify manifest.json: every file's recorded SHA-256 must match its
#      actual SHA-256.
#   3. Verify the cosign signature on modelbinding.yaml against the local
#      Sigstore trust root.
#
# Exit non-zero on any verification failure.

set -euo pipefail

BUNDLE="${1:?usage: verify-bundle.sh <bundle.tar> [<trust-root-dir>]}"
TRUST_ROOT="${2:-/etc/airgap/sigstore}"

WORK=$(mktemp -d)
trap "rm -rf $WORK" EXIT

tar -xf "$BUNDLE" -C "$WORK"

if [ ! -f "$WORK/manifest.json" ]; then
  echo "Bundle missing manifest.json" >&2
  exit 1
fi

# Verify every recorded file hash.
python3 - "$WORK" <<'EOF' || exit 1
import hashlib, json, os, sys
work = sys.argv[1]
with open(os.path.join(work, "manifest.json")) as f:
    m = json.load(f)
for entry in m["files"]:
    p = os.path.join(work, entry["name"])
    h = hashlib.sha256()
    with open(p, "rb") as fh:
        for chunk in iter(lambda: fh.read(65536), b""):
            h.update(chunk)
    actual = h.hexdigest()
    if actual != entry["sha256"]:
        print(f"HASH MISMATCH: {entry['name']} expected={entry['sha256']} actual={actual}")
        sys.exit(2)
print("manifest hashes OK")
EOF

# Verify cosign signature on modelbinding.yaml.
if command -v cosign >/dev/null; then
  cosign verify-blob \
    --offline \
    --bundle "$TRUST_ROOT" \
    --signature "$WORK/modelbinding.yaml.sig" \
    "$WORK/modelbinding.yaml"
  echo "modelbinding signature OK"
else
  echo "cosign binary not found; skipping signature verification" >&2
  exit 3
fi

echo "Bundle verification: PASS"
