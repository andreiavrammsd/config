.PHONY: all test bench lint coverage precommithook

GOLANGCI_LINT_VERSION=2.1.5

all: test lint

test:
	go test -cover -v ./...

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

lint: check-lint
	@golangci-lint run || (golangci-lint fmt && exit 1)

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt

precommithook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

check-lint:
	@if ! golangci-lint version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION)"; then \
		echo "Installing golangci-lint v$(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
			| sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION); \
	fi
