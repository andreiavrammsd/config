# Config

[![codecov](https://codecov.io/github/andreiavrammsd/config/branch/master/graph/badge.svg?token=4BV8YNIIIX)](https://app.codecov.io/github/andreiavrammsd/config)

Package config load configuration values into given struct.

The struct must be passed by reference.

Fields must be exported. Unexported fields will be ignored. They can have the `env` tag which defines the key
of the value. If no tag provided, the key will be the uppercase full path of the field (all the fields names
starting the root until current field, joined by underscore).

The `json` tag will be used for loading from JSON.

```go
package main

import (
	"fmt"
	"log"

	"github.com/andreiavrammsd/config"
)

type Config struct {
	Username string `env:"CUSTOM_USERNAME_TAG"`
	Tag      string `default:"none"`
}

func main() {
	input := []byte(`CUSTOM_USERNAME_TAG=msd # username`)

	cfg := Config{}
	if err := config.Load(&cfg).Bytes(input); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(cfg.Username)
	fmt.Println(cfg.Tag)
}
```

## Docs

[![GoDoc](https://godoc.org/github.com/andreiavrammsd/config?status.svg)](https://godoc.org/github.com/andreiavrammsd/config)

## Install

```bash
go get github.com/andreiavrammsd/config
```

## Usage

See [examples](./examples_test.go) and [tests](./config_test.go).

## Testing and QA tools for development

See [Makefile](./Makefile).
