# Config

[![build](https://github.com/andreiavrammsd/config/workflows/CI/badge.svg)](https://github.com/andreiavrammsd/config/actions/workflows/ci.yml) [![codecov](https://codecov.io/github/andreiavrammsd/config/branch/master/graph/badge.svg?token=4BV8YNIIIX)](https://app.codecov.io/github/andreiavrammsd/config) [![GoDoc](https://godoc.org/github.com/andreiavrammsd/config?status.svg)](https://godoc.org/github.com/andreiavrammsd/config)

Package `config` parses configuration values into given struct.

Requirements for configuration struct:
- A non-nil pointer to the struct must be passed.
- Fields must be exported. Unexported fields will be ignored.
- A field can have the `env` tag which defines the key of the value. If no tag provided, the key will be the uppercase full path of the field (all the fields names starting root until current field, joined by underscore).
- The `json` tag will be used for parsing from JSON.

Input sources:
- environment variables
- environment variables from files
- byte array
- json

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
	if err := config.New().FromBytes(&cfg, input); err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	fmt.Println(cfg.Username)
	fmt.Println(cfg.Tag)
}
```

## Install

```bash
go get github.com/andreiavrammsd/config
```

## Usage

See [examples](./examples_test.go) and [tests](./config_test.go).

## Testing and QA tools for development

See [Makefile](./Makefile) and [VS Code setup](.vscode).

## Known issues

- In some cases, quoted values are not parsed correctly.
```
MULTILINE_QUOTED="this is a multiline
quoted
value"

QUOTED_INCLUDING_QUOTES="{ \"name\": \"John\", \"age\": 30 }"
```
