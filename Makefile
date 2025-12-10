GO111MODULE=on
GO ?= go
GOTEST = go test -v -bench\=.
WORKDIR ?= $(shell pwd)

.PHONY: install
install:
	$(GO) install -ldflags="-s -w" -tags netgo

.PHONY: lint
lint:
	golangci-lint run --timeout=5m

.PHONY: test
test:
	$(GO) test -v ./...
