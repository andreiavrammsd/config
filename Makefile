.PHONY: all test bench lint coverage prepushhook

COVER_PROFILE=coverage.txt
GOLANGCI_LINT_VERSION=2.1.5

all: test lint

test:
	go test -cover -v ./...

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

lint: check-lint
	golangci-lint run

coverage:
	go test -v -coverprofile=$(COVER_PROFILE) -covermode=atomic ./...

coverage-report: coverage
	go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push

check-lint:
	@if ! golangci-lint version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION)"; then \
		echo "Installing golangci-lint v$(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
			| sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION); \
	fi
