package config

import (
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/converter"
	"github.com/andreiavrammsd/config/internal/parser"
)

type Config struct{}

func TestEnvFileWithParserErrorAtEnvFile(t *testing.T) {
	actual := Config{}

	loader := &Loader[Config]{
		i:          actual,
		dotEnvFile: ".env",
		parse:      func(r io.Reader, vars map[string]string) error { return errors.New("parser error with env file") },
		convert:    converter.ConvertIntoStruct[Config],
	}

	err := loader.EnvFile("testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: parser error with env file" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestEnvFileWithParserErrorBytes(t *testing.T) {
	actual := Config{}

	loader := &Loader[Config]{
		i:          actual,
		dotEnvFile: ".env",
		parse:      func(r io.Reader, vars map[string]string) error { return errors.New("parser error with bytes") },
		convert:    converter.ConvertIntoStruct[Config],
	}

	err := loader.Bytes(nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: parser error with bytes" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestEnvFileWithConverterError(t *testing.T) {
	actual := Config{}

	loader := &Loader[Config]{
		i:          actual,
		dotEnvFile: ".env",
		parse:      parser.Parse,
		convert: func(i Config, data func(string) string) error {
			return errors.New("converter error")
		},
	}

	err := loader.EnvFile("testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: converter error" {
		t.Fatal("incorrect error message:", err)
	}
}
