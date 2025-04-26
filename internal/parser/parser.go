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

func (s *Stream) isAtCommentBegin() bool {
	return s.current == '#'
}

func (s *Stream) isAtLineEnd() bool {
	return s.current == '\n' || s.current == '\r'
}

func (s *Stream) isAtEqualSign() bool {
	return s.current == '='
}

func (s *Stream) isAtSpace() bool {
	return unicode.IsSpace(s.current)
}

type Tokens struct {
	// the parser is in the variable name scope: `NAME=value #comment`
	atName bool

	// the actual variable name: `NAME=value #comment`
	name []byte

	// the parser is in the variable value scope: `name=VALUE #comment`
	atValue bool

	// the actual variable value: `name=VALUE #comment`
	value []byte

	// the parser is in the comment scope: `name=value # COMMENT`
	atComment bool
}

// appendToName adds a rune to the name array to form the variable name.
func (p *Tokens) appendToName(r rune) {
	p.name = append(p.name, byte(r))
}

// appendToValue adds a rune to the value array to form the variable value.
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
		err := p.stream.advance()

		switch {
		case err == io.EOF:
			if p.tokens.atValue {
				p.saveVar()
			}
			return nil

		case err != nil:
			return err

		case p.stream.isAtCommentBegin():
			if p.tokens.atValue {
				p.saveVar()
			}
			p.tokens.atName = false
			p.tokens.atComment = true

		case p.stream.isAtLineEnd():
			if p.tokens.atValue {
				p.saveVar()
			}
			p.tokens.atName = true
			p.tokens.atComment = false

		case p.tokens.atComment:
			continue

		case p.stream.isAtEqualSign():
			if p.tokens.atValue {
				p.tokens.appendToValue(p.stream.current)
			}
			p.tokens.atName = false
			p.tokens.atValue = true

		case p.tokens.atName:
			if !p.stream.isAtSpace() {
				p.tokens.appendToName(p.stream.current)
			}

		case p.tokens.atValue:
			p.tokens.appendToValue(p.stream.current)
		}
	}
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
