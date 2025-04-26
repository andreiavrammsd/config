package converter_test

import (
	"testing"

	"github.com/andreiavrammsd/config/internal/converter"
)

type config struct {
	S        string
	SEmpty   string
	SDefault string `default:"default value"`

	I8      int8
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

func getValue(s string) string {
	vars := make(map[string]string)
	vars["S"] = "string"
	vars["SDefault"] = ""

	vars["I8"] = "-8"
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

	return vars[s]
}

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	if actual != expected {
		t.Fatalf("%v != %v", actual, expected)
	}
}

func TestConvertIntoStruct(t *testing.T) {
	i := config{}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err != nil {
		t.Fatal("error not expected")
	}

	assertEqual(t, i.S, "string")
	assertEqual(t, i.SEmpty, "")
	assertEqual(t, i.SDefault, "default value")

	assertEqual(t, i.I8, -8)
	assertEqual(t, i.I32, -32)
	assertEqual(t, i.I64, -64)
	assertEqual(t, i.Integer, -999)

	assertEqual(t, i.UI8, 8)
	assertEqual(t, i.UI32, 32)
	assertEqual(t, i.UI64, 64)
	assertEqual(t, i.UnsignedInteger, 999)

	assertEqual(t, i.F32, 32.2345225)
	assertEqual(t, i.F64, -64.2342623678)

	assertEqual(t, i.B, true)

	assertEqual(t, string(i.Bytes), "key=value")

	assertEqual(t, i.Struct.Integer, 123)
}

func TestConvertIntoStructWithValue(t *testing.T) {
	i := struct{}{}

	err := converter.ConvertIntoStruct(i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: value passed instead of reference" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithNonStruct(t *testing.T) {
	var i *int = nil

	err := converter.ConvertIntoStruct(i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: non struct passed" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithNilStruct(t *testing.T) {
	var i *struct{} = nil

	err := converter.ConvertIntoStruct(i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: nil struct passed" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithIntParseError(t *testing.T) {
	i := struct{ Value int }{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid int value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseInt: parsing \"invalid int value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithUintParseError(t *testing.T) {
	i := struct{ Value uint }{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid uint value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseUint: parsing \"invalid uint value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithFloat32ParseError(t *testing.T) {
	i := struct{ Value float32 }{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid float32 value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseFloat: parsing \"invalid float32 value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithFloat64ParseError(t *testing.T) {
	i := struct{ Value float64 }{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid float64 value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseFloat: parsing \"invalid float64 value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithBoolParseError(t *testing.T) {
	i := struct{ Value bool }{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["VALUE"] = "invalid bool value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseBool: parsing \"invalid bool value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}

func TestConvertIntoStructWithInnerStructParseError(t *testing.T) {
	i := struct {
		Struct struct {
			Integer int
		}
	}{}

	getValue := func(s string) string {
		vars := make(map[string]string)
		vars["STRUCT_INTEGER"] = "invalid struct integer value"
		return vars[s]
	}

	err := converter.ConvertIntoStruct(&i, getValue)

	if err == nil {
		t.Fatal("error expected")
	}

	if err.Error() != "config: strconv.ParseInt: parsing \"invalid struct integer value\": invalid syntax" {
		t.Fatal("incorrect error message:", err)
	}
}
