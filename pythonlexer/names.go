package pythonlexer

import (
	"errors"
	"unicode"
)

var (
	ErrInvalidNameStart = errors.New("name starts with an invalid character")
)

func (l *Lexer) name() error {
	if !isIdentifierStart(l.file.Source[l.current]) {
		return ErrInvalidNameStart
	}

	l.current++

	for l.current < len(l.file.Source) && isIdentifierContinue(l.file.Source[l.current]) {
		l.current++
	}

	return nil
}

func isIdentifierStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentifierContinue(r rune) bool {
	return isIdentifierStart(r) || unicode.IsDigit(r)
}
