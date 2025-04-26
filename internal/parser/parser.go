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

type tokenKind byte

const (
	// Parser is in the variable name scope: `NAME=value #comment`.
	nameToken tokenKind = iota

	// Parser is in the variable value scope: `name=VALUE #comment`.
	valueToken

	// Parser is in the comment scope: `name=value # COMMENT`.
	commentToken
)

type token struct {
	kind   tokenKind
	buffer []rune
}

func (t *token) append(r rune) {
	t.buffer = append(t.buffer, r)
}

func (t token) String() string {
	return string(t.buffer)
}

type tokens struct {
	name    token
	value   token
	comment token
	current tokenKind
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
		name:    token{nameToken, nil},
		value:   token{valueToken, nil},
		comment: token{commentToken, nil},
		current: nameToken,
	}
	p.vars = vars

	for {
		err := p.stream.advance()

		switch {
		case err == io.EOF:
			if p.tokens.current == valueToken {
				p.saveVar()
			}
			return nil

		case err != nil:
			return err

		case p.stream.isAtCommentBegin():
			if p.tokens.current == valueToken {
				p.saveVar()
			}
			p.tokens.current = commentToken

		case p.stream.isAtLineEnd():
			if p.tokens.current == valueToken {
				p.saveVar()
			}
			p.tokens.current = nameToken

		case p.tokens.current == p.tokens.comment.kind:
			continue

		case p.stream.isAtEqualSign():
			if p.tokens.current == valueToken {
				p.tokens.value.append(p.stream.current)
			}
			p.tokens.current = valueToken

		case p.tokens.current == nameToken:
			if !p.stream.isAtSpace() {
				p.tokens.name.append(p.stream.current)
			}

		case p.tokens.current == valueToken:
			p.tokens.value.append(p.stream.current)
		}
	}
}

// saveVar stores the variable name and its value,
// and sets tokens to start scanning for a new variable.
func (p *Parser) saveVar() {
	p.vars[p.tokens.name.String()] = cleanVarValue(p.tokens.value.buffer)
	p.tokens.name.buffer = nil
	p.tokens.value.buffer = nil
	p.tokens.current = nameToken
}

func cleanVarValue(v []rune) string {
	return strings.Trim(strings.TrimSpace(string(v)), `"'`)
}

func New() *Parser {
	return &Parser{}
}
