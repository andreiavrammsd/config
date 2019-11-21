package config

import (
	"io/ioutil"
	"os"
)

func fromEnv(i interface{}) error {
	return parseIntoStruct(i, os.Getenv)
}

func fromEnvFile(i interface{}, files ...string) error {
	input := make([]byte, 0)
	for i := 0; i < len(files); i++ {
		data, err := ioutil.ReadFile(files[i])
		if err != nil {
			return err
		}
		input = append(input, data...)
	}

	f := func() getValue {
		vars := vars(input)

		return func(s string) string {
			return vars[s]
		}
	}

	return parseIntoStruct(i, f())
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
