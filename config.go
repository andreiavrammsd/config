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
	"io"
	"os"

	"github.com/andreiavrammsd/config/internal/interpolater"
	"github.com/andreiavrammsd/config/internal/parser"
	"github.com/andreiavrammsd/config/internal/reader"
)

// Loader provides methods to load configuration values into a struct.
type Loader[T any] struct {
	configStruct T
	dotEnvFile   string
	parse        func(r io.Reader, vars map[string]string) error
	read         func(configStruct T, data func(string) string) error
	interpolate  func(map[string]string)
}

// Load creates a Loader with given struct.
func Load[T any](config T) *Loader[T] {
	return &Loader[T]{
		configStruct: config,
		dotEnvFile:   ".env",
		parse:        parser.New().Parse,
		read:         reader.ReadToStruct[T],
		interpolate:  interpolater.New().Interpolate,
	}
}

// Env loads config into struct from environment variables.
func (l *Loader[T]) Env() error {
	return l.read(l.configStruct, os.Getenv)
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

		if err = l.parse(file, vars); err != nil {
			file.Close()
			return fmt.Errorf("config: %w", err)
		}

		file.Close()
	}

	l.interpolate(vars)

	if err := l.read(l.configStruct, func(s string) string { return vars[s] }); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return nil

}

// Bytes loads config into struct from byte array.
func (l *Loader[T]) Bytes(input []byte) error {
	return l.fromBytes(input)
}

// String loads config into struct from a string.
func (l *Loader[T]) String(input string) error {
	return l.fromBytes([]byte(input))
}

// JSON loads config into struct from json.
func (l *Loader[T]) JSON(input json.RawMessage) error {
	if err := json.Unmarshal(input, l.configStruct); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	return nil
}

func (l *Loader[T]) fromBytes(input []byte) error {
	vars := make(map[string]string)

	if err := l.parse(bytes.NewReader(input), vars); err != nil {
		return fmt.Errorf("config: %w", err)
	}

	l.interpolate(vars)

	return l.read(l.configStruct, func(s string) string { return vars[s] })
}
