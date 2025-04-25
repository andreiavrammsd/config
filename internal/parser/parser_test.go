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
	err := parser.Parse(reader, vars)

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
	err := parser.Parse(&eofReader{"a=b", 0}, vars)

	if err != nil {
		t.Error("expected no error")
	}

	assertEqual(t, vars["a"], "b")
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	err = errors.New("reader error")
	return
}

func TestParseWithReaderError(t *testing.T) {
	vars := make(map[string]string)
	err := parser.Parse(&errReader{}, vars)

	if len(vars) > 0 {
		t.Error("expected empty map")
	}

	if err == nil {
		t.Error("expected reader error")
	}

	if err.Error() != "config: cannot read from input (reader error)" {
		t.Fatal("incorrect error message:", err)
	}
}

// Benchmark_Parse-8        1934143               606.9 ns/op          4096 B/op          1 allocs/op
func Benchmark_Parse(b *testing.B) {
	reader := bytes.NewReader([]byte(environment))
	vars := make(map[string]string)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err := parser.Parse(reader, vars)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestInterpolate(t *testing.T) {
	vars := make(map[string]string)
	vars["TIMEOUT"] = "2000000000"
	vars["ABC"] = " string\\\" "
	vars["A"] = "1"
	vars["C"] = "xx"
	vars["E_NEG"] = "-1"
	vars["F32"] = "15425.2231"
	vars["F64"] = "245232212.9844448"
	vars["IsSet"] = "true"
	vars["REDIS_CONNECTION_HOST"] = " localhost "
	vars["REDIS_PORT"] = "6379"
	vars["STRUCTPTR_FIELD"] = "Val\\\"ue "
	vars["MONGO_DATABASE_HOST"] = "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb"
	vars["MONGO_DATABASE_COLLECTION_NAME"] = "us=ers"
	vars["MONGO_OTHER"] = "$A"
	vars["INTERPOLATED"] = "\\$B env_$A $ \\$B \\\\$C ${REDIS_PORT} + $"

	parser.Interpolate(vars)

	assertEqual(t, vars["TIMEOUT"], "2000000000")
	assertEqual(t, vars["ABC"], " string\\\" ")
	assertEqual(t, vars["A"], "1")
	assertEqual(t, vars["E_NEG"], "-1")
	assertEqual(t, vars["F32"], "15425.2231")
	assertEqual(t, vars["F64"], "245232212.9844448")
	assertEqual(t, vars["IsSet"], "true")
	assertEqual(t, vars["REDIS_CONNECTION_HOST"], " localhost ")
	assertEqual(t, vars["REDIS_PORT"], "6379")
	assertEqual(t, vars["STRUCTPTR_FIELD"], "Val\\\"ue ")
	assertEqual(t, vars["MONGO_DATABASE_HOST"], "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb")
	assertEqual(t, vars["MONGO_DATABASE_COLLECTION_NAME"], "us=ers")
	assertEqual(t, vars["MONGO_OTHER"], "1")
	assertEqual(t, vars["INTERPOLATED"], "$B env_1 $ $B \\xx 6379 + $")
}

// Benchmark_Interpolate-8          2369497               495.0 ns/op            80 B/op          4 allocs/op
func Benchmark_Interpolate(b *testing.B) {
	vars := make(map[string]string)
	vars["TIMEOUT"] = "2000000000"
	vars["ABC"] = " string\\\" "
	vars["A"] = "1"
	vars["C"] = "xx"
	vars["E_NEG"] = "-1"
	vars["F32"] = "15425.2231"
	vars["F64"] = "245232212.9844448"
	vars["IsSet"] = "true"
	vars["REDIS_CONNECTION_HOST"] = " localhost "
	vars["REDIS_PORT"] = "6379"
	vars["STRUCTPTR_FIELD"] = "Val\\\"ue "
	vars["MONGO_DATABASE_HOST"] = "mongodb://user:pass==@host.tld:955/?ssl=true&replicaSet=globaldb"
	vars["MONGO_DATABASE_COLLECTION_NAME"] = "us=ers"
	vars["MONGO_OTHER"] = "$A"
	vars["INTERPOLATED"] = "\\$B env_$A $ \\$B \\\\$C ${REDIS_PORT} + $"

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		parser.Interpolate(vars)
	}
}
