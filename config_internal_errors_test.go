package config

import (
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/interpolator"
	"github.com/andreiavrammsd/config/internal/parser"
	"github.com/andreiavrammsd/config/internal/reader"
)

type Configuration struct{}

func TestFromFileWithParserErrorAtEnvFile(t *testing.T) {
	config := &Config{
		parse: func(_ io.Reader, _ map[string]string) error { return errors.New("parser error with file") },
		read:  reader.ReadToStruct,
	}

	actual := Configuration{}
	err := config.FromFile(&actual, "testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "parser error with file" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromFileWithParserErrorBytes(t *testing.T) {
	config := &Config{
		parse: func(_ io.Reader, _ map[string]string) error { return errors.New("parser error with bytes") },
		read:  reader.ReadToStruct,
	}

	actual := Configuration{}
	err := config.FromBytes(&actual, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "parser error with bytes" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromFileWithReaderError(t *testing.T) {
	loader := &Config{
		parse: parser.New().Parse,
		read: func(_ any, _ func(string) string) error {
			return errors.New("reader error")
		},
		interpolate: interpolator.New().Interpolate,
	}

	actual := Configuration{}
	err := loader.FromFile(&actual, "testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "reader error" {
		t.Fatal("incorrect error message:", err)
	}
}
