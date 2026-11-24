package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type tokens struct {
	name  token
	value token
}

type Parser struct {
	vars         map[string]string
	stream       stream
	tokens       tokens
	currentToken tokenKind
	quote        rune
}

// quotePreviousOffset points to the character before the current quote.
const quotePreviousOffset = 2

// Parse consumes a reader and detects variables that it will add to the passed vars map.
func (p *Parser) Parse(r io.Reader, vars map[string]string) error {
	p.vars = vars
	p.stream = stream{reader: bufio.NewReader(r)}
	p.tokens = tokens{}
	p.currentToken = nameToken
	p.quote = 0

	for {
		err := p.stream.advance()

		switch {
		case errors.Is(err, io.EOF):
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

		case p.isAtCommentBegin():
			// Comment begins (`name=value #COMMENT`), save last variable.
			if p.atToken(valueToken) {
				p.saveVar()
			}

			p.setToken(commentToken)

		case p.isAtLineEnd():
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
			p.updateQuote()
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

func (p *Parser) isAtCommentBegin() bool {
	return p.quote == 0 && p.stream.isAtCommentBegin()
}

func (p *Parser) isAtLineEnd() bool {
	return p.quote == 0 && p.stream.isAtLineEnd()
}

func (p *Parser) updateQuote() {
	if p.stream.current != '"' && p.stream.current != '\'' {
		return
	}

	// An odd number of preceding backslashes means the quote is escaped.
	backslashes := 0
	for i := len(p.tokens.value) - quotePreviousOffset; i >= 0 && p.tokens.value[i] == '\\'; i-- {
		backslashes++
	}

	if backslashes%2 == 0 {
		switch p.quote {
		case 0:
			p.quote = p.stream.current
		case p.stream.current:
			p.quote = 0
		}
	}
}

func cleanVarValue(v []rune) string {
	return strings.Trim(strings.TrimSpace(string(v)), `"'`)
}

func New() *Parser {
	return &Parser{}
}
