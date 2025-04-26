package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

type Stream struct {
	reader *bufio.Reader
	r      rune
}

func (s *Stream) advance() (err error) {
	s.r, _, err = s.reader.ReadRune()
	return
}

func (s *Stream) isCommentBegin() bool {
	return s.r == '#'
}

func (s *Stream) isLineEnd() bool {
	return s.r == '\n' || s.r == '\r'
}

func (s *Stream) isEqualSign() bool {
	return s.r == '='
}

func (s *Stream) isSpace() bool {
	return unicode.IsSpace(s.r)
}

type Parser struct {
	vars map[string]string

	name  []byte
	value []byte

	atName    bool
	atValue   bool
	atComment bool

	stream Stream
}

func (p *Parser) Parse(r io.Reader, vars map[string]string) error {
	p.vars = vars

	p.name = nil
	p.value = nil
	p.startName()
	p.stopValue()
	p.stopComment()

	p.stream = Stream{reader: bufio.NewReader(r)}

	for {
		if err := p.stream.advance(); err != nil {
			if err == io.EOF {
				if p.isAtValue() {
					p.saveVar()
				}
				break
			}

			return err
		}

		if p.stream.isCommentBegin() {
			if p.isAtValue() {
				p.saveVar()
			}

			p.stopName()
			p.startComment()
			continue
		}

		if p.stream.isLineEnd() {
			if p.isAtValue() {
				p.saveVar()
			}

			p.startName()

			if p.isAtComment() {
				p.stopComment()
				continue
			}

			continue
		}

		if p.isAtComment() {
			continue
		}

		if p.stream.isEqualSign() {
			if p.isAtValue() {
				p.appendToValue()
			}
			p.stopName()
			p.startValue()
			continue
		}

		if p.isAtName() {
			if p.stream.isSpace() {
				continue
			}
			p.appendToName()
			continue
		}

		if p.isAtValue() {
			p.appendToValue()
		}
	}

	return nil
}

func (p *Parser) isAtComment() bool {
	return p.atComment
}

func (p *Parser) startComment() {
	p.atComment = true
}

func (p *Parser) stopComment() {
	p.atComment = false
}

func (p *Parser) isAtName() bool {
	return p.atName
}

func (p *Parser) startName() {
	p.atName = true
}

func (p *Parser) stopName() {
	p.atName = false
}

func (p *Parser) isAtValue() bool {
	return p.atValue
}

func (p *Parser) startValue() {
	p.atValue = true
}

func (p *Parser) stopValue() {
	p.atValue = false
}

func (p *Parser) appendToName() {
	p.name = append(p.name, byte(p.stream.r))
}

func (p *Parser) appendToValue() {
	p.value = append(p.value, byte(p.stream.r))
}

func (p *Parser) saveVar() {
	p.vars[string(p.name)] = varValue(p.value)
	p.name = nil
	p.value = nil
	p.atValue = false
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
