#!/bin/bash
# bundle.sh — build a signed air-gap transfer bundle.
#
# Inputs (all required):
#   --container <path>     Container image saved as docker-archive .tar
#   --model <path>         Model artifact file
#   --mlbom <path>         CycloneDX ML-BOM JSON
#   --binding <path>       ModelBinding manifest YAML
#   --binding-sig <path>   Detached cosign signature for the binding
#   --output <path>        Output bundle .tar path
#
# Bundle layout:
#   ml-model.tar           container image
#   model.<ext>            model artifact
#   ml-bom.json            CycloneDX 1.6 ML-BOM
#   modelbinding.yaml      ModelBinding CRD manifest
#   modelbinding.yaml.sig  detached cosign signature
#   manifest.json          bundle-level manifest (SHA-256 of every file)

set -euo pipefail

CONTAINER=""; MODEL=""; MLBOM=""; BINDING=""; BINDING_SIG=""; OUTPUT=""

while [ $# -gt 0 ]; do
  case "$1" in
    --container)    CONTAINER="$2"; shift 2;;
    --model)        MODEL="$2"; shift 2;;
    --mlbom)        MLBOM="$2"; shift 2;;
    --binding)      BINDING="$2"; shift 2;;
    --binding-sig)  BINDING_SIG="$2"; shift 2;;
    --output)       OUTPUT="$2"; shift 2;;
    *) echo "Unknown flag: $1" >&2; exit 1;;
  esac
done

for f in "$CONTAINER" "$MODEL" "$MLBOM" "$BINDING" "$BINDING_SIG"; do
  if [ ! -f "$f" ]; then
    echo "Required file missing: $f" >&2
    exit 1
  fi
done

WORK=$(mktemp -d)
trap "rm -rf $WORK" EXIT

cp "$CONTAINER"   "$WORK/ml-model.tar"
cp "$MODEL"       "$WORK/$(basename "$MODEL")"
cp "$MLBOM"       "$WORK/ml-bom.json"
cp "$BINDING"     "$WORK/modelbinding.yaml"
cp "$BINDING_SIG" "$WORK/modelbinding.yaml.sig"

# Bundle-level manifest with SHA-256 of every file.
{
  echo "{"
  echo "  \"created\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\","
  echo "  \"generator\": \"agm-bundle.sh/0.1.0\","
  echo "  \"files\": ["
  FIRST=1
  for f in ml-model.tar "$(basename "$MODEL")" ml-bom.json modelbinding.yaml modelbinding.yaml.sig; do
    SHA=$(sha256sum "$WORK/$f" | awk '{print $1}')
    SIZE=$(stat -c %s "$WORK/$f" 2>/dev/null || stat -f %z "$WORK/$f")
    if [ $FIRST -eq 0 ]; then echo ","; fi
    FIRST=0
    printf "    {\"name\": \"%s\", \"sha256\": \"%s\", \"size\": %s}" "$f" "$SHA" "$SIZE"
  done
  echo ""
  echo "  ]"
  echo "}"
} > "$WORK/manifest.json"

mkdir -p "$(dirname "$OUTPUT")"
tar -cf "$OUTPUT" -C "$WORK" .

echo "Bundle written: $OUTPUT"
echo "  size: $(stat -c %s "$OUTPUT" 2>/dev/null || stat -f %z "$OUTPUT") bytes"
echo "  sha256: $(sha256sum "$OUTPUT" | awk '{print $1}')"
