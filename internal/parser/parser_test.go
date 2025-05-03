package parser_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"reflect"
	"runtime"
	"testing"
	"unicode/utf8"

	"github.com/andreiavrammsd/config/internal/parser"
)

func assertVar(t *testing.T, vars map[string]string, key, expectedValue string) {
	t.Helper()

	_, file, line, _ := runtime.Caller(1)

	if value, ok := vars[key]; !ok {
		t.Fatalf("%s:%d: Key %q not found", file, line, key)
	} else if value != expectedValue {
		t.Fatalf("%s:%d: %q != %q for key %q", file, line, value, expectedValue, key)
	}
}

func TestParse(t *testing.T) {
	reader := bytes.NewReader(testdata("testdata/.env"))
	vars := make(map[string]string)

	err := parser.New().Parse(reader, vars)
	if err != nil {
		t.Error("expected no error")
	}

	expectedNumberOfVars := 67 // IS THIS OK?
	if len(vars) != expectedNumberOfVars {
		t.Fatalf("Expected %d vars, got %d", expectedNumberOfVars, len(vars))
	}

	assertVar(t, vars, "A", "1")
	assertVar(t, vars, "B", "$A")
	assertVar(t, vars, "BB", "CC")
	assertVar(t, vars, "VAR_WITH_COMMENT", "val with comment")
	assertVar(t, vars, "D", "")
	assertVar(t, vars, "D2", "")
	assertVar(t, vars, "D3", "")
	assertVar(t, vars, "E", "some value with spaces")
	assertVar(t, vars, "F", "another value with spaces")
	assertVar(t, vars, "MONGO_DATABASE_HOST", "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb")
	assertVar(t, vars, "MONGO_DATABASE_COLLECTION_NAME", "us=ers")
	assertVar(t, vars, "G", "quote 'inside' quote")
	assertVar(t, vars, "H", "quote \"inside\" quote")
	assertVar(t, vars, "I", "line1\\nline2")
	assertVar(t, vars, "J", "tab\\tseparated")
	assertVar(t, vars, "ABC", " string\\\" ")
	assertVar(t, vars, "K", "Emoji 🚀 and Unicode ü")
	assertVar(t, vars, "L", "spaced_key")
	assertVar(t, vars, "M", "spaced_value")
	assertVar(t, vars, "N", "spaced_both")
	assertVar(t, vars, "NUM", "-1")
	assertVar(t, vars, "NOT_NUM", "---1")
	assertVar(t, vars, "POS_NUM", "+1")
	assertVar(t, vars, "POS_NOT_NUM", "++1")
	// FAILS: assertVar(t, vars, "O", "#notacomment")
	assertVar(t, vars, "O2", "")
	assertVar(t, vars, "P", "key=value=another")
	assertVar(t, vars, "Q", "$UNDEFINED_VAR")
	assertVar(t, vars, "R", "$A-$B-$C")
	assertVar(t, vars, "$SPECIAL", "weird")
	assertVar(t, vars, "1NUMBER", "bad")
	assertVar(t, vars, "S", "whitespace_before_key")
	assertVar(t, vars, "T", "trailing_space")
	assertVar(t, vars, "U", "lots_of_space")
	assertVar(t, vars, "V", "first=second=third")
	assertVar(t, vars, "W", "\\uZZZZ")
	assertVar(t, vars, "X1", "true")
	assertVar(t, vars, "X2", "False")
	assertVar(t, vars, "X3", "0")
	assertVar(t, vars, "X4", "1")
	assertVar(
		t,
		vars,
		"BIG",
		"Lorem_ipsum_dolor_sit_amet_consectetur_adipiscing_elit_sed_do_eiusmod_tempor_incididunt_ut_labore_et_dolore_magna_aliqua",
	)
	assertVar(
		t,
		vars,
		"LONG",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)
	assertVar(t, vars, "Z1", "12345")
	assertVar(t, vars, "Z2", "0")
	assertVar(t, vars, "Z3", "-999")
	assertVar(t, vars, "TIMEOUT", "2000000000")
	assertVar(t, vars, "F32", "15425.2231")
	assertVar(t, vars, "F64", "245232212.9844448")
	assertVar(t, vars, "AA.key", "subvalue")
	assertVar(t, vars, "BB-key", "another_subvalue")
	assertVar(t, vars, "CC___DD", "weird_key")
	assertVar(t, vars, "EE", "[this looks like json]")
	assertVar(t, vars, "EE2", "[this looks like json]")
	assertVar(t, vars, "EE3", "[this looks like json]")
	// IS THIS OK? assertVar(t, vars, "EE4", "[this looks like json]")
	// IS THIS OK? assertVar(t, vars, "EE5", "[this looks like json]")
	// FAILS: assertVar(t, vars, "FF", "{ \"name\": \"John\", \"age\": 30 }")
	assertVar(t, vars, "ARRAY", "one,two,three")
	assertVar(t, vars, "EMPTY1", "")
	assertVar(t, vars, "EMPTY2", "")
	assertVar(t, vars, "NUM_STRING", "12345")
	// FAILS: assertVar(t, vars, "BROKEN_NEWLINE", "this is\nstill valid because quotes stay open")
	// FAILS: assertVar(t, vars, "BROKEN_NEWLINE_SINGLE_QUOTES", "this is\nstill valid because quotes stay open")
	assertVar(t, vars, "XX", "second")
	assertVar(t, vars, "INTERPOLATED", "\\$B env_$A $ \\$B \\\\$C ${REDIS_PORT} + $")
}

