package config

import (
	"io/ioutil"
)

func fromEnvFile(i interface{}, files ...string) error {
	var input []byte
	for i := 0; i < len(files); i++ {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			return err
		}
		input = append(input, data...)
	}

	return fromBytes(i, input)
}

func fromBytes(i interface{}, input []byte) error {
	f := func() getValue {
		vars := vars(input)

		return func(s string) string {
			return vars[s]
		}
	}

	return parseIntoStruct(i, f())
}
