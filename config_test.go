package config_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/andreiavrammsd/config"
)

type Struct struct {
	Field string
}

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

type envFile struct {
	AAA string
	Config
}

const testdata_file string = "testdata/.env"

func TestEnvFile(t *testing.T) {
	// Temporarily switch to testdata directory to read .env by default
	cwd, _ := os.Getwd()
	os.Chdir("testdata")
	defer func() {
		os.Chdir(cwd)
	}()

	_, expected, err := testdata(".env")
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := config.Load(&actual).EnvFile(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestEnvFileWithCustomEnvFiles(t *testing.T) {
	_, ex, err := testdata(testdata_file)
	if err != nil {
		t.Fatal(err)
	}

	expected := envFile{
		AAA:    "BBB",
		Config: ex,
	}

	actual := envFile{}
	if err := config.Load(&actual).EnvFile("testdata/.env", "testdata/.env2"); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestEnvFileWithOneFileWhichIsMissing(t *testing.T) {
	actual := envFile{}
	err := config.Load(&actual).EnvFile("somefile")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: open somefile: no such file or directory" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestEnvFileWithMultipleFilesOneMissing(t *testing.T) {
	actual := envFile{}
	err := config.Load(&actual).EnvFile("testdata/.env", "someotherfile", "testdata/.env2")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: open someotherfile: no such file or directory" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestBytes(t *testing.T) {
	input, expected, err := testdata(testdata_file)
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := config.Load(&actual).Bytes(input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestString(t *testing.T) {
	input, expected, err := testdata(testdata_file)
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := config.Load(&actual).String(string(input)); err != nil {
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
	   },
	   "Interpolated":"$B env_1 $ $B \\3 6379 + $",
	   "Default":"default value"
	}`)

	_, expected, err := testdata(testdata_file)
	if err != nil {
		t.Fatal(err)
	}

	actual := Config{}
	if err := config.Load(&actual).JSON(input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestJsonWithInvalidInput(t *testing.T) {
	input := json.RawMessage(`invalid json`)

	err := config.Load(&Config{}).JSON(input)
	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: invalid character 'i' looking for beginning of value" {
		t.Fatal("incorrect error message:", err)
	}
}

func testdata(file string) ([]byte, Config, error) {
	input, err := os.ReadFile(file)
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
