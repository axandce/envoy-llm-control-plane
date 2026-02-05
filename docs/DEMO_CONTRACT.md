# Demo Contract — envoy-llm-control-plane

## Scope
This contract defines the externally observable behavior for three demo scenarios:
- Allow + streaming (chat_stream.sh)
- Allow + non-streaming (chat_nostream.sh)
- Deny when missing API key (deny_missing_key.sh)

These behaviors must remain stable across refactors.

## Common Definitions

### Entry Point
- URL: `http://localhost:8080/v1/chat/completions`

### Authentication / Policy Input
- Required header: `x-api-key: demo-key`
- Missing `x-api-key` must result in denial (ext_authz).

### Request Shape (OpenAI-compatible)
- POST `/v1/chat/completions`
- JSON body:
  - `model` (string)
  - `messages` (array of `{role, content}`)
  - `stream` (boolean)

### Correlation / Traceability
- Logs must include per-request correlation (at least Envoy `x-request-id` if present).
- Policy decision must be observable via logs and metrics.

---

## Demo: Allow + Streaming (`demo-allow`)
### Driver
- `demo/curl/chat_stream.sh`

### Preconditions
- Header: `x-api-key: demo-key`
- Body includes `"stream": true`

### Expected Result
- HTTP status: `200`
- Response is **Server-Sent Events** (SSE):
  - Output contains one or more lines prefixed with `data:`
  - Output contains a terminal `[DONE]`
- No denial response headers/body format is present.

### Invariants
- Request is proxied to upstream (observable in upstream metric/log).
- Decision result is recorded as allow.

---

## Demo: Allow + Non-Streaming (`demo-nostream`)
### Driver
- `demo/curl/chat_nostream.sh`

### Preconditions
- Header: `x-api-key: demo-key`
- Body includes `"stream": false`

### Expected Result
- HTTP status: `200`
- Response is **single JSON response** (not SSE):
  - Output does **not** contain `data:`
  - Output contains JSON fields consistent with Chat Completions (at minimum an `id` and `choices`)

### Invariants
- Decision result recorded as allow.
- Response is not chunked SSE.

---

## Demo: Deny Missing Key (`demo-deny`)
### Driver
- `demo/curl/deny_missing_key.sh`

### Preconditions
- Request omits header `x-api-key`

### Expected Result
- HTTP status: `403` (recommended) OR stable deny code (must not change without updating this contract)
- Response is not proxied to upstream (observable via metrics/logs)
- Response is **not** SSE:
  - Output does not contain `data:`

### Invariants
- Decision result recorded as deny with reason `missing_api_key` (or equivalent stable reason)

---

## Must-Not-Break Invariants
- Entry URL and required auth header are stable for the demo.
- Deny behavior is deterministic and machine-readable.
- Streaming vs non-streaming behavior is deterministic.
- Every request produces:
  - a policy decision log line
  - a decision metric increment (allow/deny)
