// Package config load configuration values into given struct.
//
// The struct must be passed by reference.
//
// Fields must be exported. Unexported fields will be ignored. They can have the `env` tag which defines the key
// of the value. If no tag provided, the key will be the uppercase full path of the field (all the fields names
// starting the root until current field, joined by underscore).
//
// The `json` tag will be used for loading from JSON.
package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/andreiavrammsd/config/internal/converter"
	"github.com/andreiavrammsd/config/internal/parser"
)

// Loader provides methods to load configuration values into a struct
type Loader[T any] struct {
	i          T
	dotEnvFile string
}

// Load creates a Loader with given struct
func Load[T any](config T) *Loader[T] {
	return &Loader[T]{i: config, dotEnvFile: ".env"}
}

// Env loads config into struct from environment variables
func (l *Loader[T]) Env() error {
	return converter.ConvertIntoStruct(l.i, os.Getenv)
}

// EnvFile loads config into struct from environment variables in one or multiple files (dotenv).
// If no file is passed, the default is ".env".
func (l *Loader[T]) EnvFile(files ...string) error {
	if len(files) == 0 {
		files = []string{l.dotEnvFile}
	}

	vars := make(map[string]string)

	for i := 0; i < len(files); i++ {
		file, err := os.Open(files[i])
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		if err = parser.Parse(file, vars); err != nil {
			file.Close()
			return fmt.Errorf("config: %w", err)
		}

		file.Close()
	}

	parser.Interpolate(vars)

	if err := converter.ConvertIntoStruct(l.i, func(s string) string { return vars[s] }); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return nil
}

// Bytes loads config into struct from byte array
func (l *Loader[T]) Bytes(input []byte) error {
	return fromBytes(l.i, input)
}

// String loads config into struct from a string
func (l *Loader[T]) String(input string) error {
	return fromBytes(l.i, []byte(input))
}

// JSON loads config into struct from json
func (l *Loader[T]) JSON(input json.RawMessage) error {
	if err := json.Unmarshal(input, l.i); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return nil
}

func fromBytes[T any](i T, input []byte) error {
	vars := make(map[string]string)

	if err := parser.Parse(bytes.NewReader(input), vars); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	parser.Interpolate(vars)

	return converter.ConvertIntoStruct(i, func(s string) string { return vars[s] })
}
