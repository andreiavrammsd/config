package interpolator_test

import (
	"testing"

	"github.com/andreiavrammsd/config/internal/interpolator"
)

func assertEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Fatalf("%s != %s", actual, expected)
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

	interpolator.New().Interpolate(vars)

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

// Benchmark_Interpolate-8          2369497               495.0 ns/op            80 B/op          4 allocs/op.
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

	interpolator := interpolator.New()

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		interpolator.Interpolate(vars)
	}
}
