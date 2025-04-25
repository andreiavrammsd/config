package parser_test

import (
	"bytes"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/andreiavrammsd/config/internal/parser"
)

type Config struct {
	Mongo struct {
		Database struct {
			Host       string `env:"MONGO_DATABASE_HOST"`
			Collection struct {
				Name  []byte `env:"MONGO_DATABASE_COLLECTION_NAME"`
				Other byte   `env:"MONGO_OTHER"`
				X     rune   `env:"MONGO_X"`
			}
		}
	}
	Redis struct {
		Connection struct {
			Host string
			Port int `env:"REDIS_PORT"`
		}
	}
	String       string `env:"ABC" default:"ignored"`
	Struct       Struct
	StructPtr    *Struct
	D            int64
	E            int
	ENeg         int `env:"E_NEG"`
	UD           uint64
	UE           uint
	F64          float64
	Timeout      time.Duration
	C            int32
	UC           uint32
	F32          float32
	B            int16
	UB           uint16
	A            int8
	UA           uint8
	IsSet        bool
	Interpolated string
	Default      string `default:"default value"`
}

type Struct struct {
	Field string
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	err = errors.New("reader error")
	return
}

func TestWithParseReaderError(t *testing.T) {
	kv := make(map[string]string)
	err := parser.ParseVars(&errReader{}, kv)
	if len(kv) > 0 {
		t.Error("expected empty map")
	}
	if err == nil {
		t.Error("expected reader error")
	}
}

// Benchmark_parseVars-8            1663723               749 ns/op            4096 B/op          1 allocs/op
func Benchmark_parseVars(b *testing.B) {
	b.ReportAllocs()

	input, _, err := testdata()
	if err != nil {
		b.Fatal(err)
	}

	vars := make(map[string]string)
	reader := bytes.NewReader(input)
	for n := 0; n < b.N; n++ {
		err := parser.ParseVars(reader, vars)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func testdata() ([]byte, Config, error) {
	input, err := os.ReadFile("../../testdata/.env")
	if err != nil {
		return nil, Config{}, err
	}

	expected := Config{
		String: " string\\\" ",
		A:      1,
		B:      2,
		C:      3,
		D:      4,
		E:      5,
		ENeg:   -1,
		UA:     1,
		UB:     2,
		UC:     3,
		UD:     4,
		UE:     5,
		F32:    15425.2231,
		F64:    245232212.9844448,
		IsSet:  true,
		Redis: struct {
			Connection struct {
				Host string
				Port int `env:"REDIS_PORT"`
			}
		}{
			Connection: struct {
				Host string
				Port int `env:"REDIS_PORT"`
			}{
				Host: " localhost ",
				Port: 6379,
			},
		},
		Timeout: time.Second * 2,
		Struct: Struct{
			Field: "Value",
		},
		Mongo: struct {
			Database struct {
				Host       string `env:"MONGO_DATABASE_HOST"`
				Collection struct {
					Name  []byte `env:"MONGO_DATABASE_COLLECTION_NAME"`
					Other byte   `env:"MONGO_OTHER"`
					X     rune   `env:"MONGO_X"`
				}
			}
		}{Database: struct {
			Host       string `env:"MONGO_DATABASE_HOST"`
			Collection struct {
				Name  []byte `env:"MONGO_DATABASE_COLLECTION_NAME"`
				Other byte   `env:"MONGO_OTHER"`
				X     rune   `env:"MONGO_X"`
			}
		}{
			Host: "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb",
			Collection: struct {
				Name  []byte `env:"MONGO_DATABASE_COLLECTION_NAME"`
				Other byte   `env:"MONGO_OTHER"`
				X     rune   `env:"MONGO_X"`
			}{
				Name:  []byte("us=ers"),
				Other: 1,
				X:     'a',
			},
		}},
		Interpolated: "$B env_1 $ $B \\3 6379 + $",
		Default:      "default value",
	}

	return input, expected, nil
}
