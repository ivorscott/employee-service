include .env

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)

default: develop

generate:
	swagger-codegen generate -i doc/api-doc.yml -l openapi -o cmd/employee/static/swagger-ui
.PHONY: generate

fmt:
	go fmt ./...
.PHONY: fmt

lint:
	golangci-lint run
.PHONY: lint

vet:
	go vet ./...
.PHONY: vet

test: fmt lint vet
	go test --cover ./...
.PHONY: test

build: test
	go build ./cmd/employee
.PHONY: build

develop: generate
	CompileDaemon --build="go build ./cmd/employee" --log-prefix=false --command="./employee --db-disable-tls=true"
.PHONY: develop

db:
	psql $(DB_URL)
.PHONY: db

pg:
	pgcli $(DB_URL)
.PHONY: pg

# ======================================================================================================================
# Begins Migration and Seeding Helper
# ======================================================================================================================
# For the following usage -> "make migration <name>" http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),migration seed))
  name := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(name):;@:)
endif

# For the following usage -> "make up <number>", "make down <number>", "make force <number>"
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),up down force))
  num := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(num):;@:)
# "make down" without a number defaults to 1.
  ifndef num
    ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),down))
      num := 1
    endif
  endif
endif

define err_create_migration


Error: migration name is missing.
Usage: make migration <name>

$(shell echo Take a coffee break "\xE2\x98\x95")
endef

define err_force_migration


Error: migration version is missing.
Usage: make force <version>

$(shell echo Take a coffee break "\xE2\x98\x95")
endef

CHECKMARK="\xE2\x9C\x94"
BAD_INPUT="you supplied an incorrect argument"

migration:
    ifndef name
		$(error ${err_create_migration}))
    endif

	@migrate create -ext sql -dir ./res/migration -seq $(name) \
	&& echo $(CHECKMARK) Successfully created migrations!
.PHONY: migration

version:
	@migrate -path ./res/migration -database $(DB_URL) version \
	&& echo $(CHECKMARK) "Here's the current version!" \
	|| echo Did you reach the bottom? You might not be on an active version.
.PHONY: version

up:
	@migrate -path ./res/migration -verbose -database $(DB_URL) up $(num) \
	&& echo $(CHECKMARK) Successfully migrated! \
	|| echo There might not be any up migrations left or $(BAD_INPUT).
.PHONY: up

down:
	@migrate -path ./res/migration -verbose -database $(DB_URL) down $(num) \
	&& echo $(CHECKMARK) Successfully downgraded! \
	|| echo There might not be any down migrations left or $(BAD_INPUT).
.PHONY: down

downfall:
	@migrate -path ./res/migration -verbose -database $(DB_URL) down -all \
	&& echo $(CHECKMARK) Successfully applied all down migrations! \
	|| echo Use the force Luke.
.PHONY: downfall

# About force https://bit.ly/3exuENS
force:
    ifndef num
		$(error ${err_force_migration}))
    endif

	@migrate -path ./res/migration -verbose -database $(DB_URL) force $(num) \
	&& echo $(CHECKMARK) Successfully forced migration to version $(num)!
.PHONY: force

seed:
    ifndef name
		$(error seed name is missing -> make seed <name>)
    endif
	@psql -h $(POSTGRES_HOST) -U $(POSTGRES_USER) -p $(POSTGRES_PORT) $(POSTGRES_DB) -f $(name)
.PHONY: seed res/seed/*

# ======================================================================================================================
# Ends Migration and Seeding Helper
# ======================================================================================================================
