package config_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/andreiavrammsd/config"
	"github.com/andreiavrammsd/config/testdata"
)

const testdataFile string = "testdata/.env"

func TestFromFileWithDefaultFile(t *testing.T) {
	// Temporarily switch to testdata directory to read .env by default
	cwd, _ := os.Getwd() // nolint:errcheck
	os.Chdir("testdata") // nolint:errcheck
	defer func() {
		os.Chdir(cwd) // nolint:errcheck
	}()

	expected := testdata.GetExpectedOutput()

	actual := testdata.Config{}
	if err := config.New().FromFile(&actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestFromFileWithCustomFiles(t *testing.T) {
	expected := testdata.EnvFile{
		AAA:    "BBB",
		Config: testdata.GetExpectedOutput(),
	}

	actual := testdata.EnvFile{}
	if err := config.New().FromFile(&actual, "testdata/.env", "testdata/.env2"); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestFromFileWithOneMissingFile(t *testing.T) {
	actual := testdata.EnvFile{}
	err := config.New().FromFile(&actual, "somefile")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "open somefile: no such file or directory" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromFileWithMultipleFilesWhenOneIsMissing(t *testing.T) {
	actual := testdata.EnvFile{}
	err := config.New().FromFile(&actual, "testdata/.env", "someotherfile", "testdata/.env2")

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "open someotherfile: no such file or directory" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromBytes(t *testing.T) {
	input := testdata.ReadInputFile(testdataFile)
	expected := testdata.GetExpectedOutput()

	actual := testdata.Config{}
	if err := config.New().FromBytes(&actual, input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestFromBytesWithNilInput(t *testing.T) {
	err := config.New().FromBytes(&struct{}{}, nil)
	if err != nil {
		t.Fatal("error not expected")
	}
}

func TestFromJSON(t *testing.T) {
	jsonString := testdata.ReadInputFile("testdata/env.json")
	expected := testdata.GetExpectedOutput()

	var input json.RawMessage
	if err := json.Unmarshal(jsonString, &input); err != nil {
		panic(err)
	}

	actual := testdata.Config{}
	if err := config.New().FromJSON(&actual, input); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nhave: %v\nwant: %v", actual, expected)
	}
}

func TestFromJSONWithNilInput(t *testing.T) {
	err := config.New().FromJSON(&struct{}{}, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "unexpected end of JSON input" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromJSONWithInvalidInput(t *testing.T) {
	input := json.RawMessage(`invalid json`)

	err := config.New().FromJSON(&struct{}{}, input)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "invalid character 'i' looking for beginning of value" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestWithNilConfigType(t *testing.T) {
	err := config.New().FromFile(nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err != config.ErrNilPointerInput {
		t.Fatal("incorrect error message:", err)
	}
}

func TestWithValueConfigType(t *testing.T) {
	err := config.New().FromEnv(struct{}{})

	if err == nil {
		t.Fatal("error expected")
	}

	if err != config.ErrValueInput {
		t.Fatal("incorrect error message:", err)
	}
}

func TestWithNonStructConfigType(t *testing.T) {
	var i int
	err := config.New().FromBytes(&i, nil) // nil bytes?

	if err == nil {
		t.Fatal("error expected")
	}

	if err != config.ErrNonStructInput {
		t.Fatal("incorrect error message:", err)
	}
}

func TestFromJSONWithInvalidConfigType(t *testing.T) {
	var i int
	err := config.New().FromJSON(&i, nil)

	if err == nil {
		t.Fatal("error expected")
	}

	if err != config.ErrNonStructInput {
		t.Fatal("incorrect error message:", err)
	}
}
