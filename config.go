// Package `config` parses configuration values into given struct.
//
// Requirements for configuration struct:
// - A non-nil pointer to the struct must be passed.
// - Fields must be exported. Unexported fields will be ignored.
// - A field can have the `env` tag which defines the key of the value. If no tag provided, the key will be
// the uppercase full path of the field (all the fields names starting root until current field, joined by underscore).
// - The `json` tag will be used for parsing from JSON.
//
// Input sources:
// - environment variables
// - environment variables from files
// - byte array
// - json
package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/andreiavrammsd/config/internal/interpolator"
	"github.com/andreiavrammsd/config/internal/parser"
	"github.com/andreiavrammsd/config/internal/reader"
)

var ErrInvalidConfigType = errors.New("config type must be non-nil pointer to struct")

const dotEnvFile string = ".env"

// Config exposes the public API.
type Config struct {
	parse       func(r io.Reader, vars map[string]string) error
	interpolate func(map[string]string)
	read        func(configStruct any, data reader.ValueReader) error
}

// FromFile parses config into struct from one or multiple dotenv files.
// If no file given, uses .env by default.
func (c Config) FromFile(config any, files ...string) error {
	if err := validateConfigType(config); err != nil {
		return err
	}

	if len(files) == 0 {
		files = []string{dotEnvFile}
	}

	vars := make(map[string]string)

	for i := range files {
		file, err := os.Open(files[i])
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if err = c.parse(file, vars); err != nil {
			file.Close()
			return fmt.Errorf("%w", err)
		}

		file.Close()
	}

	c.interpolate(vars)

	if err := c.read(config, func(s *string) string { return vars[*s] }); err != nil {
		return err
	}

	return nil
}

// FromEnv parses config into struct from environment variables.
func (c Config) FromEnv(config any) error {
	if err := validateConfigType(config); err != nil {
		return err
	}

	return c.read(config, func(s *string) string { return os.Getenv(*s) })
}

// FromBytes parses config into struct from byte array.
func (c Config) FromBytes(config any, input []byte) error {
	if err := validateConfigType(config); err != nil {
		return err
	}

	vars := make(map[string]string)

	if err := c.parse(bytes.NewReader(input), vars); err != nil {
		return err
	}

	c.interpolate(vars)

	return c.read(config, func(s *string) string { return vars[*s] })
}

// FromJSON parses config into struct from json.
func (c Config) FromJSON(config any, input json.RawMessage) error {
	if err := validateConfigType(config); err != nil {
		return err
	}

	if err := json.Unmarshal(input, config); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// New creates the config package instance.
func New() Config {
	return Config{
		parse:       parser.New().Parse,
		interpolate: interpolator.New().Interpolate,
		read:        reader.ReadToStruct,
	}
}

func validateConfigType(config any) error {
	rv := reflect.ValueOf(config)

	if rv.Kind() != reflect.Pointer || rv.IsNil() || rv.Elem().Kind() != reflect.Struct {
		return ErrInvalidConfigType
	}

	return nil
}
