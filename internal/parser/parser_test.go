package parser_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/andreiavrammsd/config/internal/parser"
)

const environment string = ` # key=value
TIMEOUT=2000000000
ABC =" string\" "
A =1
  B  =2
C=3# key=value
D =4
E=5
E_NEG=-1
UA=1  # key=value
UB=2
# comment
UC=30

UD=40
UE=50
F32=15425.2231
F64=245232212.9844448
IsSet=true
REDIS_CONNECTION_HOST=" localhost "
REDIS_PORT=6379
STRUCT_FIELD=Value
STRUCTPTR_FIELD="Val\"ue "
MONGO_DATABASE_HOST="mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb" # db connection
MONGO_DATABASE_COLLECTION_NAME='us=ers'
MONGO_OTHER=$A
MONGO_X=97
# comment
INTERPOLATED="\$B env_$A $ \$B \\$C ${REDIS_PORT} + $"

`

func assertNotExist(t *testing.T, key string, vars map[string]string) {
	if _, ok := vars[key]; ok {
		t.Fatalf("%s not expected", key)
	}
}

func assertEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
	}
}

func TestParse(t *testing.T) {
	reader := bytes.NewReader([]byte(environment))
	vars := make(map[string]string)
	err := parser.New().Parse(reader, vars)

	if err != nil {
		t.Error("expected no error")
	}

	assertNotExist(t, "key", vars)
	assertEqual(t, vars["TIMEOUT"], "2000000000")
	assertEqual(t, vars["ABC"], " string\\\" ")
	assertEqual(t, vars["A"], "1")
	assertEqual(t, vars["B"], "2")
	assertEqual(t, vars["C"], "3")
	assertNotExist(t, "KEY3", vars)
	assertEqual(t, vars["D"], "4")
	assertEqual(t, vars["E"], "5")
	assertEqual(t, vars["E_NEG"], "-1")
	assertEqual(t, vars["UA"], "1")
	assertNotExist(t, "KEY", vars)
	assertEqual(t, vars["UB"], "2")
	assertNotExist(t, "comment", vars)
	assertEqual(t, vars["UC"], "30")
	assertEqual(t, vars["UD"], "40")
	assertEqual(t, vars["UE"], "50")
	assertEqual(t, vars["F32"], "15425.2231")
	assertEqual(t, vars["F64"], "245232212.9844448")
	assertEqual(t, vars["IsSet"], "true")
	assertEqual(t, vars["REDIS_CONNECTION_HOST"], " localhost ")
	assertEqual(t, vars["REDIS_PORT"], "6379")
	assertEqual(t, vars["STRUCT_FIELD"], "Value")
	assertEqual(t, vars["STRUCTPTR_FIELD"], "Val\\\"ue ")
	assertEqual(t, vars["MONGO_DATABASE_HOST"], "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb")
	assertEqual(t, vars["MONGO_DATABASE_COLLECTION_NAME"], "us=ers")
	assertEqual(t, vars["MONGO_OTHER"], "$A")
	assertEqual(t, vars["MONGO_X"], "97")
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

type errReader struct {
}

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
	reader := bytes.NewReader([]byte(environment))
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
