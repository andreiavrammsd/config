package config

import (
	"encoding/json"
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
	input, expected := data(t)

	kv := vars(input)
	for k, v := range kv {
		if err := os.Setenv(k, v); err != nil {
			t.Fatal(err)
		}
	}

	actual := Config{}
	if err := Load(&actual).Env(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}

	for k := range kv {
		if err := os.Unsetenv(k); err != nil {
			t.Fatal(err)
		}
	}
}

func TestEnvFile(t *testing.T) {
	_, expected := data(t)
	file := "testdata/.env"

	actual := Config{}
	if err := Load(&actual).EnvFile(file); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestBytes(t *testing.T) {
	input, expected := data(t)

	actual := Config{}
	if err := Load(&actual).Bytes(input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestString(t *testing.T) {
	input, expected := data(t)

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
	   "String":"string",
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
			 "Host":"localhost",
			 "Port":6379
		  }
	   },
	   "Timeout":2000000000,
	   "Mongo":{ 
		  "Database":{ 
			 "Host":"127.0.0.1",
			 "Collection":{ 
				"Name":"dXNlcnM=",
				"Other":1,
				"X":97
			 }
		  }
	   },
	   "Struct":{ 
		  "Field":"Value"
	   }
	}`)

	_, expected := data(t)

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

func data(t *testing.T) ([]byte, Config) {
	input, err := ioutil.ReadFile("testdata/.env")
	if err != nil {
		t.Fatal(err)
	}

	expected := Config{
		String: "string",
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
				Host: "localhost",
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
			Host: "127.0.0.1",
			Collection: struct {
				Name  []byte `env:"MONGO_DATABASE_COLLECTION_NAME"`
				Other byte   `env:"MONGO_OTHER"`
				X     rune   `env:"MONGO_X"`
			}{
				Name:  []byte("users"),
				Other: 1,
				X:     'a',
			},
		}},
	}

	return input, expected
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
