#!/usr/bin/env bash
set -euo pipefail

CFG_PATH="${1:-envoy/envoy.yaml}"

if [[ ! -f "$CFG_PATH" ]]; then
  echo "ERROR: $CFG_PATH not found" >&2
  exit 1
fi

perm="$(stat -c '%a' "$CFG_PATH" 2>/dev/null || stat -f '%Lp' "$CFG_PATH")"
gr=$(( (perm / 10) % 10 ))
or=$(( perm % 10 ))

if [[ $((gr & 4)) -eq 0 && $((or & 4)) -eq 0 ]]; then
  echo "ERROR: $CFG_PATH is not readable by group/others." >&2
  echo "Envoy often runs as a non-root user inside the container; mounted config must be readable." >&2
  echo "Fix: chmod go+r $CFG_PATH" >&2
  exit 1
fi

echo "OK: $CFG_PATH permissions look good ($perm)"
