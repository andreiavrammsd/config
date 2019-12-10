.PHONY: all test bench lint coverage prepushhook

COVER_PROFILE=cover.out
GO111MODULE=on

all: test lint

test:
	go test -cover -v ./...

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

lint: check-lint
	golint -set_exit_status ./...
	golangci-lint run

coverage:
	go test -v -coverprofile $(COVER_PROFILE) ./... && go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push

check-lint:
	@[ $(shell which golint) ] || (GO111MODULE=off && go get -u golang.org/x/lint/golint)
	@[ $(shell which golangci-lint) ] || curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
		| sh -s -- -b $(shell go env GOPATH)/bin v1.21.0
