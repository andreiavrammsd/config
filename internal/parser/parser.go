package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func ParseVars(r io.Reader, vars map[string]string) error {
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
					vars[string(name)] = varValue(value)
				}
				break
			}

			return fmt.Errorf("config: cannot read from input (%s)", err)
		}

		if r == '#' {
			if atValue {
				vars[string(name)] = varValue(value)
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
				vars[string(name)] = varValue(value)
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

func varValue(v []byte) string {
	return string(bytes.Trim(bytes.TrimSpace(v), `"'`))
}

func InterpolateVars(vars map[string]string) {
	for k, v := range vars {
		if strings.IndexByte(v, '$') == -1 {
			continue
		}

		atVar := false
		var name []byte
		var newValue []byte

		for i := 0; i < len(v); i++ {
			// Variable starts
			if v[i] == '$' {
				atVar = isAtVar(v, i)

				if i == len(v)-1 && i-1 >= 0 && v[i-1] != '\\' {
					newValue = append(newValue, v[i])
				}

				if atVar {
					continue
				}
			}

			if !atVar {
				if nextVarIsDoubleEscaped(v, i) {
					newValue = append(newValue, v[i])
					continue
				}

				if nextVarIsEscaped(v, i) {
					continue
				}

				newValue = append(newValue, v[i])
				continue
			}

			if atVar && (v[i] == '{' || v[i] == '}') {
				continue
			}

			if unicode.IsSpace(rune(v[i])) {
				newValue = append(newValue, []byte(vars[string(name)])...)
				newValue = append(newValue, v[i])
				name = nil
				atVar = false
				continue
			}

			name = append(name, v[i])
		}

		if atVar {
			newValue = append(newValue, []byte(vars[string(name)])...)
		}

		vars[k] = string(newValue)
	}
}

func isAtVar(v string, i int) (atVar bool) {
	atVar = true

	// Variable is escaped
	if i-1 >= 0 && v[i-1] == '\\' {
		atVar = false
	}

	// Variable is double escaped
	if i-2 > 0 && v[i-2] == '\\' {
		atVar = true
	}

	if i+1 < len(v) && (unicode.IsSpace(rune(v[i+1])) || v[i+1] == '"' || v[i+1] == '\'') {
		atVar = false
	}

	return
}

func nextVarIsDoubleEscaped(v string, i int) bool {
	return v[i] == '\\' && i+1 < len(v) && v[i+1] == '\\' && i+2 < len(v) && v[i+2] == '$'
}

func nextVarIsEscaped(v string, i int) bool {
	return v[i] == '\\' && i+1 < len(v) && v[i+1] == '$'
}
