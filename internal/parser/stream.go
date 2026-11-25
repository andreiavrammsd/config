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

func (s *stream) reset(r io.Reader) {
	if s.reader == nil {
		s.reader = bufio.NewReader(r)
	} else {
		s.reader.Reset(r)
	}
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
