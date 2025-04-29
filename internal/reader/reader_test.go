package reader_test

import (
	"testing"

	"github.com/andreiavrammsd/config/internal/reader"
)

type config struct {
	S        string `env:"MyString"`
	SEmpty   string
	SDefault string `default:"default value"`

	I8      int8 `env:"integer_8"`
	I16     int16
	I32     int32
	I64     int64
	Integer int

	UI8             uint8
	UI16            uint16
	UI32            uint32
	UI64            uint64
	UnsignedInteger uint

	F32 float32
	F64 float64

	B bool

	Bytes []byte

	Struct struct {
		Integer int
	}
}

func readValue(s *string) string {
	vars := make(map[string]string)
	vars["MyString"] = "string"
	vars["SDefault"] = ""

	vars["integer_8"] = "-8"
	vars["I16"] = "-16"
	vars["I32"] = "-32"
	vars["I64"] = "-64"
	vars["Integer"] = "-999"

	vars["UI8"] = "8"
	vars["UI16"] = "16"
	vars["UI32"] = "32"
	vars["UI64"] = "64"
	vars["UNSIGNEDINTEGER"] = "999"

	vars["F32"] = "32.2345225"
	vars["F64"] = "-64.2342623678"

	vars["B"] = "true"

	vars["Bytes"] = "key=value"

	vars["STRUCT_INTEGER"] = "123"

	return vars[*s]
}

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	if actual != expected {
		t.Fatalf("%v != %v", actual, expected)
	}
}

func TestReadToStruct(t *testing.T) {
	configStruct := config{}

	err := reader.ReadToStruct(&configStruct, readValue)
	if err != nil {
		t.Fatal("error not expected")
	}

	assertEqual(t, configStruct.S, "string")
	assertEqual(t, configStruct.SEmpty, "")
	assertEqual(t, configStruct.SDefault, "default value")

	assertEqual(t, configStruct.I8, -8)
	assertEqual(t, configStruct.I16, -16)
	assertEqual(t, configStruct.I32, -32)
	assertEqual(t, configStruct.I64, -64)
	assertEqual(t, configStruct.Integer, -999)

	assertEqual(t, configStruct.UI8, 8)
	assertEqual(t, configStruct.UI16, 16)
	assertEqual(t, configStruct.UI32, 32)
	assertEqual(t, configStruct.UI64, 64)
	assertEqual(t, configStruct.UnsignedInteger, 999)

	assertEqual(t, configStruct.F32, 32.2345225)
	assertEqual(t, configStruct.F64, -64.2342623678)

	assertEqual(t, configStruct.B, true)

	assertEqual(t, string(configStruct.Bytes), "key=value")

	assertEqual(t, configStruct.Struct.Integer, 123)
}

func TestReadToStructWithIntParseError(t *testing.T) {
	configStruct := struct{ Value int }{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid int value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Value (strconv.ParseInt: parsing \"invalid int value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestReadToStructWithUintParseError(t *testing.T) {
	configStruct := struct{ Value uint }{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid uint value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Value (strconv.ParseUint: parsing \"invalid uint value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestReadToStructWithFloat32ParseError(t *testing.T) {
	configStruct := struct{ Value float32 }{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid float32 value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Value (strconv.ParseFloat: parsing \"invalid float32 value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestReadToStructWithFloat64ParseError(t *testing.T) {
	configStruct := struct{ Value float64 }{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid float64 value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Value (strconv.ParseFloat: parsing \"invalid float64 value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestReadToStructWithBoolParseError(t *testing.T) {
	configStruct := struct{ Value bool }{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid bool value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Value (strconv.ParseBool: parsing \"invalid bool value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestReadToStructWithInnerStructParseError(t *testing.T) {
	configStruct := struct {
		Struct struct {
			Integer int
		}
	}{}

	readValue := func(s *string) string {
		vars := make(map[string]string)
		vars["STRUCT_INTEGER"] = "invalid struct integer value"
		return vars[*s]
	}

	err := reader.ReadToStruct(&configStruct, readValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "field Integer (strconv.ParseInt: parsing \"invalid struct integer value\": invalid syntax)" {
		t.Fatal("incorrect error message:", err)
	}
}
