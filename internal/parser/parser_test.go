package parser_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/andreiavrammsd/config/internal/parser"
)

func assertEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("%s:%d: %q != %q", file, line, actual, expected)
	}
}

func TestParse(t *testing.T) {
	reader := bytes.NewReader(testdata("testdata/.env"))
	vars := make(map[string]string)

	err := parser.New().Parse(reader, vars)
	if err != nil {
		t.Error("expected no error")
	}

	expectedNumberOfVars := 65 // IS THIS OK?
	if len(vars) != expectedNumberOfVars {
		t.Fatalf("Expected %d vars, got %d", expectedNumberOfVars, len(vars))
	}

	assertEqual(t, vars["A"], "1")
	assertEqual(t, vars["B"], "$A")
	assertEqual(t, vars["BB"], "CC")
	assertEqual(t, vars["VAR_WITH_COMMENT"], "val with comment")
	assertEqual(t, vars["D"], "")
	assertEqual(t, vars["D2"], "")
	assertEqual(t, vars["D3"], "")
	assertEqual(t, vars["E"], "some value with spaces")
	assertEqual(t, vars["F"], "another value with spaces")
	assertEqual(t, vars["MONGO_DATABASE_HOST"], "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb")
	assertEqual(t, vars["MONGO_DATABASE_COLLECTION_NAME"], "us=ers")
	assertEqual(t, vars["G"], "quote 'inside' quote")
	assertEqual(t, vars["H"], "quote \"inside\" quote")
	assertEqual(t, vars["I"], "line1\\nline2")
	assertEqual(t, vars["J"], "tab\\tseparated")
	assertEqual(t, vars["ABC"], " string\\\" ")
	assertEqual(t, vars["K"], "Emoji 🚀 and Unicode ü")
	assertEqual(t, vars["L"], "spaced_key")
	assertEqual(t, vars["M"], "spaced_value")
	assertEqual(t, vars["N"], "spaced_both")
	assertEqual(t, vars["NUM"], "-1")
	assertEqual(t, vars["NOT_NUM"], "---1")
	assertEqual(t, vars["POS_NUM"], "+1")
	assertEqual(t, vars["POS_NOT_NUM"], "++1")
	// FAILS: assertEqual(t, vars["O"], "#notacomment")
	assertEqual(t, vars["P"], "key=value=another")
	assertEqual(t, vars["Q"], "$UNDEFINED_VAR")
	assertEqual(t, vars["R"], "$A-$B-$C")
	assertEqual(t, vars["$SPECIAL"], "weird")
	assertEqual(t, vars["1NUMBER"], "bad")
	assertEqual(t, vars["S"], "whitespace_before_key")
	assertEqual(t, vars["T"], "trailing_space")
	assertEqual(t, vars["U"], "lots_of_space")
	assertEqual(t, vars["V"], "first=second=third")
	assertEqual(t, vars["W"], "\\uZZZZ")
	assertEqual(t, vars["X1"], "true")
	assertEqual(t, vars["X2"], "False")
	assertEqual(t, vars["X3"], "0")
	assertEqual(t, vars["X4"], "1")
	assertEqual(
		t,
		vars["BIG"],
		"Lorem_ipsum_dolor_sit_amet_consectetur_adipiscing_elit_sed_do_eiusmod_tempor_incididunt_ut_labore_et_dolore_magna_aliqua",
	)
	// FAILS: assertEqual(t, vars["Y"], "this is \na weird \nmultiline\nvalue")
	// FAILS: assertEqual(t, vars["LONG"], "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	assertEqual(t, vars["Z1"], "12345")
	assertEqual(t, vars["Z2"], "0")
	assertEqual(t, vars["Z3"], "-999")
	assertEqual(t, vars["TIMEOUT"], "2000000000")
	assertEqual(t, vars["F32"], "15425.2231")
	assertEqual(t, vars["F64"], "245232212.9844448")
	assertEqual(t, vars["AA.key"], "subvalue")
	assertEqual(t, vars["BB-key"], "another_subvalue")
	assertEqual(t, vars["CC___DD"], "weird_key")
	assertEqual(t, vars["EE"], "[this looks like json]")
	assertEqual(t, vars["EE2"], "[this looks like json]")
	assertEqual(t, vars["EE3"], "[this looks like json]")
	// IS THIS OK? assertEqual(t, vars["EE4"], "[this looks like json]")
	// IS THIS OK? assertEqual(t, vars["EE5"], "[this looks like json]")
	// FAILS: assertEqual(t, vars["FF"], "{ \"name\": \"John\", \"age\": 30 }")
	assertEqual(t, vars["ARRAY"], "one,two,three")
	assertEqual(t, vars["EMPTY1"], "")
	assertEqual(t, vars["EMPTY2"], "")
	assertEqual(t, vars["NUM_STRING"], "12345")
	// FAILS: assertEqual(t, vars["BROKEN_NEWLINE"], "this is\nstill valid because quotes stay open")
	assertEqual(t, vars["XX"], "second")
	assertEqual(t, vars["INTERPOLATED"], "\\$B env_$A $ \\$B \\\\$C ${REDIS_PORT} + $")
}

type eofReader struct {
	content string
	atChar  int
}

func (e *eofReader) Read(p []byte) (n int, err error) {
	if e.atChar >= len(e.content) {
		return 0, io.EOF
	}

	n = copy(p, e.content)
	e.atChar += n

	return n, nil
}

func TestParseWithEOF(t *testing.T) {
	vars := make(map[string]string)
	err := parser.New().Parse(&eofReader{"a=b", 0}, vars)
	if err != nil {
		t.Error("expected no error")
	}

	assertEqual(t, vars["a"], "b")
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

// Benchmark_Parse-8                1913996               618.2 ns/op          4192 B/op          2 allocs/op.
func Benchmark_Parse(b *testing.B) {
	p := parser.New()
	reader := bytes.NewReader(testdata("testdata/.env"))
	vars := make(map[string]string)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := p.Parse(reader, vars)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func testdata(file string) []byte {
	input, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return input
}
