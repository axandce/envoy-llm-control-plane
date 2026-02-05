#!/usr/bin/env bash
set -euo pipefail

echo "Calling without x-api-key (should be denied by ext_authz)..."

curl -sS -i \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}' \
  http://localhost:8080/v1/chat/completions | sed -n '1,40p'

echo
