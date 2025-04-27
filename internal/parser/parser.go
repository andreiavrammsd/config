package parser

import (
	"bufio"
	"io"
	"strings"
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

func (t token) String() string {
	return string(t.buffer)
}

func (t *token) append(r rune) {
	t.buffer = append(t.buffer, r)
}

type tokens struct {
	name    token
	value   token
	comment token
	current tokenKind
}

func (p *tokens) at(kind tokenKind) bool {
	return p.current == kind
}

func (p *tokens) set(kind tokenKind) {
	p.current = kind
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
			// Parsing done, save last variable.
			if p.tokens.at(valueToken) {
				p.saveVar()
			}
			return nil

		case err != nil:
			return err

		case p.stream.isAtCommentBegin():
			// Comment begins (`name=value #COMMENT`), save last variable.
			if p.tokens.at(valueToken) {
				p.saveVar()
			}
			p.tokens.set(commentToken)

		case p.stream.isAtLineEnd():
			// End of line reached, save last variable.
			if p.tokens.at(valueToken) {
				p.saveVar()
			}
			p.tokens.set(nameToken)

		case p.tokens.at(commentToken):
			// If inside a comment, just skip to next rune.
			continue

		case p.stream.isAtEqualSign():
			// If equal sign detected, start reading variable value (`name=VALUE #comment`).
			if p.tokens.at(valueToken) {
				p.tokens.value.append(p.stream.current)
			}
			p.tokens.set(valueToken)

		case p.tokens.at(nameToken):
			// Read variable name ignoring spaces (`NAME=value #comment`).
			if !p.stream.isAtSpace() {
				p.tokens.name.append(p.stream.current)
			}

		case p.tokens.at(valueToken):
			// Read variable value (`name=VALUE #comment`).
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
	p.tokens.set(nameToken)
}

func cleanVarValue(v []rune) string {
	return strings.Trim(strings.TrimSpace(string(v)), `"'`)
}

func New() *Parser {
	return &Parser{}
}