type eofReader struct {
	content string
	atChar  int
}

func (e *eofReader) Read(p []byte) (int, error) {
	if e.atChar >= len(e.content) {
		return 0, io.EOF
	}

	n := copy(p, e.content)
	e.atChar += n

	return n, nil
}

func TestParseWithEOF(t *testing.T) {
	vars := make(map[string]string)
	err := parser.New().Parse(&eofReader{"a=b", 0}, vars)
	if err != nil {
		t.Error("expected no error")
	}

	assertVar(t, vars, "a", "b")
}

type errReader struct{}

func (e *errReader) Read(_ []byte) (n int, err error) {
	err = errors.New("reader error")
	return
}

func TestParseWithReaderError(t *testing.T) {
	vars := make(map[string]string)
	err := parser.New().Parse(&errReader{}, vars)

	if len(vars) > 0 {
		t.Error("expected empty map")
	}

	if err == nil {
		t.Error("expected reader error")
	}

	if err.Error() != "reader error" {
		t.Fatal("incorrect error message:", err)
	}
}

// Benchmark_Parse-8        1623301               709.7 ns/op          4192 B/op          2 allocs/op.
func Benchmark_Parse(b *testing.B) {
	benchParser := parser.New()
	reader := bytes.NewReader(testdata("testdata/.env"))
	vars := make(map[string]string)

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		err := benchParser.Parse(reader, vars)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func FuzzParse(f *testing.F) {
	testcases := []string{string(testdata("testdata/.env")), "", " "}
	for _, tc := range testcases {
		f.Add(tc)
	}

	fuzzParser := parser.New()

	f.Fuzz(func(t *testing.T, input string) {
		varsFirst := make(map[string]string)
		if err := fuzzParser.Parse(bytes.NewReader([]byte(input)), varsFirst); err != nil {
			t.Error(err)
		}

		varsSecond := make(map[string]string)
		if err := fuzzParser.Parse(bytes.NewReader([]byte(input)), varsSecond); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(varsFirst, varsSecond) {
			t.Error("Different results")
		}

		for key, firstValue := range varsFirst {
			secondValue := varsSecond[key]

			if firstValue != secondValue {
				t.Errorf("Before: %q, after: %q", firstValue, secondValue)
			}
			if !utf8.ValidString(key) {
				t.Errorf("Parse produced invalid UTF-8 string for key: %q", key)
			}
			if !utf8.ValidString(firstValue) {
				t.Errorf("Parse produced invalid UTF-8 string for value in first map: %q", firstValue)
			}
			if !utf8.ValidString(secondValue) {
				t.Errorf("Parse produced invalid UTF-8 string for value in second map %q", secondValue)
			}
		}
	})
}

func testdata(file string) []byte {
	input, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return input
}
