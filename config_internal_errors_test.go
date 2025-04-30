package config

import (
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/interpolator"
	"github.com/andreiavrammsd/config/internal/parser"
	"github.com/andreiavrammsd/config/internal/reader"
)

func TestFromFileWithParserError(t *testing.T) {
	config := &Config{
		parse: func(_ io.Reader, _ map[string]string) error { return errors.New("parser error") },
	}

	err := config.FromFile(&struct{}{}, "testdata/.env")

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
		read: func(_ any, _ reader.ValueReader) error {
			return errors.New("reader error")
		},
	}

	err := loader.FromFile(&struct{}{}, "testdata/.env")

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

	err := config.FromBytes(&struct{}{}, nil)

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
		read: func(_ any, _ reader.ValueReader) error {
			return errors.New("reader error")
		},
	}

	err := config.FromBytes(&struct{}{}, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "reader error" {
		t.Fatal("incorrect error message:", err)
	}
}
