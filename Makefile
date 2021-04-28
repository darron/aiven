BINARY_NAME ?= aiven
CONTAINER_NAME ?= aiven
GITHUB_REPO ?= github.com/darron/aiven

BUILD_FLAGS=-X=main.GitCommit=$(GIT_COMMIT)
BUILD_COMMAND=-mod=vendor -ldflags "$(BUILD_FLAGS)" -o bin/$(BINARY_NAME) main.go
GIT_COMMIT=$(shell git rev-parse HEAD)
UNAME=$(shell uname -s | tr '[:upper:]' '[:lower:]')
THIS_FILE := $(lastword $(MAKEFILE_LIST))

all: build

deps: ## Install all dependencies.
	go mod vendor && go mod tidy

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Remove compiled binaries.
	rm -f bin/$(BINARY_NAME) || true
	rm -f bin/$(BINARY_NAME)*gz || true

docker: ## Build Docker image
	docker build . -t $(CONTAINER_NAME)

build: clean
	go build $(BUILD_COMMAND)

rebuild: clean ## Force rebuild of all packages.
	go build -a $(BUILD_COMMAND)

linux: clean ## Cross compile for linux.
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_COMMAND)

gzip: ## Compress current compiled binary.
	gzip bin/$(BINARY_NAME)
	mv bin/$(BINARY_NAME).gz bin/$(BINARY_NAME)-$(GIT_COMMIT)-$(UNAME)-amd64.gz

release: build gzip ## Full release process.

unit: ## Run unit tests.
	go test -mod=vendor -cover -race -short ./... -v

lint: ## See https://github.com/golangci/golangci-lint#install for install instructions
	golangci-lint run ./...

.PHONY: help all deps clean build gzip release unit lint