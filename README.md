# envoy-llm-control-plane

**OpenAI-compatible LLM policy enforcement using Envoy.**

This repository demonstrates how to use **Envoy** as a stable, vendor-neutral **control plane** for LLM traffic—handling authentication, coarse rate limiting, routing, and observability—while delegating all LLM-specific intelligence to a **swappable policy engine** via `ext_authz`.

> Design principle:  
> **If you can replace this with plain Envoy + a tiny auth filter and nothing breaks (besides policy), the architecture is correct.**

## Motivation

Most “LLM gateways” collapse too many responsibilities into one service:
- request handling
- policy logic
- billing
- provider abstraction
- agent orchestration

That coupling makes systems fragile, hard to audit, and painful to evolve.

This project takes the opposite approach:

- **Envoy stays boring**: transport, enforcement, observability
- **Policy stays external**: swappable, versioned, evolvable
- **Inference stays separate**: correctness and compatibility are isolated concerns

The result is an LLM ingress architecture that is:
- enterprise-friendly
- auditable
- vendor-neutral
- and easy to reason about

## What this repo is

- A **reference implementation** of Envoy as an OpenAI-compatible LLM ingress
- A demonstration of **policy enforcement via `envoy.filters.http.ext_authz`**
- A **streaming-safe** configuration for SSE / chunked responses
- A runnable local stack with:
  - Envoy
  - a minimal `ext_authz` policy stub
  - a minimal mock upstream

This is intentionally **infrastructure, not product**.

## What this repo is not

- Not an agent framework
- Not a RAG system
- Not a billing platform
- Not a full LLM “gateway” product

Those concerns belong elsewhere, behind stable interfaces.

## Architecture overview

```
Client (OpenAI SDK / curl)
        |
        v
Envoy (data plane / gateway)
        |
        |  ext_authz (gRPC)
        v
Policy Engine (stub in this repo, real engine swappable)
        |
        v
Upstream LLM (mock here; replace with vLLM, OpenAI, etc.)
```

### Key boundary

`ext_authz` is treated as a **stable contract** between Envoy and policy logic.

Envoy does not understand:
- tokens
- cost
- models
- tenants
- tools

It only enforces the decision returned by the policy engine.

That boundary is what makes the system replaceable and safe.

## OpenAI compatibility (what that means here)

Envoy is **not** re-implementing the OpenAI API.

“OpenAI-compatible” here means:
- Accepts OpenAI-shaped HTTP requests (e.g. `/v1/chat/completions`)
- Preserves request bodies and headers untouched
- Supports streaming (SSE / chunked responses)
- Forwards requests after policy approval

Any schema interpretation or provider-specific behavior belongs **downstream**, not in the control plane.

## Quickstart

### Requirements
- Docker + Docker Compose
- `make`

### Run the demo stack
```bash
make up
make demo
````

* Envoy ingress: `http://localhost:8080`
* Envoy admin UI: `http://localhost:9901`

### Streaming example

```bash
bash demo/curl/chat_stream.sh
```

You should see a streamed response pass through Envoy unchanged.

## Repo layout

```
envoy/
  envoy.yaml              # Envoy configuration (routes, ext_authz, streaming)

stubs/
  policy-engine-stub/     # Minimal ext_authz gRPC server (demo-only)
  mock-upstream/          # Minimal OpenAI-shaped upstream (demo-only)

demo/
  curl/                   # Allow/deny + streaming demo scripts
```

The stub services exist solely to make the repo runnable.
They are not intended as production implementations.

## Swapping in a real policy engine

To replace the stub policy engine with a real implementation:

1. Implement Envoy `ext_authz` gRPC
   (`envoy.service.auth.v3.Authorization`)
2. Preserve decision semantics:

   * allow / deny
   * optional header mutation
3. Point the `policy_engine` cluster at your service

**The Envoy configuration does not need to change.**

That immutability is the point.

## Design constraints (intentional)

* Envoy does **not** inspect request bodies
* Policy decisions are header- and metadata-driven
* Streaming routes disable request/stream timeouts
* Failure modes are explicit (fail-closed by default)

These constraints keep the control plane predictable and auditable.

## Who this is for

* Platform / infrastructure engineers designing LLM ingress
* Security teams reviewing LLM traffic enforcement patterns
* Architects evaluating Envoy-based control planes
* Teams that want OpenAI compatibility **without vendor lock-in**

---
