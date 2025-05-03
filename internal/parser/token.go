package parser

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
