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
	"errors"
	"fmt"
	"os"
)

const tag = "env"
const dotEnvFile = ".env"

// Loader provides methods to load configuration values into a struct
type Loader struct {
	i interface{}
}

// Load creates a Loader with given struct
func Load(i interface{}) *Loader {
	return &Loader{i: i}
}

// Env loads config into struct from environment variables
func (l *Loader) Env() error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	return parseIntoStruct(l.i, os.Getenv)
}

// EnvFile loads config into struct from environment variables in one or multiple files (dotenv).
// If no file is passed, the default is ".env".
func (l *Loader) EnvFile(files ...string) error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	if len(files) == 0 {
		files = append(files, dotEnvFile)
	}

	vars := make(map[string]string)

	for i := 0; i < len(files); i++ {
		f, err := os.Open(files[i])
		if err != nil {
			return fmt.Errorf("config: %s", err)
		}

		err = parseVars(f, vars)

		if err != nil {
			if e := f.Close(); e != nil {
				return fmt.Errorf("config: %s", e)
			}
			return fmt.Errorf("config: %s", err)
		}

		if err = f.Close(); err != nil {
			return fmt.Errorf("config: %s", err)
		}
	}

	f := func(s string) string {
		return vars[s]
	}

	return parseIntoStruct(l.i, f)
}

// Bytes loads config into struct from byte array
func (l *Loader) Bytes(input []byte) error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	return fromBytes(l.i, input)
}

// String loads config into struct from a string
func (l *Loader) String(input string) error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	return fromBytes(l.i, []byte(input))
}

// JSON loads config into struct from json
func (l *Loader) JSON(input json.RawMessage) error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	return json.Unmarshal(input, l.i)
}

func checkNilStruct(i interface{}) error {
	if i == nil {
		return errors.New("config: nil struct passed")
	}

	return nil
}

func fromBytes(i interface{}, input []byte) error {
	vars := make(map[string]string)

	if err := parseVars(bytes.NewReader(input), vars); err != nil {
		return err
	}

	f := func(s string) string {
		return vars[s]
	}

	return parseIntoStruct(i, f)
}
