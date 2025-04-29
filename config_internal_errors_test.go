package config

import (
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/interpolator"
	"github.com/andreiavrammsd/config/internal/parser"
)

type Configuration struct{}

func TestFromFileWithParserError(t *testing.T) {
	config := &Config{
		parse: func(_ io.Reader, _ map[string]string) error { return errors.New("parser error") },
	}

	actual := Configuration{}
	err := config.FromFile(&actual, "testdata/.env")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "parser error" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromFileWithReaderError(t *testing.T) {
	loader := &Config{
		parse:       parser.New().Parse,
		interpolate: interpolator.New().Interpolate,
		read: func(_ any, _ func(*string) string) error {
			return errors.New("reader error")
		},
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

func TestFromBytesWithParserError(t *testing.T) {
	config := &Config{
		parse: func(_ io.Reader, _ map[string]string) error { return errors.New("parser error") },
	}

	actual := Configuration{}
	err := config.FromBytes(&actual, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "parser error" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromBytesWithReaderError(t *testing.T) {
	config := &Config{
		parse:       parser.New().Parse,
		interpolate: interpolator.New().Interpolate,
		read: func(_ any, _ func(*string) string) error {
			return errors.New("reader error")
		},
	}

	actual := Configuration{}
	err := config.FromBytes(&actual, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "reader error" {
		t.Fatal("incorrect error message:", err)
	}
}
