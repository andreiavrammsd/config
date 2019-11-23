package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func fromEnvFile(i interface{}, files ...string) error {
	var input []byte
	for i := 0; i < len(files); i++ {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			return fmt.Errorf("config: %s", err)
		}
		input = append(input, data...)
	}

	return fromBytes(i, input)
}

func fromBytes(i interface{}, input []byte) error {
	v, err := vars(bytes.NewReader(input))
	if err != nil {
		return err
	}

	f := func(s string) string {
		return v[s]
	}

	return parseIntoStruct(i, f)
}
