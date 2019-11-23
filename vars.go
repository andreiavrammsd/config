package config

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

func vars(r io.Reader) (map[string]string, error) {
	reader := bufio.NewReader(r)
	vars := make(map[string]string)

	var k, v []byte

	atK := true
	atV := false
	atComment := false

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				if atV {
					vars[string(bytes.TrimSpace(k))] = string(bytes.Trim(bytes.TrimSpace(v), `"'`))
				}
				break
			}

			return nil, fmt.Errorf("config: cannot read input (%s)", err)
		}

		if r == '#' {
			if atV {
				vars[string(bytes.TrimSpace(k))] = string(bytes.Trim(bytes.TrimSpace(v), `"'`))
			}

			k = nil
			v = nil
			atComment = true
			atK = false
			atV = false
			continue
		}

		if r == '\n' || r == '\r' {
			if atV {
				vars[string(bytes.TrimSpace(k))] = string(bytes.Trim(bytes.TrimSpace(v), `"'`))
			}

			k = nil
			v = nil
			atK = true
			atV = false

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
			if atV {
				v = append(v, byte(r))
			}
			atK = false
			atV = true
			continue
		}

		if atK {
			if unicode.IsSpace(r) {
				continue
			}
			k = append(k, byte(r))
			continue
		}

		if atV {
			v = append(v, byte(r))
		}
	}

	return vars, nil
}
