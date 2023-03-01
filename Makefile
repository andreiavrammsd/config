.PHONY: all test bench lint coverage prepushhook

COVER_PROFILE=coverage.txt
GO111MODULE=on

all: test lint

test:
	go test -cover -v ./...

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

lint: check-lint
	golangci-lint run

coverage:
	go test -v -coverprofile=$(COVER_PROFILE) -covermode=atomic ./... && go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push

check-lint:
	@[ $(shell which golangci-lint) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.51.2
