package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Stream struct {
	reader  *bufio.Reader
	current rune
}

func (s *Stream) advance() (err error) {
	s.current, _, err = s.reader.ReadRune()
	return
}

func (s *Stream) isCommentBegin() bool {
	return s.current == '#'
}

func (s *Stream) isLineEnd() bool {
	return s.current == '\n' || s.current == '\r'
}

func (s *Stream) isEqualSign() bool {
	return s.current == '='
}

func (s *Stream) isSpace() bool {
	return unicode.IsSpace(s.current)
}

type Tokens struct {
	atName bool
	name   []byte

	atValue bool
	value   []byte

	atComment bool
}

func (p *Tokens) appendToName(r rune) {
	p.name = append(p.name, byte(r))
}

func (p *Tokens) appendToValue(r rune) {
	p.value = append(p.value, byte(r))
}

type Parser struct {
	vars   map[string]string
	stream Stream
	tokens Tokens
}

func (p *Parser) Parse(r io.Reader, vars map[string]string) error {
	p.stream = Stream{reader: bufio.NewReader(r)}
	p.tokens = Tokens{
		atName:    true,
		name:      nil,
		atValue:   false,
		value:     nil,
		atComment: false,
	}
	p.vars = vars

	for {
		if err := p.stream.advance(); err != nil {
			if err == io.EOF {
				if p.tokens.atValue {
					p.saveVar()
				}
				break
			}

			return err
		}

		if p.stream.isCommentBegin() {
			if p.tokens.atValue {
				p.saveVar()
			}

			p.tokens.atName = false
			p.tokens.atComment = true
			continue
		}

		if p.stream.isLineEnd() {
			if p.tokens.atValue {
				p.saveVar()
			}

			p.tokens.atName = true

			if p.tokens.atComment {
				p.tokens.atComment = false
				continue
			}

			continue
		}

		if p.tokens.atComment {
			continue
		}

		if p.stream.isEqualSign() {
			if p.tokens.atValue {
				p.tokens.appendToValue(p.stream.current)
			}
			p.tokens.atName = false
			p.tokens.atValue = true
			continue
		}

		if p.tokens.atName {
			if p.stream.isSpace() {
				continue
			}
			p.tokens.appendToName(p.stream.current)
			continue
		}

		if p.tokens.atValue {
			p.tokens.appendToValue(p.stream.current)
		}
	}

	return nil
}

func (p *Parser) saveVar() {
	p.vars[string(p.tokens.name)] = varValue(p.tokens.value)
	p.tokens.name = nil
	p.tokens.value = nil
	p.tokens.atValue = false
}

func varValue(v []byte) string {
	return string(bytes.Trim(bytes.TrimSpace(v), `"'`))
}

func Interpolate(vars map[string]string) {
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

func New() *Parser {
	return &Parser{}
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
