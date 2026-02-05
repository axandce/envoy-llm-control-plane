SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# -------- Settings --------
COMPOSE ?= docker compose
COMPOSE_FILE ?= docker-compose.yml
PROJECT ?=
GO_TEST_FLAGS ?= -count=1 -v -timeout=15m
ENVOY_CONFIG ?= envoy/envoy.yaml

define compose_cmd
$(COMPOSE) $(if $(PROJECT),-p $(PROJECT),) -f $(COMPOSE_FILE)
endef

# -------- Help --------
.PHONY: help
help:
	@printf "\nTargets:\n"
	@printf "  make doctor         Run reusable preflight checks (docker + compose + envoy.yaml perms)\n"
	@printf "\n"
	@printf "  make up             Start stack (compose up -d)\n"
	@printf "  make down           Stop stack (compose down -v)\n"
	@printf "  make ps             Show containers\n"
	@printf "  make logs           Follow logs (tail=200)\n"
	@printf "\n"
	@printf "  make demo           Run allow + deny demos\n"
	@printf "  make demo-allow     Streaming allow demo\n"
	@printf "  make demo-nostream  Non-streaming allow demo\n"
	@printf "  make demo-deny      Deny demo (missing key)\n"
	@printf "\n"
	@printf "  make test           Run Go tests in ./tests\n"
	@printf "  make ci             CI entrypoint (currently: make test)\n"
	@printf "\n"
	@printf "Variables:\n"
	@printf "  COMPOSE=%s\n" "$(COMPOSE)"
	@printf "  COMPOSE_FILE=%s\n" "$(COMPOSE_FILE)"
	@printf "  PROJECT=%s (optional)\n" "$(PROJECT)"
	@printf "  ENVOY_CONFIG=%s\n" "$(ENVOY_CONFIG)"
	@printf "  GO_TEST_FLAGS=%s\n" "$(GO_TEST_FLAGS)"
	@printf "\n"

# -------- Preflight / doctor --------
.PHONY: preflight
preflight:
	@scripts/preflight.sh "$(ENVOY_CONFIG)"

# Optional convenience: fix perms (not required, but handy)
.PHONY: fix-perms
fix-perms:
	@chmod go+r "$(ENVOY_CONFIG)"
	@echo "OK: set group/other read on $(ENVOY_CONFIG)"

# -------- Compose (manual ops) --------
.PHONY: up
up: doctor
	@$(compose_cmd) up -d --build --remove-orphans

.PHONY: down
down:
	@$(compose_cmd) down -v --remove-orphans

.PHONY: ps
ps:
	@$(compose_cmd) ps

.PHONY: logs
logs:
	@$(compose_cmd) logs -f --tail=200

# -------- Demos --------
.PHONY: demo
demo: demo-allow demo-deny

.PHONY: demo-allow
demo-allow:
	@bash demo/curl/chat_stream.sh

.PHONY: demo-nostream
demo-nostream:
	@bash demo/curl/chat_nostream.sh

.PHONY: demo-deny
demo-deny:
	@bash demo/curl/deny_missing_key.sh

# -------- Tests --------
.PHONY: test
test:
	@cd tests && go test ./... $(GO_TEST_FLAGS)

.PHONY: ci
ci: test
