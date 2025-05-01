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

type token []rune

func (t *token) append(r rune) {
	*t = append(*t, r)
}

func (t *token) reset() {
	*t = nil
}

type tokens struct {
	name  token
	value token
}

type Parser struct {
	vars         map[string]string
	stream       stream
	tokens       tokens
	currentToken tokenKind
}

// Parse consumes a reader and detects variables that it will add to the passed vars map.
func (p *Parser) Parse(r io.Reader, vars map[string]string) error {
	p.vars = vars
	p.stream = stream{reader: bufio.NewReader(r)}
	p.tokens = tokens{}
	p.currentToken = nameToken

	for {
		err := p.stream.advance()

		switch {
		case err == io.EOF:
			// Parsing done, save last variable.
			if p.atToken(valueToken) {
				p.saveVar()
			}
			return nil

		case err != nil:
			return err

		case p.stream.isAtEqualSign() && !p.atToken(valueToken):
			// If equal sign detected and not already scanning variable value
			// (equal sign detected first time on line),
			// the variable value starts (`name=VALUE #comment`).
			p.setToken(valueToken)

		case p.stream.isAtCommentBegin():
			// Comment begins (`name=value #COMMENT`), save last variable.
			if p.atToken(valueToken) {
				p.saveVar()
			}
			p.setToken(commentToken)

		case p.stream.isAtLineEnd():
			// End of line reached, save last variable.
			if p.atToken(valueToken) {
				p.saveVar()
			}
			p.setToken(nameToken)

		case p.atToken(commentToken):
			// If inside a comment, just skip to next rune.

		case p.atToken(nameToken):
			// Read variable name ignoring spaces (`NAME=value #comment`).
			if !p.stream.isAtSpace() {
				p.tokens.name.append(p.stream.current)
			}

		case p.atToken(valueToken):
			// Read variable value (`name=VALUE #comment`).
			p.tokens.value.append(p.stream.current)
		}
	}
}

// saveVar stores the variable name and its value,
// and sets tokens to start scanning for a new variable.
func (p *Parser) saveVar() {
	if len(p.tokens.name) > 0 {
		p.vars[string(p.tokens.name)] = cleanVarValue(p.tokens.value)
	}

	p.tokens.name.reset()
	p.tokens.value.reset()
	p.setToken(nameToken)
}

func (p *Parser) atToken(kind tokenKind) bool {
	return p.currentToken == kind
}

func (p *Parser) setToken(kind tokenKind) {
	p.currentToken = kind
}

func cleanVarValue(v []rune) string {
	return strings.Trim(strings.TrimSpace(string(v)), `"'`)
}

func New() *Parser {
	return &Parser{}
}
