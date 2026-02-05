#!/usr/bin/env bash
set -euo pipefail

CFG_PATH="${1:-envoy/envoy.yaml}"

command -v docker >/dev/null 2>&1 || { echo "ERROR: docker not found" >&2; exit 1; }
docker compose version >/dev/null 2>&1 || { echo "ERROR: docker compose not available" >&2; exit 1; }

"$(dirname "$0")/envoy-config-perms.sh" "$CFG_PATH"
echo "OK: preflight"
