package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const keyValuePattern = `(.+)=(.+)`

type getValue func(string) string

func parseIntoStruct(i interface{}, data getValue) error {
	typ := reflect.TypeOf(i)

	if typ.Kind() != reflect.Ptr {
		return errors.New("config: value passed instead of reference")
	}

	if typ.Elem().Kind() != reflect.Struct {
		return errors.New("config: non struct passed")
	}

	if err := parse(typ, reflect.ValueOf(i), data, ""); err != nil {
		return fmt.Errorf("config: %s", err)
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

func vars(input []byte) map[string]string {
	r := regexp.MustCompile(keyValuePattern)
	vars := make(map[string]string)

	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		m := r.FindSubmatch(line)
		if len(m) > 2 {
			vars[string(bytes.TrimSpace(m[1]))] = string(bytes.TrimSpace(m[2]))
		}
	}

	return vars
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
