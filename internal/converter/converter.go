package converter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	tag             = "env"
	defaultValueTag = "default"
)

type getValue = func(string) string

func ConvertIntoStruct[T any](i T, data getValue) error {
	typ := reflect.TypeOf(i)

	if typ.Kind() != reflect.Ptr {
		return errors.New("value passed instead of reference")
	}

	if typ.Elem().Kind() != reflect.Struct {
		return errors.New("non struct passed")
	}

	if reflect.ValueOf(i).IsNil() {
		return errors.New("nil struct passed")
	}

	if err := parse(typ, reflect.ValueOf(i), data, ""); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func parse(typ reflect.Type, val reflect.Value, getValue getValue, path string) error {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		path += field.Name + "_"

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if field.Type.Kind() == reflect.Struct {
			if err := parse(field.Type, val.Field(i), getValue, path); err != nil {
				return err
			}
			path = ""
		}

		value := value(&field, getValue, path)

		if value == "" {
			value = defaultValue(&field)
		}

		if value == "" {
			continue
		}
		path = ""

		fieldValue := val.Field(i)
		if err := setFieldValue(&field, fieldValue, value); err != nil {
			return err
		}
	}

	return nil
}

func value(field *reflect.StructField, getValue getValue, path string) string {
	value := getValue(key(field, path))
	if value == "" {
		value = getValue(field.Name)
	}
	return value
}

func key(field *reflect.StructField, path string) string {
	key := field.Tag.Get(tag)
	if key == "" {
		key = strings.ToUpper(strings.TrimSuffix(path, "_"))
	}
	return key
}

func defaultValue(field *reflect.StructField) string {
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
