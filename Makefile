GOLANGCI_LINT_VERSION=2.1.5

all: test lint

test:
	go test -vet=all -run=Test ./...

lint: install-lint
	@golangci-lint config verify
	@golangci-lint run || (golangci-lint fmt && exit 1)
	@govulncheck ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

precommithook:
	@echo '#!/bin/sh\n\nmake' > .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit

install-lint:
	@if ! golangci-lint version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION)"; then \
		echo "Installing golangci-lint v$(GOLANGCI_LINT_VERSION)..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
			| sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION); \
	fi
	@if ! govulncheck --version 2>/dev/null; then \
		echo "Installing govulncheck..."; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
