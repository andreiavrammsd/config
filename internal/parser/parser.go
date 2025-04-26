package parser

import (
	"bufio"
	"io"
	"unicode"
)

type stream struct {
	reader  *bufio.Reader
	current rune
}

func (s *stream) advance() (err error) {
	s.current, _, err = s.reader.ReadRune()
	return
}

func (s *stream) isAtCommentBegin() bool {
	return s.current == '#'
}

func (s *stream) isAtLineEnd() bool {
	return s.current == '\n' || s.current == '\r'
}

func (s *stream) isAtEqualSign() bool {
	return s.current == '='
}

func (s *stream) isAtSpace() bool {
	return unicode.IsSpace(s.current)
}

type tokens struct {
	// the parser is in the variable name scope: `NAME=value #comment`
	atName bool

	// the actual variable name: `NAME=value #comment`
	name []rune

	// the parser is in the variable value scope: `name=VALUE #comment`
	atValue bool

	// the actual variable value: `name=VALUE #comment`
	value []rune

	// the parser is in the comment scope: `name=value # COMMENT`
	atComment bool
}

// appendToName adds a rune to the name array to form the variable name.
func (p *tokens) appendToName(r rune) {
	p.name = append(p.name, r)
}

// appendToValue adds a rune to the value array to form the variable value.
func (p *tokens) appendToValue(r rune) {
	p.value = append(p.value, r)
}

type Parser struct {
	vars   map[string]string
	stream stream
	tokens tokens
}

// Parse consumes a reader and detects variables that it will add to the passed vars map.
func (p *Parser) Parse(r io.Reader, vars map[string]string) error {
	p.stream = stream{reader: bufio.NewReader(r)}
	p.tokens = tokens{
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

// saveVar stores the variable name and its value,
// and sets tokens to start scanning for a new variable.
func (p *Parser) saveVar() {
	p.vars[string(p.tokens.name)] = cleanVarValue(p.tokens.value)
	p.tokens.name = nil
	p.tokens.value = nil
	p.tokens.atValue = false
}

func cleanVarValue(v []rune) string {
	return strings.Trim(strings.TrimSpace(string(v)), `"'`)
}

func New() *Parser {
	return &Parser{}
}
