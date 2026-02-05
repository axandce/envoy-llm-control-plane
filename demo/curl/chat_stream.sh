#!/usr/bin/env bash
set -euo pipefail
source "$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)/scripts/harness.sh"

stack_up
wait_ready "$ENVOY_HEALTH_URL" 30

echo "=== Streaming allow demo ==="
# Keep your existing curl, but ideally reuse curl_json.
curl -N -sS \
  -H "Authorization: Bearer test" \
  -H "Content-Type: application/json" \
  "$API_BASE_URL/v1/chat/completions" \
  -d '{"model":"mock","stream":true,"messages":[{"role":"user","content":"hi"}]}'

echo
