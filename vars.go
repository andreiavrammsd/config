package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

func parseVars(r io.Reader, vars map[string]string) error {
	reader := bufio.NewReader(r)

	var name, value []byte

	atName := true
	atValue := false
	atComment := false

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if atValue {
					vars[varName(name)] = varValue(value)
				}
				break
			}

			return fmt.Errorf("config: cannot read from input (%s)", err)
		}

		if r == '#' {
			if atValue {
				vars[varName(name)] = varValue(value)
			}

			name = nil
			value = nil
			atName = false
			atValue = false
			atComment = true
			continue
		}

		if r == '\n' || r == '\r' {
			if atValue {
				vars[varName(name)] = varValue(value)
			}

			name = nil
			value = nil
			atName = true
			atValue = false

			if atComment {
				atComment = false
				continue
			}

			continue
		}

		if atComment {
			continue
		}

		if r == '=' {
			if atValue {
				value = append(value, byte(r))
			}
			atName = false
			atValue = true
			continue
		}

		if atName {
			if unicode.IsSpace(r) {
				continue
			}
			name = append(name, byte(r))
			continue
		}

		if atValue {
			value = append(value, byte(r))
		}
	}

	return nil
}

func varName(n []byte) string {
	return string(bytes.TrimSpace(n))
}

func varValue(v []byte) string {
	return string(bytes.Trim(bytes.TrimSpace(v), `"'`))
}
