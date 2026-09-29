package pythonlexer

import (
	"errors"
	"slices"
	"unicode"
)

var (
	ErrUnexpectedEOF  = errors.New("unexpected end of file")
	ErrUnexpectedChar = errors.New("unexpected character")
)

func (l *Lexer) expected(expected string) error {

	if len(expected)+l.current > len(l.file.Source) {
		return ErrUnexpectedEOF
	}

	for i, letter := range expected {
		if l.file.Source[l.current+i] != letter {
			return ErrUnexpectedChar
		}
	}

	l.current += len(expected)

	return nil
}

func (l *Lexer) expectedIgnoreCase(expected string) error {

	if len(expected)+l.current > len(l.file.Source) {
		return ErrUnexpectedEOF
	}

	for i, letter := range expected {
		if unicode.ToLower(l.file.Source[l.current+i]) != unicode.ToLower(letter) {
			return ErrUnexpectedChar
		}
	}

	l.current += len(expected)

	return nil
}

func (l *Lexer) goUntil(terminators []rune) {

	for l.current < len(l.file.Source) &&
		!slices.Contains(terminators, l.file.Source[l.current]) {
		l.current++
	}
}
