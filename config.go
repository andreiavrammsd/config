// Package config load configuration values into given struct.
// The struct must be passed by reference.
// Fields can have the `env` tag which defines the key of the value. If no tag provided, the key will be the
// uppercase full path of the field (all the fields names starting the root until current field, joined by underscore).
// The `json` tag will be used for loading from JSON.
package config

import (
	"encoding/json"
	"errors"
)

const tag = "env"
const dotEnvFile = ".env"

// Loader provides methods to load configuration values into a struct
type Loader struct {
	i interface{}
}

// Env loads config into struct from environment variables
func (l *Loader) Env() error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	return fromEnv(l.i)
}

// EnvFile loads config into struct from environment variables in files (dotenv).
func (l *Loader) EnvFile(files ...string) error {
	if err := checkNilStruct(l.i); err != nil {
		return err
	}

	if len(files) == 0 {
		files = append(files, dotEnvFile)
	}

	return fromEnvFile(l.i, files...)
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
		return errors.New("nil struct passed")
	}

	return nil
}

// Load creates a Loader with given struct
func Load(i interface{}) *Loader {
	return &Loader{i: i}
}
