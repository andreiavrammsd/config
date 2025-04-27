package config

import (
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/interpolater"
	"github.com/andreiavrammsd/config/internal/parser"
	"github.com/andreiavrammsd/config/internal/reader"
)

type Config struct{}

func TestEnvFileWithParserErrorAtEnvFile(t *testing.T) {
	actual := Config{}

	loader := &Loader[Config]{
		configStruct: actual,
		dotEnvFile:   ".env",
		parse:        func(_ io.Reader, _ map[string]string) error { return errors.New("parser error with env file") },
		read:         reader.ReadToStruct[Config],
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
		configStruct: actual,
		dotEnvFile:   ".env",
		parse:        func(_ io.Reader, _ map[string]string) error { return errors.New("parser error with bytes") },
		read:         reader.ReadToStruct[Config],
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
		configStruct: actual,
		dotEnvFile:   ".env",
		parse:        parser.New().Parse,
		read: func(_ Config, _ func(string) string) error {
			return errors.New("converter error")
		},
		interpolate: interpolater.New().Interpolate,
	}

	err := loader.EnvFile("testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: converter error" {
		t.Fatal("incorrect error message:", err)
	}
}
