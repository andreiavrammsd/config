.PHONY: all test bench lint coverage prepushhook

COVER_PROFILE=cover.out

all: test lint

GOLINT := $(shell which golint)

test:
	go test -cover -v ./...

bench:
	go test -bench=. -benchmem -v -run=Bench ./...

lint:
ifndef GOLINT
		go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./...

	@[ ! -f ./bin/golangci-lint ] && curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh \
		| sh -s -- -b ./bin v1.21.0 || true
	./bin/golangci-lint run

coverage:
	go test -v -coverprofile $(COVER_PROFILE) ./... && go tool cover -html=$(COVER_PROFILE)

prepushhook:
	echo '#!/bin/sh\n\nmake' > .git/hooks/pre-push && chmod +x .git/hooks/pre-push
