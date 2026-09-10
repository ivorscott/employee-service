# ==============================================================================
# Employee Service
#
# Getting started needs only Docker and make:
#
#   make            # build, start and seed the whole stack
#   make help       # show every command
#
# Everything that touches the database runs inside a container on the compose
# network, addressing postgres as `employee:5432`. That means no migrate
# binary, psql or pgcli on your host - and a postgres running natively on your
# machine can never be talked to by mistake.
# ==============================================================================

# Create .env automatically, so a fresh clone works with no setup step.
# GNU make remakes missing included files and restarts, so this runs first.
.env:
	@cp .env.sample .env
	@echo "==> created .env from .env.sample"

include .env

COMPOSE := docker compose

# Number of migrations to apply, e.g. `make migrate-down N=2`.
N ?=

.DEFAULT_GOAL := start

## help: show this help
help:
	@echo "Employee Service - available commands"
	@echo ""
	@grep -hE '^## ' $(firstword $(MAKEFILE_LIST)) | sed 's/^## //' \
		| awk -F': ' '{printf "  make %-16s %s\n", $$1, $$2}'
	@echo ""
.PHONY: help

# ==============================================================================
# Stack
# ==============================================================================

## start: build, start and seed the whole stack (default)
start: .env
	$(COMPOSE) up -d --build
	@echo "==> waiting for the database to be seeded..."
	@$(COMPOSE) wait seed >/dev/null 2>&1 || true
	@echo ""
	@echo "  API          http://localhost:8080/employees"
	@echo "  Grafana      http://localhost:3000"
	@echo "  Prometheus   http://localhost:9090"
	@echo "  RabbitMQ     http://localhost:15672  (guest/guest)"
	@echo ""
	@echo "  make ps   see status      make logs   follow logs"
	@echo "  make db   open a shell    make help   all commands"
	@echo ""
.PHONY: start

## stop: stop the stack, keep containers and data
stop:
	$(COMPOSE) stop
.PHONY: stop

## down: remove containers, keep data volumes
down:
	$(COMPOSE) down
.PHONY: down

## reset: remove containers AND data, then start clean
reset:
	$(COMPOSE) down -v
	@$(MAKE) start
.PHONY: reset

## ps: show the status of every service
ps:
	$(COMPOSE) ps
.PHONY: ps

## logs: follow logs (make logs S=employee-service for one service)
logs:
	$(COMPOSE) logs -f $(S)
.PHONY: logs

# ==============================================================================
# Database
# ==============================================================================

## db: open an interactive pgcli shell on the database
db:
	$(COMPOSE) run --rm pgcli
.PHONY: db

## seed: (re)load res/seed/data.sql - `make start` already does this
seed:
	$(COMPOSE) run --rm seed
.PHONY: seed

## migrate-version: print the current migration version
migrate-version:
	$(COMPOSE) run --rm migrate version
.PHONY: migrate-version

## migrate-up: apply all pending migrations (or N with N=2)
migrate-up:
	$(COMPOSE) run --rm migrate up $(N)
.PHONY: migrate-up

## migrate-down: roll back 1 migration (or N with N=2)
migrate-down:
	$(COMPOSE) run --rm migrate down $(or $(N),1)
.PHONY: migrate-down

## migrate-reset: roll back every migration
migrate-reset:
	$(COMPOSE) run --rm migrate down -all
.PHONY: migrate-reset

## migrate-force: force a dirty schema to version N, e.g. N=1
migrate-force:
	@test -n "$(N)" || { echo "usage: make migrate-force N=<version>"; exit 1; }
	$(COMPOSE) run --rm migrate force $(N)
.PHONY: migrate-force

## migrate-create: create a migration, e.g. NAME=add_widgets
migrate-create:
	@test -n "$(NAME)" || { echo "usage: make migrate-create NAME=<name>"; exit 1; }
	$(COMPOSE) run --rm migrate create -ext sql -dir /migrations -seq $(NAME)
.PHONY: migrate-create

# ==============================================================================
# Go development (needs a host Go toolchain)
# ==============================================================================

# Fail with an actionable message instead of "command not found".
define require_tool
@command -v $(1) >/dev/null 2>&1 || { \
	echo "missing host tool: $(1)"; \
	echo "install the Go dev tools with: make tools"; \
	exit 1; \
}
endef

## tools: install the host Go tools needed by develop/test
tools:
	go install github.com/githubnemo/CompileDaemon@latest
	go install github.com/vektra/mockery/v2@latest
	@command -v swagger-codegen >/dev/null 2>&1 || echo "also run: brew install swagger-codegen"
.PHONY: tools

## develop: run the service on the host with hot reload against the stack
develop:
	$(call require_tool,CompileDaemon)
	$(call require_tool,swagger-codegen)
	swagger-codegen generate -i doc/api-doc.yml -l openapi -o cmd/employee/static/swagger-ui
	CompileDaemon --build="go build ./cmd/employee" --log-prefix=false --command="./employee --db-disable-tls=true"
.PHONY: develop

## generate: regenerate mocks
generate:
	$(call require_tool,mockery)
	go generate ./...
.PHONY: generate

## test: run unit tests with coverage (starts the test database)
test: generate
	$(COMPOSE) --profile test up -d employee_test
	go test --cover ./...
.PHONY: test

## fmt: format the code
fmt:
	go fmt ./...
.PHONY: fmt

## vet: run go vet
vet:
	go vet ./...
.PHONY: vet

## lint: run golangci-lint
lint:
	$(call require_tool,golangci-lint)
	golangci-lint run
.PHONY: lint

## check: fmt, vet, lint and test
check: fmt vet lint test
.PHONY: check

## build: compile the binary
build:
	go build ./cmd/employee
.PHONY: build
