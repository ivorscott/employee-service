
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
	CompileDaemon --build="go build ./cmd/employee" --log-prefix=false --command=./employee
.PHONY: develop
