#!/usr/bin/env bash
set -euo pipefail

curl -sS \
  -H "Content-Type: application/json" \
  -H "x-api-key: demo-key" \
  -d '{
    "model":"gpt-4o-mini",
    "stream": false,
    "messages":[{"role":"user","content":"Say hello in 5 words."}]
  }' \
  http://localhost:8080/v1/chat/completions | sed -n '1,120p'

echo
