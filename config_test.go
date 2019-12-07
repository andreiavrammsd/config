package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
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
	String    string `env:"ABC"`
	Struct    Struct
	StructPtr *Struct
	D         int64
	E         int
	ENeg      int `env:"E_NEG"`
	UD        uint64
	UE        uint
	F64       float64
	Timeout   time.Duration
	C         int32
	UC        uint32
	F32       float32
	B         int16
	UB        uint16
	A         int8
	UA        uint8
	IsSet     bool
}

type Struct struct {
	Field string
}

func TestEnv(t *testing.T) {
	input, expected, err := testdata()
	if err != nil {
		t.Fatal(err)
	}

	vars := make(map[string]string)
	err = parseVars(bytes.NewReader(input), vars)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			t.Fatalf(`cannot set env variable "%s" with value "%s": "%s"`, k, v, err)
		}
	}

	actual := Config{}
	if err := Load(&actual).Env(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}

	for k := range vars {
		if err := os.Unsetenv(k); err != nil {
			t.Fatal(err)
		}
	}
}

type envFile struct {
	AAA string
	Config
}

func TestEnvFile(t *testing.T) {
	_, ex, err := testdata()
	if err != nil {
		t.Fatal(err)
	}
	file := "testdata/.env"

	expected := envFile{
		AAA:    "BBB",
		Config: ex,
	}

	actual := envFile{}
	if err := Load(&actual).EnvFile(file, "testdata/.env2"); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestBytes(t *testing.T) {
	input, expected, err := testdata()
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := Load(&actual).Bytes(input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestString(t *testing.T) {
	input, expected, err := testdata()
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := Load(&actual).String(string(input)); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestJson(t *testing.T) {
	input := json.RawMessage(`{
	   "StructPtr":null,
	   "String":" string\\\" ",
	   "A":1,
	   "B":2,
	   "C":3,
	   "D":4,
	   "E":5,
	   "ENeg":-1,
	   "UA":1,
	   "UB":2,
	   "UC":3,
	   "UD":4,
	   "UE":5,
	   "F32":15425.2231,
	   "F64":245232212.9844448,
	   "IsSet":true,
	   "Redis":{
		  "Connection":{
			 "Host":" localhost ",
			 "Port":6379
		  }
	   },
	   "Timeout":2000000000,
	   "Mongo":{
		  "Database":{
			 "Host":"mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb",
			 "Collection":{
				"Name":"dXM9ZXJz",
				"Other":1,
				"X":97
			 }
		  }
	   },
	   "Struct":{
		  "Field":"Value"
	   }
	}`)

	_, expected, err := testdata()
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := Load(&actual).JSON(input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

const parseErrorInput = "STRUCT_KEY=text"

func TestWithIntParseError(t *testing.T) {
	config := struct {
		Struct struct {
			Key int
		}
	}{}

	if err := Load(&config).String(parseErrorInput); err == nil {
		t.Error("expected parse error")
	}
}

func TestWithUintParseError(t *testing.T) {
	config := struct {
		Struct struct {
			Key uint
		}
	}{}

	if err := Load(&config).String(parseErrorInput); err == nil {
		t.Error("expected parse error")
	}
}

func TestWithFloat32ParseError(t *testing.T) {
	config := struct {
		Struct struct {
			Key float32
		}
	}{}

	if err := Load(&config).String(parseErrorInput); err == nil {
		t.Error("expected parse error")
	}
}

func TestWithFloat64ParseError(t *testing.T) {
	config := struct {
		Struct struct {
			Key float64
		}
	}{}

	if err := Load(&config).String(parseErrorInput); err == nil {
		t.Error("expected parse error")
	}
}

func TestWithBoolParseError(t *testing.T) {
	config := struct {
		Struct struct {
			Key bool
		}
	}{}

	if err := Load(&config).String(parseErrorInput); err == nil {
		t.Error("expected parse error")
	}
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	err = errors.New("reader error")
	return
}

func TestWithParseReaderError(t *testing.T) {
	kv := make(map[string]string)
	err := parseVars(&errReader{}, kv)
	if len(kv) > 0 {
		t.Error("expected empty map")
	}
	if err == nil {
		t.Error("expected reader error")
	}
}

// BenchmarkVars-8          1663723               749 ns/op            4096 B/op          1 allocs/op
func BenchmarkVars(b *testing.B) {
	b.ReportAllocs()
	input, _, err := testdata()
	if err != nil {
		b.Fatal(err)
	}

	vars := make(map[string]string)
	reader := bytes.NewReader(input)
	for n := 0; n < b.N; n++ {
		err := parseVars(reader, vars)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func testdata() ([]byte, Config, error) {
	input, err := ioutil.ReadFile("testdata/.env")
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
	}

	return input, expected, nil
}

func TestWithNilStructPassed(t *testing.T) {
	tests := []func() error{
		func() error {
			return Load(nil).Env()
		},
		func() error {
			return Load(nil).EnvFile()
		},
		func() error {
			return Load(nil).Bytes(nil)
		},
		func() error {
			return Load(nil).String("")
		},
		func() error {
			return Load(nil).JSON(nil)
		},
	}

	for _, tt := range tests {
		if tt() == nil {
			t.Fatal("expected error")
		}
	}
}

func TestWithStructPassedByValue(t *testing.T) {
	cfg := Config{}
	tests := []func() error{
		func() error {
			return Load(cfg).Env()
		},
		func() error {
			return Load(cfg).EnvFile()
		},
		func() error {
			return Load(cfg).Bytes(nil)
		},
		func() error {
			return Load(cfg).String("")
		},
		func() error {
			return Load(cfg).JSON(nil)
		},
	}

	for _, test := range tests {
		if test() == nil {
			t.Fatal("expected error")
		}
	}
}

func TestWithNonStructPassed(t *testing.T) {
	cfg := 1
	tests := []func() error{
		func() error {
			return Load(&cfg).Env()
		},
		func() error {
			return Load(&cfg).EnvFile()
		},
		func() error {
			return Load(&cfg).Bytes(nil)
		},
		func() error {
			return Load(&cfg).String("")
		},
		func() error {
			return Load(&cfg).JSON(nil)
		},
	}

	for _, test := range tests {
		if test() == nil {
			t.Fatal("expected error")
		}
	}
}
