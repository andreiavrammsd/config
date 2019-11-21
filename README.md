# Config

Package config load configuration values into given struct.

The struct must be passed by reference.

Fields can have the `env` tag which defines the key of the value. If no tag provided, the key will be the
uppercase full path of the field (all the fields names starting the root until current field, joined by underscore).

The `json` tag will be used for loading from JSON.

## Docs

[![GoDoc](https://godoc.org/github.com/andreiavrammsd/config?status.svg)](https://godoc.org/github.com/andreiavrammsd/config)

## Install

```bash
go get github.com/andreiavrammsd/config
```

## Usage

See [tests](./config_test.go).

## Testing and QA tools for development

See [Makefile](./Makefile).
