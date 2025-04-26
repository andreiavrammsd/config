package converter

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	tag             = "env"
	defaultValueTag = "default"
)

type ReadValue = func(string) string

// ConvertIntoStruct take a pointer to a struct and, for each property in the struct (recursively),
// generates a key that it passes to the given readValue function which must return the value for the property.
func ConvertIntoStruct[T any](configStruct T, readValue ReadValue) error {
	typ := reflect.TypeOf(configStruct)

	if typ.Kind() != reflect.Ptr {
		return errors.New("value passed instead of reference")
	}

	if typ.Elem().Kind() != reflect.Struct {
		return errors.New("non struct passed")
	}

	if reflect.ValueOf(configStruct).IsNil() {
		return errors.New("nil struct passed")
	}

	return parse(typ, reflect.ValueOf(configStruct), readValue, "")
}

func parse(typ reflect.Type, val reflect.Value, readValue ReadValue, path string) error {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		path += field.Name + "_"

		// Parse struct recursively.
		if field.Type.Kind() == reflect.Struct {
			if err := parse(field.Type, val.Field(i), readValue, path); err != nil {
				return err
			}
			path = ""
		}

		value := getValue(&field, readValue, path)
		if value == "" {
			continue
		}

		if err := setFieldValue(&field, val.Field(i), value); err != nil {
			return err
		}

		// After value is set, start again with a new property.
		path = ""
	}

	return nil
}

func getValue(field *reflect.StructField, readValue ReadValue, path string) (value string) {
	// Generate key and read value.
	key := generateKey(field, path)
	value = readValue(key)

	// If empty, read value from field name.
	if value == "" {
		value = readValue(field.Name)
	}

	// If empty, get default.
	if value == "" {
		value = getDefaultValue(field)
	}

	return
}

func generateKey(field *reflect.StructField, path string) (key string) {
	// Get configured key.
	key = field.Tag.Get(tag)

	// If empty, generate from path (path is property name or struct name + property name).
	if key == "" {
		key = strings.ToUpper(strings.TrimSuffix(path, "_"))
	}

	return
}

func getDefaultValue(field *reflect.StructField) string {
	return field.Tag.Get(defaultValueTag)
}

func setFieldValue(field *reflect.StructField, fieldValue reflect.Value, value string) error {
	switch field.Type.Kind() {
	case reflect.String:
		if fieldValue.CanSet() {
			fieldValue.SetString(value)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return err
		}
		fieldValue.SetInt(v)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return err
		}
		fieldValue.SetUint(v)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(v)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fieldValue.SetBool(v)
	case reflect.Slice:
		fieldValue.SetBytes([]byte(value))
	}

	return nil
}
